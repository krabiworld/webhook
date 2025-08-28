package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var (
	envProxyOnce      sync.Once
	envProxyFuncValue func(*url.URL) (*url.URL, error)
)

func FromEnvironment(req *http.Request) (*url.URL, error) {
	envProxyOnce.Do(func() {
		cfg := &Config{
			HTTPSProxy: getEnvAny("HTTPS_PROXY", "https_proxy"),
			HTTPProxy:  getEnvAny("HTTP_PROXY", "http_proxy"),
			AllProxy:   getEnvAny("ALL_PROXY", "all_proxy"),
		}
		envProxyFuncValue = cfg.ProxyFunc()
	})
	return envProxyFuncValue(req.URL)
}

type Config struct {
	HTTPSProxy string
	HTTPProxy  string
	AllProxy   string
}

type config struct {
	Config

	httpsProxy *url.URL
	httpProxy  *url.URL
	allProxy   *url.URL
}

func getEnvAny(names ...string) string {
	for _, n := range names {
		if val := os.Getenv(n); val != "" {
			return val
		}
	}
	return ""
}

func (cfg *Config) ProxyFunc() func(reqURL *url.URL) (*url.URL, error) {
	// Preprocess the Config settings for more efficient evaluation.
	cfg1 := &config{
		Config: *cfg,
	}
	cfg1.init()
	return cfg1.proxyForURL
}

func (cfg *config) proxyForURL(*url.URL) (*url.URL, error) {
	var proxy *url.URL

	if cfg.httpsProxy != nil {
		proxy = cfg.httpsProxy
	} else if cfg.httpProxy != nil {
		proxy = cfg.httpProxy
	} else if cfg.allProxy != nil {
		proxy = cfg.allProxy
	}

	if proxy == nil {
		return nil, nil
	}

	return proxy, nil
}

func parseProxy(proxy string) (*url.URL, error) {
	if proxy == "" {
		return nil, nil
	}

	proxyURL, err := url.Parse(proxy)
	if err != nil || proxyURL.Scheme == "" || proxyURL.Host == "" {
		if proxyURL, err := url.Parse("http://" + proxy); err == nil {
			return proxyURL, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
	}
	return proxyURL, nil
}

func (cfg *config) init() {
	if parsed, err := parseProxy(cfg.HTTPSProxy); err == nil {
		cfg.httpsProxy = parsed
	}
	if parsed, err := parseProxy(cfg.HTTPProxy); err == nil {
		cfg.httpProxy = parsed
	}
	if parsed, err := parseProxy(cfg.AllProxy); err == nil {
		cfg.allProxy = parsed
	}
}
