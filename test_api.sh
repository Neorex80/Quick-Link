#!/bin/bash

# URL Shortener API Test Script
# This script tests the URL shortener functionality

echo "🔗 URL Shortener API Test Script"
echo "================================="
echo ""

# Check if server is running
echo "📡 Checking if server is running..."
if ! curl -s http://localhost:8080 > /dev/null; then
    echo "❌ Server is not running. Please start it with: go run main.go"
    exit 1
fi
echo "✅ Server is running!"
echo ""

# Test 1: Valid URL shortening
echo "🧪 Test 1: Shortening a valid URL"
echo "Request: POST /shorten with https://www.google.com"
response=$(curl -s -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}')
echo "Response: $response"
echo ""

# Extract short URL for testing redirect
short_url=$(echo $response | grep -o '"short_url":"[^"]*"' | cut -d'"' -f4)
if [ ! -z "$short_url" ]; then
    echo "🔄 Testing redirect for: $short_url"
    redirect_response=$(curl -s -I "$short_url" | head -n 1)
    echo "Redirect response: $redirect_response"
    echo ""
fi

# Test 2: Invalid URL
echo "🧪 Test 2: Testing invalid URL"
echo "Request: POST /shorten with invalid URL"
response=$(curl -s -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "not-a-valid-url"}')
echo "Response: $response"
echo ""

# Test 3: Missing Content-Type
echo "🧪 Test 3: Testing missing Content-Type header"
echo "Request: POST /shorten without Content-Type header"
response=$(curl -s -X POST http://localhost:8080/shorten \
  -d '{"url": "https://example.com"}')
echo "Response: $response"
echo ""

# Test 4: Invalid JSON
echo "🧪 Test 4: Testing invalid JSON"
echo "Request: POST /shorten with malformed JSON"
response=$(curl -s -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"')
echo "Response: $response"
echo ""

# Test 5: Method not allowed
echo "🧪 Test 5: Testing wrong HTTP method"
echo "Request: GET /shorten (should be POST)"
response=$(curl -s -X GET http://localhost:8080/shorten)
echo "Response: $response"
echo ""

# Test 6: Non-existent short code
echo "🧪 Test 6: Testing non-existent short code"
echo "Request: GET /nonexistent"
response=$(curl -s -I http://localhost:8080/nonexistent | head -n 1)
echo "Response: $response"
echo ""

echo "✅ All tests completed!"
echo ""
echo "💡 To test the web interface, open: http://localhost:8080"
