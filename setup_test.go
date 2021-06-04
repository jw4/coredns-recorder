package recorder

import (
	"fmt"
	"testing"

	"github.com/coredns/caddy"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", PluginName)
	err := setup(c)
	if err != nil {
		t.Errorf("basic setup failed: %v", err)
	}

	c = caddy.NewTestController("dns", fmt.Sprintf(`%s nats core.dns nats://10.36.10.20:4222`, PluginName))
	err = setup(c)
	if err != nil {
		t.Errorf("basic setup failed: %v", err)
	}
}
