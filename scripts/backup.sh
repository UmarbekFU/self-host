#!/usr/bin/env bash

# Newsletter Platform Backup Script
# This script creates a complete backup of the newsletter platform data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BACKUP_DIR="/opt/newsletter/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="newsletter_backup_${TIMESTAMP}"
BACKUP_PATH="${BACKUP_DIR}/${BACKUP_NAME}"

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
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root to access Docker volumes and system files."
        exit 1
    fi
}

# Function to create backup directory
create_backup_dir() {
    print_status "Creating backup directory..."
    mkdir -p "$BACKUP_PATH"
    print_success "Backup directory created: $BACKUP_PATH"
}

# Function to backup database
backup_database() {
    print_status "Backing up database..."
    
    # Create database backup directory
    mkdir -p "$BACKUP_PATH/database"
    
    # Copy SQLite database
    if [ -f "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" ]; then
        cp "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" "$BACKUP_PATH/database/"
        print_success "Database backed up"
    else
        print_warning "Database file not found, skipping database backup"
    fi
}

# Function to backup configuration files
backup_config() {
    print_status "Backing up configuration files..."
    
    # Create config backup directory
    mkdir -p "$BACKUP_PATH/config"
    
    # Copy Docker Compose files
    if [ -f "docker-compose-full.yml" ]; then
        cp docker-compose-full.yml "$BACKUP_PATH/config/"
    fi
    
    # Copy environment file
    if [ -f ".env" ]; then
        cp .env "$BACKUP_PATH/config/"
    fi
    
    # Copy Caddyfile
    if [ -f "deploy/Caddyfile" ]; then
        cp deploy/Caddyfile "$BACKUP_PATH/config/"
    fi
    
    # Copy MTA configuration
    if [ -f "deploy/maddy.conf" ]; then
        cp deploy/maddy.conf "$BACKUP_PATH/config/"
    fi
    
    # Copy DKIM keys
    if [ -f "deploy/dkim.key" ]; then
        cp deploy/dkim.key "$BACKUP_PATH/config/"
    fi
    
    # Copy credentials
    if [ -f "deploy/credentials" ]; then
        cp deploy/credentials "$BACKUP_PATH/config/"
    fi
    
    print_success "Configuration files backed up"
}

# Function to export data
export_data() {
    print_status "Exporting data..."
    
    # Create data export directory
    mkdir -p "$BACKUP_PATH/data"
    
    # Export subscribers
    print_status "Exporting subscribers..."
    docker exec newsletter-app wget -qO- "http://localhost:8080/api/export/subscribers" > "$BACKUP_PATH/data/subscribers.json" 2>/dev/null || {
        print_warning "Failed to export subscribers via API, using database dump"
        sqlite3 "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" ".dump subscribers" > "$BACKUP_PATH/data/subscribers.sql"
    }
    
    # Export campaigns
    print_status "Exporting campaigns..."
    docker exec newsletter-app wget -qO- "http://localhost:8080/api/export/campaigns" > "$BACKUP_PATH/data/campaigns.json" 2>/dev/null || {
        print_warning "Failed to export campaigns via API, using database dump"
        sqlite3 "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" ".dump campaigns" > "$BACKUP_PATH/data/campaigns.sql"
    }
    
    # Export events
    print_status "Exporting events..."
    docker exec newsletter-app wget -qO- "http://localhost:8080/api/export/events" > "$BACKUP_PATH/data/events.json" 2>/dev/null || {
        print_warning "Failed to export events via API, using database dump"
        sqlite3 "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" ".dump events" > "$BACKUP_PATH/data/events.sql"
    }
    
    # Export lists
    print_status "Exporting lists..."
    docker exec newsletter-app wget -qO- "http://localhost:8080/api/export/lists" > "$BACKUP_PATH/data/lists.json" 2>/dev/null || {
        print_warning "Failed to export lists via API, using database dump"
        sqlite3 "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" ".dump lists" > "$BACKUP_PATH/data/lists.sql"
    }
    
    # Export domains
    print_status "Exporting domains..."
    docker exec newsletter-app wget -qO- "http://localhost:8080/api/export/domains" > "$BACKUP_PATH/data/domains.json" 2>/dev/null || {
        print_warning "Failed to export domains via API, using database dump"
        sqlite3 "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" ".dump domains" > "$BACKUP_PATH/data/domains.sql"
    }
    
    print_success "Data exported"
}

# Function to create full database dump
create_database_dump() {
    print_status "Creating full database dump..."
    
    # Create SQL dump
    sqlite3 "/var/lib/docker/volumes/newsletter-platform_app-data/_data/newsletter.db" ".dump" > "$BACKUP_PATH/database/full_dump.sql"
    
    print_success "Full database dump created"
}

