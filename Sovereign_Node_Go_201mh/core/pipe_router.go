package main

import (
"bufio"
"bytes"
"log"
"net/http"
"os"
"os/exec"
"strings"
)

type PipeRouter struct {
PipeFile string
OutDir   string
}

func (pr *PipeRouter) ListenAndRoute() {
for {
file, err := os.OpenFile(pr.PipeFile, os.O_RDONLY, os.ModeNamedPipe)
if err != nil {
log.Println("Pipe open error:", err)
continue
}
scanner := bufio.NewScanner(file)
for scanner.Scan() {
line := scanner.Text()
if strings.HasPrefix(line, "WORKSPACE_SYNC:") {
payload := strings.TrimPrefix(line, "WORKSPACE_SYNC:")
pr.syncToWorkspace(payload)
} else if strings.HasPrefix(line, "NAMESPACE_EXEC:") {
cmd := strings.TrimPrefix(line, "NAMESPACE_EXEC:")
pr.execProotNamespace(cmd)
} else if strings.HasPrefix(line, "GITHUB_GIST:") {
payload := strings.TrimPrefix(line, "GITHUB_GIST:")
pr.syncToGitHub(payload)
}
}
file.Close()
}
}

func (pr *PipeRouter) syncToWorkspace(data string) {
// ActiveJob pattern push to Google Apps Script Web App
log.Println("Routing to Google Workspace App Script:", data)
req, _ := http.NewRequest("POST", "http://localhost:8080/api/workspace/sync", bytes.NewBuffer([]byte(data)))
req.Header.Set("Content-Type", "application/json")
client := &http.Client{}
client.Do(req)
}

func (pr *PipeRouter) syncToGitHub(data string) {
log.Println("Routing to GitHub API:", data)
// Execution logic for GitHub API via local gh CLI
cmd := exec.Command("/data/data/com.termux/files/usr/bin/gh", "gist", "create", "-", "--desc", "Sovereign_Node_Sync")
cmd.Stdin = strings.NewReader(data)
cmd.Run()
}

func (pr *PipeRouter) execProotNamespace(cmdStr string) {
// Termux Docker alternative using proot-distro namespace isolation
log.Println("Executing in proot namespace:", cmdStr)
cmd := exec.Command("/data/data/com.termux/files/usr/bin/proot-distro", "login", "debian", "--", "bash", "-c", cmdStr)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Run()
}

func main() {
router := &PipeRouter{
PipeFile: "/data/data/com.termux/files/home/Sovereign_Node_Go/tmp/gemini_samba_bridge.fifo",
OutDir:   "/data/data/com.termux/files/home/Sovereign_Node_Go/inbound",
}
log.Println("Sovereign Pipe Router Active...")
router.ListenAndRoute()
}
