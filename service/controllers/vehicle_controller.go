package controllers

import (
	"car-comparison-service/appcontext"
	vehicleCache "car-comparison-service/cache_manager/vehicle"
	"car-comparison-service/constants"
	"car-comparison-service/contexts"
	"car-comparison-service/db/model"
	"car-comparison-service/db/repository"
	"car-comparison-service/logger"
	"car-comparison-service/service/api/request"
	"car-comparison-service/utils"
	"context"
	"github.com/google/uuid"
	"sync"
)

type Vehicle struct {
	DbClient           repository.IVehicle
	VehicleSuggestions vehicleCache.IVehicleSuggestions
}

type IVehicle interface {
	GetVehiclesByModelName(ctx context.Context, modelName string) ([]*model.VehicleWithAttachmentInformation, error)
	GetVehicleInfoById(ctx context.Context, id uuid.UUID) ([]*model.VehicleWithFeatures, error)
	GetVehicleSuggestions(ctx context.Context, id uuid.UUID) ([]model.VehicleSuggestionResult, error)
	GetVehicleComparison(ctx context.Context, req request.VehicleComparisonRequest) (map[string][]interface{}, error)
	CreateVehicle(ctx context.Context, req request.CreateVehicleRequest) (*model.Vehicle, error)
	AddVehicleAttachments(ctx context.Context, vehicleId uuid.UUID,
		vehicleAttachmentsRequest []request.BulkAddVehicleAttachmentsRequest) ([]*model.VehicleAttachment, error)
	AddVehicleFeatures(ctx context.Context, vehicleId uuid.UUID,
		vehicleFeaturesRequest []request.BulkAddVehicleFeaturesRequest) ([]*model.VehicleFeatures, error)
}

var (
	VehicleControllerDoOnce sync.Once
	VehicleController       Vehicle
)

func InitializeVehicleController() Vehicle {
	VehicleControllerDoOnce.Do(func() {
		VehicleController = Vehicle{
			DbClient:           appcontext.GetDbClient(),
			VehicleSuggestions: vehicleCache.NewVehicleSuggestionsManager(),
		}
	})
	return VehicleController
}

