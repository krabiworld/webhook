package client

import (
	"fmt"
	"time"
	"webhook/config"
	"webhook/structs"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

var client *fasthttp.Client

func Init() {
	client = &fasthttp.Client{
		ReadTimeout:                   time.Millisecond * 500,
		WriteTimeout:                  time.Millisecond * 500,
		MaxConnDuration:               time.Hour,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}

	proxy := config.Get().Proxy
	if proxy != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialerTimeout(proxy, time.Second*2)
	}
}

func ExecuteWebhook(eventResult *structs.Webhook, creds structs.Credentials) error {
	url := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", creds.ID, creds.Token)

	body, err := sonic.Marshal(eventResult)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBodyRaw(body)

	resp := fasthttp.AcquireResponse()
	if err := client.Do(req, resp); err != nil {
		return err
	}

	fasthttp.ReleaseRequest(req)

	if resp.StatusCode() != fasthttp.StatusNoContent {
		return fmt.Errorf(string(resp.Body()))
	}

	fasthttp.ReleaseResponse(resp)

	return nil
}
