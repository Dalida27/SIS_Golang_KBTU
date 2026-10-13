package task

type Manager struct {
	tasks     map[int]*Task
	lastID    int
	processor Processor
}

func NewManager(processor Processor) *Manager {
	return &Manager{
		tasks:     make(map[int]*Task),
		lastID:    0,
		processor: processor,
	}
}

func (m *Manager) CreateTask(scores []int) (*Task, error) {
	err := validateScores(scores)
	if err != nil {
		return nil, err
	}

	scoresCopy := make([]int, len(scores))
	copy(scoresCopy, scores)

	m.lastID++

	newTask := &Task{
		ID:     m.lastID,
		Scores: scoresCopy,
		Status: StatusPending,
	}
	m.tasks[newTask.ID] = newTask

	return newTask, nil
}

func (m *Manager) GetTask(id int) (*Task, error) {
	foundTask, ok := m.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	return foundTask, nil
}

func (m *Manager) ProcessTask(id int) error {
	t, err := m.GetTask(id)
	if err != nil {
		return err
	}

	err = t.changeStatus(StatusProcessing)
	if err != nil {
		return err
	}

	t.Result = m.processor.Process(t.Scores)

	err = t.changeStatus(StatusCompleted)
	if err != nil {
		return err
	}

	return nil
}
