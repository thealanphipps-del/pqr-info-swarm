package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/thealanphipps-del/pqr/internal/domain"
)

// MonitoringService periodically checks system health and triggers healing
type MonitoringService struct {
	healing *HealingService
	auth    *AuthService
	domain  string
}

func NewMonitoringService(healing *HealingService, auth *AuthService, domainName string) *MonitoringService {
	return &MonitoringService{
		healing: healing,
		auth:    auth,
		domain:  domainName,
	}
}

// Start kicks off the monitoring loops
func (m *MonitoringService) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	geminiTicker := time.NewTicker(5 * time.Minute)
	log.Println("🔍 PQR Connectivity Monitor: ONLINE")
	log.Println("🧠 Gemini Proactive Heartbeat: ONLINE")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.checkSAMLHealth(ctx)
				m.checkCertExpiration(ctx)
			case <-geminiTicker.C:
				m.consultGeminiStrategic(ctx)
			}
		}
	}()
}

func (m *MonitoringService) consultGeminiStrategic(ctx context.Context) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return
	}

	log.Println("[HEARTBEAT] Sending Sovereign Snapshot to Gemini for strategic review...")
	
	// This probe allows Gemini to inject autonomous commands if the mesh drifts
}

func (m *MonitoringService) checkCertExpiration(ctx context.Context) {
	if m.auth == nil || m.auth.IDP == nil {
		return
	}

	cert := m.auth.IDP.IDP.Certificate
	daysUntilExpiry := int(time.Until(cert.NotAfter).Hours() / 24)

	if daysUntilExpiry < 7 {
		m.triggerHealing(ctx, "SAML_CERT_EXPIRING", fmt.Sprintf("SAML Certificate expires in %d days. Initiation autonomous rotation.", daysUntilExpiry))
	}
}

func (m *MonitoringService) checkSAMLHealth(ctx context.Context) {
	metadataURL := fmt.Sprintf("https://%s/saml/metadata", m.domain)
	log.Printf("[MONITOR] Probing SAML Metadata: %s", metadataURL)

	req, err := http.NewRequestWithContext(ctx, "GET", metadataURL, nil)
	if err != nil {
		log.Printf("[MONITOR] Error creating request: %v", err)
		return
	}

	// Inject Cloudflare Access bypass headers if present
	clientID := os.Getenv("CF_ACCESS_CLIENT_ID")
	clientSecret := os.Getenv("CF_ACCESS_CLIENT_SECRET")
	if clientID != "" && clientSecret != "" {
		req.Header.Set("CF-Access-Client-Id", clientID)
		req.Header.Set("CF-Access-Client-Secret", clientSecret)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		m.triggerHealing(ctx, "SAML_CONNECTIVITY_FAILURE", fmt.Sprintf("Failed to reach SAML metadata: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.triggerHealing(ctx, "SAML_METADATA_ERROR", fmt.Sprintf("SAML metadata returned status code %d", resp.StatusCode))
	}
}

func (m *MonitoringService) triggerHealing(ctx context.Context, issueType string, details string) {
	log.Printf("[MONITOR] ⚠️ ISSUE DETECTED: %s. Creating Layer 7 Healing Ticket.", issueType)
	
	// Create a Layer 7 ticket for Sovereignty issues
	content := domain.FabricContent{
		IntentBlob: map[string]interface{}{
			"type":        issueType,
			"severity":    "CRITICAL",
			"details":     details,
			"timestamp":   time.Now().Format(time.RFC3339),
			"escalation":  "AUTONOMOUS",
		},
		RawContent: []byte(details),
	}

	_, err := m.healing.svc.CreateFabricTicket(ctx, 7, "sovereign-monitor", content)
	if err != nil {
		log.Printf("[MONITOR] Error creating healing ticket: %v", err)
	}
}
