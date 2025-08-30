package proxy

import (
	"net/http"
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
)

type Proxy struct {
	proxyAddr string
	proxyFunc func(*url.URL) (*url.URL, error)
}

func New() *Proxy {
	proxy := parse()

	var proxyAddr string
	var proxyFunc func(*url.URL) (*url.URL, error)

	if proxy != nil {
		proxyAddr = proxy.String()
		proxyFunc = func(url *url.URL) (*url.URL, error) {
			return proxy, nil
		}
	} else {
		proxyFunc = func(url *url.URL) (*url.URL, error) {
			return nil, nil
		}
	}

	return &Proxy{proxyAddr: proxyAddr, proxyFunc: proxyFunc}
}

func (p *Proxy) Addr() string {
	return p.proxyAddr
}

func (p *Proxy) Func(req *http.Request) (*url.URL, error) {
	return p.proxyFunc(req.URL)
}

func parse() *url.URL {
	if parsed := parseProxy(getEnvAny("HTTPS_PROXY", "https_proxy")); parsed != nil {
		return parsed
	}
	if parsed := parseProxy(getEnvAny("HTTP_PROXY", "http_proxy")); parsed != nil {
		return parsed
	}
	if parsed := parseProxy(getEnvAny("ALL_PROXY", "all_proxy")); parsed != nil {
		return parsed
	}
	return nil
}

func getEnvAny(names ...string) string {
	for _, n := range names {
		if val := os.Getenv(n); val != "" {
			return val
		}
	}
	return ""
}

func parseProxy(proxy string) *url.URL {
	if proxy == "" {
		return nil
	}

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Error().Err(err).Str("url", proxy).Msg("failed to parse url")
		return nil
	}

	return proxyURL
}
