package storage

import (
	"github.com/ksputo/k8s-teamhack/internal/storage/driver/memory"
	"github.com/ksputo/k8s-teamhack/internal/storage/model"
)

type TaskStorage interface {
	Insert(task model.Task) error
	GetByID(taskID string) (*model.Task, error)
}

func NewFromConfig() (TaskStorage, error) {
	return nil, nil
}

func NewMemoryStorage() TaskStorage {
	return memory.NewTaskStorage()
}
