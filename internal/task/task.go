package task

import (
	"errors"
	"fmt"
	"time"
)

type TaskState int

const (
	Pending TaskState = iota
	Running
	Completed
	Failed
)

type Task struct {
	ID          int       `json:"ID"`
	Status      TaskState `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
	Description string    `json:"Description"`
}

func (t *Task) Do() string {
	t.Status = Completed
	return fmt.Sprintf("Task ID:%d done(%d)", t.ID, t.Status)
}

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
