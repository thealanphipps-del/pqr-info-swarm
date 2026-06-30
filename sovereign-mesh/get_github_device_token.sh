#!/usr/bin/env bash
# Helper script to perform GitHub Device Flow using the profile in device_login_profile.json
# Prerequisite: jq installed (sudo apt-get install -y jq)

PROFILE="/home/aellok/sovereign-mesh/device_login_profile.json"

if [[ ! -f "$PROFILE" ]]; then
  echo "Profile not found at $PROFILE"
  exit 1
fi

CLIENT_ID=$(jq -r '.client_id' "$PROFILE")
CLIENT_SECRET=$(jq -r '.client_secret' "$PROFILE")
SCOPES=$(jq -r '.scopes | join(" ")' "$PROFILE")

# Step 1: Request device code
DEVICE_RESPONSE=$(curl -s -X POST "https://github.com/login/device/code" \
  -d "client_id=$CLIENT_ID&scope=$SCOPES")

DEVICE_CODE=$(echo "$DEVICE_RESPONSE" | jq -r '.device_code')
USER_CODE=$(echo "$DEVICE_RESPONSE" | jq -r '.user_code')
VERIFICATION_URI=$(echo "$DEVICE_RESPONSE" | jq -r '.verification_uri')
INTERVAL=$(echo "$DEVICE_RESPONSE" | jq -r '.interval')

echo "Please open a browser and go to: $VERIFICATION_URI"
echo "Enter the user code: $USER_CODE"

echo "Waiting for you to authorize..."

# Step 2: Poll for token
while true; do
  TOKEN_RESPONSE=$(curl -s -X POST "https://github.com/login/oauth/access_token" \
    -d "client_id=$CLIENT_ID&device_code=$DEVICE_CODE&grant_type=urn:ietf:params:oauth:grant-type:device_code" \
    -H "Accept: application/json")

  ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
  ERROR=$(echo "$TOKEN_RESPONSE" | jq -r '.error')

  if [[ "$ACCESS_TOKEN" != "null" && -n "$ACCESS_TOKEN" ]]; then
    echo "✅ Access token obtained!"
    echo "Access Token: $ACCESS_TOKEN"
    # Optionally store token securely
    exit 0
  fi

  if [[ "$ERROR" == "authorization_pending" ]]; then
    sleep $INTERVAL
    continue
  else
    echo "❌ Error: $ERROR"
    exit 1
  fi
done
