package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime"
	"net/http"
	"strings"
	"webhook/config"
	"webhook/context"
	"webhook/events"
	"webhook/structs/discord"

	"github.com/rs/zerolog/log"
)

const (
	githubEvent           = "X-GitHub-Event"
	githubSignature       = "X-Hub-Signature-256"
	githubSignaturePrefix = "sha256="
	githubUserAgentPrefix = "GitHub-Hookshot/"
)

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		log.Error().Err(err).Msg("Invalid media type")
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

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.Error(w, "Path must be in format /:id/:token", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read body")
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close body")
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
			log.Error().Err(err).Msg("Failed to decode signature")
			http.Error(w, "Invalid signature hex", http.StatusBadRequest)
			return
		}

		if !hmac.Equal(expectedMAC, receivedSig) {
			http.Error(w, "Signature mismatch", http.StatusUnauthorized)
			return
		}
	}

	ctx := context.NewContext(r.URL.Query())

	w.WriteHeader(http.StatusNoContent)
	go events.Parse(eventHeader, body, ctx, discord.Credentials{ID: parts[0], Token: parts[1]})
}
