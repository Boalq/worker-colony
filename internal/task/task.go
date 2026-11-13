package task

import (
	"errors"
	"time"
)

type TaskState int

// "Enums" for the state of the Tasks
const (
	Pending TaskState = iota
	Running
	Completed
	Failed
)

// This is a Mock "Task" Struct, thus the main focus is on the
// actual mechanism of the Worker and Http request handling
// JSON tags for encoding and decoding so it is suited for any API use
type Task struct {
	ID          int       `json:"ID"`
	Status      TaskState `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
	Description string    `json:"Description"`
}

// Transition Machine with Rules on what transitions are
// allowed and what not
func (t *Task) TransitionTo(newStatus TaskState) error {
	switch {
	case t.Status == Pending && newStatus == Running:
		t.Status = newStatus
		return nil
	case t.Status == Running && (newStatus == Completed || newStatus == Failed):
		t.Status = newStatus
		return nil
	default:
		return errors.New("Invalid Transition")
	}
}
