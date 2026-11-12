package handler

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	memorystorage "github.com/Boalq/worker-colony/internal/memory-storage"
	"github.com/Boalq/worker-colony/internal/task"
)

type TaskHandler struct {
	storage memorystorage.InMemoryStorage
	queue   chan task.Task
	counter atomic.Int64
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		storage: *memorystorage.NewInMemoryStorage(),
		queue:   make(chan task.Task, 3),
		counter: atomic.Int64{},
	}
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	id := h.counter.Add(1)
	t := task.Task{ID: int(id), Status: task.Pending, CreatedAt: time.Now()}

	h.storage.SaveTask(t)
	h.queue <- t

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "task rn in the making\n")
}
