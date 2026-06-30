package service

import (
	"bufio"
	"context"
	"log"
	"os/exec"
	"regexp"
	"sync"
)

type TunnelAgent struct {
	mu           sync.RWMutex
	cmd          *exec.Cmd
	failoverURL  string
	isRunning    bool
	stopChan     chan struct{}
}

var (
	cfURLRegex = regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)
)

func NewTunnelAgent() *TunnelAgent {
	return &TunnelAgent{
		stopChan: make(chan struct{}),
	}
}

// StartFailoverTunnel launches trycloudflare quick tunnel asynchronously
func (t *TunnelAgent) StartFailoverTunnel(ctx context.Context) error {
	t.mu.Lock()
	if t.isRunning {
		t.mu.Unlock()
		return nil
	}
	t.isRunning = true
	t.mu.Unlock()

	log.Println("[TUNNEL-AGENT] Spin-up Initiated: TryCloudflare failover route...")

	// Spawn cloudflared tunnel pointing to the local Go server port 8196
	cmd := exec.CommandContext(ctx, "/usr/local/bin/cloudflared", "tunnel", "--url", "http://localhost:8196")
	
	stdout, err := cmd.StderrPipe() // cloudflared output logs are typically on stderr
	if err != nil {
		t.mu.Lock()
		t.isRunning = false
		t.mu.Unlock()
		return err
	}

	if err := cmd.Start(); err != nil {
		t.mu.Lock()
		t.isRunning = false
		t.mu.Unlock()
		return err
	}

	t.mu.Lock()
	t.cmd = cmd
	t.mu.Unlock()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if match := cfURLRegex.FindString(line); match != "" {
				t.mu.Lock()
				t.failoverURL = match
				t.mu.Unlock()
				log.Printf("[TUNNEL-AGENT] ✓ Active TryCloudflare failover tunnel established: %s", match)
			}
		}
	}()

	go func() {
		// Wait for command termination
		_ = cmd.Wait()
		t.mu.Lock()
		t.isRunning = false
		t.failoverURL = ""
		t.mu.Unlock()
		log.Println("[TUNNEL-AGENT] TryCloudflare tunnel terminated.")
	}()

	return nil
}

// GetFailoverURL returns the active trycloudflare failover address
func (t *TunnelAgent) GetFailoverURL() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.failoverURL
}

// Stop terminates the active trycloudflare tunnel
func (t *TunnelAgent) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.cmd != nil && t.cmd.Process != nil {
		_ = t.cmd.Process.Kill()
	}
	t.isRunning = false
	t.failoverURL = ""
}
