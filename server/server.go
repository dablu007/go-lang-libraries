package server

import (
	"flow/logger"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func Init() {
	r := NewRouter()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	err := r.Run(viper.GetString("server.port"))
	if err != nil {
		logger.SugarLogger.Error("Server not able to startup with error: ", err)
	}
}
