package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/supu-io/messages"
	"github.com/supu-io/payload"
)

// Move ...
func move(m *messages.UpdateIssue) error {
	t := getTransition(m)
	if t == nil {
		return errors.New("Invalid transition")
	}
	callHooks(m, t)

	return nil
}

func getTransition(m *messages.UpdateIssue) *messages.Transition {
	from := m.Issue.Status
	to := m.Status
	for _, t := range m.Workflow.Transitions {
		if t.From == from && t.To == to {
			return &t
		}
	}

	return nil
}

func callHooks(m *messages.UpdateIssue, t *messages.Transition) {
	p := getPayload(m)
	for _, h := range t.Hooks {
		hook(h, p)
	}
}

func hook(hook messages.Hook, p string) {
	req, err := http.NewRequest("POST", hook.URL, nil)
	// TODO: We will need at some point to support tokens
	// req.Header.Add("X-AUTH-TOKEN", "token")
	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		log.Println("Couldn't connect to the server")
	}
}

func getPayload(m *messages.UpdateIssue) string {
	from := m.Issue.Status
	to := m.Status

	c := m.Config
	i := m.Issue
	g := payload.Github{Token: &c.Github.Token}
	conf := payload.Config{Github: &g}

	t := payload.Transition{From: &from, To: &to}

	s := getAllStatus(m.Workflow)

	issue := payload.Issue{ID: &i.ID}
	p := payload.Payload{
		Config:     &conf,
		Transition: &t,
		Status:     &s,
		Issue:      &issue,
	}

	body, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(body)
}

func getAllStatus(w messages.Workflow) []string {
	status := []string{}
	for _, t := range w.Transitions {
		addFrom := true
		addTo := true
		for _, s := range status {
			if s == t.To {
				addTo = false
			}
			if s == t.From {
				addFrom = false
			}
		}
		if addFrom == true {
			status = append(status, string(t.From))
		}
		if addTo == true {
			status = append(status, string(t.To))
		}
	}

	return status
}
