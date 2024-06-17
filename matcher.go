package hostmatcher

import (
	"net"
	"strings"
)

// Matcher represents the matching rule
type Matcher interface {
	Match(host string) bool
}

func NewMatcher(hosts []string) Matcher {
	var m hostMatcher
	for _, host := range hosts {
		host = strings.TrimSpace(host)
		index := strings.Index(host, "/")
		if index == -1 {
			match, err := newHostMatcher(host)
			if err != nil {
				continue
			}
			switch match.(type) {
			case allMatch:
				return allMatcher{}
			case cidrMatch, ipMatch:
				m.ipMatchers = append(m.ipMatchers, match)
			case domainMatch:
				m.domainMatchers = append(m.domainMatchers, match)
			}
		} else {
			match, err := newHostMatcher(host[:index])
			if err != nil {
				continue
			}
			pathMatch := newPathMatcher(host[index+1:])
			switch match.(type) {
			case allMatch:
				m.ipAndPathMatchers = append(m.ipAndPathMatchers, hostAdnPathMatcher{
					matcher:     match,
					pathMatcher: pathMatch,
				})
				m.domainAndPathMatchers = append(m.domainAndPathMatchers, hostAdnPathMatcher{
					matcher:     match,
					pathMatcher: pathMatch,
				})
			case cidrMatch, ipMatch:
				m.ipAndPathMatchers = append(m.ipAndPathMatchers, hostAdnPathMatcher{
					matcher:     match,
					pathMatcher: pathMatch,
				})
			case domainMatch:
				m.domainAndPathMatchers = append(m.domainAndPathMatchers, hostAdnPathMatcher{
					matcher:     match,
					pathMatcher: pathMatch,
				})
			}
		}
	}
	return &m
}

type hostMatcher struct {
	ipMatchers     []matcher
	domainMatchers []matcher

	ipAndPathMatchers     []hostAdnPathMatcher
	domainAndPathMatchers []hostAdnPathMatcher
}

type hostAdnPathMatcher struct {
	matcher     matcher
	pathMatcher pathMatcher
}

func (m *hostMatcher) Match(addr string) bool {
	index := strings.Index(addr, "/")
	if index == -1 {
		return m.matchHost(addr)
	}
	return m.matchHostAndPath(addr[:index], addr[index+1:])
}

func (m *hostMatcher) matchHost(addr string) bool {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}

	ip := net.ParseIP(host)

	if ip != nil {
		for _, m := range m.ipMatchers {
			if m.match(addr, port, ip) {
				return true
			}
		}
	} else {
		addr = idnaASCII(strings.ToLower(host))
		for _, m := range m.domainMatchers {
			if m.match(addr, port, ip) {
				return true
			}
		}
	}
	return false
}

func (m *hostMatcher) matchHostAndPath(addr string, path string) bool {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}

	ip := net.ParseIP(host)
	if ip != nil {
		for _, m := range m.ipAndPathMatchers {
			if m.matcher.match(addr, port, ip) && m.pathMatcher.match(path) {
				return true
			}
		}
	} else {
		addr = idnaASCII(strings.ToLower(host))
		for _, m := range m.domainAndPathMatchers {
			if m.matcher.match(addr, port, ip) && m.pathMatcher.match(path) {
				return true
			}
		}
	}
	return false
}
