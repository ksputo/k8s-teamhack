package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ksputo/k8s-teamhack/internal/storage"
	"github.com/ksputo/k8s-teamhack/internal/task/api"
	"github.com/ksputo/k8s-teamhack/internal/task/complexity"
	"github.com/ksputo/k8s-teamhack/internal/task/service"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	DbInMemory           bool   `envconfig:"default=true"`
	ComplexityServiceURL string `envconfig:"default=http://localhost:3000/complexity"`
	Host                 string `envconfig:"default=localhost"`
	Port                 string `envconfig:"default=8080"`
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Info("Starting TaskService")

	var cfg Config
	err := envconfig.InitWithPrefix(&cfg, "APP")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("initiated config")
	var db storage.TaskStorage
	if cfg.DbInMemory {
		logger.Info("db in memory")
		db = storage.NewMemoryStorage()
	} else {
		s, err := storage.NewFromConfig()
		if err != nil {
			logger.Fatal(err)
		}
		db = s
	}
	logger.Info("initiated db")

	complexityProvider := complexity.NewClient(cfg.ComplexityServiceURL)

	taskService := service.NewTaskService(db, complexityProvider)

	h := api.NewHandler(taskService, logger)
	r := mux.NewRouter()
	h.AttachRoutes(r)
	logger.Info("initiated here")

	err = http.ListenAndServe(cfg.Host+":"+cfg.Port, r)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Server started on %s:%s", cfg.Host, cfg.Port)
}
