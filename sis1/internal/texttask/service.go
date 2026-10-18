package texttask

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidText  = errors.New("invalid text")
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidState = errors.New("invalid task state")
	ErrNilProcessor = errors.New("processor is required")
)

type Service struct {
	tasks     map[int]Task
	nextID    int
	processor Processor
}

func NewService(processor Processor) (*Service, error) {
	if processor == nil {
		return nil, ErrNilProcessor
	}

	return &Service{
		tasks:     make(map[int]Task),
		nextID:    1,
		processor: processor,
	}, nil
}

func (service *Service) Create(text string) (Task, error) {
	if !utf8.ValidString(text) {
		return Task{}, fmt.Errorf("%w: text must be valid UTF-8", ErrInvalidText)
	}
	if strings.TrimSpace(text) == "" {
		return Task{}, fmt.Errorf("%w: text must contain a non-whitespace character", ErrInvalidText)
	}

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
	task, exists := service.tasks[id]
	if !exists {
		return Task{}, fmt.Errorf("%w: id=%d", ErrTaskNotFound, id)
	}

	return task, nil
}

func (service *Service) Process(id int) (Task, error) {
	task, err := service.Get(id)
	if err != nil {
		return Task{}, err
	}
	if task.State != Pending {
		return Task{}, fmt.Errorf("%w: task %d is %s, expected %s", ErrInvalidState, id, task.State, Pending)
	}

	task.State = Processing
	service.tasks[id] = task

	task.Result = service.processor.Process(task.Text)
	task.State = Completed
	service.tasks[id] = task

	return task, nil
}
