package client

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"
	"webhook/proxy"
	"webhook/structs/discord"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

const baseURL = "https://discord.com/api"

var client *http.Client

func Init() {
	p := proxy.New()

	dial := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 p.Func,
			DialContext:           dial.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	log.Info().Str("proxy", p.Addr()).Msg("Client initialized")
}

func ExecuteWebhook(eventResult *discord.Webhook, creds discord.Credentials) error {
	url := fmt.Sprintf("%s/webhooks/%s/%s", baseURL, creds.ID, creds.Token)

	body, err := sonic.Marshal(eventResult)
	if err != nil {
		return fmt.Errorf("sonic.Marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return fmt.Errorf("discord api error: %s", buf.String())
	}

	return nil
}
