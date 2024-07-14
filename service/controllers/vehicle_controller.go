package controllers

import (
	"car-comparison-service/appcontext"
	vehicleCache "car-comparison-service/cache_manager/vehicle"
	"car-comparison-service/constants"
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
	DbClient repository.IVehicle
}

type IVehicle interface {
	GetVehiclesByModelName(ctx context.Context, modelName string) ([]*model.VehicleWithAttachmentInformation, error)
	GetVehicleInfoById(ctx context.Context, id uuid.UUID) ([]*model.VehicleWithFeatures, error)
	GetVehicleSuggestions(ctx context.Context, id uuid.UUID) ([]model.VehicleSuggestionResult, error)
	GetVehicleComparison(ctx context.Context, req request.VehicleComparisonRequest) (map[string][]interface{}, error)
}

var (
	VehicleControllerDoOnce sync.Once
	VehicleController       Vehicle
)

func InitializeVehicleController() Vehicle {
	VehicleControllerDoOnce.Do(func() {
		VehicleController = Vehicle{
			DbClient: appcontext.GetDbClient(),
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
	vehicle, err := vc.DbClient.GetVehiclesById(ctx, id)
	if err != nil {
		return nil, err
	}

	vehicleSuggestionCache := vehicleCache.CreateSuggestionVehicle(ctx, vehicle.Id)
	cachedSuggestions, err := vehicleSuggestionCache.GetVehicleSuggestionsDetails()
	if err != nil {
		logger.Get(ctx).Errorf("Error in fetching cached suggestions for id: %s, err: %v", id, err.Error())

		suggestionsFactory := SuggestionsFactory{}
		suggestionsControllerObj := suggestionsFactory.GetSuggestionsController(vehicle.Type)
		suggestedVehicles, err := suggestionsControllerObj.ExecuteRules(ctx, repository.DbClient().DB, vehicle)
		if err != nil {
			return nil, err
		}

		if err := vehicleSuggestionCache.SetVehicleSuggestionsDetails(suggestedVehicles); err != nil {
			logger.Get(ctx).Errorf("Error in caching suggestions for id: %s, err: %v", id, err.Error())
		}
		return suggestedVehicles, nil
	}
	return cachedSuggestions, nil
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
