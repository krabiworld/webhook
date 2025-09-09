package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
	"webhook/internal/structs/discord"
)

const baseURL = "https://discord.com/api"

var client = &http.Client{}

func init() {
	// Add support for socks5 in http_proxy and all_proxy
	for _, key := range []string{"HTTP_PROXY", "ALL_PROXY"} {
		var val string
		if v, ok := os.LookupEnv(key); ok {
			val = v
		} else {
			val = os.Getenv(strings.ToLower(key))
		}

		if strings.HasPrefix(val, "socks5://") {
			_ = os.Setenv("HTTPS_PROXY", val)
			break
		}
	}
}

func ExecuteWebhook(eventResult *discord.Webhook, creds discord.Credentials) error {
	url := fmt.Sprintf("%s/webhooks/%s/%s", baseURL, creds.ID, creds.Token)

	body, err := json.Marshal(eventResult)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response

	for i := 0; i < 10; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		slog.Warn("request failed, retrying...", "attempt", i+1, "err", err.Error())
		time.Sleep(time.Second)
	}

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Failed to close response body", "err", err.Error())
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return fmt.Errorf("discord api error: %s", buf.String())
	}

	return nil
}
