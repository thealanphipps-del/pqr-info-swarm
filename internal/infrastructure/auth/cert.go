package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

// GenerateSelfSignedCert generates a temporary RSA key and certificate for SAML
func GenerateSelfSignedCert(commonName string) (*rsa.PrivateKey, *x509.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"PQR Sovereign Node"},
			CommonName:   commonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	return priv, cert, nil
}

// EncodeCertToPEM encodes the certificate to PEM format
func EncodeCertToPEM(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
}

// EncodeKeyToPEM encodes the private key to PEM format
func EncodeKeyToPEM(key *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
}

// LoadCertFromFiles reads a PEM cert and key from disk
func LoadCertFromFiles(certPath, keyPath string) (*rsa.PrivateKey, *x509.Certificate, error) {
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}

	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil {
		return nil, nil, fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("failed to decode key PEM")
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		// Try PKCS8 if PKCS1 fails
		pk8Key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, nil, err
		}
		rsaKey, ok := pk8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, fmt.Errorf("not an RSA private key")
		}
		key = rsaKey
	}

	return key, cert, nil
}
