package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Status godoc
// @Summary Show the status of the API
// @Description get the status of the API
// @Tags status
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router /status [get]
func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
