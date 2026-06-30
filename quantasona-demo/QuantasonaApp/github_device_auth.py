import urllib.request
import urllib.parse
import json
import time
import sys

CLIENT_ID = "178ee8a7de8d3f6a22aa"  # Official GitHub CLI Client ID
SCOPES = "repo read:org write:org"

def main():
    print("Initiating GitHub Device Authorization Flow...")
    
    # 1. Request device and user codes
    url = "https://github.com/login/device/code"
    data = urllib.parse.urlencode({
        "client_id": CLIENT_ID,
        "scope": SCOPES
    }).encode("utf-8")
    
    req = urllib.request.Request(url, data=data, headers={
        "Accept": "application/json",
        "Content-Type": "application/x-www-form-urlencoded"
    })
    
    try:
        with urllib.request.urlopen(req) as response:
            res = json.loads(response.read().decode("utf-8"))
    except Exception as e:
        print(f"Error connecting to GitHub: {e}")
        return

    device_code = res.get("device_code")
    user_code = res.get("user_code")
    verification_uri = res.get("verification_uri")
    expires_in = res.get("expires_in")
    interval = res.get("interval", 5)

    print("\n" + "="*50)
    print(f"1. Open your standard browser and go to:\n   {verification_uri}")
    print(f"\n2. Enter the following authorization code:\n   {user_code}")
    print("="*50 + "\n")
    print("Waiting for authorization (polling GitHub)... Press Ctrl+C to cancel.")

    # 2. Poll for token
    token_url = "https://github.com/login/oauth/access_token"
    token_data = urllib.parse.urlencode({
        "client_id": CLIENT_ID,
        "device_code": device_code,
        "grant_type": "urn:ietf:params:oauth:grant-type:device_code"
    }).encode("utf-8")

    token_req = urllib.request.Request(token_url, data=token_data, headers={
        "Accept": "application/json",
        "Content-Type": "application/x-www-form-urlencoded"
    })

    start_time = time.time()
    while time.time() - start_time < expires_in:
        try:
            with urllib.request.urlopen(token_req) as response:
                token_res = json.loads(response.read().decode("utf-8"))
        except Exception as e:
            print(f"\nError polling: {e}")
            time.sleep(interval)
            continue

        error = token_res.get("error")
        if error:
            if error == "authorization_pending":
                # User hasn't authorized yet, sleep and retry
                sys.stdout.write(".")
                sys.stdout.flush()
                time.sleep(interval)
                continue
            elif error == "slow_down":
                interval += 5
                time.sleep(interval)
                continue
            else:
                print(f"\nAuthorization failed: {token_res.get('error_description')}")
                return
        
        access_token = token_res.get("access_token")
        if access_token:
            print("\n\nSuccessfully Authenticated!")
            save_token(access_token)
            return

        time.sleep(interval)

    print("\nAuthorization timed out.")

def save_token(token):
    try:
        token_file = "C:/Users/theal/QuantasonaApp/github_token.txt"
        with open(token_file, "w") as f:
            f.write(token)
        print(f"GitHub token saved locally to {token_file}.")
    except Exception as e:
        print(f"Error saving token to file: {e}")

if __name__ == "__main__":
    main()
