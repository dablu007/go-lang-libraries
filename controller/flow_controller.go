package controller

import (
	"flow/service"
	"flow/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FlowController struct {
	flowService service.FlowService
}

/*
This method fetches and returns all the flows associated with given merchant.
*/
func (u FlowController) GetFlows() gin.HandlerFunc {
	//Add log
	fn := func(c *gin.Context) {
		merchantId := c.Query("merchantId")
		if utility.IsValidUUID(merchantId) {
			//Add log
			flows := u.flowService.GetMerchantFlows(merchantId)
			if len(flows.FlowResponses) > 0 {
				c.JSON(http.StatusOK, gin.H{"flows": flows})
				return
			} else {
				c.JSON(http.StatusNotFound, gin.H{})
				return
			}
		}
		//Add log
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	return fn
}
