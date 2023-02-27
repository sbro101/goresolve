package goresolve

import (
	"github.com/miekg/dns"
)

// Hosts struct
type Data struct {
	Hostname     string   `json:"hostname,omitempty"`
	IPv4         []string `json:"ipv4,omitempty"`
	IPv6         []string `json:"ipv6,omitempty"`
	CNAME        string   `json:"cname,omitempty"`
	Error        bool     `json:"error,omitempty"`
	ErrorMessage string   `json:"errormessage,omitempty"`
}

// Hostname function
func Hostname(hostname string, nameserver string) *Data {
	r := new(Data)

	r.Hostname = hostname

	cname, err := GetCNAME(r.Hostname, nameserver)
	if err != nil {
		r.Error = true
		r.ErrorMessage = err.Error()
		return r
	}

	if len(cname) > 0 {
		r.CNAME = cname
		return r
	}

	ar, err := GetA(r.Hostname, nameserver)
	if err != nil {
		r.Error = true
		r.ErrorMessage = err.Error()
		return r
	}
	r.IPv4 = ar

	aaaar, err := GetAAAA(r.Hostname, nameserver)
	if err != nil {
		r.Error = true
		r.ErrorMessage = err.Error()
		return r
	}
	r.IPv6 = aaaar

	return r

}

// GetCNAME function
func GetCNAME(hostname string, nameserver string) (string, error) {
	var cname string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeCNAME)
	c := new(dns.Client)
	m.MsgHdr.RecursionDesired = true
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return "none", err
	}
	for _, rin := range in.Answer {
		if r, ok := rin.(*dns.CNAME); ok {
			cname = r.Target
		}
	}
	return cname, nil
}

// GetA function
func GetA(hostname string, nameserver string) ([]string, error) {
	var record []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeA)
	c := new(dns.Client)
	m.MsgHdr.RecursionDesired = true
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return nil, err
	}
	for _, rin := range in.Answer {
		if r, ok := rin.(*dns.A); ok {
			record = append(record, r.A.String())
		}
	}

	return record, nil
}

// GetAAAA function
func GetAAAA(hostname string, nameserver string) ([]string, error) {
	var record []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeAAAA)
	c := new(dns.Client)
	m.MsgHdr.RecursionDesired = true
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return nil, err
	}
	for _, rin := range in.Answer {
		if r, ok := rin.(*dns.AAAA); ok {
			record = append(record, r.AAAA.String())
		}
	}

	return record, nil
}
