package service

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SAMLMetadataIndex represents the indexed SAML metadata
type SAMLMetadataIndex struct {
	SHA256             string    `json:"sha256"`
	LastFetch          time.Time `json:"last_fetch"`
	NextScheduled      time.Time `json:"next_scheduled"`
	Version            int       `json:"version"`
	CachedXML          string    `json:"cached_xml"`
}

// SAMLCertIndex represents the indexed SAML certificate
type SAMLCertIndex struct {
	SHA256          string    `json:"sha256"`
	Expiration      time.Time `json:"expiration"`
	RotationWindow  int       `json:"rotation_window_days"`
	VaultVersion    int       `json:"vault_version"`
	InMemoryPointer string    `json:"in_memory_pointer"`
}

// SAMLRotationState represents the active certificate rotation state
type SAMLRotationState struct {
	LastRotation   time.Time `json:"last_rotation"`
	NextRotation   time.Time `json:"next_rotation"`
	RotationReason string    `json:"rotation_reason"`
	CertVersion    int       `json:"cert_version"`
}

// SAMLAgent orchestrates the caching, indexing, and drift detection for SAML
type SAMLAgent struct {
	mu            sync.RWMutex
	auth          *AuthService
	metadataIndex *SAMLMetadataIndex
	certIndex     *SAMLCertIndex
	rotationState *SAMLRotationState
	
	indexPath     string
	certPath      string
	statePath     string
}

func NewSAMLAgent(auth *AuthService) *SAMLAgent {
	// Use standard paths with local fallbacks if system directories are unwritable
	idxDir := "/var/pqr/index"
	stateDir := "/var/pqr/state"
	
	if err := os.MkdirAll(idxDir, 0755); err != nil {
		log.Printf("[SAML-AGENT] /var/pqr/index unwritable, falling back to local certs/index")
		idxDir = "certs/index"
		os.MkdirAll(idxDir, 0755)
	}
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		log.Printf("[SAML-AGENT] /var/pqr/state unwritable, falling back to local certs/state")
		stateDir = "certs/state"
		os.MkdirAll(stateDir, 0755)
	}

	return &SAMLAgent{
		auth:      auth,
		indexPath: filepath.Join(idxDir, "saml_metadata.index"),
		certPath:  filepath.Join(idxDir, "saml_cert.index"),
		statePath: filepath.Join(stateDir, "saml_rotation.state"),
	}
}

// LoadFromDisk retrieves caches from the file indexes
func (a *SAMLAgent) LoadFromDisk() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Load Metadata Index
	if data, err := os.ReadFile(a.indexPath); err == nil {
		var idx SAMLMetadataIndex
		if err := json.Unmarshal(data, &idx); err == nil {
			a.metadataIndex = &idx
			log.Printf("[SAML-AGENT] Loaded cached Metadata Index version %d", idx.Version)
		}
	}

	// Load Cert Index
	if data, err := os.ReadFile(a.certPath); err == nil {
		var idx SAMLCertIndex
		if err := json.Unmarshal(data, &idx); err == nil {
			a.certIndex = &idx
			log.Printf("[SAML-AGENT] Loaded cached SAML Certificate Index")
		}
	}

	// Load Rotation State
	if data, err := os.ReadFile(a.statePath); err == nil {
		var state SAMLRotationState
		if err := json.Unmarshal(data, &state); err == nil {
			a.rotationState = &state
			log.Printf("[SAML-AGENT] Loaded active SAML Rotation State")
		}
	}

	return nil
}

// SaveToDisk persists the current index/cache states to disk
func (a *SAMLAgent) SaveToDisk() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.metadataIndex != nil {
		if data, err := json.MarshalIndent(a.metadataIndex, "", "  "); err == nil {
			_ = os.WriteFile(a.indexPath, data, 0644)
		}
	}

	if a.certIndex != nil {
		if data, err := json.MarshalIndent(a.certIndex, "", "  "); err == nil {
			_ = os.WriteFile(a.certPath, data, 0644)
		}
	}

	if a.rotationState != nil {
		if data, err := json.MarshalIndent(a.rotationState, "", "  "); err == nil {
			_ = os.WriteFile(a.statePath, data, 0644)
		}
	}

	return nil
}

// GetMetadata returns cached XML metadata without doing expensive disk/network operations
func (a *SAMLAgent) GetMetadata() (string, error) {
	a.mu.RLock()
	// If cached metadata exists and TTL has not expired (e.g. 300s)
	if a.metadataIndex != nil && time.Since(a.metadataIndex.LastFetch) < 300*time.Second {
		xml := a.metadataIndex.CachedXML
		a.mu.RUnlock()
		return xml, nil
	}
	a.mu.RUnlock()

	// If expired or missing, regenerate in-memory metadata via auth service
	if a.auth == nil || a.auth.IDP == nil {
		return "", fmt.Errorf("auth service metadata provider unavailable")
	}

	// Retrieve XML from crewjam/saml provider and marshal it
	metaXML, err := xml.MarshalIndent(a.auth.IDP.IDP.Metadata(), "", "  ")
	if err != nil {
		return "", err
	}

	h := sha256.Sum256(metaXML)
	hashStr := hex.EncodeToString(h[:])

	a.mu.Lock()
	version := 1
	if a.metadataIndex != nil {
		version = a.metadataIndex.Version
		if a.metadataIndex.SHA256 != hashStr {
			version++
		}
	}

	a.metadataIndex = &SAMLMetadataIndex{
		SHA256:        hashStr,
		LastFetch:     time.Now(),
		NextScheduled: time.Now().Add(300 * time.Second),
		Version:       version,
		CachedXML:     string(metaXML),
	}
	a.mu.Unlock()

	_ = a.SaveToDisk()
	return string(metaXML), nil
}

// UpdateCertIndex caches certificate state
func (a *SAMLAgent) UpdateCertIndex(cert *x509.Certificate) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if cert == nil {
		return
	}

	block := pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	pemBytes := pem.EncodeToMemory(&block)
	h := sha256.Sum256(pemBytes)
	hashStr := hex.EncodeToString(h[:])

	a.certIndex = &SAMLCertIndex{
		SHA256:          hashStr,
		Expiration:      cert.NotAfter,
		RotationWindow:  7,
		VaultVersion:    1,
		InMemoryPointer: fmt.Sprintf("%p", cert),
	}
	
	if a.rotationState == nil {
		a.rotationState = &SAMLRotationState{
			LastRotation:   time.Now(),
			NextRotation:   cert.NotAfter.AddDate(0, 0, -7),
			RotationReason: "Initialization",
			CertVersion:    1,
		}
	}
}

// DetectDrift checks if certificate needs immediate healing/rotation
func (a *SAMLAgent) DetectDrift() (bool, string) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.certIndex == nil {
		return true, "Missing certificate index information"
	}

	daysUntilExpiry := int(time.Until(a.certIndex.Expiration).Hours() / 24)
	if daysUntilExpiry < 7 {
		return true, fmt.Sprintf("Certificate expiring in %d days (threshold < 7 days)", daysUntilExpiry)
	}

	return false, ""
}
