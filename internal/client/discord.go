package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"webhook/internal/structs/discord"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

const baseURL = "https://discord.com/api"

var client *http.Client

func Init() {
	// Add support for socks5 in http_proxy and all_proxy
	if proxy := proxyFromEnv("HTTP_PROXY", "ALL_PROXY"); proxy != "" {
		_ = os.Setenv("HTTPS_PROXY", proxy)
	}

	client = &http.Client{}

	log.Info().Msg("Client initialized")
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

func getenvInsensitive(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return os.Getenv(strings.ToLower(key))
}

func proxyFromEnv(keys ...string) string {
	for _, key := range keys {
		if val := getenvInsensitive(key); strings.HasPrefix(val, "socks5://") {
			return val
		}
	}
	return ""
}
