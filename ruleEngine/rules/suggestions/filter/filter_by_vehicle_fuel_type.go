package filter

import (
	"car-comparison-service/db/model"
	"car-comparison-service/logger"
	"car-comparison-service/orm"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules"
	"car-comparison-service/utils"
	"context"
	"gorm.io/gorm/clause"
)

func filterByVehicleFuelType(ctx context.Context, qe *orm.QueryEngine, re *ruleEngine.RuleEngineExecutor) (*orm.QueryEngine, error) {
	var vehicleTable *clause.Table
	vehicleTable, err := ruleEngine.GetCacheValueHelper[clause.Table](re, rules.VehicleTable)
	if err != nil {
		logger.Log.Error(ctx, err, "Rule Id - filterByVehicleFuelType, error fetching VehicleTable")
		return nil, err
	}

	result, err := ruleEngine.GetCacheValueHelper[[]model.VehicleSuggestionResult](re, rules.VehicleSuggestions)
	if err != nil {
		return nil, err
	}
	vehicleSuggestionIds := make([]string, 0)
	for _, suggestedVehicle := range *result {
		vehicleSuggestionIds = append(vehicleSuggestionIds, suggestedVehicle.Id.String())
	}

	fuelTypeValue, err := ruleEngine.GetCacheValueHelper[*model.FuelType](re, rules.VehicleFuelTypeVariable)
	if err != nil {
		logger.Log.Error(ctx, err, "Rule Id - filterByVehicleFuelType, error fetching fuelTypeValue")
		return qe, err
	}

	qe.Where(
		orm.Eq(orm.Column(*vehicleTable, "fuel_type"), *fuelTypeValue),
	)

	if len(vehicleSuggestionIds) > 0 {
		qe.Where(orm.In(orm.Column(*vehicleTable, "id"), utils.TypeCastToInterfaceSlice(vehicleSuggestionIds)))
	}
	return qe, nil
}

func VehicleFuelTypeFilter() *ruleEngine.DbRule {
	return ruleEngine.CreateDbRule("vehicle_fuel_type_filter").
		AddTask(filterByVehicleFuelType)
}
