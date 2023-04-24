//go:build with_proxyprovider

package proxyprovider

import (
	"context"
	mDNS "github.com/miekg/dns"
	dns "github.com/sagernet/sing-dns"
	dnsQUIC "github.com/sagernet/sing-dns/quic"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	"net"
	"net/netip"
	"net/url"
	"strings"
	"sync"
	"time"
)

const defaultDNS = "223.5.5.5:53"

func (p *ProxyProvider) initDNS() error {
	switch {
	case strings.HasPrefix(p.options.DNS, "udp://"):
		addr := strings.Replace(p.options.DNS, "udp://", "", 1)

		ip, err := netip.ParseAddr(strings.Trim(addr, "[]"))
		if err == nil {
			addr = net.JoinHostPort(ip.String(), "53")
			dnsTransport, err := dns.NewUDPTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			if host == "" {
				return E.New("invalid dns address: ", addr)
			}

			ip, err = netip.ParseAddr(host)
			if err != nil {
				return E.New("invalid dns address: ", addr)
			}

			addr = net.JoinHostPort(ip.String(), port)
			dnsTransport, err := dns.NewUDPTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		return E.New("invalid dns address: ", addr)
	case strings.HasPrefix(p.options.DNS, "tcp://"):
		addr := strings.Replace(p.options.DNS, "tcp://", "", 1)

		ip, err := netip.ParseAddr(strings.Trim(addr, "[]"))
		if err == nil {
			addr = net.JoinHostPort(ip.String(), "53")
			dnsTransport, err := dns.NewTCPTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			if host == "" {
				return E.New("invalid dns address: ", addr)
			}

			ip, err = netip.ParseAddr(host)
			if err != nil {
				return E.New("invalid dns address: ", addr)
			}

			addr = net.JoinHostPort(ip.String(), port)
			dnsTransport, err := dns.NewTCPTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		return E.New("invalid dns address: ", addr)
	case strings.HasPrefix(p.options.DNS, "tls://"):
		addr := strings.Replace(p.options.DNS, "tls://", "", 1)

		ip, err := netip.ParseAddr(strings.Trim(addr, "[]"))
		if err == nil {
			addr = net.JoinHostPort(ip.String(), "853")
			dnsTransport, err := dns.NewTLSTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			if host == "" {
				return E.New("invalid dns address: ", addr)
			}

			ip, err = netip.ParseAddr(host)
			if err != nil {
				return E.New("invalid dns address: ", addr)
			}

			addr = net.JoinHostPort(ip.String(), port)
			dnsTransport, err := dns.NewTLSTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		return E.New("invalid dns address: ", addr)
	case strings.HasPrefix(p.options.DNS, "https://"):
		u, err := url.Parse(p.options.DNS)
		if err != nil {
			return E.New("invalid dns address: ", p.options.DNS)
		}

		port := u.Port()
		if port == "" {
			port = "443"
		}

		ip, err := netip.ParseAddr(u.Hostname())
		if err != nil {
			return E.New("invalid dns address: ", p.options.DNS)
		}

		u.Host = net.JoinHostPort(ip.String(), port)

		u.Fragment = ""

		dnsTransport := dns.NewHTTPSTransport("proxy-provider-dns", p.dialer, u.String())
		p.dnsTransport = dnsTransport
		return nil
	case strings.HasPrefix(p.options.DNS, "h3://"):
		urlStr := strings.Replace(p.options.DNS, "h3://", "https://", 1)

		u, err := url.Parse(urlStr)
		if err != nil {
			return E.New("invalid dns address: ", p.options.DNS)
		}

		port := u.Port()
		if port == "" {
			port = "443"
		}

		ip, err := netip.ParseAddr(u.Hostname())
		if err != nil {
			return E.New("invalid dns address: ", p.options.DNS)
		}

		u.Host = net.JoinHostPort(ip.String(), port)

		u.Fragment = ""
		u.Scheme = "https"

		dnsTransport := dnsQUIC.NewHTTP3Transport("proxy-provider-dns", p.dialer, u.String())
		p.dnsTransport = dnsTransport
		return nil
	case strings.HasPrefix(p.options.DNS, "quic://"):
		addr := strings.Replace(p.options.DNS, "quic://", "", 1)

		ip, err := netip.ParseAddr(strings.Trim(addr, "[]"))
		if err == nil {
			addr = net.JoinHostPort(ip.String(), "784")
			dnsTransport, err := dnsQUIC.NewTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			if host == "" {
				return E.New("invalid dns address: ", addr)
			}

			ip, err = netip.ParseAddr(host)
			if err != nil {
				return E.New("invalid dns address: ", addr)
			}

			addr = net.JoinHostPort(ip.String(), port)
			dnsTransport, err := dnsQUIC.NewTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		return E.New("invalid dns address: ", addr)
	case p.options.DNS == "":
		dnsTransport, err := dns.NewUDPTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(defaultDNS))
		if err != nil {
			return E.Cause(err, "create dns transport")
		}
		p.dnsTransport = dnsTransport
		return nil
	default:
		ip, err := netip.ParseAddr(p.options.DNS)
		if err == nil {
			addr := net.JoinHostPort(ip.String(), "53")
			dnsTransport, err := dns.NewUDPTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		host, port, err := net.SplitHostPort(p.options.DNS)
		if err == nil {
			if host == "" {
				return E.New("invalid dns address: ", p.options.DNS)
			}

			ip, err = netip.ParseAddr(host)
			if err != nil {
				return E.New("invalid dns address: ", p.options.DNS)
			}

			addr := net.JoinHostPort(ip.String(), port)
			dnsTransport, err := dns.NewUDPTransport("proxy-provider-dns", p.ctx, p.dialer, M.ParseSocksaddr(addr))
			if err != nil {
				return E.Cause(err, "create dns transport")
			}
			p.dnsTransport = dnsTransport
			return nil
		}

		return E.New("invalid dns address: ", p.options.DNS)
	}
}

func (p *ProxyProvider) queryWrapper(domain string, t uint16) ([]netip.Addr, error) {
	dnsMsg := new(mDNS.Msg)
	dnsMsg.SetQuestion(mDNS.Fqdn(domain), mDNS.TypeA)

	ctx, cancel := context.WithTimeout(p.ctx, time.Second*20)
	defer cancel()
	respMsg, err := p.dnsTransport.Exchange(ctx, dnsMsg)
	if err != nil {
		return nil, err
	}

	ips := make([]netip.Addr, 0)

	for _, answer := range respMsg.Answer {
		if a, ok := answer.(*mDNS.A); ok {
			nIP, _ := netip.ParseAddr(a.A.String())
			ips = append(ips, nIP)
		}
	}

	if len(ips) == 0 {
		return nil, E.New("no ip found")
	}

	return ips, nil
}

func (p *ProxyProvider) queryA(domain string) ([]netip.Addr, error) {
	return p.queryWrapper(domain, mDNS.TypeA)
}

func (p *ProxyProvider) queryAAAA(domain string) ([]netip.Addr, error) {
	return p.queryWrapper(domain, mDNS.TypeAAAA)
}

func (p *ProxyProvider) query(domain string) ([]netip.Addr, error) {
	ch := make(chan []netip.Addr, 2)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ips, err := p.queryA(domain)
		if err != nil {
			return
		}
		ch <- ips
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ips, err := p.queryAAAA(domain)
		if err != nil {
			return
		}
		ch <- ips
	}()
	wg.Wait()
	ips := make([]netip.Addr, 0)
	for {
		select {
		case ip := <-ch:
			ips = append(ips, ip...)
		default:
			if len(ips) == 0 {
				return nil, E.New("no ip found")
			}

			return ips, nil
		}
	}
}
