#!/bin/bash

# API Testing Script untuk Haslaw Services
# Usage: ./test-api.sh [server-ip]

SERVER_IP=${1:-"103.179.57.14"}
BASE_URL="http://${SERVER_IP}:8080"

echo "🚀 Testing Haslaw API at ${BASE_URL}"
echo "=================================="

# Test 1: Basic connectivity
echo "📡 Test 1: Basic Connectivity"
if curl -s --connect-timeout 5 "${BASE_URL}/" > /dev/null; then
    echo "✅ Server is reachable"
else
    echo "❌ Server is not reachable"
    exit 1
fi

# Test 2: Health check
echo -e "\n🏥 Test 2: Health Check"
HEALTH_RESPONSE=$(curl -s -w "%{http_code}" "${BASE_URL}/health" -o /dev/null)
if [ "$HEALTH_RESPONSE" = "200" ]; then
    echo "✅ Health check passed"
else
    echo "⚠️  Health check returned: $HEALTH_RESPONSE"
fi

# Test 3: API endpoints
echo -e "\n🔍 Test 3: API Endpoints"

# Test root endpoint
echo "Testing root endpoint..."
curl -s -X GET "${BASE_URL}/" | head -c 200
echo ""

# Test auth endpoints (if available)
echo -e "\nTesting auth endpoints..."
curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}' | head -c 200
echo ""

# Test admin endpoints (if available)
echo -e "\nTesting admin endpoints..."
curl -s -X GET "${BASE_URL}/admin/profile" \
  -H "Content-Type: application/json" | head -c 200
echo ""

# Test 4: Response time
echo -e "\n⏱️  Test 4: Response Time"
TIME_RESPONSE=$(curl -o /dev/null -s -w "%{time_total}" "${BASE_URL}/")
echo "Response time: ${TIME_RESPONSE} seconds"

echo -e "\n✅ Testing completed!"
