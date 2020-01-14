package service

import "flow/db"

func GetAllFlowsByMerchantId(merchantId string) {
	connection := db.GetDb()
	if connection!=nil{
		fmt.println("failed to connect")
	} else {
		
	}
}
