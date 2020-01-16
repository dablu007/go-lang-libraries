package service

import (
	"encoding/json"
	"flow/cache"
	"flow/db"
	"flow/model"
	"flow/model/response_dto"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func GetAllFlowsByMerchantId(merchantId string, tenantId string, channel string) response_dto.FlowResponseDto {
	fmt.Println("Trying to fetch from redis cache.")
	redisClient := cache.GetRedisClient()
	var merchantFlows []uint8
	var response response_dto.FlowResponseDto
	if redisClient != nil {
		flows, err := redisClient.Get(merchantId).Result()
		fmt.Println("flows", flows)
		if err != nil {
			fmt.Println("Failed to fetch from redis with error: ", err)
		}
		if flows == "" {
			fmt.Println("Fetching all from db")
			flow := FetchAllMerchantFlowsFromDB(merchantId)
			merchantFlows = *flow.ModuleVersions
			if merchantFlows != nil {
				flowResponseDto, err := ParseMerchantFlowsJsonValues(flow)
				if err != nil {

					fmt.Println(err)
				}
				redisClient.Set(merchantId, flowResponseDto, 0)
				response = flowResponseDto
			}
		}
	} else {
		fmt.Println("Failed to connect with Redis.", redisClient)
	}
	return response
}

func FetchAllMerchantFlowsFromDB(merchantId string) model.Flow {
	fmt.Println("fetching from db")
	dbConnection := db.GetDB()
	var flow model.Flow
	if dbConnection == nil {
		fmt.Print("failed to connect")
	} else {
		flow := &model.Flow{MerchantId: uuid.FromStringOrNil(merchantId)}
		dbConnection.Debug().First(&flow)
	}
	fmt.Print(flow.Id, flow.Name)
	return flow
}

func ParseMerchantFlowsJsonValues(flow model.Flow) (response_dto.FlowResponseDto, error) {
	//Todo: write logic to parse entire json
	var flowResponseDto response_dto.FlowResponseDto
	flowResponseDto.Name = flow.Name
	flowResponseDto.Version = flow.Version
	flowResponseDto.Type = flow.Type
	flowResponseDto.Status = flow.Status
	var moduleVersionResponses []response_dto.ModuleVersionResponseDto
	var moduleList []uint8
	dbConnection := db.GetDB()
	moduleList = *flow.ModuleVersions
	for i := 0; i < len(moduleList); i++ {
		var moduleVersionResponseDto response_dto.ModuleVersionResponseDto
		module := &model.Module{Id: int(moduleList[i])}
		dbConnection.Find(&module)
		moduleVersionResponseDto.Name = module.Name
		moduleVersionResponseDto.Status = module.Status
		moduleVersion := &model.ModuleVersion{ModuleId: int(moduleList[i])}
		dbConnection.Find(&moduleVersion)
		moduleVersionResponseDto.Version = moduleVersion.Version
		moduleVersionResponseDto.ExternalId = moduleVersion.ExternalId
		moduleVersionResponseDto.Properties = moduleVersion.Properties

		//Todo: Find sections
		var sectionList []int
		err := json.Unmarshal([]byte(moduleVersion.SectionVersions), sectionList)
		if err != nil {

		}
		for j := 0; j < len(sectionList); j++ {
			var sectionVersionResponseDto response_dto.SectionVersionsResponseDto
			section := &model.Section{Id: sectionList[i]}
			dbConnection.Find(&section)
			sectionVersionResponseDto.Name = section.Name
			sectionVersionResponseDto.Status = section.Status
			sectionVersion := &model.SectionVersion{SectionId: sectionList[i]}
			dbConnection.Find(&sectionVersion)
			sectionVersionResponseDto.ExternalId = sectionVersion.ExternalId
			sectionVersionResponseDto.Version = sectionVersion.Version

			//Todo: Find Fields
			for j := 0; j < len(sectionList); j++ {
				var fieldVersionsResponseDto response_dto.FieldVersionsResponseDto
				field := &model.Field{SectionId: sectionList[i]}
				dbConnection.Find(&field)
				fieldVersionsResponseDto.Name = field.Name
				fieldVersionsResponseDto.Status = field.Status
				fieldVersion := &model.FieldVersion{FieldId: field.Id}
				dbConnection.Find(&fieldVersion)
				fieldVersionsResponseDto.ExternalId = fieldVersion.ExternalId
				fieldVersionsResponseDto.Version = fieldVersion.Version

				sectionVersionResponseDto.Fields = append(sectionVersionResponseDto.Fields, fieldVersionsResponseDto)
			}

			moduleVersionResponseDto.Sections = append(moduleVersionResponseDto.Sections, sectionVersionResponseDto)
		}

		moduleVersionResponses = append(moduleVersionResponses, moduleVersionResponseDto)
	}

	flowResponseDto.Modules = moduleVersionResponses

	return flowResponseDto, nil
}
