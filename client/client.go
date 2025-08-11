package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
	"webhook/structs"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/proxy"
)

const baseURL = "https://discord.com/api"

var client *http.Client

func Init() {
	dial := proxy.FromEnvironment().(proxy.ContextDialer)
	client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dial.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
		},
		Timeout: 10 * time.Second,
	}
}

func ExecuteWebhook(eventResult *structs.Webhook, creds structs.Credentials) error {
	url := fmt.Sprintf("%s/webhooks/%s/%s", baseURL, creds.ID, creds.Token)

	body, err := sonic.Marshal(eventResult)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
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
