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
	dial := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 proxy.FromEnvironment,
			DialContext:           dial.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
}

func ExecuteWebhook(eventResult *discord.Webhook, creds discord.Credentials) error {
	url := fmt.Sprintf("%s/webhooks/%s/%s", baseURL, creds.ID, creds.Token)

	body, err := sonic.Marshal(eventResult)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return fmt.Errorf("discord webhook error: %s", buf.String())
	}

	return nil
}
