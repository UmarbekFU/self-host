#!/bin/bash

# Test script for the newsletter platform API
echo "Testing Newsletter Platform API..."

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s http://localhost:8081/api/health || echo "Health endpoint not available"

echo -e "\n2. Testing domains endpoint..."
curl -s http://localhost:8081/api/domains || echo "Domains endpoint not available"

echo -e "\n3. Testing campaigns endpoint..."
curl -s http://localhost:8081/api/campaigns || echo "Campaigns endpoint not available"

echo -e "\n4. Testing root endpoint..."
curl -s http://localhost:8081/ | head -10

echo -e "\n5. Testing static files..."
curl -s http://localhost:8081/static/ | head -5

echo -e "\nAPI testing complete!"
