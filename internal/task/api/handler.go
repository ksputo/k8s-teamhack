package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ksputo/k8s-teamhack/internal/storage/model"
	"github.com/ksputo/k8s-teamhack/internal/task/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type apiHandler struct {
	taskService *service.TaskService
	log         logrus.FieldLogger
}

func (h *apiHandler) AttachRoutes(rtr *mux.Router) {
	rtr.HandleFunc("/task/{task_id}", h.createTask).Methods(http.MethodPost)
	rtr.HandleFunc("/task/{task_id}", h.getTask).Methods(http.MethodGet)
}

func NewHandler(ts *service.TaskService, log logrus.FieldLogger) *apiHandler {
	return &apiHandler{
		taskService: ts,
		log:         log,
	}
}

func (h *apiHandler) createTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskID := vars["task_id"]
	logger := h.log.WithField("taskID", taskID)

	logger.Info("creating task")

	var task model.Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		h.log.Error(errors.Wrapf(err, "while decoding request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.taskService.CreateTask(task)
	if err != nil {
		h.log.Error(errors.Wrapf(err, "while creating task"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Info("successfully created task")
	w.WriteHeader(http.StatusCreated)
}

func (h *apiHandler) getTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	taskID := vars["task_id"]

	logger := h.log.WithField("taskID", taskID)

	logger.Info("getting task")

	task, err := h.taskService.GetTask(taskID)
	if err != nil {
		h.log.Error(errors.Wrapf(err, "while getting task"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		h.log.Warnf("could not write response %s", string(res))
	}
}
