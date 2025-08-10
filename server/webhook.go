package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"webhook/config"
	"webhook/parser"
	"webhook/structs"

	"github.com/valyala/fasthttp"
)

func webhook(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		return
	}

	contentType := ctx.Request.Header.Peek("Content-Type")
	if len(contentType) == 0 || !bytes.HasPrefix(contentType, []byte("application/json")) {
		ctx.SetStatusCode(fasthttp.StatusUnsupportedMediaType)
		return
	}

	eventHeader := ctx.Request.Header.Peek("X-GitHub-Event")
	if len(eventHeader) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Missing event")
		return
	}

	userAgentHeader := ctx.Request.Header.Peek("User-Agent")
	if len(userAgentHeader) == 0 || !bytes.HasPrefix(userAgentHeader, []byte("GitHub-Hookshot/")) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Incorrect user agent")
		return
	}

	parts := strings.Split(string(ctx.Path()), "/")

	if len(parts) != 3 || parts[1] == "" || parts[2] == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Path must be in format /:id/:token")
		return
	}

	if config.Get().Secret != "" {
		signatureHeader := ctx.Request.Header.Peek("X-Hub-Signature-256")
		if len(signatureHeader) == 0 {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("Missing signature")
			return
		}

		body := ctx.PostBody()

		mac := hmac.New(sha256.New, []byte(config.Get().Secret))
		mac.Write(body)
		expectedMAC := mac.Sum(nil)

		const prefix = "sha256="
		sig := string(signatureHeader)
		if len(sig) <= len(prefix) || sig[:len(prefix)] != prefix {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString("Invalid signature format")
			return
		}

		receivedSig, err := hex.DecodeString(sig[len(prefix):])
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString("Invalid signature hex")
			return
		}

		if !hmac.Equal(expectedMAC, receivedSig) {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("Signature mismatch")
			return
		}
	}

	creds := structs.Credentials{
		ID:    parts[1],
		Token: parts[2],
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	eventCopy := string(bytes.Clone(eventHeader))
	bodyCopy := bytes.Clone(ctx.PostBody())
	go parser.Parse(eventCopy, bodyCopy, creds)
}
