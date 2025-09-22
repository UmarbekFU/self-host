#!/bin/bash

# Newsletter Platform Installer
# One-command installer for self-hosted newsletter platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="/opt/newsletter"
DOCKER_COMPOSE_VERSION="2.20.0"
APP_VERSION="v1.0.0"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "This script should not be run as root for security reasons."
        log_info "Please run as a regular user with sudo privileges."
        exit 1
    fi
}

check_dependencies() {
    log_info "Checking system dependencies..."
    
    # Check if running on supported OS
    if [[ "$OSTYPE" != "linux-gnu"* ]]; then
        log_error "This installer only supports Linux systems."
        exit 1
    fi
    
    # Check if curl is available
    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed."
        exit 1
    fi
    
    # Check if wget is available
    if ! command -v wget &> /dev/null; then
        log_error "wget is required but not installed."
        exit 1
    fi
    
    log_success "System dependencies check passed."
}

install_docker() {
    log_info "Installing Docker..."
    
    if command -v docker &> /dev/null; then
        log_info "Docker is already installed."
        return
    fi
    
    # Install Docker
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh
    
    # Add current user to docker group
    sudo usermod -aG docker $USER
    
    log_success "Docker installed successfully."
    log_warning "Please log out and log back in for group changes to take effect."
}

install_docker_compose() {
    log_info "Installing Docker Compose..."
    
    if command -v docker-compose &> /dev/null; then
        log_info "Docker Compose is already installed."
        return
    fi
    
    # Download Docker Compose
    sudo curl -L "https://github.com/docker/compose/releases/download/v${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    
    # Create symlink
    sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    
    log_success "Docker Compose installed successfully."
}

create_system_user() {
    log_info "Creating system user for newsletter platform..."
    
    if id "newsletter" &>/dev/null; then
        log_info "User 'newsletter' already exists."
    else
        sudo useradd -r -s /bin/false -d $INSTALL_DIR newsletter
        log_success "User 'newsletter' created."
    fi
}

get_server_ip() {
    # Try to get public IP
    SERVER_IP=$(curl -s ifconfig.me || curl -s ipinfo.io/ip || curl -s icanhazip.com || echo "YOUR_SERVER_IP")
    echo $SERVER_IP
}

prompt_configuration() {
    log_info "Please provide the following configuration:"
    echo
    
    # Admin domain
    read -p "Admin domain (e.g., panel.example.com): " ADMIN_DOMAIN
    if [[ -z "$ADMIN_DOMAIN" ]]; then
        log_error "Admin domain is required."
        exit 1
    fi
    
    # Sending domain
    read -p "Sending domain (e.g., news.example.com): " SENDING_DOMAIN
    if [[ -z "$SENDING_DOMAIN" ]]; then
        log_error "Sending domain is required."
        exit 1
    fi
    
    # License key
    read -p "License key: " LICENSE_KEY
    if [[ -z "$LICENSE_KEY" ]]; then
        log_error "License key is required."
        exit 1
    fi
    
    # Server IP
    SERVER_IP=$(get_server_ip)
    read -p "Server IP address [$SERVER_IP]: " INPUT_IP
    if [[ -n "$INPUT_IP" ]]; then
        SERVER_IP="$INPUT_IP"
    fi
    
    # Email
    read -p "Admin email address: " ADMIN_EMAIL
    if [[ -z "$ADMIN_EMAIL" ]]; then
        log_error "Admin email is required."
        exit 1
    fi
    
    log_success "Configuration collected."
}

generate_dkim_keys() {
    log_info "Generating DKIM keys..."
    
    # Generate DKIM keys using OpenSSL
    DKIM_SELECTOR="newsletter"
    DKIM_PRIVATE_KEY_FILE="$INSTALL_DIR/mta/keys/dkim.key"
    DKIM_PUBLIC_KEY_FILE="$INSTALL_DIR/mta/keys/dkim.pub"
    
    # Create keys directory
    sudo mkdir -p "$INSTALL_DIR/mta/keys"
    
    # Generate private key
    sudo openssl genrsa -out "$DKIM_PRIVATE_KEY_FILE" 2048
    sudo chmod 600 "$DKIM_PRIVATE_KEY_FILE"
    
    # Generate public key
    sudo openssl rsa -in "$DKIM_PRIVATE_KEY_FILE" -pubout -out "$DKIM_PUBLIC_KEY_FILE"
    sudo chmod 644 "$DKIM_PUBLIC_KEY_FILE"
    
    # Extract public key for DNS record
    DKIM_PUBLIC_KEY=$(sudo openssl rsa -in "$DKIM_PRIVATE_KEY_FILE" -pubout -outform DER | openssl base64 -A)
    
    log_success "DKIM keys generated."
}

generate_tls_cert() {
    log_info "Generating self-signed TLS certificate..."
    
    TLS_CERT_FILE="$INSTALL_DIR/mta/keys/tls.crt"
    TLS_KEY_FILE="$INSTALL_DIR/mta/keys/tls.key"
    
    sudo openssl req -x509 -newkey rsa:4096 -keyout "$TLS_KEY_FILE" -out "$TLS_CERT_FILE" -days 365 -nodes \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=$SENDING_DOMAIN"
    
    sudo chmod 600 "$TLS_KEY_FILE"
    sudo chmod 644 "$TLS_CERT_FILE"
    
    log_success "TLS certificate generated."
}

create_directories() {
    log_info "Creating installation directories..."
    
    sudo mkdir -p "$INSTALL_DIR"/{app,deploy,mta/keys,scripts,backups}
    sudo chown -R newsletter:newsletter "$INSTALL_DIR"
    
    log_success "Directories created."
}

