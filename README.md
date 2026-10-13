# SIS Project — Exam Score Analyzer

Checkpoint 1: Core Go. A small synchronous, in-memory Go application.
Only the Go standard library is used.

## What the application does

A task is a list of exam scores (from 0 to 100) that must be analyzed.

| Part | In this project |
|---|---|
| Input | a list of scores, for example [85, 42, 100, 67, 50] |
| Processing | count the scores, find min and max, calculate the average, count how many students passed (score >= 50) |
| Result | count: 5, min: 42, max: 100, average: 68.80, passed: 4 |

Each task has a unique integer ID (1, 2, 3, ...), its input scores, a status and a result.
All tasks are stored in memory in a map.

## Validation rules

- the list of scores must not be empty;
- every score must be from 0 to 100.

Invalid input is rejected with an error, and the task is not created.

## Task lifecycle

```
pending  ──►  processing  ──►  completed
```

| Status | Meaning |
|---|---|
| pending | the task is created, but not processed yet |
| processing | the task is being processed right now |
| completed | the task is processed, the result is ready |

Allowed transitions:

| From | To | When |
|---|---|---|
| pending | processing | processing of the task starts |
| processing | completed | the result is calculated and saved in the task |

Any other transition is invalid and returns ErrInvalidTransition, and the status does not change.
For example, a completed task cannot be processed again.

## Errors

| Error | When |
|---|---|
| ErrEmptyScores | the list of scores is empty |
| ErrScoreOutOfRange | a score is less than 0 or greater than 100 |
| ErrTaskNotFound | there is no task with this ID |
| ErrInvalidTransition | the status change is not allowed (for example, processing a completed task) |

## Project structure

```
sis-project/
├── go.mod
├── main.go            entry point, shows how the application works
├── README.md
└── task/
    ├── errors.go      all application errors
    ├── task.go        Task, Status, Result, input validation, status transitions
    ├── processor.go   Processor interface and ScoreProcessor (the processing logic)
    ├── manager.go     Manager: creates, stores, finds and processes tasks
    └── task_test.go   unit tests
```

## How to run

```
go run .
```

## How to test

```
go test ./...
```

## Checks before submission

```
gofmt -l .
go vet ./...
go test ./...
```

gofmt -l . must print nothing.
