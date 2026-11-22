package texttask

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func receive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for a goroutine")
		var zero T
		return zero
	}
}

func createTask(t *testing.T, service *Service, text string) Task {
	t.Helper()
	task, err := service.Create(text)
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func checkTask(t *testing.T, service *Service, original Task, state State, result Result) {
	t.Helper()
	want := original
	want.State, want.Result = state, result
	got, err := service.Get(original.ID)
	if err != nil || got != want {
		t.Fatalf("Get(%d) = (%+v, %v), want %+v", original.ID, got, err, want)
	}
}

func TestWorkerPoolOverlapAndQueue(t *testing.T) {
	// A worker reports entry, then waits. Two entry reports prove overlap.
	started := make(chan string, 6)
	gate := make(chan struct{})
	release := sync.OnceFunc(func() { close(gate) })
	service, err := NewService(observingProcessor{observe: func(text string) {
		started <- text
		<-gate
	}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(service.Shutdown)
	t.Cleanup(release)

	texts := []string{"Go", "Go\nGo", "a b", "abc", "x\ny\n", "last task"}
	want := []Result{{1, 2, 1}, {2, 5, 2}, {2, 3, 1}, {1, 3, 1}, {2, 4, 3}, {2, 9, 1}}
	tasks := make([]Task, len(texts))
	for i, text := range texts {
		tasks[i] = createTask(t, service, text)
	}
	for _, task := range tasks[:2] {
		returned := make(chan error, 1)
		go func() { returned <- service.Submit(task.ID) }()
		if err := receive(t, returned); err != nil {
			t.Fatal(err)
		}
	}
	first, second := receive(t, started), receive(t, started)
	running := map[string]bool{first: true, second: true}
	if !running[texts[0]] || !running[texts[1]] {
		t.Fatalf("unexpected overlapping inputs: %q and %q", first, second)
	}
	for _, task := range tasks[:2] {
		checkTask(t, service, task, Processing, Result{})
		if err := service.Submit(task.ID); !errors.Is(err, ErrInvalidState) {
			t.Fatalf("resubmit processing task: %v", err)
		}
	}
	for _, task := range tasks[2:5] {
		if err := service.Submit(task.ID); err != nil {
			t.Fatal(err)
		}
		checkTask(t, service, task, Queued, Result{})
	}
	if err := service.Submit(tasks[2].ID); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("resubmit queued task: %v", err)
	}
	if err := service.Submit(tasks[5].ID); !errors.Is(err, ErrQueueFull) {
		t.Fatalf("submit to full queue: %v", err)
	}
	checkTask(t, service, tasks[5], Pending, Result{})

	release()
	receive(t, started) // A queued job now occupies a worker, freeing a slot.
	if err := service.Submit(tasks[5].ID); err != nil {
		t.Fatalf("retry: %v", err)
	}
	service.Shutdown()
	for i, task := range tasks {
		checkTask(t, service, task, Completed, want[i])
	}
}

func TestShutdownDrainsActiveAndQueuedTasks(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		gate := make(chan struct{})
		service, err := NewService(observingProcessor{observe: func(string) { <-gate }})
		if err != nil {
			t.Fatal(err)
		}
		release := sync.OnceFunc(func() { close(gate) })
		defer service.Shutdown()
		defer release()

		var accepted []Task
		for _, text := range []string{"a", "b", "c"} {
			task := createTask(t, service, text)
			if err := service.Submit(task.ID); err != nil {
				t.Fatal(err)
			}
			accepted = append(accepted, task)
		}
		synctest.Wait() // Both workers are held inside the processor.
		checkTask(t, service, accepted[2], Queued, Result{})
		returned := make(chan struct{})
		go func() {
			service.Shutdown()
			close(returned)
		}()
		synctest.Wait() // Shutdown has closed admission and is waiting for workers.
		fresh := createTask(t, service, "new")
		if err := service.Submit(fresh.ID); !errors.Is(err, ErrStopped) {
			t.Fatalf("submission during shutdown: %v", err)
		}
		select {
		case <-returned:
			t.Fatal("Shutdown returned with unfinished work")
		default:
		}
		release()
		<-returned
		for _, task := range accepted {
			checkTask(t, service, task, Completed, Result{1, 1, 1})
		}
	})
}

func TestSimultaneousSubmissionsOfOneTask(t *testing.T) {
	var calls atomic.Int32
	service, err := NewService(observingProcessor{observe: func(string) { calls.Add(1) }})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(service.Shutdown)
	task := createTask(t, service, "Go")
	if err := service.Submit(999); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("missing task: %v", err)
	}
	const callers = 32
	begin := make(chan struct{})
	outcomes := make(chan error, callers)
	for i := 0; i < callers; i++ {
		go func() {
			<-begin
			outcomes <- service.Submit(task.ID)
		}()
	}
	close(begin)
	accepted := 0
	for i := 0; i < callers; i++ {
		err := receive(t, outcomes)
		if err == nil {
			accepted++
		} else if !errors.Is(err, ErrInvalidState) {
			t.Fatalf("duplicate submission: %v", err)
		}
	}
	service.Shutdown()
	if accepted != 1 || calls.Load() != 1 {
		t.Fatalf("accepted %d times, processed %d times", accepted, calls.Load())
	}
	checkTask(t, service, task, Completed, Result{1, 2, 1})
}

func TestConcurrentCreateReadSubmitAndShutdown(t *testing.T) {
	service := newTestService(t)
	seed := createTask(t, service, "already accepted")
	if err := service.Submit(seed.ID); err != nil {
		t.Fatal(err)
	}
	const clients = 48
	type outcome struct {
		task Task
		err  error
	}
	outcomes := make(chan outcome, clients)
	begin := make(chan struct{})
	for i := 0; i < clients; i++ {
		go func() {
			<-begin
			task, err := service.Create(fmt.Sprintf("client %d", i))
			if err != nil {
				outcomes <- outcome{err: err}
				return
			}
			_, err = service.Get(seed.ID)
			if err == nil {
				err = service.Submit(task.ID)
			}
			outcomes <- outcome{task, err}
		}()
	}
	stopped := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		go func() {
			<-begin
			service.Shutdown()
			stopped <- struct{}{}
		}()
	}
	close(begin)
	receive(t, stopped)
	receive(t, stopped)
	checkTask(t, service, seed, Completed, Result{2, 16, 1})
	seen := map[int]bool{seed.ID: true}
	for i := 0; i < clients; i++ {
		o := receive(t, outcomes)
		if o.task.ID == 0 || seen[o.task.ID] {
			t.Fatalf("invalid or duplicate ID %d (err=%v)", o.task.ID, o.err)
		}
		seen[o.task.ID] = true
		switch {
		case o.err == nil:
			checkTask(t, service, o.task, Completed, (TextAnalyzer{}).Process(o.task.Text))
		case errors.Is(o.err, ErrStopped), errors.Is(o.err, ErrQueueFull):
			checkTask(t, service, o.task, Pending, Result{})
		default:
			t.Fatalf("unexpected submission error: %v", o.err)
		}
	}
	fresh := createTask(t, service, "after shutdown")
	if err := service.Submit(fresh.ID); !errors.Is(err, ErrStopped) {
		t.Fatalf("submit after shutdown: %v", err)
	}
	checkTask(t, service, fresh, Pending, Result{})
}
