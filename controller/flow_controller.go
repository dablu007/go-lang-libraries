package controller

import (
	"flow/repository"
	"flow/utility"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FlowController struct{}

/*
This method fetches and returns all the flows associated with given merchant.
*/
func (u FlowController) GetAllFlows() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		merchantId := c.Param("merchantId")
		if utility.IsValidUUID(merchantId) {
			flows, err := repository.GetAllFlowByMerchantId(merchantId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving data", "error": err})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved data!", "flows": flows})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	return fn
}
