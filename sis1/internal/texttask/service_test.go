package texttask

import (
	"errors"
	"testing"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	service, err := NewService(TextAnalyzer{})
	if err != nil {
		t.Fatalf("NewService() failed: %v", err)
	}
	return service
}

func TestNewServiceRejectsNilProcessor(t *testing.T) {
	service, err := NewService(nil)
	if service != nil || !errors.Is(err, ErrNilProcessor) {
		t.Fatalf("NewService(nil) = (%v, %v), want (nil, ErrNilProcessor)", service, err)
	}
}

func TestCreateAndRetrieveTask(t *testing.T) {
	service := newTestService(t)
	created, err := service.Create("  Hello Go  ")
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	want := Task{ID: 1, Text: "  Hello Go  ", State: Pending}
	if created != want {
		t.Fatalf("Create() = %+v, want %+v", created, want)
	}

	stored, err := service.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if stored != created {
		t.Fatalf("Get() = %+v, want %+v", stored, created)
	}
}

func TestCreateRejectsInvalidText(t *testing.T) {
	testCases := []struct {
		name string
		text string
	}{
		{name: "empty", text: ""},
		{name: "spaces", text: "   "},
		{name: "tabs and newlines", text: "\t\r\n"},
		{name: "Unicode whitespace", text: "\u00a0\u2003"},
		{name: "invalid UTF-8", text: string([]byte{0xff})},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := newTestService(t)
			rejected, err := service.Create(testCase.text)
			if !errors.Is(err, ErrInvalidText) {
				t.Fatalf("Create(%q) error = %v, want ErrInvalidText", testCase.text, err)
			}
			if rejected != (Task{}) {
				t.Fatalf("invalid input returned a task: %+v", rejected)
			}
			if _, err := service.Get(1); !errors.Is(err, ErrTaskNotFound) {
				t.Fatalf("invalid input was stored: Get(1) error = %v", err)
			}
			valid, err := service.Create("Go")
			if err != nil || valid.ID != 1 {
				t.Fatalf("valid task after rejection = (%+v, %v), want ID 1", valid, err)
			}
		})
	}
}

func TestMissingTask(t *testing.T) {
	service := newTestService(t)
	testCases := []struct {
		name      string
		operation func(int) (Task, error)
	}{
		{name: "retrieve", operation: service.Get},
		{name: "process", operation: service.Process},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := testCase.operation(999)
			if !errors.Is(err, ErrTaskNotFound) || got != (Task{}) {
				t.Fatalf("operation(999) = (%+v, %v), want zero task and ErrTaskNotFound", got, err)
			}
		})
	}
}

func TestProcessTask(t *testing.T) {
	service := newTestService(t)
	created, err := service.Create("Hello Go world")
	if err != nil {
		t.Fatal(err)
	}
	completed, err := service.Process(created.ID)
	if err != nil {
		t.Fatalf("Process() failed: %v", err)
	}
	want := Task{
		ID:     created.ID,
		Text:   created.Text,
		State:  Completed,
		Result: Result{WordCount: 3, CharacterCount: 14, LineCount: 1},
	}
	if completed != want {
		t.Fatalf("Process() = %+v, want %+v", completed, want)
	}
	stored, err := service.Get(created.ID)
	if err != nil || stored != want {
		t.Fatalf("stored task = (%+v, %v), want %+v", stored, err, want)
	}
}

func TestProcessRejectsCompletedTask(t *testing.T) {
	service := newTestService(t)
	created, err := service.Create("Go")
	if err != nil {
		t.Fatal(err)
	}
	completed, err := service.Process(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Process(created.ID); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("second Process() error = %v, want ErrInvalidState", err)
	}
	stored, err := service.Get(created.ID)
	if err != nil || stored != completed {
		t.Fatalf("rejected operation changed task: (%+v, %v)", stored, err)
	}
}

type observingProcessor struct {
	observe func(string)
}

func (processor observingProcessor) Process(text string) Result {
	processor.observe(text)
	return (TextAnalyzer{}).Process(text)
}

func TestStateDuringProcessing(t *testing.T) {
	var service *Service
	var taskID int
	observed := false
	processor := observingProcessor{
		observe: func(text string) {
			observed = true
			stored, err := service.Get(taskID)
			if err != nil || stored.State != Processing || stored.Result != (Result{}) || text != "Go" {
				t.Fatalf("task during processing = (%+v, %v), input=%q", stored, err, text)
			}
			if _, err := service.Process(taskID); !errors.Is(err, ErrInvalidState) {
				t.Fatalf("processing a busy task: error = %v, want ErrInvalidState", err)
			}
		},
	}
	var err error
	service, err = NewService(processor)
	if err != nil {
		t.Fatal(err)
	}
	created, err := service.Create("Go")
	if err != nil {
		t.Fatal(err)
	}
	taskID = created.ID
	completed, err := service.Process(taskID)
	if err != nil || !observed || completed.State != Completed {
		t.Fatalf("Process() = (%+v, %v), observed=%t", completed, err, observed)
	}
}

func TestIndependentTasks(t *testing.T) {
	service := newTestService(t)
	first, err := service.Create("Hello Go world")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Create("Привет\nмир")
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == second.ID {
		t.Fatalf("duplicate ID: %d", first.ID)
	}
	firstCompleted, err := service.Process(first.ID)
	if err != nil {
		t.Fatal(err)
	}
	secondPending, err := service.Get(second.ID)
	if err != nil || secondPending != second {
		t.Fatalf("processing first task changed second: (%+v, %v)", secondPending, err)
	}
	secondCompleted, err := service.Process(second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if secondCompleted.State != Completed || secondCompleted.Text != second.Text || secondCompleted.Result != (Result{WordCount: 2, CharacterCount: 10, LineCount: 2}) {
		t.Fatalf("unexpected second result: %+v", secondCompleted)
	}
	storedFirst, err := service.Get(first.ID)
	if err != nil || storedFirst != firstCompleted {
		t.Fatalf("processing second task changed first: (%+v, %v)", storedFirst, err)
	}
}

func TestReturnedTaskIsACopy(t *testing.T) {
	service := newTestService(t)
	created, err := service.Create("Go")
	if err != nil {
		t.Fatal(err)
	}
	id := created.ID
	created.ID = 999
	created.Text = "changed"
	created.State = Completed
	stored, err := service.Get(id)
	if err != nil || stored != (Task{ID: id, Text: "Go", State: Pending}) {
		t.Fatalf("modifying Create() result changed storage: (%+v, %v)", stored, err)
	}
	completed, err := service.Process(id)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := service.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	snapshot.Result.WordCount = 999
	completed.Text = "changed again"
	stored, err = service.Get(id)
	if err != nil || stored.Text != "Go" || stored.Result != (Result{WordCount: 1, CharacterCount: 2, LineCount: 1}) {
		t.Fatalf("modifying Get()/Process() result changed storage: (%+v, %v)", stored, err)
	}
}
