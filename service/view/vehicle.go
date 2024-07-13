package view

import (
	"car-comparison-service/constants"
	"car-comparison-service/db/model"
	"car-comparison-service/service/api/response"
	"car-comparison-service/utils"
	"github.com/google/uuid"
)

func CreateVehicleSearchResponse(vehicles []*model.VehicleWithAttachmentInformation) []response.VehicleResponse {
	vehicleIdToAttachmentsMap := make(map[uuid.UUID][]response.VehicleAttachment)
	vehicleIdToDetailsMap := make(map[uuid.UUID]response.VehicleResponse)
	for _, vehicle := range vehicles {
		if vehicle.AttachmentId != nil {
			vehicleIdToAttachmentsMap[utils.GetValFromPtr(vehicle.Id)] = append(vehicleIdToAttachmentsMap[utils.GetValFromPtr(vehicle.Id)],
				response.VehicleAttachment{
					Id:        vehicle.AttachmentId,
					Path:      vehicle.Path,
					MediaType: vehicle.MediaType,
				})
		}
		vehicleIdToDetailsMap[utils.GetValFromPtr(vehicle.Id)] = response.VehicleResponse{
			Id:                vehicle.Id,
			Brand:             vehicle.Brand,
			Model:             vehicle.Model,
			ManufacturingYear: vehicle.ManufacturingYear,
			Price:             vehicle.Price,
			Type:              utils.NewPtr(string(utils.GetValFromPtr(vehicle.Type))),
			Attachments:       nil,
		}
	}

	responseList := make([]response.VehicleResponse, 0, 1)
	for vehicleId, vehicleDetails := range vehicleIdToDetailsMap {
		if _, ok := vehicleIdToAttachmentsMap[vehicleId]; ok {
			vehicleDetails.Attachments = vehicleIdToAttachmentsMap[vehicleId]
		}
		responseList = append(responseList, vehicleDetails)
	}

	return responseList
}

func CreateVehicleComparisonResponse(values map[string][]interface{}) response.VehicleComparisonResponse {
	ids := values[constants.VehicleId]
	delete(values, constants.VehicleId)
	return response.VehicleComparisonResponse{
		Ids:             ids,
		ComparisonTable: values,
	}
}
