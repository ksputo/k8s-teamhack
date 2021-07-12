package complexity

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ksputo/k8s-teamhack/internal/storage/model"
	"github.com/pkg/errors"
)

type client struct {
	baseUrl string
}

func NewClient(baseURL string) *client {
	return &client{
		baseUrl: baseURL,
	}
}

func (c *client) Get(duration string) (string, error) {
	var taskDuration model.Task
	taskDuration.Duration = duration

	payload, err := json.Marshal(taskDuration)
	if err != nil {
		return "", errors.Wrap(err, "while creating HTTP request")
	}

	req, err := http.NewRequest(http.MethodPost, c.baseUrl, bytes.NewBuffer(payload))
	if err != nil {
		return "", errors.Wrap(err, "while creating HTTP request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "while making HTTP call")
	}
	defer resp.Body.Close()

	bodyRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "while reading HTTP response body")
	}
	var result model.Task
	if err = json.Unmarshal(bodyRaw, &result); err != nil {
		return "", errors.Wrap(err, "while decoding HTTP response")
	}

	return result.Complexity, nil
}
