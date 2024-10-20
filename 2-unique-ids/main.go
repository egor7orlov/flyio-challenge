package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var msgBody map[string]any
		if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
			return err
		}

		body := make(map[string]any)

		body["type"] = "generate_ok"
		body["id"] = uuid.New().String()
		body["in_reply_to"] = msgBody["msg_id"]

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
