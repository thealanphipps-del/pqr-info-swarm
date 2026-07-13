package main

import (
"bufio"
"fmt"
"os"
"time"
)

func main() {
fmt.Println("[GMODEM_UPGRADE] GMUDD interactive router ignited (LOCAL_LOOPBACK).")
file, err := os.Open("gmudd_inbox.log")
if err != nil {
fmt.Println("Error (1): Cannot open IPC log:", err)
return
}
defer file.Close()

reader := bufio.NewReader(file)
for {
line, err := reader.ReadString('\n')
if err == nil {
fmt.Print("[MASTER_RECEIVED]: ", line)
} else {
time.Sleep(2 * time.Second) // 2000ms Poll
}
}
}
