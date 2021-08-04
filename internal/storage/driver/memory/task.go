package memory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ksputo/k8s-teamhack/internal/storage/model"
)

type memoryStore struct {
	mu    sync.Mutex
	tasks map[string]model.Task
}

func NewTaskStorage() *memoryStore {
	return &memoryStore{
		tasks: make(map[string]model.Task),
	}
}

func (s *memoryStore) Insert(task model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task
	return nil
}

func (s *memoryStore) GetByID(taskID string) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, exists := s.tasks[taskID]
	if !exists {
		return nil, errors.New(fmt.Sprintf("task with id %s does not exist", taskID))
	}

	return &t, nil
}
