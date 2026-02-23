#!/bin/bash
# Comprehensive API Test Script
# Tests all RDP API endpoints

TOKEN="Bearer test-token-123"
BASE_URL="http://localhost:8080/api/v1"

echo "=== RDP API Comprehensive Test ==="
echo ""

# Test Health Endpoint
echo "1. Testing Health Endpoint..."
curl -s "$BASE_URL/health" | jq .
echo ""

# Test Forum APIs
echo "2. Testing Forum APIs..."
echo "   2a. List Boards:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/boards" | jq '.data | length'
echo ""

echo "   2b. List Posts:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/posts" | jq '.data | length'
echo ""

echo "   2c. Get Single Post:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/posts/01HPVJ8ZM00000000000000301" | jq '.data.title'
echo ""

echo "   2d. Get Post Replies:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/posts/01HPVJ8ZM00000000000000301/replies" | jq '.data | length'
echo ""

# Test Knowledge APIs
echo "3. Testing Knowledge APIs..."
echo "   3a. List Knowledge:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/knowledge" | jq '.code'
echo ""

echo "   3b. List Categories:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/categories" | jq '.code'
echo ""

# Test Analytics APIs
echo "4. Testing Analytics APIs..."
echo "   4a. Dashboard:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/analytics/dashboard" | jq '.code'
echo ""

# Test Project APIs
echo "5. Testing Project APIs..."
echo "   5a. List Projects:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/projects" | jq '.message'
echo ""

# Test Zotero APIs
echo "6. Testing Zotero APIs..."
echo "   6a. List Items:"
curl -s -H "Authorization: $TOKEN" "$BASE_URL/zotero/items" | jq '.code'
echo ""

echo "=== Test Complete ==="
