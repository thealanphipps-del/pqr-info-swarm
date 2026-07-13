package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"rt_api/internal/mcp"
	"fmt"
)

func main() {
	router := &mcp.TPipeRouter{
		CodeSide: "/data/data/com.termux/files/home/Sovereign_Node_Go",
		OSSide:   "/data/data/com.termux/files/usr/etc",
	}

	// SPAWN LOOPBACK TELNET IN GOROUTINE
	go mcp.StartTelnetExplorer("2323", router)

	r := gin.Default()
	
	// ZEROCODE MCP ENDPOINT FOR GEMINI INFERENCE
	r.POST("/api/v2/mcp/goreleaser", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "mcp_bridge_active",
			"tpipe_code": router.CodeSide,
			"tpipe_os": router.OSSide,
			"immutable": true,
		})
	})

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	fmt.Printf("[MASTER] Binding REST on %s, Telnet on 2323\n", port)
	r.Run(":" + port)
}
