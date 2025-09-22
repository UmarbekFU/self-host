#!/usr/bin/env bash

# Self-Hosted Newsletter Platform Installer
# This script sets up the complete newsletter platform with Docker, Caddy, and MTA

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="newsletter-platform"
APP_USER="newsletter"
APP_DIR="/opt/newsletter"
DOCKER_COMPOSE_FILE="docker-compose-full.yml"
ENV_FILE=".env"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_error "This script should not be run as root for security reasons."
        print_error "Please run as a regular user with sudo privileges."
        exit 1
    fi
}

# Function to check if user has sudo privileges
check_sudo() {
    if ! sudo -n true 2>/dev/null; then
        print_error "This script requires sudo privileges."
        print_error "Please ensure your user has sudo access."
        exit 1
    fi
}

# Function to detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [ -f /etc/debian_version ]; then
            OS="debian"
        elif [ -f /etc/redhat-release ]; then
            OS="redhat"
        else
            OS="linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    else
        OS="unknown"
    fi
    print_status "Detected OS: $OS"
}

# Function to install Docker
install_docker() {
    if command -v docker &> /dev/null; then
        print_success "Docker is already installed"
        return
    fi
    
    print_status "Installing Docker..."
    
    case $OS in
        "debian")
            sudo apt-get update
            sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
            curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
            echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
            sudo apt-get update
            sudo apt-get install -y docker-ce docker-ce-cli containerd.io
            ;;
        "redhat")
            sudo yum install -y yum-utils
            sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
            sudo yum install -y docker-ce docker-ce-cli containerd.io
            sudo systemctl start docker
            sudo systemctl enable docker
            ;;
        "macos")
            print_warning "Please install Docker Desktop for Mac from https://www.docker.com/products/docker-desktop"
            print_warning "After installation, run this script again."
            exit 1
            ;;
        *)
            print_error "Unsupported operating system. Please install Docker manually."
            exit 1
            ;;
    esac
    
    # Add current user to docker group
    sudo usermod -aG docker $USER
    print_success "Docker installed successfully"
    print_warning "Please log out and log back in for Docker group changes to take effect."
}

# Function to install Docker Compose
install_docker_compose() {
    if command -v docker-compose &> /dev/null; then
        print_success "Docker Compose is already installed"
        return
    fi
    
    print_status "Installing Docker Compose..."
    
    # Get latest version
    COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
    
    sudo curl -L "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    
    print_success "Docker Compose installed successfully"
}

# Function to create application user
create_app_user() {
    if id "$APP_USER" &>/dev/null; then
        print_success "User $APP_USER already exists"
    else
        print_status "Creating application user: $APP_USER"
        sudo useradd -r -s /bin/false -d $APP_DIR $APP_USER
        print_success "User $APP_USER created"
    fi
}

# Function to get server IP
get_server_ip() {
    # Try to get external IP
    EXTERNAL_IP=$(curl -s ifconfig.me 2>/dev/null || curl -s ipinfo.io/ip 2>/dev/null || echo "")
    
    if [ -n "$EXTERNAL_IP" ]; then
        SERVER_IP=$EXTERNAL_IP
    else
        # Fallback to local IP
        SERVER_IP=$(hostname -I | awk '{print $1}')
    fi
    
    print_status "Detected server IP: $SERVER_IP"
}

