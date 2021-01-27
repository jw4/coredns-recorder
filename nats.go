package recorder

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/coredns/coredns/request"
	"github.com/nats-io/nats.go"
)

// Errors.
var (
	ErrInvalidTopic    = errors.New("invalid subject / topic name")
	ErrInvalidEndpoint = errors.New("invalid endpoint")
)

type natsOutput struct {
	Topic    string
	Endpoint string
	conn     *nats.Conn
}

func (o *natsOutput) Record(r request.Request) error {
	if err := o.Connect(); err != nil {
		return err
	}

	now := time.Now()
	message := fmt.Sprintf(
		`{"ip":%q,"name":%q,"class":%q,"time":%q,"ts":%d}`,
		r.IP(), r.Name(), r.Type(), now.Format(time.RFC3339), now.Unix())

	return o.conn.Publish(o.Topic, []byte(message))
}

func (o *natsOutput) Connect() error {
	if o.conn != nil {
		if o.conn.IsConnected() {
			return nil
		}
	}

	nc, err := nats.Connect(o.Endpoint)
	if err != nil {
		return err
	}

	o.conn = nc

	return nil
}

func (o *natsOutput) Flush() error {
	if o.conn != nil {
		if err := o.conn.Flush(); err != nil {
			return err
		}
		o.conn.Close()
		o.conn = nil
	}

	return nil
}

func (o *natsOutput) Validate() error {
	if strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '.':
			return r
		default:
			return -1
		}
	}, o.Topic) != o.Topic {
		return ErrInvalidTopic
	}

	endpoints := strings.Split(o.Endpoint, ",")
	for _, endpoint := range endpoints {
		if _, err := url.Parse(strings.TrimSpace(endpoint)); err != nil {
			return ErrInvalidEndpoint
		}
	}

	return nil
}
