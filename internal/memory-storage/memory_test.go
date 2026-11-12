package memorystorage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Boalq/worker-colony/internal/memory-storage"
	"github.com/Boalq/worker-colony/internal/task"
)

func TestSaveTask(t *testing.T) {
	mem := memorystorage.NewInMemoryStorage()
	
	mem.SaveTask(task.Task{ID: 1, Status: task.Pending, CreatedAt: time.Now()})
	mem.SaveTask(task.Task{ID: 1, Status: task.Pending, CreatedAt: time.Now()})
	mem.SaveTask(task.Task{ID: 2, Status: task.Pending, CreatedAt: time.Now()})
	mem.SaveTask(task.Task{ID: 3, Status: task.Pending, CreatedAt: time.Now()})

	fmt.Print(mem.String())
}

func TestUpdateTask(t *testing.T) {
	mem := memorystorage.NewInMemoryStorage()

	mem.SaveTask(task.Task{ID: 1, Status: task.Pending, CreatedAt: time.Now()})
	err := mem.UpdateTask(1, task.Running)

	if err != nil{
		t.Error(err)
	}

	fmt.Print(mem.String())
}