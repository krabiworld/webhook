# Webhook

A beautiful GitHub webhook proxy for Discord

## Installation

```shell
docker run --name webhook -d -p 8080:8080 ghcr.io/krabiworld/webhook
```

## Configuration

See [.env.example](.env.example) for example

## Proxy support

All standard environment variables are supported, such as `HTTP_PROXY`, `HTTPS_PROXY` and `ALL_PROXY`. SOCKS5 only works in `ALL_PROXY`
