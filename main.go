package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"webhook/events"
)

type Discord struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

const githubEvent = "X-GitHub-Event"

var discordWebhookURL = os.Getenv("DISCORD_WEBHOOK")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		event := r.Header.Get(githubEvent)
		if strings.TrimSpace(event) == "" {
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Print(err)
			}
		}(r.Body)

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}

		go parseEvent(event, bodyBytes)
		w.WriteHeader(204)
	})

	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseEvent(event string, data []byte) {
	switch event {
	case "push":
		e := events.Push{}
		err := parseJSON(data, &e)
		if err != nil {
			log.Print(err)
			return
		}

		builder := strings.Builder{}

		for _, c := range e.Commits {
			builder.WriteString(fmt.Sprintf("[`%s`](<%s>) %s\n", c.Id[:7], c.Url, c.Message))
		}

		branch := strings.TrimPrefix(e.Ref, "refs/heads/")

		builder.WriteString(fmt.Sprintf(
			"\n- [%s](<%s>) on [%s](<%s>)/[%s](<%s>)",
			e.Pusher.Name,
			e.Sender.Url,
			e.Repository.Name,
			e.Repository.Url,
			branch,
			e.Repository.Url+"/tree/"+branch,
		))

		executeWebhook(builder.String(), e.Pusher.Name, e.Sender.AvatarUrl)
	}
}

func executeWebhook(content, username, avatar string) {
	body := Discord{
		Content:   content,
		Username:  username,
		AvatarURL: avatar,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Print(err)
	}

	if resp.StatusCode != 204 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}

		log.Printf("Discord Error: %s", respBody)
	}
}

func parseJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
