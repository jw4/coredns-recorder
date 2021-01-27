// Package recorder provides a CoreDNS plugin to record queries.
package recorder

import (
	"context"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// PluginName is the CoreDNS name for this plugin.
const PluginName = "recorder"

// Plugin provides a CoreDNS plugin to record to redis.
type Plugin struct {
	Next plugin.Handler
	out  output
	log  clog.P
}

// ServeDNS implements the plugin interface.
func (p Plugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if err := p.out.Record(request.Request{W: w, Req: r}); err != nil {
		p.log.Error(err)
	}

	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

// Name implements the plugin interface.
func (p Plugin) Name() string { return PluginName }

// Ready implements the plugin interface.
func (p Plugin) Ready() bool { return true }
