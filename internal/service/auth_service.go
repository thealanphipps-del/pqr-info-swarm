package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlidp"
	"github.com/thealanphipps-del/pqr/internal/domain"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/auth"
)

// AuthService handles SAML Identity Provider (IdP) logic
type AuthService struct {
	IDP  *samlidp.Server
	repo domain.UserRepository
}

// NewAuthService creates a new SAML IdP service
func NewAuthService(repo domain.UserRepository, baseURL string, key *rsa.PrivateKey, cert *x509.Certificate) (*AuthService, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	// Initialize the IdP server
	idpServer := &samlidp.Server{
		Store: &samlidp.MemoryStore{},
	}
	
	// Hydrate the internal IdentityProvider structure
	idpServer.IDP.Certificate = cert
	idpServer.IDP.Key = key
	idpServer.IDP.SSOURL = *u
	idpServer.IDP.SSOURL.Path = "/saml/sso"
	idpServer.IDP.MetadataURL = *u
	idpServer.IDP.MetadataURL.Path = "/saml/metadata"

	return &AuthService{
		IDP:  idpServer,
		repo: repo,
	}, nil
}

// HandleMetadata serves the SAML IdP metadata
func (s *AuthService) HandleMetadata(w http.ResponseWriter, r *http.Request) {
	s.IDP.ServeHTTP(w, r)
}

// HandleSSO handles the SAML Single Sign-On request
func (s *AuthService) HandleSSO(w http.ResponseWriter, r *http.Request) {
	s.IDP.ServeHTTP(w, r)
}

// AddUser adds a user to the IdP store (temporary for testing)
func (s *AuthService) AddUser(username, password, email, displayName string) error {
	// In a real implementation, this would sync with the database
	// For now, we'll add it to the IdP's internal store if it's using one, 
	// or just rely on the database during the actual login flow.
	return nil
}

// GetIdPHandler returns the http handler for SAML endpoints
func (s *AuthService) GetIdPHandler() http.Handler {
	return s.IDP
}

// RotateCertificates generates a new self-signed cert and updates the IdP
func (s *AuthService) RotateCertificates(ctx context.Context, commonName string) (*rsa.PrivateKey, *x509.Certificate, error) {
	privKey, cert, err := auth.GenerateSelfSignedCert(commonName)
	if err != nil {
		return nil, nil, err
	}

	// Update the live IdP server
	s.IDP.IDP.Key = privKey
	s.IDP.IDP.Certificate = cert
	
	// Re-initialize the IdP handler to pick up new certs
	// Note: samlidp.Server.Handler() creates a new handler each time
	
	return privKey, cert, nil
}
