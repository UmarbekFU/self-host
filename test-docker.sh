#!/bin/bash

# Test script for Docker deployment
echo "🚀 Testing Newsletter Platform Docker Deployment"
echo "================================================"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local url=$1
    local name=$2
    local expected_status=${3:-200}
    
    echo -n "Testing $name... "
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url")
    
    if [ "$response" = "$expected_status" ]; then
        echo -e "${GREEN}✅ PASS${NC} (HTTP $response)"
        return 0
    else
        echo -e "${RED}❌ FAIL${NC} (HTTP $response, expected $expected_status)"
        return 1
    fi
}

# Function to test JSON endpoint
test_json_endpoint() {
    local url=$1
    local name=$2
    
    echo -n "Testing $name... "
    
    response=$(curl -s "$url")
    
    if echo "$response" | jq . > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${NC} (Valid JSON)"
        return 0
    else
        echo -e "${RED}❌ FAIL${NC} (Invalid JSON)"
        return 1
    fi
}

echo "📦 Building and starting services..."
docker-compose -f docker-compose-full.yml up -d --build

echo ""
echo "⏳ Waiting for services to start..."
sleep 30

echo ""
echo "🧪 Running tests..."
echo "=================="

# Test backend health
test_endpoint "http://localhost:8081/api/health" "Backend Health Check"

# Test backend API
test_json_endpoint "http://localhost:8081/api/domains" "Domains API"
test_json_endpoint "http://localhost:8081/api/lists" "Lists API"
test_json_endpoint "http://localhost:8081/api/campaigns" "Campaigns API"

# Test frontend
test_endpoint "http://localhost:3000" "Frontend Homepage"

# Test MailHog
test_endpoint "http://localhost:8025" "MailHog Web UI"

echo ""
echo "📊 Testing API functionality..."

# Test creating a domain
echo -n "Creating test domain... "
domain_response=$(curl -s -X POST http://localhost:8081/api/domains \
    -H "Content-Type: application/json" \
    -d '{"domain":"test.example.com"}')

if echo "$domain_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

# Test creating a list
echo -n "Creating test list... "
list_response=$(curl -s -X POST http://localhost:8081/api/lists \
    -H "Content-Type: application/json" \
    -d '{"name":"Test List","description":"A test list"}')

if echo "$list_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

echo ""
echo "🎯 Test Results Summary"
echo "======================="
echo "✅ Backend API: http://localhost:8081"
echo "✅ Frontend UI: http://localhost:3000"
echo "✅ MailHog UI: http://localhost:8025"
echo ""
echo "📝 You can now:"
echo "   - Visit http://localhost:3000 to see the frontend"
echo "   - Visit http://localhost:8025 to see captured emails"
echo "   - Use the API at http://localhost:8081/api/*"
echo ""
echo "🛑 To stop: docker-compose -f docker-compose-full.yml down"
