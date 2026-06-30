#!/bin/bash
set -e

# Syncs OS-SPARK SovereignIdentity tokens from HashiCorp Vault
# Requires VAULT_ADDR and VAULT_TOKEN to be exported in the environment.

echo "=> Syncing Sovereign Identity from HashiCorp Vault..."

if [ -z "$VAULT_ADDR" ] || [ -z "$VAULT_TOKEN" ]; then
    echo "ERROR: VAULT_ADDR or VAULT_TOKEN not set."
    exit 1
fi

# Fetch the Sovereign Identity token and inject it into the Systemd environment file
vault kv get -field=identity_token secret/sovereign/mesh > /etc/sovereign-kernel/identity.token
echo "SOVEREIGN_IDENTITY_TOKEN=$(cat /etc/sovereign-kernel/identity.token)" > /etc/sovereign-kernel/vault.env

chmod 600 /etc/sovereign-kernel/vault.env
chmod 600 /etc/sovereign-kernel/identity.token

echo "=> Identity synced successfully. Restarting Sovereign Kernel..."
systemctl restart sovereign-kernel
