package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/supu-io/payload"
)

// Move ...
func move(m *Move) error {
	t := getTransition(m)
	if t == nil {
		return errors.New("Invalid transition")
	}
	callHooks(m, t)

	return nil
}

func getTransition(m *Move) *Transition {
	from := m.Issue.Status
	to := *m.Status
	for _, t := range m.Workflow.Transitions {
		if t.From == from && t.To == to {
			return &t
		}
	}

	return nil
}

func callHooks(m *Move, t *Transition) {
	p := getPayload(m)
	for _, h := range t.Hooks {
		hook(h, p)
	}
}

func hook(hook Hook, p string) {
	req, err := http.NewRequest("POST", hook.URL, nil)
	// TODO: We will need at some point to support tokens
	// req.Header.Add("X-AUTH-TOKEN", "token")
	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		log.Println("Couldn't connect to the server")
	}
}

func getPayload(m *Move) string {
	from := m.Issue.Status
	to := *m.Status

	c := m.Config
	i := m.Issue
	g := payload.Github{Token: &c.Github.Token}
	conf := payload.Config{Github: &g}

	t := payload.Transition{From: &from, To: &to}

	s := []string{"", ""}

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
