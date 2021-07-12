package model

type Task struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Duration    string `json:"duration,omitempty"`
	Complexity  string `json:"complexity,omitempty"`
	Done        bool   `json:"done,omitempty"`
}
