package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	topic := os.Getenv("NATS_TOPIC")
	urls := os.Getenv("NATS_URLS")

	var out io.WriteCloser = &writer{}
	if len(os.Args) > 1 {
		out = &writer{filename: os.Args[1]}
	}

	defer out.Close()

	nc, err := nats.Connect(urls)
	if err != nil {
		log.Fatalf("connecting to %s: %v", urls, err)
	}

	defer nc.Close()

	sub, err := nc.SubscribeSync(topic)
	if err != nil {
		log.Fatalf("subscribing to %s: %v", topic, err)
	}

	for {
		m, err := sub.NextMsg(time.Minute * 10)
		if err != nil {
			log.Fatalf("retrieving next message: %v", err)
		}

		d := &message{}
		if err = json.NewDecoder(bytes.NewReader(m.Data)).Decode(d); err != nil {
			log.Fatalf("decoding next message: %v", err)
		}
		fmt.Fprintf(out, "%s\n", d)
	}
}
