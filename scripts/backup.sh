#!/bin/bash

# Newsletter Platform Backup Script
# Creates a complete backup of the platform including database, configs, and keys

set -e

# Configuration
INSTALL_DIR="/opt/newsletter"
BACKUP_DIR="$INSTALL_DIR/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="newsletter_backup_$TIMESTAMP"
BACKUP_PATH="$BACKUP_DIR/$BACKUP_NAME"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

check_installation() {
    if [[ ! -d "$INSTALL_DIR" ]]; then
        log_error "Newsletter platform not found at $INSTALL_DIR"
        exit 1
    fi
}

create_backup_directory() {
    log_info "Creating backup directory..."
    mkdir -p "$BACKUP_PATH"
    log_success "Backup directory created: $BACKUP_PATH"
}

backup_database() {
    log_info "Backing up database..."
    
    # Stop the application to ensure clean database state
    log_info "Stopping application..."
    cd "$INSTALL_DIR"
    docker-compose -f deploy/docker-compose.yml stop app
    
    # Copy database file
    if [[ -f "$INSTALL_DIR/app-data/newsletter.db" ]]; then
        cp "$INSTALL_DIR/app-data/newsletter.db" "$BACKUP_PATH/"
        log_success "Database backed up"
    else
        log_warning "Database file not found"
    fi
    
    # Restart the application
    log_info "Restarting application..."
    docker-compose -f deploy/docker-compose.yml start app
}

backup_configurations() {
    log_info "Backing up configurations..."
    
    # Copy Docker Compose files
    cp -r "$INSTALL_DIR/deploy" "$BACKUP_PATH/"
    
    # Copy environment file
    if [[ -f "$INSTALL_DIR/.env" ]]; then
        cp "$INSTALL_DIR/.env" "$BACKUP_PATH/"
    fi
    
    # Copy DNS records
    if [[ -f "$INSTALL_DIR/dns-records.txt" ]]; then
        cp "$INSTALL_DIR/dns-records.txt" "$BACKUP_PATH/"
    fi
    
    log_success "Configurations backed up"
}

backup_keys() {
    log_info "Backing up keys and certificates..."
    
    # Create keys directory
    mkdir -p "$BACKUP_PATH/keys"
    
    # Copy DKIM keys
    if [[ -d "$INSTALL_DIR/mta/keys" ]]; then
        cp -r "$INSTALL_DIR/mta/keys"/* "$BACKUP_PATH/keys/"
        log_success "Keys and certificates backed up"
    else
        log_warning "Keys directory not found"
    fi
}

backup_logs() {
    log_info "Backing up logs..."
    
    # Create logs directory
    mkdir -p "$BACKUP_PATH/logs"
    
    # Get container logs
    cd "$INSTALL_DIR"
    docker-compose -f deploy/docker-compose.yml logs > "$BACKUP_PATH/logs/docker-compose.log" 2>&1 || true
    
    # Get individual service logs
    docker-compose -f deploy/docker-compose.yml logs app > "$BACKUP_PATH/logs/app.log" 2>&1 || true
    docker-compose -f deploy/docker-compose.yml logs mta > "$BACKUP_PATH/logs/mta.log" 2>&1 || true
    docker-compose -f deploy/docker-compose.yml logs proxy > "$BACKUP_PATH/logs/proxy.log" 2>&1 || true
    
    log_success "Logs backed up"
}

create_backup_manifest() {
    log_info "Creating backup manifest..."
    
    cat > "$BACKUP_PATH/MANIFEST.txt" << EOF
Newsletter Platform Backup
=========================
Created: $(date)
Version: v1.0.0

Contents:
- newsletter.db: SQLite database
- deploy/: Docker Compose and configuration files
- keys/: DKIM keys and TLS certificates
- logs/: Application logs
- .env: Environment configuration
- dns-records.txt: DNS configuration

Restore Instructions:
1. Stop the newsletter platform
2. Copy all files to the installation directory
3. Restart the platform
4. Verify DNS records are still valid

Backup Size: $(du -sh "$BACKUP_PATH" | cut -f1)
EOF
    
    log_success "Backup manifest created"
}

compress_backup() {
    log_info "Compressing backup..."
    
    cd "$BACKUP_DIR"
    tar -czf "${BACKUP_NAME}.tar.gz" "$BACKUP_NAME"
    rm -rf "$BACKUP_NAME"
    
    BACKUP_FILE="${BACKUP_NAME}.tar.gz"
    BACKUP_SIZE=$(du -sh "$BACKUP_FILE" | cut -f1)
    
    log_success "Backup compressed: $BACKUP_FILE ($BACKUP_SIZE)"
}

cleanup_old_backups() {
    log_info "Cleaning up old backups..."
    
    # Keep only the last 10 backups
    cd "$BACKUP_DIR"
    ls -t newsletter_backup_*.tar.gz | tail -n +11 | xargs -r rm -f
    
    log_success "Old backups cleaned up"
}

print_backup_info() {
    echo
    log_success "Backup completed successfully!"
    echo
    echo "Backup file: $BACKUP_DIR/$BACKUP_FILE"
    echo "Backup size: $BACKUP_SIZE"
    echo
    echo "To restore this backup:"
    echo "1. Stop the platform: docker-compose -f $INSTALL_DIR/deploy/docker-compose.yml down"
    echo "2. Extract backup: tar -xzf $BACKUP_DIR/$BACKUP_FILE -C $INSTALL_DIR"
    echo "3. Start the platform: docker-compose -f $INSTALL_DIR/deploy/docker-compose.yml up -d"
    echo
}

# Main backup process
main() {
    echo "Newsletter Platform Backup"
    echo "========================="
    echo
    
    check_installation
    create_backup_directory
    backup_database
    backup_configurations
    backup_keys
    backup_logs
    create_backup_manifest
    compress_backup
    cleanup_old_backups
    print_backup_info
}

# Run main function
main "$@"
