package controller

import (
	"flow/service"
	"flow/utility"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FlowController struct{}

/*
This method fetches and returns all the flows associated with given merchant.
*/
func (u FlowController) GetAllFlows() gin.HandlerFunc {
	fmt.Print("Reached controller method")
	fn := func(c *gin.Context) {
		merchantId := c.Query("merchantId")
		tenantId := c.Query("tenantId")
		channel := c.Query("X-Channel")
		if utility.IsValidUUID(merchantId) {
			fmt.Print("merchantId is valid", merchantId)
			flows := service.GetAllFlowsByMerchantId(merchantId, tenantId, channel)
			c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved data!", "flows": flows})
			return
		}
		fmt.Print("InValid merchantId.")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	return fn
}

func (u FlowController) TeaPot() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hi, I am a Teapot."})
	}
	return fn
}