func (vc Vehicle) GetVehiclesByModelName(ctx context.Context, modelName string) ([]*model.VehicleWithAttachmentInformation, error) {
	vehicles, err := vc.DbClient.GetVehiclesByModel(ctx, modelName)
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (vc Vehicle) GetVehicleInfoById(ctx context.Context, id uuid.UUID) ([]*model.VehicleWithFeatures, error) {
	vehicle, err := vc.DbClient.GetVehicleWithFeaturesById(ctx, id)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (vc Vehicle) GetVehicleSuggestions(ctx context.Context, id uuid.UUID) ([]model.VehicleSuggestionResult, error) {
	vehicle, err := vc.DbClient.GetVehicleInfoById(ctx, id)
	if err != nil {
		return nil, err
	}

	key := vc.VehicleSuggestions.CreateKey(ctx, id.String())
	cachedSuggestions, err := vc.VehicleSuggestions.GetVehicleSuggestionsDetails(ctx, key)
	if err != nil {
		logger.Get(ctx).Errorf("Error in fetching cached suggestions for id: %s, err: %v", id, err.Error())

		suggestionsFactory := SuggestionsFactory{}
		suggestionsControllerObj := suggestionsFactory.GetSuggestionsController(vehicle.Type)
		suggestedVehicles, err := suggestionsControllerObj.ExecuteRules(ctx, repository.DbClient().DB, vehicle)
		if err != nil {
			return nil, err
		}

		newCtx := contexts.Copy(ctx)
		go vc.cacheSuggestedVehicles(newCtx, key, suggestedVehicles)
		return suggestedVehicles, nil
	}
	return cachedSuggestions, nil
}

func (vc Vehicle) cacheSuggestedVehicles(ctx context.Context, key string, suggestedVehicles []model.VehicleSuggestionResult) {
	if err := vc.VehicleSuggestions.SetVehicleSuggestionsDetails(ctx, key, suggestedVehicles); err != nil {
		logger.Get(ctx).Errorf("Error in caching suggestions for key: %s, err: %v", key, err.Error())
	}
}

func (vc Vehicle) GetVehicleComparison(ctx context.Context, req request.VehicleComparisonRequest) (map[string][]interface{}, error) {
	vehicles, err := vc.DbClient.GetVehiclesByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	keyFeaturesMap := vc.getVehicleKeyFeaturesMap()
	vehicleFeatureToValuesMap := make(map[string][]interface{})
	for _, vehicle := range vehicles {
		vehicleMap := utils.StructToMap(vehicle)
		for featureName, value := range vehicleMap {
			if _, ok := keyFeaturesMap[featureName]; ok {
				vehicleFeatureToValuesMap[featureName] = append(vehicleFeatureToValuesMap[featureName], value)
			}
		}
	}

	if req.HideCommonFeatures {
		vehicleFeatureToValuesMap = vc.filterCommonFeatures(vehicleFeatureToValuesMap)
	}

	return vehicleFeatureToValuesMap, nil
}

func (vc Vehicle) CreateVehicle(ctx context.Context, req request.CreateVehicleRequest) (*model.Vehicle, error) {
	vehicle, err := vc.DbClient.CreateVehicle(ctx, &model.Vehicle{
		DbId:              model.DbId{Id: utils.NewPtr(uuid.New())},
		Model:             req.Model,
		Brand:             req.Brand,
		ManufacturingYear: req.ManufacturingYear,
		Type:              req.Type,
		Price:             req.Price,
		FuelType:          req.FuelType,
		Mileage:           req.Mileage,
	})
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (vc Vehicle) filterCommonFeatures(vehicleFeatureToValuesMap map[string][]interface{}) map[string][]interface{} {
	filteredVehicleFeatureToValuesMap := make(map[string][]interface{})
	for featureName, values := range vehicleFeatureToValuesMap {
		if utils.ContainsSameValues(values) {
			continue
		}
		filteredVehicleFeatureToValuesMap[featureName] = values
	}
	return filteredVehicleFeatureToValuesMap
}

func (vc Vehicle) getVehicleKeyFeaturesMap() map[string]bool {
	keyFeaturesMap := make(map[string]bool)
	keyFeatures := constants.VehicleKeyFeatures
	for _, feature := range keyFeatures {
		keyFeaturesMap[feature] = true
	}
	return keyFeaturesMap
}

func (vc Vehicle) AddVehicleAttachments(ctx context.Context, vehicleId uuid.UUID,
	vehicleAttachmentsRequest []request.BulkAddVehicleAttachmentsRequest) ([]*model.VehicleAttachment, error) {
	vehicleAttachments := make([]*model.VehicleAttachment, 0, 1)
	for _, vehicleAttachmentRequest := range vehicleAttachmentsRequest {
		vehicleAttachments = append(vehicleAttachments, &model.VehicleAttachment{
			DbId:      model.DbId{Id: utils.NewPtr(uuid.New())},
			Name:      vehicleAttachmentRequest.Name,
			Path:      vehicleAttachmentRequest.Path,
			MediaType: vehicleAttachmentRequest.MediaType,
			VehicleId: utils.NewPtr(vehicleId),
		})
	}
	resp, err := vc.DbClient.BulkAddVehicleAttachments(ctx, vehicleAttachments)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (vc Vehicle) AddVehicleFeatures(ctx context.Context, vehicleId uuid.UUID,
	vehicleFeaturesRequest []request.BulkAddVehicleFeaturesRequest) ([]*model.VehicleFeatures, error) {
	vehicleFeatures := make([]*model.VehicleFeatures, 0, 1)
	for _, vehicleFeatureRequest := range vehicleFeaturesRequest {
		vehicleFeatures = append(vehicleFeatures, &model.VehicleFeatures{
			DbId:      model.DbId{Id: utils.NewPtr(uuid.New())},
			Key:       vehicleFeatureRequest.Key,
			Value:     vehicleFeatureRequest.Value,
			VehicleId: utils.NewPtr(vehicleId),
		})
	}
	resp, err := vc.DbClient.BulkAddVehicleFeatures(ctx, vehicleFeatures)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
