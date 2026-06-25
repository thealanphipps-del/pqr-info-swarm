package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultSecretManager struct {
	client *api.Client
}

func NewVaultSecretManager() (*VaultSecretManager, error) {
	config := api.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")
	if config.Address == "" {
		config.Address = "http://localhost:8200"
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	token := os.Getenv("VAULT_TOKEN")
	if token != "" {
		client.SetToken(token)
	}

	return &VaultSecretManager{client: client}, nil
}

func (v *VaultSecretManager) StoreSAMLKey(ctx context.Context, keyPEM, certPEM []byte) error {
	data := map[string]interface{}{
		"private_key": string(keyPEM),
		"certificate": string(certPEM),
	}

	_, err := v.client.Logical().Write("secret/data/saml", map[string]interface{}{
		"data": data,
	})
	return err
}

func (v *VaultSecretManager) GetSAMLKey(ctx context.Context) (keyPEM, certPEM []byte, err error) {
	secret, err := v.client.Logical().Read("secret/data/saml")
	if err != nil {
		return nil, nil, err
	}
	if secret == nil || secret.Data["data"] == nil {
		return nil, nil, fmt.Errorf("saml secret not found")
	}

	data := secret.Data["data"].(map[string]interface{})
	return []byte(data["private_key"].(string)), []byte(data["certificate"].(string)), nil
}

func (v *VaultSecretManager) VerifyIdentity(ctx context.Context) error {
	// Look up the identity of the token to verify it's valid
	_, err := v.client.Auth().Token().LookupSelf()
	if err != nil {
		return fmt.Errorf("token lookup failed: %w", err)
	}
	return nil
}
