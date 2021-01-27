// Package recorder provides a CoreDNS plugin to record queries.
package recorder

import (
	"context"
	"fmt"
	"os"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// PluginName is the CoreDNS name for this plugin.
const PluginName = "recorder"

// Plugin provides a CoreDNS plugin to record to redis.
type Plugin struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin interface.
func (p Plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	fmt.Fprintf(os.Stderr, `{"ip":%q,"name":%q,"class":%q}`, state.IP(), state.Name(), state.Type())
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

// Name implements the plugin interface.
func (p Plugin) Name() string { return PluginName }

// Ready implements the plugin interface.
func (p Plugin) Ready() bool { return true }

func setup(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return plugin.Error(PluginName, c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler { return Plugin{Next: next} })

	return nil
}
