package goresolve

import (
	"errors"

	"github.com/miekg/dns"
)

// Hosts struct
type Data struct {
	Hostname     string   `json:"hostname"`
	IPv4         []string `json:"ipv4,omitempty"`
	IPv6         []string `json:"ipv6,omitempty"`
	CNAME        string   `json:"cname,omitempty"`
	Error        bool     `json:"error,omitempty"`
	ErrorMessage string   `json:"errormessage,omitempty"`
}

func Hostname(hostname, nameserver string) (*Data, error) {
	d := new(Data)

	d.Hostname = hostname

	cname, err := GetCNAME(hostname, nameserver)
	if err != nil {
		return setError(d, err), nil
	}

	// If there is a CNAME record, return it and do not look up A/AAAA records.
	if len(cname) > 0 {
		d.CNAME = cname
		return d, nil
	}

	// Look up A records for the given hostname.
	d.IPv4, err = getIPRecords(hostname, nameserver, dns.TypeA)
	if err != nil {
		return setError(d, err), nil
	}

	// Look up AAAA records for the given hostname.
	d.IPv6, err = getIPRecords(hostname, nameserver, dns.TypeAAAA)
	if err != nil {
		return setError(d, err), nil
	}

	return d, nil
}

// GetCNAME returns the CNAME record for the given hostname using the
// specified nameserver.
//
// The returned string will be empty if there is no CNAME record for the given
// hostname. If an error occurs, the returned error will be non-nil.
func GetCNAME(hostname, nameserver string) (string, error) {
	if hostname == "" {
		return "", errors.New("empty hostname")
	}
	if nameserver == "" {
		return "", errors.New("empty nameserver")
	}

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(hostname), dns.TypeCNAME)

	client := &dns.Client{}
	msg.MsgHdr.RecursionDesired = true

	resp, _, err := client.Exchange(msg, nameserver+":53")
	if err != nil {
		return "", err
	}

	// Iterate over the returned DNS records and find the CNAME record.
	for _, ans := range resp.Answer {
		cname, ok := ans.(*dns.CNAME)
		if !ok {
			continue
		}

		// The CNAME record must have a valid target.
		if cname.Target == "" {
			return "", errors.New("empty CNAME target")
		}

		return cname.Target, nil
	}

	// If no CNAME record is found, return NO errors.
	return "", nil
}

// getIPRecords performs a DNS query to retrieve the IP addresses for the
// given hostname. The recordType parameter specifies the type of IP address
// to retrieve (either A or AAAA records). The returned error is non-nil if
// there is a problem with the DNS query. The returned slice of strings will
// be empty if there is no IP address for the given hostname.
func getIPRecords(hostname, nameserver string, recordType uint16) ([]string, error) {
	// Sanity check the input parameters.
	if hostname == "" {
		return nil, errors.New("empty hostname")
	}
	if nameserver == "" {
		return nil, errors.New("empty nameserver")
	}

	// Create a new DNS message.
	var records []string
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(hostname), recordType)

	// Create a new DNS client.
	client := &dns.Client{}
	msg.MsgHdr.RecursionDesired = true

	// Perform the DNS query.
	resp, _, err := client.Exchange(msg, nameserver+":53")
	if err != nil {
		// If there is a problem with the DNS query, return an error.
		return nil, err
	}

	// Iterate over the returned DNS records and find the IP address records.
	for _, answer := range resp.Answer {
		switch recordType {
		case dns.TypeA:
			// We are looking for A records.
			a, ok := answer.(*dns.A)
			if !ok {
				// If this is not an A record, skip it.
				continue
			}
			if a.A == nil {
				// If the A record is empty, skip it.
				continue
			}
			// Append the IP address to the list of results.
			records = append(records, a.A.String())
		case dns.TypeAAAA:
			// We are looking for AAAA records.
			aaaa, ok := answer.(*dns.AAAA)
			if !ok {
				// If this is not an AAAA record, skip it.
				continue
			}
			if aaaa.AAAA == nil {
				// If the AAAA record is empty, skip it.
				continue
			}
			// Append the IP address to the list of results.
			records = append(records, aaaa.AAAA.String())
		default:
			// If we encounter an unknown record type, return an error.
			return nil, errors.New("unknown record type")
		}
	}

	return records, nil
}

// setError sets the Error field of the given Data struct to true and sets the
// ErrorMessage field to the given error message.
func setError(data *Data, err error) *Data {
	data.Error = true
	data.ErrorMessage = err.Error()
	return data
}
