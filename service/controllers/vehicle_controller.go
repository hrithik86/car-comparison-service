package controllers

import (
	vehicleCache "car-comparison-service/cache_manager/vehicle"
	"car-comparison-service/constants"
	"car-comparison-service/db/model"
	"car-comparison-service/db/repository"
	"car-comparison-service/logger"
	"car-comparison-service/service/api/request"
	"car-comparison-service/utils"
	"context"
	"github.com/google/uuid"
)

type VehicleController struct {
	db repository.CarComparisonServiceDb
}

func NewVehicleController() *VehicleController {
	return &VehicleController{db: repository.DbClient()}
}

func (vc *VehicleController) GetVehiclesByModelName(ctx context.Context, modelName string) ([]*model.VehicleWithAttachmentInformation, error) {
	vehicles, err := vc.db.GetVehiclesByModel(ctx, modelName)
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (vc *VehicleController) GetVehicleById(ctx context.Context, id uuid.UUID) (*model.Vehicle, error) {
	vehicle, err := vc.db.GetVehiclesById(ctx, id)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (vc *VehicleController) GetVehicleSuggestions(ctx context.Context, id uuid.UUID) ([]model.VehicleSuggestionResult, error) {
	vehicle, err := vc.db.GetVehiclesById(ctx, id)
	if err != nil {
		return nil, err
	}

	vehicleSuggestionCache := vehicleCache.CreateSuggestionVehicle(ctx, vehicle.Id)
	cachedSuggestions, err := vehicleSuggestionCache.GetVehicleSuggestionsDetails()
	if err != nil {
		logger.Get(ctx).Errorf("Error in fetching cached suggestions for id: %s, err: %v", id, err.Error())

		suggestionsFactory := SuggestionsFactory{}
		suggestionsControllerObj := suggestionsFactory.GetSuggestionsController(vehicle.Type)
		suggestedVehicles, err := suggestionsControllerObj.ExecuteRules(ctx, vc.db.DB, vehicle)
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

func (vc *VehicleController) GetVehicleComparison(ctx context.Context, req request.VehicleComparisonRequest) (map[string][]interface{}, error) {
	vehicles, err := vc.db.GetVehiclesByIds(ctx, req.Ids)
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

func (vc *VehicleController) filterCommonFeatures(vehicleFeatureToValuesMap map[string][]interface{}) map[string][]interface{} {
	filteredVehicleFeatureToValuesMap := make(map[string][]interface{})
	for featureName, values := range vehicleFeatureToValuesMap {
		if utils.ContainsSameValues(values) {
			continue
		}
		filteredVehicleFeatureToValuesMap[featureName] = values
	}
	return filteredVehicleFeatureToValuesMap
}

func (vc *VehicleController) getVehicleKeyFeaturesMap() map[string]bool {
	keyFeaturesMap := make(map[string]bool)
	keyFeatures := constants.VehicleKeyFeatures
	for _, feature := range keyFeatures {
		keyFeaturesMap[feature] = true
	}
	return keyFeaturesMap
}
