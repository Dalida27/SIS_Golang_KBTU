package task

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
)

type Result struct {
	Count   int
	Min     int
	Max     int
	Average float64
	Passed  int
}

type Task struct {
	ID     int
	Scores []int
	Status Status
	Result Result
}

func validateScores(scores []int) error {
	if len(scores) == 0 {
		return ErrEmptyScores
	}

	for _, score := range scores {
		if score < 0 || score > 100 {
			return ErrScoreOutOfRange
		}
	}

	return nil
}

func (t *Task) changeStatus(newStatus Status) error {
	if t.Status == StatusPending && newStatus == StatusProcessing {
		t.Status = newStatus
		return nil
	}

	if t.Status == StatusProcessing && newStatus == StatusCompleted {
		t.Status = newStatus
		return nil
	}

	return ErrInvalidTransition
}
