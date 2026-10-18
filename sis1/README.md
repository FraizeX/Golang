# Text Analyzer — SIS Checkpoint 1

A small synchronous Go application that creates text analysis tasks, stores them
in memory, and calculates word, Unicode character, and line counts.


## Requirements and scope

- Go **1.27**; verified locally with Go 1.27.1.
- Standard library only; no additional installations or dependencies.
- Tasks and results exist only while the program runs. IDs start at 1 for each service instance.
- Processing is synchronous. A service instance is intended for sequential use.
- The executable runs a demonstration with two inputs and expected error cases.

## Run

Open a terminal in this directory, which contains this project's `go.mod`:

```powershell
go run .
```

From the parent `golang` directory, first run `cd sis1`. This folder is an
independent Go module and can be published as its own repository.

To try other texts, change the `inputs` slice in `main.go`, then run the program again.

Selected output:

```text
Processed task #1: state=completed, words=3, characters=14, lines=1
Processed task #2: state=completed, words=4, characters=23, lines=2
Expected lookup error: task not found: id=999
Expected state error: invalid task state: task 1 is completed, expected pending
```

The `Expected ... error` messages demonstrate handled errors and are part of a
successful run.

## Counting rules

| Metric | Definition | Example |
| --- | --- | --- |
| Words | Nonempty tokens separated by Unicode whitespace, using `strings.Fields` | `hello,world!` is one word |
| Characters | Unicode code points (runes), including whitespace and punctuation | `Привет мир` has 10 runes |
| Lines | Number of `\n` separators plus one for nonempty text | `Go\n` has two lines |

Characters are runes, not bytes or grapheme clusters. A combined emoji may contain
multiple runes. Windows `\r\n` counts as one line separator and two characters.
The original input is preserved. Empty, whitespace-only, and invalid UTF-8 inputs
are rejected when creating a task. The standalone analyzer returns a zero result
for an empty string; task creation applies the validation rules.

## Task lifecycle

```text
Create(valid input) -> pending -> processing -> completed
```

Only a `pending` task can be processed. The service stores `processing` immediately
before invoking the processor, then saves the result and `completed` state.
Reprocessing a `processing` or `completed` task returns `ErrInvalidState`.
Because processing is synchronous, a regular call returns after completion;
`TestStateDuringProcessing` observes the intermediate state inside the processor.
The analyzer has no expected processing failures after input validation.

## Project structure

```text
sis1/
  go.mod
  .gitignore
  main.go
  internal/texttask/
    task.go
    processor.go
    service.go
    processor_test.go
    service_test.go
  README.md
  EXPLANATION.md
  GITHUB_GUIDE.md
```

| File | Responsibility |
| --- | --- |
| `main.go` | Demonstrates creation, processing, retrieval, and handled errors |
| `task.go` | Defines `Task`, `Result`, and lifecycle states |
| `processor.go` | Defines the `Processor` interface and `TextAnalyzer` implementation |
| `service.go` | Validates input, generates IDs, stores tasks, controls transitions |
| `*_test.go` | Checks analysis rules, task behavior, errors, and independent results |

The service stores task **values** in `map[int]Task`. Returned tasks are copies;
changing a returned task cannot change stored data. The `Processor` interface
separates analysis from task storage and permits an observing processor in tests.

## Checks

Run from this directory:

```powershell
gofmt -l .
go vet ./...
go test ./...
go build -o ./bin/text-analyzer.exe .
```

`gofmt -l .` should print nothing; the other commands should finish successfully.
For verbose tests or a coverage report:

```powershell
go test -v ./...
go test -cover ./...
```

## Checkpoint coverage

| Behavior | Implementation / test |
| --- | --- |
| Create, validate, reject invalid input | `Create`, `TestCreateAndRetrieveTask`, `TestCreateRejectsInvalidText` |
| Unique IDs and in-memory storage | `nextID`, `tasks`, `TestIndependentTasks` |
| Retrieve and report missing tasks | `Get`, `TestMissingTask` |
| Calculate and store correct results | `Process`, `TestProcessTask`, `TestTextAnalyzer` |
| Track lifecycle and reject invalid operations | `TestStateDuringProcessing`, `TestProcessRejectsCompletedTask` |
| Keep task data independent | `TestIndependentTasks`, `TestReturnedTaskIsACopy` |
| Use an interface and receiver methods | `Processor`, `TextAnalyzer.Process`, `*Service` methods |
| Standard-library and table-driven tests | `processor_test.go`, `service_test.go` |


