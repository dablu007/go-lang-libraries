package controller

import (
	"fmt"
	"github.com/dablu007/go-lang-libraries/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	logger.SugarLogger.Info("Success!")
	fmt.Print("Success!")
	c.String(http.StatusOK, "Working!")
}