# Function to prompt for configuration
prompt_config() {
    echo
    print_status "Newsletter Platform Configuration"
    echo "======================================"
    
    # App domain
    while true; do
        read -p "Enter your admin panel domain (e.g., panel.example.com): " APP_DOMAIN
        if [[ $APP_DOMAIN =~ ^[a-zA-Z0-9][a-zA-Z0-9.-]*[a-zA-Z0-9]$ ]]; then
            break
        else
            print_error "Please enter a valid domain name"
        fi
    done
    
    # Sending domain
    while true; do
        read -p "Enter your sending domain (e.g., news.example.com): " SENDING_DOMAIN
        if [[ $SENDING_DOMAIN =~ ^[a-zA-Z0-9][a-zA-Z0-9.-]*[a-zA-Z0-9]$ ]]; then
            break
        else
            print_error "Please enter a valid domain name"
        fi
    done
    
    # License key
    while true; do
        read -p "Enter your license key: " LICENSE_KEY
        if [ -n "$LICENSE_KEY" ]; then
            break
        else
            print_error "License key is required"
        fi
    done
    
    # Email address for admin
    while true; do
        read -p "Enter admin email address: " ADMIN_EMAIL
        if [[ $ADMIN_EMAIL =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
            break
        else
            print_error "Please enter a valid email address"
        fi
    done
    
    # Generate random password for admin
    ADMIN_PASSWORD=$(openssl rand -base64 32)
    
    print_success "Configuration collected"
}

# Function to generate DKIM keys
generate_dkim_keys() {
    print_status "Generating DKIM keys for $SENDING_DOMAIN..."
    
    # Create temporary directory for keys
    TEMP_DIR=$(mktemp -d)
    
    # Generate private key
    openssl genrsa -out "$TEMP_DIR/dkim_private.pem" 2048
    
    # Generate public key
    openssl rsa -in "$TEMP_DIR/dkim_private.pem" -pubout -out "$TEMP_DIR/dkim_public.pem"
    
    # Extract public key in DNS format
    DKIM_PUBLIC_KEY=$(openssl rsa -in "$TEMP_DIR/dkim_private.pem" -pubout -outform DER 2>/dev/null | openssl base64 -A)
    DKIM_SELECTOR="newsletter"
    
    # Store keys
    DKIM_PRIVATE_KEY=$(cat "$TEMP_DIR/dkim_private.pem")
    DKIM_PUBLIC_KEY_DNS="v=DKIM1; k=rsa; p=$DKIM_PUBLIC_KEY"
    
    # Cleanup
    rm -rf "$TEMP_DIR"
    
    print_success "DKIM keys generated"
}

# Function to create environment file
create_env_file() {
    print_status "Creating environment configuration..."
    
    cat > $ENV_FILE << EOF
# Newsletter Platform Configuration
APP_DOMAIN=$APP_DOMAIN
SENDING_DOMAIN=$SENDING_DOMAIN
LICENSE_KEY=$LICENSE_KEY
ADMIN_EMAIL=$ADMIN_EMAIL
ADMIN_PASSWORD=$ADMIN_PASSWORD
SERVER_IP=$SERVER_IP

# Database
DATABASE_URL=sqlite:///var/app/newsletter.db

# SMTP Configuration (internal to Docker network)
SMTP_HOST=mta
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=

# DKIM Configuration
DKIM_SELECTOR=$DKIM_SELECTOR
DKIM_PRIVATE_KEY='$DKIM_PRIVATE_KEY'
DKIM_PUBLIC_KEY='$DKIM_PUBLIC_KEY_DNS'

# App Configuration
PORT=8080
LOG_LEVEL=info
EOF

    print_success "Environment file created: $ENV_FILE"
}

# Function to create Docker Compose file
create_docker_compose() {
    print_status "Creating Docker Compose configuration..."
    
    cat > $DOCKER_COMPOSE_FILE << EOF
version: '3.8'

services:
  # Reverse Proxy with TLS
  proxy:
    image: caddy:2-alpine
    container_name: newsletter-proxy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./deploy/Caddyfile:/etc/caddy/Caddyfile
      - caddy-data:/data
      - caddy-config:/config
    environment:
      - APP_DOMAIN=$APP_DOMAIN
      - SENDING_DOMAIN=$SENDING_DOMAIN
    restart: unless-stopped
    depends_on:
      - app

  # Newsletter Application
  app:
    build:
      context: ./app
      dockerfile: Dockerfile
    container_name: newsletter-app
    environment:
      - DATABASE_URL=sqlite:///var/app/newsletter.db
      - PORT=8080
      - LICENSE_KEY=$LICENSE_KEY
      - SMTP_HOST=mta
      - SMTP_PORT=587
      - DKIM_SELECTOR=$DKIM_SELECTOR
      - DKIM_PRIVATE_KEY='$DKIM_PRIVATE_KEY'
    volumes:
      - app-data:/var/app
    restart: unless-stopped
    depends_on:
      - mta
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Mail Transfer Agent (MTA)
  mta:
    image: foxcpp/maddy:latest
    container_name: newsletter-mta
    volumes:
      - ./deploy/mta:/data
      - ./deploy/maddy.conf:/etc/maddy/maddy.conf
    environment:
      - SENDING_DOMAIN=$SENDING_DOMAIN
      - DKIM_SELECTOR=$DKIM_SELECTOR
      - DKIM_PRIVATE_KEY='$DKIM_PRIVATE_KEY'
    restart: unless-stopped
    ports:
      - "25:25"   # SMTP
      - "587:587" # Submission
    healthcheck:
      test: ["CMD", "maddy", "ctl", "status"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  app-data:
    driver: local
  caddy-data:
    driver: local
  caddy-config:
    driver: local

networks:
  default:
    name: newsletter-network
EOF

    print_success "Docker Compose file created: $DOCKER_COMPOSE_FILE"
}

# Function to create Caddyfile
create_caddyfile() {
    print_status "Creating Caddyfile..."
    
    mkdir -p deploy
    
    cat > deploy/Caddyfile << EOF
# Newsletter Platform Caddyfile
{$APP_DOMAIN} {
    # Admin Panel
    reverse_proxy app:8080 {
        header_up Host {host}
        header_up X-Real-IP {remote}
        header_up X-Forwarded-For {remote}
        header_up X-Forwarded-Proto {scheme}
    }
    
    # Security headers
    header {
        # Enable HSTS
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        
        # Prevent clickjacking
        X-Frame-Options "SAMEORIGIN"
        
        # Prevent MIME type sniffing
        X-Content-Type-Options "nosniff"
        
        # XSS protection
        X-XSS-Protection "1; mode=block"
        
        # Referrer policy
        Referrer-Policy "strict-origin-when-cross-origin"
    }
    
    # Logging
    log {
        output file /var/log/caddy/access.log
        format single_field common_log
    }
}

# Health check endpoint
{$APP_DOMAIN}/health {
    respond "OK" 200
}
EOF

    print_success "Caddyfile created"
}

# Function to create MTA configuration
create_mta_config() {
    print_status "Creating MTA configuration..."
    
    cat > deploy/maddy.conf << EOF
# Maddy MTA Configuration for Newsletter Platform

# Global configuration
{
    hostname = "mail.{$SENDING_DOMAIN}"
    primary_domain = "{$SENDING_DOMAIN}"
    
    # TLS configuration
    tls = "tls://0.0.0.0:465"
    tls = "tls://0.0.0.0:587"
    
    # Submission port
    submission = "tls://0.0.0.0:587"
    
    # Authentication
    auth {
        plain /etc/maddy/credentials
    }
    
    # Storage
    storage = "sqlite3:///data/maddy.db"
    
    # Queue
    queue = "memory"
    
    # Logging
    log {
        level = "info"
        output = "stdout"
    }
}

# SMTP server
(smtp) tcp://0.0.0.0:25 {
    limits {
        all rate 10 1s
        all concurrency 5
    }
    
    # DKIM signing
    dkim {
        domain {$SENDING_DOMAIN}
        selector {$DKIM_SELECTOR}
        key file /etc/maddy/dkim.key
    }
    
    # Bounce handling
    default_destination {
        deliver_to &local_routing
    }
}

# Local delivery
(local_routing) {
    deliver_to &remote_queue
}

# Remote delivery queue
(remote_queue) {
    destination {
        deliver_to &remote_smtp
    }
}

# Remote SMTP delivery
(remote_smtp) {
    limits {
        destination rate 5 1s
        destination concurrency 2
    }
    
    mx_auth {
        dane
        mta_sts
    }
}

# Bounce processing
(bounce_processor) {
    # Process bounces and send to webhook
    deliver_to &bounce_webhook
}

# Bounce webhook
(bounce_webhook) {
    # Send bounces to application webhook
    webhook {
        url "http://app:8080/api/hooks/bounce"
        method "POST"
    }
}
EOF

    # Create credentials file
    cat > deploy/credentials << EOF
# MTA credentials for internal use
newsletter:$(openssl rand -base64 32)
EOF

    # Create DKIM key file
    echo "$DKIM_PRIVATE_KEY" > deploy/dkim.key
    
    print_success "MTA configuration created"
}

# Function to create DNS records file
create_dns_records() {
    print_status "Creating DNS records file..."
    
    cat > dns-records.txt << EOF
# DNS Records for Newsletter Platform
# Add these records to your DNS provider

# SPF Record
Type: TXT
Name: $SENDING_DOMAIN
Value: v=spf1 a mx ip4:$SERVER_IP ~all

# DKIM Record
Type: TXT
Name: $DKIM_SELECTOR._domainkey.$SENDING_DOMAIN
Value: $DKIM_PUBLIC_KEY_DNS

# DMARC Record
Type: TXT
Name: _dmarc.$SENDING_DOMAIN
Value: v=DMARC1; p=quarantine; rua=mailto:dmarc@$SENDING_DOMAIN; ruf=mailto:dmarc@$SENDING_DOMAIN

# PTR Record (Reverse DNS)
# Set up reverse DNS for your server IP ($SERVER_IP) to point to: mail.$SENDING_DOMAIN
# This must be configured with your hosting provider

# Admin Panel (Optional - for subdomain)
Type: A
Name: $APP_DOMAIN
Value: $SERVER_IP

# Sending Domain (Optional - for subdomain)
Type: A
Name: $SENDING_DOMAIN
Value: $SERVER_IP
EOF
    
    print_success "DNS records file created: dns-records.txt"
}

# Function to start services
start_services() {
    print_status "Starting newsletter platform services..."
    
    # Build and start services
    docker-compose -f $DOCKER_COMPOSE_FILE up -d --build
    
    print_success "Services started successfully"
}

# Function to wait for services
wait_for_services() {
    print_status "Waiting for services to be ready..."
    
    # Wait for app to be healthy
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if docker-compose -f $DOCKER_COMPOSE_FILE ps app | grep -q "healthy"; then
            print_success "Application is ready"
            break
        fi
        
        attempt=$((attempt + 1))
        print_status "Waiting for application... (attempt $attempt/$max_attempts)"
        sleep 10
    done
    
    if [ $attempt -eq $max_attempts ]; then
        print_error "Application failed to start properly"
        print_error "Check logs with: docker-compose -f $DOCKER_COMPOSE_FILE logs"
        exit 1
    fi
}

# Function to display final instructions
display_final_instructions() {
    echo
    print_success "Newsletter Platform Installation Complete!"
    echo "=============================================="
    echo
    print_status "Next Steps:"
    echo "1. Add the DNS records from 'dns-records.txt' to your DNS provider"
    echo "2. Wait for DNS propagation (usually 5-15 minutes)"
    echo "3. Access your admin panel at: https://$APP_DOMAIN"
    echo "4. Login with:"
    echo "   Email: $ADMIN_EMAIL"
    echo "   Password: $ADMIN_PASSWORD"
    echo
    print_status "Important Files:"
    echo "- DNS Records: dns-records.txt"
    echo "- Environment: $ENV_FILE"
    echo "- Docker Compose: $DOCKER_COMPOSE_FILE"
    echo
    print_status "Useful Commands:"
    echo "- View logs: docker-compose -f $DOCKER_COMPOSE_FILE logs"
    echo "- Restart services: docker-compose -f $DOCKER_COMPOSE_FILE restart"
    echo "- Stop services: docker-compose -f $DOCKER_COMPOSE_FILE down"
    echo "- Update platform: docker-compose -f $DOCKER_COMPOSE_FILE pull && docker-compose -f $DOCKER_COMPOSE_FILE up -d"
    echo
    print_warning "Keep your admin password safe! You can change it in the admin panel."
    echo
}

# Main installation function
main() {
    echo "Self-Hosted Newsletter Platform Installer"
    echo "========================================="
    echo
    
    # Pre-flight checks
    check_root
    check_sudo
    detect_os
    
    # Install dependencies
    install_docker
    install_docker_compose
    
    # Create application user
    create_app_user
    
    # Get server information
    get_server_ip
    
    # Get configuration from user
    prompt_config
    
    # Generate DKIM keys
    generate_dkim_keys
    
    # Create configuration files
    create_env_file
    create_docker_compose
    create_caddyfile
    create_mta_config
    create_dns_records
    
    # Start services
    start_services
    wait_for_services
    
    # Display final instructions
    display_final_instructions
}

# Run main function
main "$@"