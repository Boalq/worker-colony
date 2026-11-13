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

// TaskHandler struct that has
// A temporary storage to keep a collection of tasks
// Queue for buffered concurrent tasks
// Counter for ID assigments of incoming tasks
type TaskHandler struct {
	storage *memorystorage.InMemoryStorage
	queue   chan task.Task
	counter atomic.Int64
}

// "Constructor" for a new Task Handler 
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		storage: memorystorage.NewInMemoryStorage(),
		queue:   make(chan task.Task, 3),
		counter: atomic.Int64{},
	}
}

// Adding a Task to the TaskHandler Struct
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

	// Http response for the User to verify if the task has been
	// stored in the memory and queued up in the Channel
	select {
	case h.queue <- t:
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "task is being processed => ID: %d\n", id)
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "queue full, try again later\n")
	}
}

// A Worker (bee) instance with an Designaited ID Number
// for logging and monitoring concurrent work on the Tasks
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

// Getter Method for a Task (usually not that prefered)
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
