package texttask

type State string

const (
	Pending    State = "pending"
	Queued     State = "queued"
	Processing State = "processing"
	Completed  State = "completed"
)

type Result struct {
	WordCount      int
	CharacterCount int
	LineCount      int
}

type Task struct {
	ID     int
	Text   string
	State  State
	Result Result
}
