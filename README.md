# Webhook

![Build](https://github.com/krabiworld/webhook/actions/workflows/test.yml/badge.svg)
![License](https://img.shields.io/github/license/krabiworld/webhook)
![GitHub stars](https://img.shields.io/github/stars/krabiworld/webhook?style=social)

A beautiful GitHub webhook proxy for Discord

## Installation

**Docker:**

```shell
docker run --name webhook -d -p 8080:8080 ghcr.io/krabiworld/webhook
```

**Docker Compose:**

```yaml
services:
  webhook:
    image: "ghcr.io/krabiworld/webhook"
    ports:
      - "8080:8080"
```

**Kubernetes**

```shell
helm install webhook oci://ghcr.io/krabiworld/webhook
```

See [helm/values.yaml](helm/values.yaml)

## Usage

1. Create a webhook in Discord
   1. Go to Discord channel settings -> Integrations
   2. Click Create Webhook, give it a name
   3. Copy the webhook URL
2. Configure GitHub
   1. Go to `https://github.com/<username>/<repository>/settings/hooks/new`
   2. In **Payload URL**, paste copied Discord Webhook URL and replace `https://discord.com/api/webhooks` with your webhook proxy URL
    ```text
    Original Discord URL:
    https://discord.com/api/webhooks/123456/abcdef
    
    Replace with proxy URL:
    https://webhook.your-domain.com/123456/abcdef
    ```
   3. In **Content Type**, select `application/json`
   4. In **Secret**, paste your secret (if you set environment variable `SECRET`)
   5. In **Which events would you like to trigger this webhook?**, select `Send me everything.`
   6. Click Add webhook

## Configuration

All environment variables are optional.

| Variable                    | Description                                                                                  | Default value | Example                                                       |
|-----------------------------|----------------------------------------------------------------------------------------------|---------------|---------------------------------------------------------------|
| LOG_LEVEL                   | How much detail to log. Options: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic` | `info`        | `debug`                                                       |
| LOG_MODE                    | Log output format: machine-friendly `json` or human-friendly `pretty`.                       | `json`        | `pretty`                                                      |
| ADDR                        | Address and port the server listens on.                                                      | `:8080`       | `127.0.0.1:9000`                                              |
| SECRET                      | Shared secret to validate GitHub signatures. Leave empty to disable checks.                  |               | `random-string`                                               |
| STORAGE_BACKEND             | Where to keep temporary cache. Options: `memory` (in-RAM) or `redis`.                        | `memory`      | `redis`                                                       |
| REDIS_URL                   | Redis connection string (used if `STORAGE_BACKEND=redis`).                                   |               | `redis://<user>:<pass>@localhost:6379/<db>`                   |
| HAPPY_EMOJI                 | Emoji shown when someone stars the repo.                                                     |               | 🔥 or `<:foxtada:1399709119304306746>`                        |
| SUCCESS_EMOJI               | Emoji for successful workflows or checks.                                                    |               | ✨ or `<:catgood:1399709119304306747>`                         |
| FAILURE_EMOJI               | Emoji for failed workflows or checks.                                                        |               | 😭 or `<:catscream:1399709119304306748>`                      |
| DISABLED_EVENTS             | Comma-separated list of GitHub events to ignore completely.                                  |               | `release,fork`                                                |
| IGNORE_PRIVATE_REPOSITORIES | Skip events from private repositories.                                                       | `false`       | `true`                                                        |
| IGNORED_REPOSITORIES        | Comma-separated list of repos to ignore.                                                     |               | `torvalds/linux,rust-lang/rust`                               |
| IGNORED_WORKFLOWS           | List of GitHub Actions workflows to ignore globally.                                         |               | `"CodeQL,Automatic Dependency Submission,Dependabot Updates"` |

## Endpoints

- GET `/health`
- POST `/:id/:token`

Query parameters:

| Parameter        | Description                         | Example                |
|------------------|-------------------------------------|------------------------|
| ignoredEvents    | Comma-separated events to ignore    | check_run,workflow_run |
| ignoredChecks    | Comma-separated checks to ignore    | Cloudflare,Vercel      |
| ignoredWorkflows | Comma-separated workflows to ignore | CodeQL,Lint            |

## Implemented events

- `check_run`
- `fork`
- `issue_comment`
- `issues`
- `pull_request`
- `push`
- `release`
- `star`
- `workflow_run`

## Proxy support

All standard environment variables are supported, such as `HTTP_PROXY`, `HTTPS_PROXY` and `ALL_PROXY`. SOCKS5 only works in `ALL_PROXY`
