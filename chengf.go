package chengf

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("chengf")

// Chengf is an chengf plugin to show how to write a plugin.
type Chengf struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface. This method gets called when chengf is used
// in a Server.
func (e Chengf) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// This function could be simpler. I.e. just fmt.Println("chengf") here, but we want to show
	// a slightly more complex chengf as to make this more interesting.
	// Here we wrap the dns.ResponseWriter in a new ResponseWriter and call the next plugin, when the
	// answer comes back, it will print "chengf".

	// Debug log that we've have seen the query. This will only be shown when the debug plugin is loaded.
	log.Debug("Received response")

	// Wrap.
	pw := NewResponsePrinter(w)

	// Export metric with the server label set to the current server handling the request.
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	// Call next plugin (if any).
	return plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
}

// Name implements the Handler interface.
func (e Chengf) Name() string { return "chengf" }

// ResponsePrinter wrap a dns.ResponseWriter and will write chengf to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "chengf" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	log.Info("chengf")
	a := new(dns.Msg)
	var rr dns.RR
	rr = new(dns.A)
	rr.(*dns.A).A = net.ParseIP("192.168.0.111").To4()
	a.Extra = []dns.RR{rr}
	return r.ResponseWriter.WriteMsg(res)
}
