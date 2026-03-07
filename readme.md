# Deadman's Snitch (Go Heartbeat Monitor)

A lightweight, concurrent heartbeat monitoring tool written in Go. It monitors services via **TCP** or **HTTP/HTTPS** and sends alerts to a **Discord Webhook** when a service becomes unreachable.

## Features

* **Multi-Protocol Support**: Easily monitor `tcp://` and `https://` endpoints.
* **Discord Integration**: Real-time incident reporting with timestamps and error details.
* **Smart Backoff**: Prevents notification spam by increasing the check interval (up to 60 minutes) during persistent outages.
* **Fault Tolerance**: Configurable `REPORTCOUNT` to avoid \"flapping\" alerts (only alerts after N consecutive failures).
* **Highly Concurrent**: Uses Go routines to monitor multiple hosts simultaneously without blocking.

---

## Configuration

The application is configured entirely via Environment Variables.

| Variable | Description | Example | Default |
| :--- | :--- | :--- | :--- |
| `HOSTS` | Comma-separated list of targets. | `https://api.test.com,tcp://1.1.1.1:53` | **Required** |
| `DISCORD_WEBHOOK` | Your Discord Webhook URL. | `https://discord.com/api/webhooks/...` | **Required** |
| `REPORTCOUNT` | Number of failed tries before alerting. | `5` | `3` |

---

## Deployment with Docker Compose

Since the `docker-compose.yml` is already included in the repository, you can start the monitor with a single command.

### 1. Update Environment
Ensure your `docker-compose.yml` (or `.env` file) contains your specific configuration.
