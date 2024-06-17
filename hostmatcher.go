package hostmatcher

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/idna"
)

func newHostMatcher(host string) (matcher, error) {
	host = strings.ToLower(host)
	if len(host) == 0 {
		return nil, fmt.Errorf("%q is no host", host)
	}

	if host == "*" {
		return allMatch{}, nil
	}

	// IPv4/CIDR, IPv6/CIDR
	if _, pnet, err := net.ParseCIDR(host); err == nil {
		return cidrMatch{cidr: pnet}, nil
	}

	// IPv4:port, [IPv6]:port
	phost, pport, err := net.SplitHostPort(host)
	if err == nil {
		if len(phost) == 0 {
			return nil, fmt.Errorf("%q is no host", host)
		}
		if phost[0] == '[' && phost[len(phost)-1] == ']' {
			phost = phost[1 : len(phost)-1]
		}
	} else {
		phost = host
	}
	// IPv4, IPv6
	if pip := net.ParseIP(phost); pip != nil {
		return ipMatch{ip: pip, port: pport}, nil
	}

	if len(phost) == 0 {
		return nil, fmt.Errorf("%q is no host", host)
	}

	// domain.com or domain.com:80
	// foo.com matches bar.foo.com
	// .domain.com or .domain.com:port
	// *.domain.com or *.domain.com:port
	if strings.HasPrefix(phost, "*.") {
		phost = phost[1:]
	}
	matchHost := false
	if phost[0] != '.' {
		matchHost = true
		phost = "." + phost
	}
	phost = idnaASCII(phost)
	return domainMatch{host: phost, port: pport, matchHost: matchHost}, nil
}

func idnaASCII(host string) string {
	h, err := idna.Lookup.ToASCII(host)
	if err != nil {
		return host
	}
	return h
}

// matcher represents the matching rule
type matcher interface {
	// match returns true if the host and optional port or ip and optional port
	// are allowed
	match(host, port string, ip net.IP) bool
}

// allMatcher matches on all possible inputs
type allMatcher struct{}

func (allMatcher) Match(host string) bool {
	return true
}

type allMatch struct{}

func (allMatch) match(host, port string, ip net.IP) bool {
	return true
}

type cidrMatch struct {
	cidr *net.IPNet
}

func (m cidrMatch) match(host, port string, ip net.IP) bool {
	return m.cidr.Contains(ip)
}

type ipMatch struct {
	ip   net.IP
	port string
}

func (m ipMatch) match(host, port string, ip net.IP) bool {
	if m.ip.Equal(ip) {
		return m.port == "" || m.port == port
	}
	return false
}

type domainMatch struct {
	host string
	port string

	matchHost bool
}

func (m domainMatch) match(host, port string, ip net.IP) bool {
	if strings.HasSuffix(host, m.host) || (m.matchHost && host == m.host[1:]) {
		return m.port == "" || m.port == port
	}
	return false
}
