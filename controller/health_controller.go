package controller

import (
	"boiler-plate/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	logger.SugarLogger.Info("Success!")
	c.String(http.StatusOK, "Working!")
}
