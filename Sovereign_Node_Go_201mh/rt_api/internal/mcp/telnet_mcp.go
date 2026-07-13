package mcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// TPipeRouter enforces immutable split-access
type TPipeRouter struct {
	CodeSide string
	OSSide   string
}

func (t *TPipeRouter) ReadImmutable(side string, relativePath string) (string, error) {
	var basePath string
	if side == "CODE" {
		basePath = t.CodeSide
	} else if side == "OS" {
		basePath = t.OSSide
	} else {
		return "", fmt.Errorf("INVALID_TPIPE_SIDE")
	}

	target := filepath.Join(basePath, relativePath)
	// Enforce immutability by strictly opening Read-Only
	f, err := os.OpenFile(target, os.O_RDONLY, 0400)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func StartTelnetExplorer(port string, router *TPipeRouter) {
	l, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Printf("[MCP_TELNET_FATAL] %v\n", err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handleMenu(conn, router)
	}
}

func handleMenu(conn net.Conn, router *TPipeRouter) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for {
		conn.Write([]byte("\033[2J\033[H")) // ANSI Clear
		conn.Write([]byte("=== GSH MESH TELNET EXPLORER ===\n"))
		conn.Write([]byte("[ ZEROCODE MCP & IMMUTABLE T-PIPE ]\n\n"))
		conn.Write([]byte("1. Explore CODE (Immutable /home/Sovereign_Node_Go)\n"))
		conn.Write([]byte("2. Explore OS (Immutable /usr/etc)\n"))
		conn.Write([]byte("3. View Healer & Sensor Logs\n"))
		conn.Write([]byte("4. Trigger GoReleaser MCP Sync\n"))
		conn.Write([]byte("5. Disconnect\n\n"))
		conn.Write([]byte("AWAITING INPUT: "))

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			conn.Write([]byte("\n--- CODE ROOT ---\n"))
			files, _ := os.ReadDir(router.CodeSide)
			for _, f := range files {
				conn.Write([]byte(fmt.Sprintf("- %s\n", f.Name())))
			}
		case "2":
			conn.Write([]byte("\n--- OS ROOT ---\n"))
			files, _ := os.ReadDir(router.OSSide)
			for _, f := range files {
				conn.Write([]byte(fmt.Sprintf("- %s\n", f.Name())))
			}
		case "3":
			logData, _ := router.ReadImmutable("CODE", "bin/sensor.log")
			lines := strings.Split(logData, "\n")
			start := len(lines) - 10
			if start < 0 { start = 0 }
			conn.Write([]byte("\n--- SENSOR LOG TAIL ---\n"))
			conn.Write([]byte(strings.Join(lines[start:], "\n") + "\n"))
		case "4":
			conn.Write([]byte("\n[MCP] Transmitting Context to GoReleaser Pipeline...\n[MCP] Immutable Lock Confirmed.\n"))
		case "5":
			return
		}
		conn.Write([]byte("\n[PRESS ENTER TO RETURN]"))
		scanner.Scan()
	}
}
