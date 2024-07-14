package view

import (
	"car-comparison-service/constants"
	"car-comparison-service/db/model"
	"car-comparison-service/service/api/response"
	"car-comparison-service/utils"
	"github.com/google/uuid"
)

func CreateVehicleSearchResponse(vehicles []*model.VehicleWithAttachmentInformation) []response.VehicleWithAttachmentsResponse {
	vehicleIdToAttachmentsMap := make(map[uuid.UUID][]response.VehicleAttachment)
	vehicleIdToDetailsMap := make(map[uuid.UUID]response.VehicleWithAttachmentsResponse)
	for _, vehicle := range vehicles {
		if vehicle.AttachmentId != nil {
			vehicleIdToAttachmentsMap[utils.GetValFromPtr(vehicle.Id)] = append(vehicleIdToAttachmentsMap[utils.GetValFromPtr(vehicle.Id)],
				response.VehicleAttachment{
					Id:        vehicle.AttachmentId,
					Path:      vehicle.Path,
					MediaType: vehicle.MediaType,
				})
		}

		vehicleIdToDetailsMap[utils.GetValFromPtr(vehicle.Id)] = response.VehicleWithAttachmentsResponse{
			VehicleBaseResponse: response.VehicleBaseResponse{
				Id:                vehicle.Id,
				Brand:             vehicle.Brand,
				Model:             vehicle.Model,
				ManufacturingYear: vehicle.ManufacturingYear,
				Price:             vehicle.Price,
				MileageType:       vehicle.Mileage,
				FuelType:          utils.NewPtr(string(utils.GetValFromPtr(vehicle.FuelType))),
				Type:              utils.NewPtr(string(utils.GetValFromPtr(vehicle.Type))),
			},
			Attachments: nil,
		}
	}

	responseList := make([]response.VehicleWithAttachmentsResponse, 0, 1)
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

func CreateVehicleFeaturesResponse(vehicles []*model.VehicleWithFeatures) response.VehicleWithFeaturesResponse {
	vehicleBaseInfo := vehicles[0]
	featuresList := make([]response.VehicleFeatures, 0, 1)
	for _, vehicle := range vehicles {
		if vehicle.FeatureId != nil {
			featuresList = append(featuresList, response.VehicleFeatures{
				FeatureId: vehicle.FeatureId,
				Key:       vehicle.Key,
				Value:     vehicle.Value,
			})
		}
	}

	return response.VehicleWithFeaturesResponse{
		VehicleBaseResponse: response.VehicleBaseResponse{
			Id:                vehicleBaseInfo.Id,
			Brand:             vehicleBaseInfo.Brand,
			Model:             vehicleBaseInfo.Model,
			ManufacturingYear: vehicleBaseInfo.ManufacturingYear,
			Price:             vehicleBaseInfo.Price,
			Type:              utils.NewPtr(string(utils.GetValFromPtr(vehicleBaseInfo.Type))),
			FuelType:          utils.NewPtr(string(utils.GetValFromPtr(vehicleBaseInfo.FuelType))),
			MileageType:       vehicleBaseInfo.Mileage,
		},
		Features: featuresList,
	}
}
