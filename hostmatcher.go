package hostmatcher

import (
	"net"
	"strings"

	"golang.org/x/net/idna"
)

// Matcher represents the matching rule
type Matcher interface {
	Match(host string) bool
}

func NewMatcher(hosts []string) Matcher {
	var m hostMatcher
	for _, p := range hosts {
		p = strings.ToLower(strings.TrimSpace(p))
		if len(p) == 0 {
			continue
		}

		if p == "*" {
			return allMatch{}
		}

		// IPv4/CIDR, IPv6/CIDR
		if _, pnet, err := net.ParseCIDR(p); err == nil {
			m.ipMatchers = append(m.ipMatchers, cidrMatch{cidr: pnet})
			continue
		}

		// IPv4:port, [IPv6]:port
		phost, pport, err := net.SplitHostPort(p)
		if err == nil {
			if len(phost) == 0 {
				// There is no host part, likely the entry is malformed; ignore.
				continue
			}
			if phost[0] == '[' && phost[len(phost)-1] == ']' {
				phost = phost[1 : len(phost)-1]
			}
		} else {
			phost = p
		}
		// IPv4, IPv6
		if pip := net.ParseIP(phost); pip != nil {
			m.ipMatchers = append(m.ipMatchers, ipMatch{ip: pip, port: pport})
			continue
		}

		if len(phost) == 0 {
			// There is no host part, likely the entry is malformed; ignore.
			continue
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
		m.domainMatchers = append(m.domainMatchers, domainMatch{host: phost, port: pport, matchHost: matchHost})
	}
	return &m
}

func idnaASCII(host string) string {
	h, err := idna.Lookup.ToASCII(host)
	if err != nil {
		return host
	}
	return h
}

type hostMatcher struct {
	ipMatchers     []matcher
	domainMatchers []matcher
}

func (m *hostMatcher) Match(addr string) bool {
	if len(addr) == 0 {
		return true
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(host)

	if ip != nil {
		for _, m := range m.ipMatchers {
			if m.match(addr, port, ip) {
				return true
			}
		}
	} else {
		addr = idnaASCII(strings.ToLower(strings.TrimSpace(host)))
		for _, m := range m.domainMatchers {
			if m.match(addr, port, ip) {
				return true
			}
		}
	}
	return false
}

// matcher represents the matching rule
type matcher interface {
	// match returns true if the host and optional port or ip and optional port
	// are allowed
	match(host, port string, ip net.IP) bool
}

// allMatch matches on all possible inputs
type allMatch struct{}

func (allMatch) Match(host string) bool {
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
