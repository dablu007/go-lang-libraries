package server

import (
	"flow/controller"
	"flow/db"
	"flow/db/repository"
	"flow/service"
	"flow/utility"
)

var validator *utility.RequestValidator
var mapUtil *utility.MapUtil
var dbService *db.DBService
var journeyRepo repository.JourneyRepository
var fieldRepo repository.FieldRepository
var moduleRepo repository.ModuleRepository
var sectionRepo repository.SectionRepository
var journeyServiceUtil *service.JourneyServiceUtil
var journeyService *service.JourneyService
var flowController *controller.JourneyController


func init(){
	validator = utility.NewRequestValidator()
	mapUtil = utility.NewMapUtil()
	dbService = new(db.DBService)
	journeyRepo = repository.NewJourneyRepository()
	fieldRepo = repository.NewFieldRepository()
	moduleRepo = repository.NewModuleRepository()
	sectionRepo = repository.NewSectionRepository()
	journeyServiceUtil = service.NewJourneyServiceUtil(mapUtil, dbService, journeyRepo, fieldRepo, moduleRepo, sectionRepo)
	journeyService = service.NewJourneyService(journeyServiceUtil, validator, moduleRepo, journeyRepo)
	flowController = controller.NewJourneyController(journeyService, validator)
}
