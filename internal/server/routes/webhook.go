package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strings"
	"webhook/internal/config"
	"webhook/internal/events"
	"webhook/internal/structs/discord"
)

const (
	githubEvent           = "X-GitHub-Event"
	githubSignature       = "X-Hub-Signature-256"
	githubSignaturePrefix = "sha256="
	githubUserAgentPrefix = "GitHub-Hookshot/"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		w.Header().Set("Accept-Post", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	eventHeader := r.Header.Get(githubEvent)
	if eventHeader == "" {
		http.Error(w, "Missing event", http.StatusBadRequest)
		return
	}

	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" || !strings.HasPrefix(userAgent, githubUserAgentPrefix) {
		http.Error(w, "Incorrect user agent", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read body", "err", err.Error())
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			slog.Error("Failed to close body", "err", err.Error())
		}
	}()

	if secret := config.Get().Secret; secret != "" {
		sig := r.Header.Get(githubSignature)
		if sig == "" {
			http.Error(w, "Missing signature", http.StatusUnauthorized)
			return
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expectedMAC := mac.Sum(nil)

		if !strings.HasPrefix(sig, githubSignaturePrefix) {
			http.Error(w, "Invalid signature format", http.StatusBadRequest)
			return
		}

		receivedSig, err := hex.DecodeString(sig[len(githubSignaturePrefix):])
		if err != nil {
			slog.Error("Failed to decode signature", "err", err.Error())
			http.Error(w, "Invalid signature hex", http.StatusBadRequest)
			return
		}

		if !hmac.Equal(expectedMAC, receivedSig) {
			http.Error(w, "Signature mismatch", http.StatusUnauthorized)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
	go events.Parse(eventHeader, body, discord.Credentials{ID: r.PathValue("id"), Token: r.PathValue("token")})
}
