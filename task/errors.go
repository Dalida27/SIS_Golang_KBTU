package task

import "errors"

var (
	ErrEmptyScores       = errors.New("scores list is empty")
	ErrScoreOutOfRange   = errors.New("every score must be from 0 to 100")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid status transition")
)
