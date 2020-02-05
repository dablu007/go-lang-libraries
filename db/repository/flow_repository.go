package repository

import (
	"flow/db"
	"flow/enum"
	"flow/model"
	"flow/utility"
)

type FlowRepository interface {
	FindByExternalId(flowExternalId string) model.Flow
	FindActiveFlowsByFlowContext(merchantId string, tenantId string, channelId string) []model.Flow
}

type FlowRepositoryImpl struct {
	MapUtil   utility.MapUtil
	DBService db.DBService
}

func NewFlowRepositoryImpl(DBService db.DBService) *FlowRepositoryImpl {
	repo := &FlowRepositoryImpl{
		DBService: DBService,
	}
	return repo
}

func (f FlowRepositoryImpl) FindByExternalId(flowExternalId string) model.Flow {
	var flow model.Flow
	dbConnection := f.DBService.GetDB()
	dbConnection.Where(" external_id = ? ", flowExternalId).Find(&flow)
	return flow
}

func (f FlowRepositoryImpl) FindActiveFlowsByFlowContext(merchantId string, tennatId string, channelId string) []model.Flow {
	dbConnection := f.DBService.GetDB()
	var flows []model.Flow
	if dbConnection == nil {
		return flows
	}
	dbConnection.Where("flow_context->>'MerchantId' = ? and flow_context->>'TenantId' = ? and flow_context->>'ChannelId' = ? and status = ? and deleted_on is NULL", merchantId, tennatId, channelId, enum.Active).Find(&flows)
	return flows
}
