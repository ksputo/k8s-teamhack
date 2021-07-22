package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ksputo/k8s-teamhack/internal/storage/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

const (
	EasyComplexity           string        = "easy"
	EasyComplexityDuration   time.Duration = time.Hour
	MediumComplexity         string        = "medium"
	MediumComplexityDuration time.Duration = 2 * time.Hour
	HardComplexity           string        = "hard"
	HardComplexityDuration   time.Duration = 3 * time.Hour
	UnknownComplexity        string        = "unknown"
)

type Config struct {
	Port string `envconfig:"default=3000"`
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Info("Starting ComplexityService")

	var cfg Config
	err := envconfig.InitWithPrefix(&cfg, "APP")
	if err != nil {
		logger.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/complexity", func(w http.ResponseWriter, req *http.Request) {
		var taskDuration model.Task
		if err := json.NewDecoder(req.Body).Decode(&taskDuration); err != nil {
			logger.Error(errors.Wrapf(err, "while decoding request"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		timeDuration, err := time.ParseDuration(taskDuration.Duration)
		if err != nil {
			logger.Error(errors.Wrapf(err, "while parsing task duration"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		taskDuration.Complexity = countComplexity(timeDuration)
		logger.Infof("complexity for task duration %q is %q", taskDuration.Duration, taskDuration.Complexity)

		res, err := json.Marshal(taskDuration)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(res)
		if err != nil {
			logger.Warnf("could not write response %s", string(res))
		}
	}).Methods(http.MethodPost)

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		logger.Fatal(err)
	}
}

func countComplexity(d time.Duration) string {
	if d < EasyComplexityDuration {
		return EasyComplexity
	}
	if d > EasyComplexityDuration && d <= MediumComplexityDuration {
		return MediumComplexity
	}
	if d > MediumComplexityDuration && d <= HardComplexityDuration {
		return MediumComplexity
	}
	if d > HardComplexityDuration {
		return HardComplexity
	}

	return UnknownComplexity
}
