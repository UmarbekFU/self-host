# Self-Hosted Newsletter Platform

A self-hosted web application for sending newsletter campaigns with solid deliverability basics, simple setup, and no external services.

## Features

- One-command installer with Docker Compose
- Admin web UI for campaign management
- Outbound email pipeline with DKIM signing
- DNS verification and deliverability checks
- Analytics and tracking
- Full data export capabilities
- No third-party dependencies

## Quick Start

```bash
curl -sSL https://raw.githubusercontent.com/your-repo/install.sh | bash
```

## Architecture

- **Backend**: Go HTTP API server
- **Frontend**: SvelteKit admin UI
- **Database**: SQLite (with Postgres option)
- **MTA**: Maddy for SMTP delivery
- **Proxy**: Caddy for TLS termination
- **Templates**: MJML for responsive emails

## License

One-time purchase license. See LICENSE file for details.
