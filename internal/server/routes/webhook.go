package routes

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gohook/internal/config"
	"gohook/internal/parser"
	"gohook/internal/structs/discord"
	"mime"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	githubEvent           = "X-GitHub-Event"
	githubSignature       = "X-Hub-Signature-256"
	githubSignaturePrefix = "sha256="
	githubUserAgentPrefix = "GitHub-Hookshot/"
)

func Webhook(ctx *fasthttp.RequestCtx) {
	mediaType, _, err := mime.ParseMediaType(string(ctx.Request.Header.Peek(fasthttp.HeaderContentType)))
	if err != nil || mediaType != "application/json" {
		// w.Header().Set("Accept-Post", "application/json")
		// w.WriteHeader(http.StatusUnsupportedMediaType)
		ctx.SetStatusCode(fasthttp.StatusUnsupportedMediaType)
		return
	}

	eventHeader := ctx.Request.Header.Peek(githubEvent)
	if len(eventHeader) == 0 {
		ctx.Error("Missing event", fasthttp.StatusBadRequest)
		return
	}

	userAgent := string(ctx.Request.Header.Peek("User-Agent"))
	if userAgent == "" || !strings.HasPrefix(userAgent, githubUserAgentPrefix) {
		ctx.Error("Incorrect user agent", fasthttp.StatusBadRequest)
		return
	}

	parts := strings.Split(string(ctx.Path()), "/")

	if len(parts) != 3 || parts[1] == "" || parts[2] == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Path must be in format /:id/:token")
		return
	}

	body := ctx.PostBody()

	if secret := config.Get().Secret; secret != "" {
		sig := string(ctx.Request.Header.Peek(githubSignature))
		if sig == "" {
			ctx.Error("Missing signature", fasthttp.StatusUnauthorized)
			return
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expectedMAC := mac.Sum(nil)

		if !strings.HasPrefix(sig, githubSignaturePrefix) {
			ctx.Error("Invalid signature format", fasthttp.StatusBadRequest)
			return
		}

		receivedSig, err := hex.DecodeString(sig[len(githubSignaturePrefix):])
		if err != nil {
			fmt.Println("Failed to decode signature:", err)
			ctx.Error("Invalid signature hex", fasthttp.StatusBadRequest)
			return
		}

		if !hmac.Equal(expectedMAC, receivedSig) {
			ctx.Error("Signature mismatch", fasthttp.StatusUnauthorized)
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)

	eventCopy := string(bytes.Clone(eventHeader))
	bodyCopy := bytes.Clone(body)

	go parser.Parse(eventCopy, bodyCopy, discord.Credentials{ID: parts[1], Token: parts[2]})
}
