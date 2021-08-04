package service

import (
	"github.com/ksputo/k8s-teamhack/internal/storage"
	"github.com/ksputo/k8s-teamhack/internal/storage/model"
	"github.com/pkg/errors"
)

type TaskService struct {
	storage            storage.TaskStorage
	complexityProvider ComplexityProvider
}

type ComplexityProvider interface {
	Get(taskDuration string) (string, error)
}

func NewTaskService(s storage.TaskStorage, dp ComplexityProvider) *TaskService {
	return &TaskService{storage: s, complexityProvider: dp}
}

func (s TaskService) CreateTask(task model.Task) error {
	if task.Complexity == "" {
		complexity, err := s.getComplexity(task.Duration)
		if err != nil {
			return errors.Wrapf(err, "while getting task complexity")
		}
		task.Complexity = complexity
	}

	err := s.storage.Insert(task)
	if err != nil {
		return errors.Wrapf(err, "while inserting task to db")
	}

	return nil
}

func (s TaskService) getComplexity(taskDuration string) (string, error) {
	taskComplexity, err := s.complexityProvider.Get(taskDuration)
	if err != nil {
		return "", err
	}
	return taskComplexity, nil
}

func (s TaskService) GetTask(id string) (model.Task, error) {
	task, err := s.storage.GetByID(id)
	if err != nil {
		return model.Task{}, errors.Wrapf(err, "while getting task from db")
	}

	return *task, nil
}
