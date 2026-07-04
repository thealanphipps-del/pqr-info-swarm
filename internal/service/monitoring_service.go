package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/thealanphipps-del/pqr/internal/domain"
	"golang.org/x/sys/windows"
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
				m.checkDiskSpace(ctx)
				
				metrics, _ := CollectHardwareMetrics()
				m.checkGPUOverload(ctx, metrics)
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

func (m *MonitoringService) checkDiskSpace(ctx context.Context) {
	var freeBytesAvailableToCaller, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	err := windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr("C:\\"),
		&freeBytesAvailableToCaller,
		&totalNumberOfBytes,
		&totalNumberOfFreeBytes)

	if err != nil {
		log.Printf("[MONITOR] Error checking disk space: %v", err)
		return
	}

	freeMB := freeBytesAvailableToCaller / (1024 * 1024)
	if freeMB < 2048 { // Less than 2GB
		msg := fmt.Sprintf("C:\\ drive has critically low space: %d MB remaining. Initiating autonomous cleanup.", freeMB)
		m.triggerHealing(ctx, "LOW_DISK_SPACE", msg)
		m.performAutonomousCleanup(ctx)
	}
}

func (m *MonitoringService) performAutonomousCleanup(ctx context.Context) {
	log.Println("[HEAL] Executing Autonomous Cleanup: Purging Go build cache...")
	cmd := exec.Command("go", "clean", "-cache")
	if err := cmd.Run(); err != nil {
		log.Printf("[HEAL] Error clearing Go cache: %v", err)
	}

	log.Println("[HEAL] Executing Autonomous Cleanup: Purging Temp folder...")
	psCmd := `Get-ChildItem -Path $env:TEMP -Recurse | Remove-Item -Force -Recurse -ErrorAction SilentlyContinue`
	cmd2 := exec.Command("powershell", "-Command", psCmd)
	if err := cmd2.Run(); err != nil {
		log.Printf("[HEAL] Error clearing Temp folder: %v", err)
	}
	
	log.Println("[HEAL] Autonomous Cleanup complete.")
}

var gpu0OverloadCounter int
var gpu1OverloadCounter int
var gpu2OverloadCounter int // NPU

func (m *MonitoringService) checkGPUOverload(ctx context.Context, metrics HardwareMetrics) {
	m.trackGPU(ctx, "gpu0", metrics.GPU0Percent, &gpu0OverloadCounter)
	m.trackGPU(ctx, "gpu1", metrics.GPU1Percent, &gpu1OverloadCounter)
	m.trackGPU(ctx, "gpu2", metrics.GPU2Percent, &gpu2OverloadCounter) // NPU
}

func (m *MonitoringService) trackGPU(ctx context.Context, label string, percent float64, counter *int) {
	if percent >= 99.0 {
		*counter++
	} else {
		*counter = 0
		return
	}

	// 3 minutes = 3 iterations (1/min)
	if *counter == 3 {
		m.triggerHealing(ctx, label+"_OVERLOAD_WARN", label+" overload detected for 3 minutes")
		m.broadcastAlert(label + " at 100% for 3 minutes - investigating cause")
		m.searchForGPUHog(label)
	}

	// 12 iterations = 12 minutes
	if *counter >= 12 {
		m.broadcastAlert(label + " overload unresolved after 12 minutes - mission critical")
		m.triggerHealing(ctx, label+"_OVERLOAD_CRIT", label+" overload unresolved after 12 minutes")
	}
}

func (m *MonitoringService) searchForGPUHog(label string) {
	// TODO: integrate with nvidia-smi or Windows GPU APIs
	// For now, placeholder:
	log.Printf("[monitor] searching for %s hog process...\n", label)
}

func (m *MonitoringService) broadcastAlert(msg string) {
	// TODO: integrate with cockpit notifications
	log.Println("[ALERT]", msg)
}
