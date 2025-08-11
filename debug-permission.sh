#!/bin/bash

# Test script untuk debug permission
# Usage: ./test-news-create.sh

API_BASE="http://localhost:8080"

echo "🔍 Testing News Creation Permission"
echo "=================================="

# Test 1: Login dulu
echo "1. Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"superadmin","password":"superadmin123"}')

echo "Login Response:"
echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token' 2>/dev/null)

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Failed to get token from login"
    exit 1
fi

echo -e "\n✅ Token obtained: ${TOKEN:0:20}..."

# Test 2: Check profile untuk verify role
echo -e "\n2. Checking user profile..."
PROFILE_RESPONSE=$(curl -s -X GET "$API_BASE/auth/profile" \
  -H "Authorization: Bearer $TOKEN")

echo "Profile Response:"
echo "$PROFILE_RESPONSE" | jq '.' 2>/dev/null || echo "$PROFILE_RESPONSE"

# Test 3: Try to create news
echo -e "\n3. Testing news creation..."
NEWS_RESPONSE=$(curl -s -X POST "$API_BASE/admin/news" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "news_title": "Test News",
    "content": "Test content for the news article",
    "status": "Posted",
    "category": "test",
    "image": "https://example.com/test-image.jpg"
  }')

echo "News Creation Response:"
echo "$NEWS_RESPONSE" | jq '.' 2>/dev/null || echo "$NEWS_RESPONSE"

echo -e "\n✅ Debug test completed!"
