variable "rust_log" {
  type    = string
  default = "info"
}

variable "secret" {
  type = string
}

job "foxogram-webhook" {
  datacenters = ["dc1"]

  update {
    max_parallel     = 1
    min_healthy_time = "5s"
    healthy_deadline = "30s"
    stagger          = "5s"
    auto_revert      = true
  }

  group "webhook" {
    network {
      port "http" {
        to = 8080
      }
    }

    task "webhook" {
      driver = "docker"

      config {
        image        = "foxogram/webhook:local"
        network_mode = "foxogram"
        labels = {
          "traefik.enable"                                                  = "true"
          "traefik.http.routers.foxogram-webhook.rule"                      = "Host(`webhook.foxogram.su`)"
          "traefik.http.routers.foxogram-webhook.tls.certresolver"          = "letsencrypt"
          "traefik.http.services.foxogram-webhook.loadbalancer.server.port" = "8080"
          "traefik.http.routers.foxogram-webhook.middlewares"               = "ratelimit@file"
        }
      }

      env {
        RUST_LOG = var.rust_log
        SECRET   = var.secret
      }

      service {
        name = "webhook"

        check {
          address_mode   = "driver"
          port           = "http"
          name           = "health"
          type           = "http"
          path           = "/health"
          interval       = "10s"
          timeout        = "2s"
          initial_status = "critical"
        }
      }

      restart {
        attempts = 3
        interval = "10m"
        delay    = "15s"
        mode     = "fail"
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
