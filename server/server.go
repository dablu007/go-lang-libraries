package server

import (
	"flow/config"
	"flow/logger"

	"github.com/gin-gonic/gin"
)

func Init() {
	serverConfig := config.GetConfig()
	r := NewRouter()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	err := r.Run(serverConfig.GetString("server.port"))
	if err != nil {
		logger.SugarLogger.Error("Server not able to startup with error: ", err)
	}
}
