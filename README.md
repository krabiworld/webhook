# Gohook

![Build](https://github.com/krabiworld/gohook/actions/workflows/build.yml/badge.svg)
![License](https://img.shields.io/github/license/krabiworld/gohook)
![GitHub stars](https://img.shields.io/github/stars/krabiworld/gohook?style=social)

A lightweight GitHub webhook proxy for Discord.

## Installation

### Docker

```shell
docker run --name gohook -d -p 8080:8080 ghcr.io/krabiworld/gohook
```

### Docker Compose

```yaml
services:
  gohook:
    image: "ghcr.io/krabiworld/gohook"
    ports:
      - "8080:8080"
```

### Kubernetes

```shell
helm repo add gohook https://krabiworld.github.io/gohook
helm install gohook gohook/gohook
```

Or you can use OCI repository:

```shell
helm install gohook oci://ghcr.io/krabiworld/gohook
```

See [helm/values.yaml](helm/values.yaml)

### Precompiled binaries

You can download precompiled binaries for Linux, macOS and Windows from the [GitHub Releases](https://github.com/krabiworld/gohook/releases) page.

### Building from source

To build gohook from source, you only need Go (the version specified in [go.mod](go.mod) or later).

Start by cloning the repository:

```shell
git clone https://github.com/krabiworld/gohook.git
cd gohook
```

Then build the binary and run it:

```shell
go build ./cmd/gohook
./gohook
```

## Usage

1. Create a webhook in Discord
   1. Go to Discord channel settings -> Integrations
   2. Click Create Webhook, give it a name
   3. Copy the webhook URL
2. Configure GitHub
   1. Go to `https://github.com/<username>/<repository>/settings/hooks/new`
   2. In **Payload URL**, paste the copied Discord Webhook URL and replace `https://discord.com/api/webhooks` with your Gohook URL
    ```text
    Original Discord URL:
    https://discord.com/api/webhooks/123456/abcdef
    
    Replace with Gohook URL:
    https://gohook.your-domain.com/123456/abcdef
    ```
   3. In **Content Type**, select `application/json`
   4. In **Secret**, paste your secret (if you set environment variable `SECRET`)
   5. In **Which events would you like to trigger this webhook?**, select `Send me everything.`
   6. Click Add webhook

## Configuration

All environment variables are optional.

| Variable        | Description                                                                 | Default value | Example                                  |
|-----------------|-----------------------------------------------------------------------------|---------------|------------------------------------------|
| `LOG_LEVEL`     | How much detail to log. Options: `debug`, `info`, `warn`, `error`.          | `info`        | `debug`                                  |
| `ADDR`          | Address and port the server listens on.                                     | `:8080`       | `127.0.0.1:9000`                         |
| `SECRET`        | Shared secret to validate GitHub signatures. Leave empty to disable checks. |               | `random-string`                          |
| `HAPPY_EMOJI`   | Emoji displayed when someone stars the repository.                          |               | 🔥 or `<:foxtada:1399709119304306746>`   |
| `SUCCESS_EMOJI` | Emoji displayed for successful workflows or checks.                         |               | ✨ or `<:catgood:1399709119304306747>`    |
| `FAILURE_EMOJI` | Emoji displayed for failed workflows or checks.                             |               | 😭 or `<:catscream:1399709119304306748>` |

## Events

| API name        | UI name            | Supported actions                                    | Notes                                                                                         |
|-----------------|--------------------|------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| `check_run`     | Check runs         | `completed`                                          | Ignores these checks: `Dependabot`, `GitHub Actions` and `GitHub Advanced Security`           |
| `fork`          | Forks              |                                                      |                                                                                               |
| `issue_comment` | Issue comments     | All                                                  |                                                                                               |
| `issues`        | Issues             | All                                                  |                                                                                               |
| `public`        | Visibility changes |                                                      |                                                                                               |
| `pull_request`  | Pull requests      | All except `labeled` and `synchronize`               |                                                                                               |
| `push`          | Pushes             |                                                      |                                                                                               |
| `release`       | Releases           | `published`                                          |                                                                                               |
| `repository`    | Repositories       | `archived`, `privatized`, `renamed` and `unarchived` |                                                                                               |
| `star`          | Stars              | `created`                                            |                                                                                               |
| `workflow_run`  | Workflow runs      | `completed`                                          | Ignores these workflows: `CodeQL`, `Dependabot Updates` and `Automatic Dependency Submission` |

## Endpoints

- GET: `/health`
- POST: `/:id/:token`

## Proxy support

All standard environment variables are supported, such as `HTTP_PROXY`, `HTTPS_PROXY` and `ALL_PROXY`.
