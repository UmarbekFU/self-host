# Self-Hosted Newsletter Platform

A complete, production-ready self-hosted newsletter platform that gives you full control over your email marketing without relying on external services.

## 🚀 Features

### Core Functionality
- **Self-hosted**: No external dependencies or third-party services
- **One-command setup**: Automated installer with Docker and DNS configuration
- **Modern UI**: Clean, responsive admin interface built with SvelteKit
- **Email Templates**: MJML-based responsive email templates
- **Campaign Management**: Create, schedule, and track email campaigns
- **Subscriber Management**: Import, export, and manage subscriber lists
- **Analytics**: Track opens, clicks, bounces, and unsubscribes

### Deliverability & Security
- **Built-in Deliverability Checks**: SPF, DKIM, DMARC, PTR, and TLS validation
- **DKIM Signing**: Automatic email authentication
- **Bounce Handling**: Automatic bounce processing and suppression
- **Rate Limiting**: Provider-aware throttling and warmup
- **Content Validation**: Spam filter checks and content warnings

### Technical Features
- **Docker-based**: Easy deployment and scaling
- **TLS Termination**: Automatic HTTPS with Let's Encrypt
- **Backup System**: Complete data export and restore
- **API-first**: RESTful API for all operations
- **Queue System**: Reliable background job processing

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Caddy Proxy   │────│  Go Application │────│  Maddy MTA      │
│   (TLS/HTTPS)   │    │  (API + Admin)  │    │  (SMTP + DKIM)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
    ┌─────────┐            ┌─────────┐            ┌─────────┐
    │ Browser │            │ SQLite  │            │ Internet│
    │ (Admin) │            │ (Data)  │            │ (Email) │
    └─────────┘            └─────────┘            └─────────┘
```

### Components
- **Caddy**: Reverse proxy with automatic TLS certificates
- **Go Application**: REST API, admin UI, and background workers
- **Maddy MTA**: SMTP server with DKIM signing and bounce handling
- **SQLite**: Embedded database (PostgreSQL support planned)
- **SvelteKit**: Modern frontend framework

## 🚀 Quick Start

### Prerequisites
- Linux server (Ubuntu 20.04+ recommended)
- Domain name pointing to your server
- Root or sudo access

### Installation

1. **Clone and run the installer:**
   ```bash
   git clone https://github.com/your-org/newsletter-platform.git
   cd newsletter-platform
   chmod +x scripts/install.sh
   ./scripts/install.sh
   ```

2. **Follow the setup wizard:**
   - Enter your admin panel domain (e.g., `panel.example.com`)
   - Enter your sending domain (e.g., `news.example.com`)
   - Provide your license key
   - Set admin email and password

3. **Configure DNS records:**
   The installer will generate a `dns-records.txt` file with all required DNS records:
   - SPF record
   - DKIM record
   - DMARC record
   - PTR record (reverse DNS)

4. **Access your admin panel:**
   Visit `https://your-panel-domain.com` and log in with your credentials.

### First Steps

1. **Verify Domain Setup**: Check that all DNS records are properly configured
2. **Create Mailing Lists**: Set up your subscriber lists
3. **Import Subscribers**: Upload CSV files or add subscribers manually
4. **Create Campaigns**: Use templates or create custom emails
5. **Send Test Emails**: Verify everything works before sending to your list

## 📋 System Requirements

### Minimum
- 1 CPU core
- 1GB RAM
- 10GB storage
- Ubuntu 20.04+ or similar Linux distribution

### Recommended
- 2 CPU cores
- 4GB RAM
- 50GB storage
- Ubuntu 22.04 LTS

## 🔧 Configuration

### Environment Variables

The installer creates a `.env` file with the following configuration:

```bash
# Application
APP_DOMAIN=panel.example.com
SENDING_DOMAIN=news.example.com
LICENSE_KEY=your-license-key
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=generated-password

# Database
DATABASE_URL=sqlite:///var/app/newsletter.db

# SMTP (internal)
SMTP_HOST=mta
SMTP_PORT=587

# DKIM
DKIM_SELECTOR=newsletter
DKIM_PRIVATE_KEY=generated-private-key
```

