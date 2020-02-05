package controller

import (
	"flow/logger"
	"flow/service"
	"flow/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JourneyController struct {
	journeyService   service.JourneyService
	requestValidator utility.RequestValidator
}

/*
This method fetches and returns all the flows associated with given merchant.
*/
func (u JourneyController) GetJourneys() gin.HandlerFunc {
	methodName := "GetJourneys:"
	fn := func(c *gin.Context) {
		merchantId := c.Query("merchantId")
		tenantId := c.Query("tenantId")
		channelId := c.Query("channelId")
		logger.SugarLogger.Info(methodName, "Recieved request to get all the Journeys associated with merchant: ", merchantId)
		if u.requestValidator.IsValidRequest(merchantId, tenantId, channelId) {
			logger.SugarLogger.Info(methodName, "Request is validated")
			flows := u.journeyService.GetJourneys(merchantId, tenantId, channelId)
			if len(flows.JourneyResponses) > 0 {
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

func (u JourneyController) GetJourneyById() gin.HandlerFunc {
	methodName := "GetJourneyById:"
	fn := func(c *gin.Context) {
		journeyId := c.Param("journeyId")
		logger.SugarLogger.Info(methodName, "Recieved request to get Journey by JourneyId ", journeyId)
		if len(journeyId) <= 0 {
			logger.SugarLogger.Info(methodName, " journey id passed is empty or null ", journeyId)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		flow := u.journeyService.GetJourneyById(journeyId)
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
