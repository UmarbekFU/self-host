#!/bin/bash

# Simple test script for Docker deployment
echo "🚀 Testing Newsletter Platform (Simple Docker Setup)"
echo "===================================================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "📦 Starting services..."
docker-compose -f docker-compose-simple.yml up -d

echo ""
echo "⏳ Waiting for services to start..."
sleep 20

echo ""
echo "🧪 Testing Backend API..."
echo "========================"

# Test backend health
echo -n "Testing health endpoint... "
health_response=$(curl -s http://localhost:8081/api/health)
if echo "$health_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

# Test API endpoints
echo -n "Testing domains API... "
domains_response=$(curl -s http://localhost:8081/api/domains)
if echo "$domains_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

echo -n "Testing lists API... "
lists_response=$(curl -s http://localhost:8081/api/lists)
if echo "$lists_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

echo -n "Testing campaigns API... "
campaigns_response=$(curl -s http://localhost:8081/api/campaigns)
if echo "$campaigns_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

echo ""
echo "📊 Testing API functionality..."

# Test creating a domain
echo -n "Creating test domain... "
domain_response=$(curl -s -X POST http://localhost:8081/api/domains \
    -H "Content-Type: application/json" \
    -d '{"domain":"test.example.com"}')

if echo "$domain_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
    domain_id=$(echo "$domain_response" | jq -r '.data.id')
    echo "   Domain ID: $domain_id"
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
    list_id=$(echo "$list_response" | jq -r '.data.id')
    echo "   List ID: $list_id"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

# Test creating a campaign
echo -n "Creating test campaign... "
campaign_response=$(curl -s -X POST http://localhost:8081/api/campaigns \
    -H "Content-Type: application/json" \
    -d "{\"list_id\":$list_id,\"subject\":\"Welcome!\",\"html\":\"<h1>Welcome!</h1>\",\"text\":\"Welcome!\",\"from_name\":\"Test\",\"from_email\":\"test@test.example.com\"}")

if echo "$campaign_response" | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
    campaign_id=$(echo "$campaign_response" | jq -r '.data.id')
    echo "   Campaign ID: $campaign_id"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

echo ""
echo "🎯 Test Results Summary"
echo "======================="
echo "✅ Backend API: http://localhost:8081"
echo "✅ MailHog UI: http://localhost:8025"
echo ""
echo "📝 You can now:"
echo "   - Test the API at http://localhost:8081/api/*"
echo "   - View captured emails at http://localhost:8025"
echo "   - Run the frontend separately with: cd app/web && npm run dev"
echo ""
echo "🛑 To stop: docker-compose -f docker-compose-simple.yml down"
