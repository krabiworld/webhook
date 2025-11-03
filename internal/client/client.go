package client

import (
	"encoding/json"
	"fmt"
	"gohook/internal/config"
	"gohook/internal/structs/discord"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

const baseURL = "https://discord.com/api"

var client *fasthttp.Client

func Init() {
	client = &fasthttp.Client{
		ReadTimeout:              time.Millisecond * 500,
		WriteTimeout:             time.Millisecond * 500,
		MaxConnDuration:          time.Hour,
		NoDefaultUserAgentHeader: true,
	}

	proxy := config.Get().Proxy
	if proxy != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialerTimeout(proxy, time.Second*2)
	}
}

func ExecuteWebhook(eventResult *discord.Webhook, creds discord.Credentials) error {
	url := fmt.Sprintf("%s/webhooks/%s/%s", baseURL, creds.ID, creds.Token)

	body, err := json.Marshal(eventResult)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBodyRaw(body)

	resp := fasthttp.AcquireResponse()

	for i := range 10 {
		err = client.Do(req, resp)
		if err == nil {
			break
		}
		fmt.Println("Request failed, retrying... Attempt", i+1, "Error", err)
		time.Sleep(time.Second)
	}

	fasthttp.ReleaseRequest(req)

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusNoContent {
		return fmt.Errorf("discord api error: %s", string(resp.Body()))
	}

	return nil
}
