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

func filterByVehicleType(ctx context.Context, qe *orm.QueryEngine, re *ruleEngine.RuleEngineExecutor) (*orm.QueryEngine, error) {
	var vehicleTable *clause.Table
	vehicleTable, err := ruleEngine.GetCacheValueHelper[clause.Table](re, rules.VehicleTable)
	if err != nil {
		logger.Log.Error(ctx, err, "Rule Id - filterByVehicleType, error fetching VehicleTable")
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

	typeValue, err := ruleEngine.GetCacheValueHelper[*model.VehicleType](re, rules.VehicleTypeVariable)
	if err != nil {
		logger.Log.Error(ctx, err, "Rule Id - filterByVehicleType, error fetching typeValue")
		return qe, err
	}

	qe.Where(
		orm.Eq(orm.Column(*vehicleTable, "type"), *typeValue),
	)

	if len(vehicleSuggestionIds) > 0 {
		qe.Where(orm.In(orm.Column(*vehicleTable, "id"), utils.TypeCastToInterfaceSlice(vehicleSuggestionIds)))
	}
	return qe, nil
}

func VehicleTypeFilter() *ruleEngine.DbRule {
	return ruleEngine.CreateDbRule("vehicle_type_filter").
		AddTask(filterByVehicleType)
}
