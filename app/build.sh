#!/bin/bash

# Build script for Newsletter Platform
# Builds the SvelteKit frontend and copies it to the Go static directory

set -e

# Colors for output
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

# Configuration
WEB_DIR="web"
STATIC_DIR="static"
BUILD_DIR="build"

log_info "Building Newsletter Platform..."

# Check if web directory exists
if [[ ! -d "$WEB_DIR" ]]; then
    log_warning "Web directory not found. Creating basic structure..."
    mkdir -p "$WEB_DIR/src/routes"
    mkdir -p "$WEB_DIR/static"
fi

# Install dependencies if package.json exists
if [[ -f "$WEB_DIR/package.json" ]]; then
    log_info "Installing frontend dependencies..."
    cd "$WEB_DIR"
    
    if [[ ! -d "node_modules" ]]; then
        npm install
    else
        log_info "Dependencies already installed"
    fi
    
    # Build the frontend
    log_info "Building SvelteKit frontend..."
    npm run build
    
    # Copy build output to static directory
    log_info "Copying build output to static directory..."
    if [[ -d "$BUILD_DIR" ]]; then
        rm -rf "../$STATIC_DIR"
        cp -r "$BUILD_DIR" "../$STATIC_DIR"
        log_success "Frontend built and copied to static directory"
    else
        log_warning "Build directory not found. Creating placeholder static files..."
        mkdir -p "../$STATIC_DIR"
        cat > "../$STATIC_DIR/index.html" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Newsletter Platform</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 2rem; background: #f8fafc; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 2rem; border-radius: 0.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        h1 { color: #1e293b; margin-bottom: 1rem; }
        p { color: #64748b; line-height: 1.6; }
        .status { background: #dcfce7; color: #166534; padding: 1rem; border-radius: 0.375rem; margin: 1rem 0; }
        .warning { background: #fef3c7; color: #92400e; padding: 1rem; border-radius: 0.375rem; margin: 1rem 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Newsletter Platform</h1>
        <div class="status">
            <strong>✓ Backend API is running</strong>
        </div>
        <div class="warning">
            <strong>⚠ Frontend not built</strong><br>
            The SvelteKit frontend needs to be built. Run the build script to compile the admin interface.
        </div>
        <p>
            This is a self-hosted newsletter platform. The backend API is running and ready to use.
            You can access the API endpoints at <code>/api/*</code>.
        </p>
        <h2>API Endpoints</h2>
        <ul>
            <li><code>GET /api/health</code> - Health check</li>
            <li><code>POST /api/domains</code> - Create domain</li>
            <li><code>GET /api/domains</code> - List domains</li>
            <li><code>POST /api/campaigns</code> - Create campaign</li>
            <li><code>GET /api/campaigns</code> - List campaigns</li>
        </ul>
    </div>
</body>
</html>
EOF
        log_success "Placeholder static files created"
    fi
    
    cd ..
else
    log_warning "No package.json found. Creating placeholder static files..."
    mkdir -p "$STATIC_DIR"
    cat > "$STATIC_DIR/index.html" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Newsletter Platform</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 2rem; background: #f8fafc; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 2rem; border-radius: 0.5rem; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
        h1 { color: #1e293b; margin-bottom: 1rem; }
        p { color: #64748b; line-height: 1.6; }
        .status { background: #dcfce7; color: #166534; padding: 1rem; border-radius: 0.375rem; margin: 1rem 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Newsletter Platform</h1>
        <div class="status">
            <strong>✓ Backend API is running</strong>
        </div>
        <p>
            This is a self-hosted newsletter platform. The backend API is running and ready to use.
            You can access the API endpoints at <code>/api/*</code>.
        </p>
        <h2>API Endpoints</h2>
        <ul>
            <li><code>GET /api/health</code> - Health check</li>
            <li><code>POST /api/domains</code> - Create domain</li>
            <li><code>GET /api/domains</code> - List domains</li>
            <li><code>POST /api/campaigns</code> - Create campaign</li>
            <li><code>GET /api/campaigns</code> - List campaigns</li>
        </ul>
    </div>
</body>
</html>
EOF
    log_success "Placeholder static files created"
fi

log_success "Build completed successfully!"
