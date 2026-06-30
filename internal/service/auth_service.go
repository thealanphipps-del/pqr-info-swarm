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
	IDP   *samlidp.Server
	repo  domain.UserRepository
	Agent *SAMLAgent
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

	authService := &AuthService{
		IDP:  idpServer,
		repo: repo,
	}

	authService.Agent = NewSAMLAgent(authService)
	_ = authService.Agent.LoadFromDisk()
	authService.Agent.UpdateCertIndex(cert)
	_ = authService.Agent.SaveToDisk()

	return authService, nil
}

// HandleMetadata serves the SAML IdP metadata using cached XML in SAMLAgent
func (s *AuthService) HandleMetadata(w http.ResponseWriter, r *http.Request) {
	if s.Agent != nil {
		meta, err := s.Agent.GetMetadata()
		if err == nil {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(meta))
			return
		}
	}
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
	user := samlidp.User{
		Name:       username,
		Groups:     []string{"users"},
		Email:      email,
		CommonName: displayName,
	}
	if store, ok := s.IDP.Store.(*samlidp.MemoryStore); ok {
		return store.Put("/users/"+username, user)
	}
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
	
	if s.Agent != nil {
		s.Agent.UpdateCertIndex(cert)
		_ = s.Agent.SaveToDisk()
	}
	
	return privKey, cert, nil
}
