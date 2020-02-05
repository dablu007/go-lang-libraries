package repository

import (
	"flow/db"
	"flow/enum"
	"flow/model"
	"flow/utility"
)

type JourneyRepository interface {
	FindByExternalId(flowExternalId string) model.Journey
	FindActiveJourneysByJourneyContext(merchantId string, tenantId string, channelId string) []model.Journey
}

type JourneyRepositoryImpl struct {
	MapUtil   utility.MapUtil
	DBService db.DBService
}

func NewJourneyRepository() *JourneyRepositoryImpl {
	repo := &JourneyRepositoryImpl{
		MapUtil:utility.MapUtil{},
		DBService:db.DBService{},
	}
	return repo
}

func NewFlowRepositoryImpl(DBService db.DBService) *JourneyRepositoryImpl {
	repo := &JourneyRepositoryImpl{
		DBService: DBService,
	}
	return repo
}

func (f JourneyRepositoryImpl) FindByExternalId(flowExternalId string) model.Journey {
	var journey model.Journey
	dbConnection := f.DBService.GetDB()
	dbConnection.Where(" external_id = ? ", flowExternalId).Find(&journey)
	return journey
}

func (f JourneyRepositoryImpl) FindActiveJourneysByJourneyContext(merchantId string, tennatId string, channelId string) []model.Journey {
	dbConnection := f.DBService.GetDB()
	var journeys []model.Journey
	if dbConnection == nil {
		return journeys
	}
	dbConnection.Where("flow_context->>'MerchantId' = ? and flow_context->>'TenantId' = ? and flow_context->>'ChannelId' = ? and status = ? and deleted_on is NULL", merchantId, tennatId, channelId, enum.Active).Find(&journeys)
	return journeys
}
