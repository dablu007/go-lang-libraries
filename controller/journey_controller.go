package controller

import (
	"encoding/json"
	"flow/auth"
	"flow/logger"
	"flow/service"
	"flow/utility"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
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
		var scopes = "internal_services"
		if !auth.ValidateScope(token, scopes) {
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
		var scopes = "internal_services"
		if !auth.ValidateScope(token, scopes) {
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
				c.JSON(http.StatusNotFound, gin.H{})
				return
			}
		}
		flow := u.journeyService.GetJourneyById(journeyId)
		if len(flow.Name) > 0 {
			c.JSON(http.StatusOK, flow)
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
	}
	return fn
}

func (u JourneyController) GetJourneyListForJourneyIds() gin.HandlerFunc {
	methodName := "GetJourneyById:"
	fn := func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.SugarLogger.Info(methodName, "Invalid body passed")
			c.JSON(http.StatusBadRequest, "Invalid Body passed")
			return
		}
		var journeyIds []string
		err = json.Unmarshal(body,&journeyIds)
		if err != nil {
			logger.SugarLogger.Info(methodName, "Unable to parse the body")
			c.JSON(http.StatusBadRequest, "Invalid Body passed")
			return
		}
		token := c.Request.Header.Get("Authorization")
		logger.SugarLogger.Info(methodName, "Recieved request to get Journey by JourneyId ", journeyIds)
		var scopes = "internal_services"
		if !auth.ValidateScope(token,scopes) {
			logger.SugarLogger.Info(methodName, "Invalid scope passed for fetching data ")
			c.JSON(http.StatusUnauthorized, "Invalid Scope")
			return
		}
		if len(journeyIds) <= 0 {
			logger.SugarLogger.Info(methodName, " journey id passed is empty or null ", journeyIds)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		flow := u.journeyService.GetJourneyDetailsListForJourneyIds(journeyIds)
		if len(flow) > 0 {
			c.JSON(http.StatusOK, flow)
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

	}
	return fn
}

