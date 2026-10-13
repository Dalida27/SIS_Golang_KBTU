package task

import (
	"errors"
	"testing"
)

func newTestManager() *Manager {
	return NewManager(ScoreProcessor{PassScore: 50})
}

func TestCreateTask_Success(t *testing.T) {
	m := newTestManager()

	created, err := m.CreateTask([]int{70, 80, 90})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}
	if created.Status != StatusPending {
		t.Errorf("expected status %s, got %s", StatusPending, created.Status)
	}
	if len(created.Scores) != 3 {
		t.Errorf("expected 3 scores, got %d", len(created.Scores))
	}
}

func TestCreateTask_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		scores  []int
		wantErr error
	}{
		{"empty list", []int{}, ErrEmptyScores},
		{"nil list", nil, ErrEmptyScores},
		{"negative score", []int{50, -5}, ErrScoreOutOfRange},
		{"score above 100", []int{101, 70}, ErrScoreOutOfRange},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newTestManager()

			created, err := m.CreateTask(tc.scores)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}
			if created != nil {
				t.Errorf("expected no task, got task with ID %d", created.ID)
			}

			_, err = m.GetTask(1)
			if !errors.Is(err, ErrTaskNotFound) {
				t.Errorf("invalid task was saved in memory")
			}
		})
	}
}

func TestCreateTask_CopiesScores(t *testing.T) {
	m := newTestManager()

	scores := []int{60, 70}
	created, err := m.CreateTask(scores)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	scores[0] = 0

	if created.Scores[0] != 60 {
		t.Errorf("task scores were changed from outside: got %v", created.Scores)
	}
}

func TestGetTask_Success(t *testing.T) {
	m := newTestManager()

	created, err := m.CreateTask([]int{60, 75})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := m.GetTask(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
	if len(got.Scores) != 2 || got.Scores[0] != 60 || got.Scores[1] != 75 {
		t.Errorf("expected scores [60 75], got %v", got.Scores)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	m := newTestManager()

	got, err := m.GetTask(42)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
	}
	if got != nil {
		t.Errorf("expected no task, got task with ID %d", got.ID)
	}
}

func TestProcessTask_Success(t *testing.T) {
	m := newTestManager()

	created, err := m.CreateTask([]int{40, 60, 80, 100})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = m.ProcessTask(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := m.GetTask(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Status != StatusCompleted {
		t.Errorf("expected status %s, got %s", StatusCompleted, got.Status)
	}

	want := Result{Count: 4, Min: 40, Max: 100, Average: 70, Passed: 3}
	if got.Result != want {
		t.Errorf("expected result %+v, got %+v", want, got.Result)
	}
}

func TestScoreProcessor_Process(t *testing.T) {
	p := ScoreProcessor{PassScore: 50}

	tests := []struct {
		name   string
		scores []int
		want   Result
	}{
		{"one score", []int{75}, Result{Count: 1, Min: 75, Max: 75, Average: 75, Passed: 1}},
		{"everybody passed", []int{50, 70, 90}, Result{Count: 3, Min: 50, Max: 90, Average: 70, Passed: 3}},
		{"nobody passed", []int{10, 20, 30}, Result{Count: 3, Min: 10, Max: 30, Average: 20, Passed: 0}},
		{"lowest and highest scores", []int{0, 100}, Result{Count: 2, Min: 0, Max: 100, Average: 50, Passed: 1}},
		{"average with a fraction", []int{50, 51}, Result{Count: 2, Min: 50, Max: 51, Average: 50.5, Passed: 2}},
		{"empty list", []int{}, Result{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := p.Process(tc.scores)
			if got != tc.want {
				t.Errorf("expected %+v, got %+v", tc.want, got)
			}
		})
	}
}

func TestProcessTask_AlreadyCompleted(t *testing.T) {
	m := newTestManager()

	created, err := m.CreateTask([]int{55, 65})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = m.ProcessTask(created.ID)
	if err != nil {
		t.Fatalf("first processing: expected no error, got %v", err)
	}

	err = m.ProcessTask(created.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("second processing: expected error %v, got %v", ErrInvalidTransition, err)
	}

	if created.Status != StatusCompleted {
		t.Errorf("expected status to stay %s, got %s", StatusCompleted, created.Status)
	}
}

func TestProcessTask_NotFound(t *testing.T) {
	m := newTestManager()

	err := m.ProcessTask(7)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
	}
}

func TestChangeStatus(t *testing.T) {
	tests := []struct {
		name    string
		from    Status
		to      Status
		wantErr bool
	}{
		{"pending to processing", StatusPending, StatusProcessing, false},
		{"processing to completed", StatusProcessing, StatusCompleted, false},
		{"pending to completed", StatusPending, StatusCompleted, true},
		{"processing to pending", StatusProcessing, StatusPending, true},
		{"completed to processing", StatusCompleted, StatusProcessing, true},
		{"completed to pending", StatusCompleted, StatusPending, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tk := Task{Status: tc.from}

			err := tk.changeStatus(tc.to)

			if tc.wantErr {
				if !errors.Is(err, ErrInvalidTransition) {
					t.Errorf("expected error %v, got %v", ErrInvalidTransition, err)
				}
				if tk.Status != tc.from {
					t.Errorf("status must stay %s, got %s", tc.from, tk.Status)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if tk.Status != tc.to {
					t.Errorf("expected status %s, got %s", tc.to, tk.Status)
				}
			}
		})
	}
}

func TestMultipleTasks(t *testing.T) {
	m := newTestManager()

	first, err := m.CreateTask([]int{10, 20, 30})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	second, err := m.CreateTask([]int{90, 100})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if first.ID == second.ID {
		t.Fatalf("tasks must have different IDs, both have %d", first.ID)
	}

	err = m.ProcessTask(second.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if first.Status != StatusPending {
		t.Errorf("first task: expected status %s, got %s", StatusPending, first.Status)
	}
	if first.Result.Count != 0 {
		t.Errorf("first task must have no result yet, got %+v", first.Result)
	}

	err = m.ProcessTask(first.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	gotFirst, err := m.GetTask(first.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	gotSecond, err := m.GetTask(second.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wantFirst := Result{Count: 3, Min: 10, Max: 30, Average: 20, Passed: 0}
	wantSecond := Result{Count: 2, Min: 90, Max: 100, Average: 95, Passed: 2}

	if gotFirst.Result != wantFirst {
		t.Errorf("first task: expected %+v, got %+v", wantFirst, gotFirst.Result)
	}
	if gotSecond.Result != wantSecond {
		t.Errorf("second task: expected %+v, got %+v", wantSecond, gotSecond.Result)
	}
}
