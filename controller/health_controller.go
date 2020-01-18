package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"flow/logger"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	logger.SugarLogger.Info("Success!")
	c.String(http.StatusOK, "Working!")
}
