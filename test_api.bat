@echo off
REM URL Shortener API Test Script for Windows
REM This script tests the URL shortener functionality

echo 🔗 URL Shortener API Test Script
echo =================================
echo.

REM Check if server is running
echo 📡 Checking if server is running...
curl -s http://localhost:8080 >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Server is not running. Please start it with: go run main.go
    pause
    exit /b 1
)
echo ✅ Server is running!
echo.

REM Test 1: Valid URL shortening
echo 🧪 Test 1: Shortening a valid URL
echo Request: POST /shorten with https://www.google.com
curl -s -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d "{\"url\": \"https://www.google.com\"}"
echo.
echo.

REM Test 2: Invalid URL
echo 🧪 Test 2: Testing invalid URL
echo Request: POST /shorten with invalid URL
curl -s -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d "{\"url\": \"not-a-valid-url\"}"
echo.
echo.

REM Test 3: Missing Content-Type
echo 🧪 Test 3: Testing missing Content-Type header
echo Request: POST /shorten without Content-Type header
curl -s -X POST http://localhost:8080/shorten -d "{\"url\": \"https://example.com\"}"
echo.
echo.

REM Test 4: Invalid JSON
echo 🧪 Test 4: Testing invalid JSON
echo Request: POST /shorten with malformed JSON
curl -s -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d "{\"url\": \"https://example.com\""
echo.
echo.

REM Test 5: Method not allowed
echo 🧪 Test 5: Testing wrong HTTP method
echo Request: GET /shorten (should be POST)
curl -s -X GET http://localhost:8080/shorten
echo.
echo.

REM Test 6: Non-existent short code
echo 🧪 Test 6: Testing non-existent short code
echo Request: GET /nonexistent
curl -s -I http://localhost:8080/nonexistent
echo.
echo.

echo ✅ All tests completed!
echo.
echo 💡 To test the web interface, open: http://localhost:8080
echo.
pause
