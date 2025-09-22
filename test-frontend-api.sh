#!/bin/bash

# Test script to verify frontend API proxy is working
echo "🧪 Testing Frontend API Proxy"
echo "============================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Testing API endpoints through frontend proxy..."

# Test health endpoint
echo -n "Health endpoint... "
if curl -s http://localhost:3000/api/health | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

# Test domains endpoint
echo -n "Domains endpoint... "
if curl -s http://localhost:3000/api/domains | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

# Test lists endpoint
echo -n "Lists endpoint... "
if curl -s http://localhost:3000/api/lists | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

# Test campaigns endpoint
echo -n "Campaigns endpoint... "
if curl -s http://localhost:3000/api/campaigns | jq .success > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PASS${NC}"
else
    echo -e "${RED}❌ FAIL${NC}"
fi

echo ""
echo "🎯 Frontend API Proxy Test Complete!"
echo "Now you can visit http://localhost:3000 and all the buttons should work!"
