package proxy

import (
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	envProxyOnce      sync.Once
	envProxyFuncValue func(*url.URL) (*url.URL, error)
)

func FromEnvironment(req *http.Request) (*url.URL, error) {
	envProxyOnce.Do(func() {
		cfg := &config{}
		cfg.init()
		envProxyFuncValue = cfg.proxyForURL
	})
	return envProxyFuncValue(req.URL)
}

type config struct {
	proxy *url.URL
}

func (cfg *config) proxyForURL(*url.URL) (*url.URL, error) {
	if cfg.proxy != nil {
		return cfg.proxy, nil
	}

	return nil, nil
}

func (cfg *config) init() {
	if parsed := parseProxy(getEnvAny("HTTPS_PROXY", "https_proxy")); parsed != nil {
		cfg.proxy = parsed
		return
	}
	if parsed := parseProxy(getEnvAny("HTTP_PROXY", "http_proxy")); parsed != nil {
		cfg.proxy = parsed
		return
	}
	if parsed := parseProxy(getEnvAny("ALL_PROXY", "all_proxy")); parsed != nil {
		cfg.proxy = parsed
		return
	}
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
