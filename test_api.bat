@echo off
REM Test script for Blog API
REM Make sure the API is running on localhost:8080

set BASE_URL=http://localhost:8080

echo Testing Blog API...
echo ==================

REM Test health check
echo 1. Testing health check...
curl -s "%BASE_URL%/health"
echo.

REM Test user registration
echo 2. Testing user registration...
curl -s -X POST "%BASE_URL%/api/v1/auth/register" ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"testuser\", \"email\": \"test@example.com\", \"password\": \"password123\"}"
echo.

REM Test user login
echo 3. Testing user login...
curl -s -X POST "%BASE_URL%/api/v1/auth/login" ^
  -H "Content-Type: application/json" ^
  -d "{\"email\": \"test@example.com\", \"password\": \"password123\"}"
echo.

echo API testing completed!
