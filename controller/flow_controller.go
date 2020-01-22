package controller

import (
	"flow/logger"
	"flow/service"
	"flow/utility"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FlowController struct {
	flowService      service.FlowService
	requestValidator utility.RequestValidator
}

/*
This method fetches and returns all the flows associated with given merchant.
*/
func (u FlowController) GetFlows() gin.HandlerFunc {
	methodName := "GetFlows:"
	fn := func(c *gin.Context) {
		merchantId := c.Query("merchantId")
		tenantId := c.Query("tenantId")
		channelId := c.Query("channelId")
		logger.SugarLogger.Info(methodName, "Recieved request to get all the flows associated with merchant: ", merchantId)
		if u.requestValidator.IsValidRequest(merchantId, tenantId, channelId) {
			flows := u.flowService.GetFlows(merchantId, tenantId, channelId)
			if len(flows.FlowResponses) > 0 {
				c.JSON(http.StatusOK, gin.H{"flows": flows})
				return
			} else {
				c.JSON(http.StatusNotFound, gin.H{})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	return fn
}