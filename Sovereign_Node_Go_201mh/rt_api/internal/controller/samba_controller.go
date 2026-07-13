package controller

import (
	"github.com/gin-gonic/gin"
	"rt_api/internal/samba"
	"net/http"
)

func SambaIOHandler(c *gin.Context) {
	var req struct {
		Path  string `json:"path"`
		Data  string `json:"data"`
		Write bool   `json:"write"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// AUTH CREDENTIALS FETCHED FROM 39.MH ENV VARS
	cfg := samba.SambaConfig{
		Addr:  "10.0.39.1:445",
		User:  "alan_phipps",
		Pass:  "KEEP_SECRET", 
		Share: "Sovereign_Project",
	}

	res, err := samba.ExecuteSambaIO(cfg, req.Path, []byte(req.Data), req.Write)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payload": res})
}
