package netbox

import (
	"context"

	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

// Structure for the plugin
type Netbox struct {
	Url   string
	Token string
}

// The Name() method which is used so other plugins can check if a certain plugin is loaded. The method just returns the string netbox.
func (n Netbox) Name() string { return "netbox" }

// ServeDNS implements the plugin.Handler interface. This methos is called for every dns request
func (n Netbox) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	reqcurl := n.Url + CreateUrl(state.QType(), state.Name())
	log.Info("Request received from " + GetOutboundIP().String() + " for " + state.Name() + " with type " + Saytype(state.QType()))

	body := Client(reqcurl, n.Token)
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	var response = JsonToStruct(body)
	var answers = CreateResponse(state.QType(), state.QName(), response)

	// Message creation
	m := new(dns.Msg)
	// We define the message as type "reply".
	m.SetReply(r)
	m.Authoritative = true
	//Adding the answer to the answer table
	m.Answer = answers
	//Writing the answer
	w.WriteMsg(m)
	log.Info("Response sent for " + state.Name() + " with type: " + Saytype(state.QType()))
	return dns.RcodeSuccess, nil
}
