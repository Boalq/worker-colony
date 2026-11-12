package memorystorage

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/Boalq/worker-colony/internal/task"
)

type InMemoryStorage struct {
	mu    sync.RWMutex
	tasks map[int]task.Task
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[int]task.Task),
	}
}

func (m *InMemoryStorage) SaveTask(newTask task.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, task := range m.tasks {
		if task == newTask {
			return
		}
	}

	m.tasks[newTask.ID] = newTask
}

func (m *InMemoryStorage) GetTask(id int) (task.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, task := range m.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return task.Task{}, errors.New("Task Not Found")
}

func (m *InMemoryStorage) UpdateTask(id int, status task.TaskState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if ok == false {
		return errors.New("Could not find the Task")
	}

	err := task.TransitionTo(status)
	if err != nil {
		return err
	}

	m.tasks[id] = task
	return nil
}

func (m *InMemoryStorage) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var buf bytes.Buffer
	buf.WriteString("A Memory Storage:\n")
	for task := range m.tasks {
		fmt.Fprintf(&buf, "%v\n", m.tasks[task])
	}

	return buf.String()
}
