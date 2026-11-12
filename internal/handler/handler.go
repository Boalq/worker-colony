package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	memorystorage "github.com/Boalq/worker-colony/internal/memory-storage"
	"github.com/Boalq/worker-colony/internal/task"
)

type TaskHandler struct {
	storage *memorystorage.InMemoryStorage
	queue   chan task.Task
	counter atomic.Int64
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		storage: memorystorage.NewInMemoryStorage(),
		queue:   make(chan task.Task, 3),
		counter: atomic.Int64{},
	}
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	id := h.counter.Add(1)
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	t := task.Task{ID: int(id), Status: task.Pending, CreatedAt: time.Now(), Description: string(body)}

	h.storage.SaveTask(t)

	select {
	case h.queue <- t:
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "task rn in the making\n")
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "queue full, try again later\n")
	}
}

func (h *TaskHandler) StartWorker(WorkerID int) {
	for tsk := range h.queue {
		err := h.storage.UpdateTask(tsk.ID, task.Running)
		if err != nil {
			fmt.Printf("Error: %s while working on the task\n", err)
			continue
		}

		fmt.Printf("Worker %d working on task %d\n", WorkerID, tsk.ID)

		time.Sleep(1 * time.Second)

		err = h.storage.UpdateTask(tsk.ID, task.Completed)
		if err != nil {
			fmt.Printf("Error: %s while working on the task\n", err)
			continue
		}

		fmt.Printf("Worker %d finished task %d\nDescription: %s\n", WorkerID, tsk.ID, tsk.Description)
	}
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	asked_id := r.PathValue("id")
	id, _ := strconv.Atoi(asked_id)

	t, err := h.storage.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
