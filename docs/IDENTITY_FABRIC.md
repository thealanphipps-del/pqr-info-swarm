# 🛡️ PQR Identity Fabric

The Identity Fabric ensures that only authorized entities can interact with the Sovereign Node, while allowing the Node itself to interact with the global enterprise ecosystem.

## 🔐 SAML Identity Provider (IdP)
The node acts as a standalone SAML IdP.
- **Metadata URL**: `https://pqr.info/saml/metadata`
- **SSO URL**: `https://pqr.info/saml/sso`
- **Certificate**: Self-signed RSA-2048 (Stored in Vault).

### 🔄 Autonomous Certificate Rotation
The **MonitoringService** checks the SAML certificate health every minute. If the certificate is within 7 days of expiration, the **HealingService** triggers an autonomous rotation:
1. Generate new Key/Cert pair.
2. Update Vault.
3. Reload the live `AuthService` in-memory.

## 🧱 Cloudflare Access Bypass
The `pqr.info` domain is protected by Cloudflare Access. To allow the **Healing Agents** to perform health checks from the outside, they use **Service Tokens**.

### Credentials (Stored in Vault)
- **CF-Access-Client-Id**: `c98ca7026f54305b05cd24975a3ce6d2.access`
- **CF-Access-Client-Secret**: `ebf3177d992adb0c3db7b088fb5b9e3d83e96649fb9bc5b86a25301af5c8e744`

### Usage
Every internal request to `pqr.info` from the Monitoring Service automatically injects these headers to "pierce" the Access wall for forensic probing.
