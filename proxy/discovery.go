package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func findOllamaServer(port string) (string, error) {
	// Find local IP
	var localIP string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				break
			}
		}
	}

	if localIP == "" {
		return "", fmt.Errorf("could not determine local IP")
	}

	parts := strings.Split(localIP, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid local IP format")
	}

	subnet := fmt.Sprintf("%s.%s.%s.", parts[0], parts[1], parts[2])
	
	fmt.Printf("Scanning subnet %sx for port %s...\n", subnet, port)

	var wg sync.WaitGroup
	resultChan := make(chan string, 1)
	doneChan := make(chan struct{})

	for i := 1; i < 255; i++ {
		targetIP := fmt.Sprintf("%s%d", subnet, i)
		if targetIP == localIP {
			continue // Skip ourselves if we know we don't run it locally
		}
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			target := net.JoinHostPort(ip, port)
			conn, err := net.DialTimeout("tcp", target, 500*time.Millisecond)
			if err == nil {
				conn.Close()
				select {
				case resultChan <- ip:
				default:
				}
			}
		}(targetIP)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case ip := <-resultChan:
		return ip, nil
	case <-doneChan:
		return "", fmt.Errorf("no server found on port %s", port)
	case <-time.After(2 * time.Second):
		return "", fmt.Errorf("scan timed out")
	}
}
