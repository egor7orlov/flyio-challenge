package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var mu sync.Mutex

func main() {
	n := maelstrom.NewNode()

	numSet := make(map[float64]any)

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var msgBody map[string]any
		if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
			return err
		}

		num := msgBody["message"].(float64)
		resBody := make(map[string]any)

		mu.Lock()
		numSet[num] = struct{}{}
		mu.Unlock()

		resBody["type"] = "broadcast_ok"

		return n.Reply(msg, resBody)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var msgBody map[string]any
		if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
			return err
		}

		resBody := make(map[string]any)

		mu.Lock()
		resBody["messages"] = getSetValues(&numSet)
		mu.Unlock()

		resBody["type"] = "read_ok"

		return n.Reply(msg, resBody)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var msgBody map[string]any
		if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
			return err
		}

		resBody := make(map[string]any)
		resBody["type"] = "topology_ok"

		return n.Reply(msg, resBody)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func getSetValues(set *map[float64]any) []float64 {
	values := make([]float64, 0, len(*set))

	for key := range *set {
		values = append(values, key)
	}

	return values
}
