package main

// Github ...
type Github struct {
	Token string `json:"token"`
}

// Config ...
type Config struct {
	Github *Github `json:"github, omitempty"`
}

// Hook ...
type Hook struct {
	URL string `json:"url"`
}

// Transition ...
type Transition struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Hooks []Hook `json:"hooks"`
}

// Workflow ...
type Workflow struct {
	Transitions []Transition `json:"transitions"`
}

// Issue
type Issue struct {
	ID     string `json:"string"`
	Number string `json:"string"`
	Repo   string `json:"repo"`
	Owner  string `json:"owner"`
	Status string `json:"status"`
}

// Move ...
type Move struct {
	Issue    *Issue  `json:"issue"`
	Status   *string `json:"status"`
	Config   `json:"config"`
	Workflow `json:"workflow"`
}

// ErrorMessage representation of a json error message
type ErrorMessage struct {
	Error string `json:"error"`
}