### DNS Configuration

The platform requires several DNS records for proper email deliverability:

1. **SPF Record**: Authorizes your server to send emails
2. **DKIM Record**: Provides email authentication
3. **DMARC Record**: Defines email policy and reporting
4. **PTR Record**: Reverse DNS lookup for your server IP

## 📊 Usage

### Creating Campaigns

1. Navigate to **Campaigns** in the admin panel
2. Click **Create Campaign**
3. Fill in campaign details:
   - Subject line
   - From name and email
   - HTML and text content
   - Target mailing list
4. Use the **Test Send** feature to verify your email
5. **Schedule** or **Send Now**

### Managing Subscribers

1. Go to **Lists** to manage your mailing lists
2. **Import** subscribers from CSV files
3. **Export** subscriber data
4. View subscriber status and engagement

### Monitoring Deliverability

1. Check the **Domains** page for DNS configuration status
2. Monitor bounce rates and complaints
3. Review campaign analytics and engagement metrics

## 🛠️ Development

### Backend Development

```bash
cd app
go mod download
go run cmd/server/main.go
```

### Frontend Development

```bash
cd app/web
npm install
npm run dev
```

### Building for Production

```bash
# Build Go application
cd app
go build -o newsletter cmd/server/main.go

# Build frontend
cd app/web
npm run build
```

## 📦 Docker Deployment

### Using Docker Compose

```bash
# Production deployment
docker-compose -f docker-compose-full.yml up -d

# Development with MailHog
docker-compose -f docker-compose-simple.yml up -d
```

### Manual Docker Build

```bash
# Build application image
docker build -t newsletter-platform ./app

# Run with external MTA
docker run -d \
  --name newsletter-app \
  -p 8080:8080 \
  -e DATABASE_URL=sqlite:///var/app/newsletter.db \
  -e LICENSE_KEY=your-license-key \
  newsletter-platform
```

## 🔒 Security

### Built-in Security Features
- HTTPS enforcement via Caddy
- CSRF protection
- Secure session management
- Input validation and sanitization
- Rate limiting on public endpoints

### Best Practices
- Keep your system updated
- Use strong passwords
- Regularly backup your data
- Monitor logs for suspicious activity
- Keep DKIM keys secure

## 📈 Monitoring & Maintenance

### Health Checks
- Application health: `https://your-domain.com/api/health`
- Database connectivity
- SMTP server status
- DNS configuration validation

### Backup & Restore

```bash
# Create backup
sudo ./scripts/backup.sh

# Restore from backup
tar -xzf newsletter_backup_YYYYMMDD_HHMMSS.tar.gz
# Follow instructions in MANIFEST.txt
```

### Logs

```bash
# View application logs
docker-compose logs -f app

# View MTA logs
docker-compose logs -f mta

# View proxy logs
docker-compose logs -f proxy
```

## 🚨 Troubleshooting

### Common Issues

1. **DNS Records Not Working**
   - Wait for DNS propagation (up to 48 hours)
   - Verify records with online DNS checkers
   - Check for typos in domain names

2. **Emails Going to Spam**
   - Verify all DNS records are correct
   - Check bounce rates and complaints
   - Review email content for spam triggers

3. **Application Won't Start**
   - Check Docker logs: `docker-compose logs app`
   - Verify environment variables
   - Ensure ports are not in use

4. **Database Issues**
   - Check disk space
   - Verify file permissions
   - Restore from backup if needed

### Getting Help

- Check the logs for error messages
- Review the troubleshooting section
- Open an issue on GitHub with detailed information

## 🔄 Updates

### Updating the Platform

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
docker-compose -f docker-compose-full.yml down
docker-compose -f docker-compose-full.yml up -d --build
```

### Database Migrations

The application automatically runs database migrations on startup. No manual intervention required.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 🙏 Acknowledgments

- [Maddy](https://github.com/foxcpp/maddy) - Modern MTA
- [Caddy](https://caddyserver.com/) - Web server with automatic HTTPS
- [SvelteKit](https://kit.svelte.dev/) - Frontend framework
- [MJML](https://mjml.io/) - Email template framework

---

**Built with ❤️ for the self-hosted community**