package recorder

import (
	"fmt"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

func setup(c *caddy.Controller) error {
	var out output = &defaultOutput{}
	c.Next()
	if c.NextArg() {
		switch c.Val() {
		case "nats":
			args := c.RemainingArgs()
			if len(args) != 2 {
				err := fmt.Errorf("expected <topic> <endpoint>, got %s: %w", strings.Join(args, " "), c.ArgErr())
				return plugin.Error(PluginName, err)
			}

			out = &natsOutput{Topic: args[0], Endpoint: args[1]}
			if err := out.(*natsOutput).Validate(); err != nil {
				return plugin.Error(PluginName, fmt.Errorf("invalid nats config: %v: %w", err, c.ArgErr()))
			}
		default:
			return plugin.Error(PluginName, c.ArgErr())
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Plugin{
			Next: next,
			out:  out,
			log:  clog.NewWithPlugin(PluginName),
		}
	})

	c.OnStartup(func() error { return out.Connect() })
	c.OnShutdown(func() error { return out.Flush() })

	return nil
}

func init() { plugin.Register(PluginName, setup) }
