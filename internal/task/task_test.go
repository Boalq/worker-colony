package task_test

import (
	"testing"
	"time"

	"github.com/Boalq/worker-colony/internal/task"
)

func TestTransitionTo(t *testing.T) {
	t.Run("Check for Invalid Transition", func(t *testing.T) {
		
		tsk := task.Task{ID: 1, Status: task.Pending, CreatedAt: time.Now(), Description: "getting some milk"}		
		err := tsk.TransitionTo(task.Completed)

		if err == nil{
			t.Errorf("Have %s want %s", err, "Invalid Transition")
		}
	})
	t.Run("Check for Transition from Running to Failed (exception case)", func(t *testing.T) {

		tsk := task.Task{ID: 1, Status: task.Running, CreatedAt: time.Now(), Description: "getting some milk"}		
		err := tsk.TransitionTo(task.Failed)

		if err != nil{
			t.Errorf("Have error: %s, expected no errors", err)
		}
	})
}