# Function to backup Docker volumes
backup_docker_volumes() {
    print_status "Backing up Docker volumes..."
    
    # Create volumes backup directory
    mkdir -p "$BACKUP_PATH/volumes"
    
    # Backup app data volume
    if docker volume inspect newsletter-platform_app-data >/dev/null 2>&1; then
        docker run --rm -v newsletter-platform_app-data:/data -v "$BACKUP_PATH/volumes":/backup alpine tar czf /backup/app-data.tar.gz -C /data .
        print_success "App data volume backed up"
    else
        print_warning "App data volume not found"
    fi
    
    # Backup Caddy data volume
    if docker volume inspect newsletter-platform_caddy-data >/dev/null 2>&1; then
        docker run --rm -v newsletter-platform_caddy-data:/data -v "$BACKUP_PATH/volumes":/backup alpine tar czf /backup/caddy-data.tar.gz -C /data .
        print_success "Caddy data volume backed up"
    else
        print_warning "Caddy data volume not found"
    fi
}

# Function to create backup manifest
create_manifest() {
    print_status "Creating backup manifest..."
    
    cat > "$BACKUP_PATH/MANIFEST.txt" << EOF
Newsletter Platform Backup
=========================
Created: $(date)
Backup Name: $BACKUP_NAME
Backup Path: $BACKUP_PATH

Contents:
- database/: Database files and dumps
- config/: Configuration files
- data/: Exported data (JSON/SQL)
- volumes/: Docker volume backups
- MANIFEST.txt: This file

Database:
- newsletter.db: SQLite database file
- full_dump.sql: Complete database dump

Configuration:
- docker-compose-full.yml: Docker Compose configuration
- .env: Environment variables
- Caddyfile: Caddy reverse proxy configuration
- maddy.conf: MTA configuration
- dkim.key: DKIM private key
- credentials: MTA credentials

Data Exports:
- subscribers.json/sql: Subscriber data
- campaigns.json/sql: Campaign data
- events.json/sql: Event tracking data
- lists.json/sql: Mailing list data
- domains.json/sql: Domain configuration

Volumes:
- app-data.tar.gz: Application data volume
- caddy-data.tar.gz: Caddy data volume

Restore Instructions:
1. Extract the backup: tar -xzf $BACKUP_NAME.tar.gz
2. Restore database: cp database/newsletter.db /var/lib/docker/volumes/newsletter-platform_app-data/_data/
3. Restore configuration: cp config/* ./
4. Restore volumes: docker run --rm -v newsletter-platform_app-data:/data -v \$(pwd)/volumes:/backup alpine tar xzf /backup/app-data.tar.gz -C /data
5. Restart services: docker-compose -f docker-compose-full.yml up -d
EOF

    print_success "Backup manifest created"
}

# Function to compress backup
compress_backup() {
    print_status "Compressing backup..."
    
    cd "$BACKUP_DIR"
    tar -czf "${BACKUP_NAME}.tar.gz" "$BACKUP_NAME"
    rm -rf "$BACKUP_NAME"
    
    print_success "Backup compressed: ${BACKUP_NAME}.tar.gz"
}

# Function to cleanup old backups
cleanup_old_backups() {
    print_status "Cleaning up old backups..."
    
    # Keep only the last 7 backups
    cd "$BACKUP_DIR"
    ls -t newsletter_backup_*.tar.gz | tail -n +8 | xargs -r rm -f
    
    print_success "Old backups cleaned up"
}

# Function to display backup information
display_backup_info() {
    echo
    print_success "Backup completed successfully!"
    echo "=================================="
    echo
    print_status "Backup Details:"
    echo "- Name: $BACKUP_NAME"
    echo "- Location: $BACKUP_DIR/${BACKUP_NAME}.tar.gz"
    echo "- Size: $(du -h "$BACKUP_DIR/${BACKUP_NAME}.tar.gz" | cut -f1)"
    echo
    print_status "To restore this backup:"
    echo "1. Extract: tar -xzf $BACKUP_DIR/${BACKUP_NAME}.tar.gz"
    echo "2. Follow the instructions in MANIFEST.txt"
    echo
    print_status "Backup contents:"
    echo "- Database: SQLite file and SQL dumps"
    echo "- Configuration: All config files and keys"
    echo "- Data: Exported subscribers, campaigns, events"
    echo "- Volumes: Docker volume backups"
    echo
}

# Main backup function
main() {
    echo "Newsletter Platform Backup Script"
    echo "================================="
    echo
    
    # Pre-flight checks
    check_root
    
    # Create backup directory
    create_backup_dir
    
    # Perform backups
    backup_database
    backup_config
    export_data
    create_database_dump
    backup_docker_volumes
    
    # Create manifest
    create_manifest
    
    # Compress backup
    compress_backup
    
    # Cleanup old backups
    cleanup_old_backups
    
    # Display information
    display_backup_info
}

# Run main function
main "$@"