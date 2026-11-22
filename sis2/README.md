# Text Analyzer · SIS Checkpoint 2

Text Analyzer creates in-memory tasks and counts words, Unicode code points and lines. This checkpoint extends the original `sis1` module with an asynchronous worker pool. It uses only the Go standard library.

## Run and check

From the repository root:

```sh
cd sis2
go run .
gofmt -l .
go vet ./...
go test ./...
go test -race ./...
```

The module requires Go 1.27.0 or later. The race detector on Windows/amd64 requires cgo and a compatible MinGW-w64 C compiler on PATH. The demo submits two texts, waits for the workers, then prints one result per ID. For the first input, `Hello Go world`, expect 3 words, 14 characters and 1 line. The final output order follows the input IDs. Workers can finish in either order.

## Counting rules from Checkpoint 1

| Metric | Rule |
| --- | --- |
| Words | `strings.Fields` splits on Unicode whitespace. Punctuation stays in its token. |
| Characters | `utf8.RuneCountInString` counts Unicode code points, including whitespace and punctuation. |
| Lines | For nonempty text, count newline characters (`\n`) and add one. A trailing newline introduces an empty last line. |

`Create` rejects invalid UTF-8 and whitespace-only text. It preserves the original input and assigns an integer ID starting at 1. Tasks exist for the lifetime of the service.

## API and states

| Operation | Behavior |
| --- | --- |
| `NewService(TextAnalyzer{})` | Starts two reusable worker goroutines and allocates a buffered queue for three waiting jobs. The caller owns the service and calls `Shutdown` when finished. |
| `Create(text)` | Stores a new task in `pending`. Creation alone does not schedule work. |
| `Submit(id)` | Enqueues a pending task and returns on admission. A worker may already have started or finished when the caller resumes. |
| `Get(id)` | Returns a value snapshot of the task and its result at the time of the read. |
| `Process(id)` | The synchronous API from Checkpoint 1 submits through the same queue, then waits on a private reply channel for that job. It uses the same admission errors and concurrency limit as `Submit`. |
| `Shutdown()` | Stops admission, drains accepted jobs, and waits for both workers. Repeated and concurrent calls are safe. `Get` and `Create` still work after shutdown. |

The normal lifecycle is **pending → queued → processing → completed**. Only pending tasks can be submitted. A missing ID returns `ErrTaskNotFound`. A repeat submission returns `ErrInvalidState`. When all three queue slots are occupied and no worker can receive immediately, submission returns `ErrQueueFull`. The task stays pending and may be submitted again. Once shutdown begins closing admission, a submission returns `ErrStopped`. The stock analyzer has no processing error path after validation.

## Concurrency design

The queue carries a task ID and an optional reply channel. A worker receives a job, sets its state to processing, runs the analyzer on a local task copy, then publishes the result and completed state together. The analyzer is stateless and safe to call from both workers. A replacement processor must also support concurrent independent calls.

One `sync.RWMutex` protects the task map, ID allocation, task transitions and the stopped flag. `Get` uses the read lock. Submission holds the write lock across validation, the nonblocking send and recording the queued state. A worker can receive the job immediately, but it needs that same lock before changing state. Keeping the analysis outside the lock lets both workers compute and other callers read.

Only the service closes the queue. `Shutdown` sets the flag and closes the channel under the same lock used for admission, releases the lock, then waits on a `sync.WaitGroup`. Workers range over the channel until all buffered jobs have been processed. The wait group counts workers for their entire lifetime.

## Source and tests

| File | Responsibility |
| --- | --- |
| `main.go` | Runnable example of submission, retrieval, errors and shutdown |
| `internal/texttask/task.go`, `processor.go` | Data types and the original analyzer |
| `internal/texttask/service.go` | In-memory store, admission, worker loop and shutdown |
| `internal/texttask/processor_test.go`, `service_test.go` | Original analysis, validation, state and snapshot tests |
| `internal/texttask/pool_test.go` | Overlap, full queue and retry, duplicate submissions, concurrent access and shutdown |

To watch the overlap test: `go test -run TestWorkerPoolOverlapAndQueue -v ./internal/texttask`. It holds the first two tasks on a channel until both workers have entered the processor. A separate test uses `testing/synctest` to observe shutdown with both workers busy and a third job still queued.

