package main

import (
"crypto/sha256"
"fmt"
"io"
"os"
"os/exec"
"path/filepath"
"time"
)

// Sentry context utilizing pointer receivers
type Sentry struct {
TargetIP string
CaseDir  string
AdbPath  string
LogFile  *os.File
}

func (s *Sentry) InitLogger() error {
logPath := filepath.Join(s.CaseDir, "MASTER_CUSTODY_LOG.txt")
f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
if err != nil {
return err
}
s.LogFile = f
s.Audit("SENTRY_INITIALIZED", "Chain of Custody Protocol Active")
return nil
}

func (s *Sentry) Audit(action, details string) {
timestamp := time.Now().UTC().Format(time.RFC3339)
entry := fmt.Sprintf("[%s] | UID: 203 | ACTION: %s | DETAILS: %s\n", timestamp, action, details)
s.LogFile.WriteString(entry)
fmt.Print(entry)
}

func (s *Sentry) HashAndStore(filename, data string) {
outPath := filepath.Join(s.CaseDir, "telemetry", filename)
os.WriteFile(outPath, []byte(data), 0600)

h := sha256.New()
h.Write([]byte(data))
hash := fmt.Sprintf("%x", h.Sum(nil))

s.Audit("EVIDENCE_SEALED", fmt.Sprintf("File: %s | SHA256: %s", filename, hash))
}

func (s *Sentry) ExtractTelemetry(uri, filename string) {
s.Audit("EXTRACTION_START", "URI: "+uri)
cmd := exec.Command(s.AdbPath, "-s", s.TargetIP, "shell", "content", "query", "--uri", uri)
out, err := cmd.CombinedOutput()
if err != nil {
s.Audit("EXTRACTION_ERROR", fmt.Sprintf("URI: %s | Err: %v", uri, err))
return
}
s.HashAndStore(filename, string(out))
}

func (s *Sentry) ExtractDumpsys(service, filename string) {
s.Audit("DUMPSYS_START", "Service: "+service)
cmd := exec.Command(s.AdbPath, "-s", s.TargetIP, "shell", "dumpsys", service)
out, err := cmd.CombinedOutput()
if err != nil {
s.Audit("DUMPSYS_ERROR", fmt.Sprintf("Service: %s | Err: %v", service, err))
return
}
s.HashAndStore(filename, string(out))
}

func (s *Sentry) GenerateManifests() {
s.Audit("MANIFEST_START", "Scanning User 0 and User 150 partitions")

// User 0
cmd0 := exec.Command(s.AdbPath, "-s", s.TargetIP, "shell", "find", "/sdcard", "-type", "f", "2>/dev/null")
out0, _ := cmd0.CombinedOutput()
s.HashAndStore("../manifests/USER_0_MANIFEST.txt", string(out0))

// User 150 (Secure Folder boundary attempt)
cmd150 := exec.Command(s.AdbPath, "-s", s.TargetIP, "shell", "find", "/storage/emulated/150", "-type", "f", "2>/dev/null")
out150, _ := cmd150.CombinedOutput()
s.HashAndStore("../manifests/USER_150_MANIFEST.txt", string(out150))

// Cross-Reference Hook for Google Takeout Logic
s.Audit("CROSS_REF_READY", "Takeout ingestion endpoints primed for next phase.")
}

func main() {
agent := &Sentry{
TargetIP: os.Args[1],
CaseDir:  os.Args[2],
AdbPath:  os.Args[3],
}

if err := agent.InitLogger(); err != nil {
panic(err)
}
defer agent.LogFile.Close()

agent.ExtractTelemetry("content://call_log/calls", "CALL_LOGS_FORENSIC.txt")
agent.ExtractTelemetry("content://sms", "SMS_FORENSIC.txt")
agent.ExtractDumpsys("location", "GPS_LOCATION_TELEMETRY.txt")
agent.ExtractDumpsys("telephony.registry", "TOWER_PING_REGISTRY.txt")
agent.GenerateManifests()
}
