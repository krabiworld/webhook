# Webhook

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
    environment:
      - SECRET=random-string
```

**Kubernetes**

```shell
helm install webhook oci://ghcr.io/krabiworld/webhook
```

## Configuration

See [.env.example](.env.example) for example

## Endpoints

- `/:id/:token?ignoreWorkflows=CodeQL`
- `/health`

## Proxy support

All standard environment variables are supported, such as `HTTP_PROXY`, `HTTPS_PROXY` and `ALL_PROXY`. SOCKS5 only works in `ALL_PROXY`
