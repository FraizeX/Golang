package texttask

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"
)

var (
	ErrInvalidText  = errors.New("invalid text")
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidState = errors.New("invalid task state")
	ErrNilProcessor = errors.New("processor is required")
	ErrQueueFull    = errors.New("task queue is full")
	ErrStopped      = errors.New("worker pool is stopped")
)

const (
	WorkerCount   = 2
	QueueCapacity = 3
)

type job struct {
	id   int
	done chan Task // Only the synchronous Process caller needs a reply.
}

type Service struct {
	mu        sync.RWMutex
	tasks     map[int]Task
	nextID    int
	processor Processor
	queue     chan job
	stopped   bool
	workers   sync.WaitGroup
}

func NewService(processor Processor) (*Service, error) {
	if processor == nil {
		return nil, ErrNilProcessor
	}

	service := &Service{
		tasks:     make(map[int]Task),
		nextID:    1,
		processor: processor,
		queue:     make(chan job, QueueCapacity),
	}
	service.workers.Add(WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		go service.worker()
	}
	return service, nil
}

func (service *Service) Create(text string) (Task, error) {
	if !utf8.ValidString(text) {
		return Task{}, fmt.Errorf("%w: text must be valid UTF-8", ErrInvalidText)
	}
	if strings.TrimSpace(text) == "" {
		return Task{}, fmt.Errorf("%w: text must contain a non-whitespace character", ErrInvalidText)
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	task := Task{
		ID:    service.nextID,
		Text:  text,
		State: Pending,
	}
	service.tasks[task.ID] = task
	service.nextID++

	return task, nil
}

func (service *Service) Get(id int) (Task, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()
	task, exists := service.tasks[id]
	if !exists {
		return Task{}, fmt.Errorf("%w: id=%d", ErrTaskNotFound, id)
	}

	return task, nil
}

func (service *Service) Process(id int) (Task, error) {
	done := make(chan Task, 1)
	if err := service.enqueue(job{id: id, done: done}); err != nil {
		return Task{}, err
	}
	return <-done, nil
}

// Submit returns once the job is accepted. Get can observe its progress.
func (service *Service) Submit(id int) error {
	return service.enqueue(job{id: id})
}

func (service *Service) enqueue(j job) error {
	service.mu.Lock()
	defer service.mu.Unlock()
	if service.stopped {
		return ErrStopped
	}
	task, exists := service.tasks[j.id]
	if !exists {
		return fmt.Errorf("%w: id=%d", ErrTaskNotFound, j.id)
	}
	if task.State != Pending {
		return fmt.Errorf("%w: task %d is %s, expected %s", ErrInvalidState, j.id, task.State, Pending)
	}
	select {
	case service.queue <- j:
		task.State = Queued
		service.tasks[j.id] = task
		return nil
	default:
		return ErrQueueFull
	}
}

func (service *Service) worker() {
	defer service.workers.Done()
	for j := range service.queue {
		service.mu.Lock()
		task := service.tasks[j.id]
		task.State = Processing
		service.tasks[j.id] = task
		service.mu.Unlock()

		task.Result = service.processor.Process(task.Text)
		task.State = Completed

		service.mu.Lock()
		service.tasks[j.id] = task
		service.mu.Unlock()
		if j.done != nil {
			j.done <- task
		}
	}
}

// Shutdown stops admission, drains accepted work, and joins the workers.
func (service *Service) Shutdown() {
	service.mu.Lock()
	if !service.stopped {
		service.stopped = true
		close(service.queue)
	}
	service.mu.Unlock()
	service.workers.Wait()
}
