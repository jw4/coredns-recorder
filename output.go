package recorder

import (
	"fmt"
	"os"
	"time"

	"github.com/coredns/coredns/request"
)

type output interface {
	Record(r request.Request) error
	Connect() error
	Flush() error
}

type defaultOutput struct{}

func (o defaultOutput) Record(r request.Request) error {
	now := time.Now()
	message := fmt.Sprintf(
		`{"ip":%q,"name":%q,"class":%q,"time":%q,"ts":%d}`,
		r.IP(), r.Name(), r.Type(), now.Format(time.RFC3339), now.Unix())

	_, err := fmt.Fprintf(os.Stderr, "%s: %s\n", PluginName, message)
	return err
}

func (o defaultOutput) Connect() error { return nil }
func (o defaultOutput) Flush() error   { return nil }
