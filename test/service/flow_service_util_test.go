package service

import (
	"flow/config"
	"flow/logger"
	"flow/model"
	"flow/service"
	mock_repository "flow/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func init() {
	service := "finbox-integration"
	environment := "dev"
	config.Init(service, environment)
	logger.InitLogger()
}

func TestFetchFlowByIdFromDB(t *testing.T) {
	var controller = gomock.NewController(t)
	defer controller.Finish()

	var flow model.Journey
	flow.Id = 1
	flow.CreatedOn = time.Now()

	var journeyRepository = mock_repository.NewMockJourneyRepository(controller)
	var externalId = "123"
	journeyRepository.EXPECT().FindByExternalId(externalId).Return(flow)

	var journeyService = &service.JourneyServiceUtil{
		JourneyRepository: journeyRepository,
	}
	var flowActual = journeyService.FetchJourneyByIdFromDB(externalId)
	assert.Equal(t, flowActual.Id, flow.Id)


	journeyRepository.EXPECT().FindByExternalId("").Return(model.Journey{})

	journeyService = &service.JourneyServiceUtil{
		JourneyRepository: journeyRepository,
	}
	flowActual = journeyService.FetchJourneyByIdFromDB("")
	assert.Equal(t, len(flowActual.Name), 0)
}
