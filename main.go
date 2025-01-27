package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"webhook/events"
)

type Discord struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type Credentials struct {
	ID    string
	Token string
}

const (
	discordBaseURL = "https://discord.com/api"
	githubEvent    = "X-GitHub-Event"
)

func main() {
	http.HandleFunc("/{id}/{token}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		creds := Credentials{
			ID:    r.PathValue("id"),
			Token: r.PathValue("token"),
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

		go parseEvent(event, bodyBytes, creds)
		w.WriteHeader(204)
	})

	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseEvent(event string, data []byte, creds Credentials) {
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

		executeWebhook(builder.String(), e.Pusher.Name, e.Sender.AvatarUrl, creds)
	case "workflow_run":
		e := events.WorkflowRun{}
		err := parseJSON(data, &e)
		if err != nil {
			log.Print(err)
			return
		}

		if e.Action != "completed" {
			return
		}

		content := fmt.Sprintf("Workflow %s on [%s](<%s>)/[%s](<%s>)",
			e.WorkflowRun.Conclusion,
			e.Repository.Name,
			e.Repository.Url,
			e.WorkflowRun.HeadBranch,
			e.Repository.Url+"/tree/"+e.WorkflowRun.HeadBranch,
		)

		executeWebhook(content, e.Workflow.Name, e.Sender.AvatarUrl, creds)
	}
}

func executeWebhook(content, username, avatar string, creds Credentials) {
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

	url := fmt.Sprintf("%s/webhooks/%s/%s", discordBaseURL, creds.ID, creds.Token)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
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
