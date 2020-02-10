package controller

import (
	"flow/auth"
	"flow/logger"
	"flow/service"
	"flow/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JourneyController struct {
	journeyService   *service.JourneyService
	requestValidator *utility.RequestValidator
}

/*
This method fetches and returns all the flows associated with given merchant.
*/

func NewJourneyController(journeyService *service.JourneyService, validator *utility.RequestValidator) *JourneyController {
	controller := &JourneyController{
		journeyService:   journeyService,
		requestValidator: validator,
	}
	return controller
}
func (u JourneyController) GetJourneys() gin.HandlerFunc {
	methodName := "GetJourneys:"
	fn := func(c *gin.Context) {
		merchantId := c.Query("merchantId")
		tenantId := c.Query("tenantId")
		channelId := c.Query("channelId")
		token := c.Request.Header.Get("Authorization")
		logger.SugarLogger.Info(methodName, "Recieved request to get all the Journeys associated with merchant: ", merchantId)

		if !auth.ValidateScope(token) {
			logger.SugarLogger.Info(methodName, "Invalid scope passed for fetching data ")
			c.JSON(http.StatusUnauthorized, "Invalid Scope")
			return
		}
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
		var isNested = c.Query("isNested")
		var nestedValue, err = strconv.ParseBool(isNested)
		token := c.Request.Header.Get("Authorization")
		logger.SugarLogger.Info(methodName, "Recieved request to get Journey by JourneyId ", journeyId)
		if !auth.ValidateScope(token) {
			logger.SugarLogger.Info(methodName, "Invalid scope passed for fetching data ")
			c.JSON(http.StatusUnauthorized, "Invalid Scope")
			return
		}
		if len(journeyId) <= 0 {
			logger.SugarLogger.Info(methodName, " journey id passed is empty or null ", journeyId)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		if err != nil && len(isNested) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		if nestedValue == false {
			flow := u.journeyService.GetJourneyDetailsAsList(journeyId)
			if len(flow.Name) > 0 {
				c.JSON(http.StatusOK, flow)
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			}
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
