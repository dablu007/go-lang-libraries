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
			logger.SugarLogger.Info(methodName, "Request is validated")
			flows := u.flowService.GetFlows(merchantId, tenantId, channelId)
			if len(flows.FlowResponses) > 0 {
				c.JSON(http.StatusOK, flows)
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

func (u FlowController) GetFlowById() gin.HandlerFunc  {
	methodName := "GetFlowById:"
	fn := func(c *gin.Context) {
		flowId := c.Param("flowId")
		logger.SugarLogger.Info(methodName, "Recieved request to get flow by flowId ", flowId)
		if len(flowId) <= 0 {
			logger.SugarLogger.Info(methodName, " Flow id passed is empty or null ", flowId)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		flow := u.flowService.GetFlowById(flowId)
		if len(flow.Name) > 0 {
			c.JSON(http.StatusOK, flow)
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
	}
	return fn
}