copy_files() {
    log_info "Copying application files..."
    
    # Copy deployment files
    sudo cp -r deploy/* "$INSTALL_DIR/deploy/"
    sudo cp -r scripts/* "$INSTALL_DIR/scripts/"
    
    # Set permissions
    sudo chown -R newsletter:newsletter "$INSTALL_DIR"
    sudo chmod +x "$INSTALL_DIR/scripts"/*.sh
    
    log_success "Files copied."
}

generate_configs() {
    log_info "Generating configuration files..."
    
    # Generate .env file
    cat > "$INSTALL_DIR/.env" << EOF
# Newsletter Platform Configuration
ADMIN_DOMAIN=$ADMIN_DOMAIN
SENDING_DOMAIN=$SENDING_DOMAIN
LICENSE_KEY=$LICENSE_KEY
SERVER_IP=$SERVER_IP
ADMIN_EMAIL=$ADMIN_EMAIL
DKIM_SELECTOR=newsletter
MADDY_HOSTNAME=$SENDING_DOMAIN
EOF
    
    # Generate Caddyfile
    sed -i "s/\${ADMIN_DOMAIN:panel.example.com}/$ADMIN_DOMAIN/g" "$INSTALL_DIR/deploy/Caddyfile"
    sed -i "s/\${SENDING_DOMAIN:news.example.com}/$SENDING_DOMAIN/g" "$INSTALL_DIR/deploy/Caddyfile"
    
    # Generate Maddy config
    sed -i "s/\${env.SENDING_DOMAIN}/$SENDING_DOMAIN/g" "$INSTALL_DIR/deploy/mta/maddy.conf"
    sed -i "s/\${env.MADDY_HOSTNAME}/$SENDING_DOMAIN/g" "$INSTALL_DIR/deploy/mta/maddy.conf"
    sed -i "s/\${env.DKIM_SELECTOR}/newsletter/g" "$INSTALL_DIR/deploy/mta/maddy.conf"
    
    log_success "Configuration files generated."
}

build_and_deploy() {
    log_info "Building and deploying application..."
    
    cd "$INSTALL_DIR"
    
    # Build the application
    log_info "Building Go application..."
    docker build -t newsletter:latest ../app/
    
    # Start services
    log_info "Starting services..."
    docker-compose -f deploy/docker-compose.yml up -d
    
    log_success "Application deployed successfully."
}

generate_dns_records() {
    log_info "Generating DNS records..."
    
    # Get DKIM public key
    DKIM_PUBLIC_KEY=$(sudo openssl rsa -in "$INSTALL_DIR/mta/keys/dkim.key" -pubout -outform DER | openssl base64 -A)
    
    # Generate DNS records file
    cat > "$INSTALL_DIR/dns-records.txt" << EOF
# DNS Records for Newsletter Platform
# Add these records to your DNS provider

# SPF Record
Type: TXT
Name: $SENDING_DOMAIN
Value: v=spf1 a mx ip4:$SERVER_IP ~all

# DKIM Record
Type: TXT
Name: newsletter._domainkey.$SENDING_DOMAIN
Value: v=DKIM1; k=rsa; p=$DKIM_PUBLIC_KEY

# DMARC Record
Type: TXT
Name: _dmarc.$SENDING_DOMAIN
Value: v=DMARC1; p=quarantine; rua=mailto:dmarc@$SENDING_DOMAIN

# PTR Record (set by your hosting provider)
# Reverse DNS should point to: mail.$SENDING_DOMAIN

# Admin Domain (CNAME or A record)
Type: A
Name: $ADMIN_DOMAIN
Value: $SERVER_IP
EOF
    
    log_success "DNS records generated in $INSTALL_DIR/dns-records.txt"
}

wait_for_services() {
    log_info "Waiting for services to start..."
    
    # Wait for app to be ready
    for i in {1..30}; do
        if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
            log_success "Application is ready!"
            break
        fi
        log_info "Waiting for application... ($i/30)"
        sleep 10
    done
    
    if ! curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
        log_error "Application failed to start. Check logs with: docker-compose logs"
        exit 1
    fi
}

print_success_message() {
    echo
    log_success "Newsletter Platform installed successfully!"
    echo
    echo "Next steps:"
    echo "1. Add the DNS records from $INSTALL_DIR/dns-records.txt to your DNS provider"
    echo "2. Wait for DNS propagation (usually 5-15 minutes)"
    echo "3. Access the admin panel at: https://$ADMIN_DOMAIN"
    echo "4. Complete the domain verification process"
    echo
    echo "Useful commands:"
    echo "  View logs: docker-compose -f $INSTALL_DIR/deploy/docker-compose.yml logs"
    echo "  Restart: docker-compose -f $INSTALL_DIR/deploy/docker-compose.yml restart"
    echo "  Stop: docker-compose -f $INSTALL_DIR/deploy/docker-compose.yml down"
    echo
    echo "Backup script: $INSTALL_DIR/scripts/backup.sh"
    echo
}

# Main installation process
main() {
    echo "Newsletter Platform Installer"
    echo "============================="
    echo
    
    check_root
    check_dependencies
    install_docker
    install_docker_compose
    create_system_user
    prompt_configuration
    create_directories
    copy_files
    generate_dkim_keys
    generate_tls_cert
    generate_configs
    build_and_deploy
    generate_dns_records
    wait_for_services
    print_success_message
}

# Run main function
main "$@"
