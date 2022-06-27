package netbox

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

// Structure for the records, this structure is used to access to the results
type Response struct {
	Records []Records `json:"results"`
}

// Structure for the values, this structure is used to access to the value for the desired record
type Records struct {
	Value string `json:"value"`
}

//This method is called for every Api call. This fonction calls the right url depending on the type and the hostname
func CreateUrl(qtype uint16, name string) string {

	hostname := ""
	hostname = strings.TrimSuffix(name, ".")

	switch qtype {
	case dns.TypeA:
		return "/api/plugins/netbox-dns/records/?type=A&name=" + hostname

	case dns.TypeAAAA:
		return "/api/plugins/netbox-dns/records/?type=AAAA&name=" + hostname

	case dns.TypeMX:
		return "/api/plugins/netbox-dns/records/?type=MX&name=" + hostname

	case dns.TypeTXT:
		return "/api/plugins/netbox-dns/records/?type=TXT&name=" + hostname

	case dns.TypeCNAME:
		return "/api/plugins/netbox-dns/records/?type=CNAME&name=" + hostname

	case dns.TypeSOA:
		return "/api/plugins/netbox-dns/records/?type=SOA&name=" + hostname

	case dns.TypeNS:
		return "/api/plugins/netbox-dns/records/?type=NS&name=" + hostname

	case dns.TypeSRV:
		return "/api/plugins/netbox-dns/records/?type=SRV&name=" + hostname

	case dns.TypePTR:
		return "/api/plugins/netbox-dns/records/?type=PTR&name=" + hostname

	case dns.TypeSPF:
		return "/api/plugins/netbox-dns/records/?type=SPF&name=" + hostname
	}
	return ""
}

//This fonction converts the uint16 to string
func Saytype(qtype uint16) string {

	switch qtype {
	case dns.TypeA:
		return "A"

	case dns.TypeAAAA:
		return "AAAA"

	case dns.TypeMX:
		return "MX"

	case dns.TypeTXT:
		return "TXT"

	case dns.TypeCNAME:
		return "CNAME"

	case dns.TypeSOA:
		return "SOA"

	case dns.TypeNS:
		return "NS"

	case dns.TypeSRV:
		return "SRV"

	case dns.TypePTR:
		return "PTR"

	case dns.TypeSPF:
		return "SPF"
	}
	return ""
}

// This fonction converts JSON to a structure
func JsonToStruct(body []byte) Response {
	var response Response
	json.Unmarshal(body, &response)
	return response
}

// This fonction returns the right answer depending on the type and the hostname. To get the right answer, this fonction will compare the values for each records.
func CreateResponse(qtype uint16, name string, response Response) []dns.RR {
	// on cree un tbl de response
	answers := []dns.RR{}

	for _, r := range response.Records {
		switch qtype {
		case dns.TypeA:
			// on cree une reponse
			rr := new(dns.A)
			// on cree un header
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET}
			// on parse l'IP dans un format IPV4
			rr.A = net.ParseIP(string(r.Value)).To4()
			// on ajoute la réponse au tableau
			answers = append(answers, rr)

		case dns.TypeAAAA:
			rr := new(dns.AAAA)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET}
			// on parse l'IP dans un format IPV6
			rr.AAAA = net.ParseIP(string(r.Value)).To16()
			answers = append(answers, rr)

		case dns.TypeMX:
			rr := new(dns.MX)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeMX, Class: dns.ClassINET}
			//on récupère le string
			rr.Mx = r.Value + "."
			answers = append(answers, rr)

		case dns.TypeTXT:
			rr := new(dns.TXT)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeTXT, Class: dns.ClassINET}
			rr.Txt = []string{r.Value + "."}
			answers = append(answers, rr)

		case dns.TypeCNAME:
			rr := new(dns.CNAME)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET}
			rr.Target = r.Value + "."
			answers = append(answers, rr)

		case dns.TypeSOA:
			rr := new(dns.SOA)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeSOA, Class: dns.ClassINET}
			rr.Ns = r.Value + "."
			rr.Mbox = r.Value + "."
			answers = append(answers, rr)

		case dns.TypeNS:
			rr := new(dns.NS)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeNS, Class: dns.ClassINET}
			rr.Ns = r.Value + "."
			answers = append(answers, rr)

		case dns.TypeSRV:
			rr := new(dns.SRV)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeSRV, Class: dns.ClassINET}
			rr.Target = r.Value + "."
			answers = append(answers, rr)

		case dns.TypePTR:
			rr := new(dns.PTR)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypePTR, Class: dns.ClassINET}
			rr.Ptr = r.Value + "."
			answers = append(answers, rr)

		case dns.TypeSPF:
			rr := new(dns.SPF)
			rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeSPF, Class: dns.ClassINET}
			rr.Txt = []string{r.Value + "."}
			answers = append(answers, rr)
		}
	}
	if len(answers) == 0 {

		log.Fatal(name + " does not exist in Netbox")
	}
	return answers
}

// This fonction returns the local IP address
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
