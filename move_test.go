package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func getMessage() *Move {
	source := "fixtures/move.json"
	absPath, _ := filepath.Abs(source)
	file, err := os.Open(absPath)
	log.Printf("Reading config from: %s", source)
	if err != nil {
		log.Panic("error:", err)
	}

	decoder := json.NewDecoder(file)
	m := Move{}
	err = decoder.Decode(&m)
	if err != nil {
		log.Println("Definition file is invalid")
		log.Panic("error:", err)
	}

	return &m
}

var done = false
var wg sync.WaitGroup

func mockServer(route string, method string, status int, output string) *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		s := output
		if s == "" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			s = buf.String()
		}
		w.WriteHeader(status)
		w.Header().Set("X-Auth-Token", "")
		fmt.Fprint(w, s)
		println("supu")
		wg.Done()
	}).Methods(method)

	return httptest.NewServer(r)
}

func TestInvalidTransition(t *testing.T) {
	Convey("Given received transition is invalid", t, func() {
		Convey("When send the message", func() {
			done = false
			server := mockServer("/", "POST", 200, "hello")
			m := getMessage()
			for i, t := range m.Workflow.Transitions {
				m.Workflow.Transitions[i].Hooks = append(t.Hooks, Hook{URL: server.URL})
			}
			m.Issue.Status = "done"
			move(m)
		})
	})
}

func TestValidTransitionWithoutHooks(t *testing.T) {
	Convey("Given received transition is invalid", t, func() {
		Convey("When send the message", func() {
			done = false
			m := getMessage()
			move(m)
		})
	})
}

func TestValidTransitionXX(t *testing.T) {
	Convey("Given received transition is invalid", t, func() {
		Convey("When send the message", func() {
			done = false
			server := mockServer("/", "POST", 200, "hello")
			m := getMessage()
			wg.Add(1)
			for i, t := range m.Workflow.Transitions {
				m.Workflow.Transitions[i].Hooks = append(t.Hooks, Hook{URL: server.URL})
			}
			move(m)
			wg.Wait()
		})
	})
}
