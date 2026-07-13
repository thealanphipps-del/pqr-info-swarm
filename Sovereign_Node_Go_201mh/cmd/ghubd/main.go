package main

import (
"bufio"
"encoding/base64"
"fmt"
"os"
"strings"
)

const (
GHUB_PATH = "/data/data/com.termux/files/home/Sovereign_Node_Go/ghub/"
PIPE_PATH = "/data/data/com.termux/files/home/.gemini_strike_pipe"
)

func main() {
fmt.Println("[GHUB] Remote Listener Active. Awaiting Gemini Offloads...")

for {
file, err := os.OpenFile(PIPE_PATH, os.O_RDONLY, os.ModeNamedPipe)
if err != nil {
continue
}

scanner := bufio.NewScanner(file)
for scanner.Scan() {
line := scanner.Text()
if strings.HasPrefix(line, "GHUB_PUSH|") {
parts := strings.Split(line, "|")
if len(parts) < 3 {
continue
}
filename := parts[1]
encodedData := parts[2]

data, err := base64.StdEncoding.DecodeString(encodedData)
if err != nil {
fmt.Printf("[GHUB] Error decoding %s: %v\n", filename, err)
continue
}

err = os.WriteFile(GHUB_PATH+filename, data, 0644)
if err == nil {
fmt.Printf("[GHUB] Offload Successful: %s (%d bytes)\n", filename, len(data))
}
}
}
file.Close()
}
}
