#!/bin/bash
# This script is used to test the auth API.

set -e

PORT=3000
URL="http://localhost:$PORT"
AUTH="$URL/api/auth"

req() {
  curl -s -X POST -H "Content-Type: application/json" -d "$1" $2
}

# Create a new user & log in

CREDENTIALS='{"username":"testuser","password":"testpassword","email":"test123@google.com"}'

TOKEN_1=$(req "$CREDENTIALS" "$AUTH/signup" | jq -r '.token')
TOKEN_2=$(req "$CREDENTIALS" "$AUTH/login" | jq -r '.token')

if [ "$TOKEN_1" != "$TOKEN_2" ]; then
  echo "Tokens do not match ($TOKEN_1 != $TOKEN_2)"
  exit 1
fi

# Test the token

TOKEN_TEST=$(req "{\"token\":\"$TOKEN_1\",\"username\":\"testuser\"}" "$AUTH/confirm_identity")

if [ "$TOKEN_TEST" != "{\"ok\":\"true\"}" ]; then
  echo "Token test failed"
  exit 1
fi

echo "Ok!"
