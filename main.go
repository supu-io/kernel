package main

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/nats-io/nats"
	"github.com/supu-io/messages"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

	nc.Subscribe("workflow.move", func(msg *nats.Msg) {
		m := messages.UpdateIssue{}
		json.Unmarshal(msg.Data, &m)
		if err := move(&m); err != nil {
			e := ErrorMessage{Error: err.Error()}
			nc.Publish(msg.Reply, ToJSON(e))
		} else {
			nc.Publish(msg.Reply, ToJSON(m.Issue))
		}
	})
	runtime.Goexit()
}

// ToJSON represents an object as a json []byte
func ToJSON(i interface{}) []byte {
	json, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	return json
}
