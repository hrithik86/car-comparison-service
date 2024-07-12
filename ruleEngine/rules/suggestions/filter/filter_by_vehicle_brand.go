package filter

import (
	"car-comparison-service/db/model"
	"car-comparison-service/logger"
	"car-comparison-service/orm"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules"
	"context"
	"gorm.io/gorm/clause"
)

func filterByVehicleBrand(ctx context.Context, qe *orm.QueryEngine, re *ruleEngine.RuleEngineExecutor) (*orm.QueryEngine, error) {
	var vehicleTable *clause.Table
	vehicleTable, err := ruleEngine.GetCacheValueHelper[clause.Table](re, rules.VehicleTable)
	if err != nil {
		logger.Log.Error(ctx, err, "Rule Id - filterByVehicleBrand, error fetching VehicleTable")
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

	brandValue, err := ruleEngine.GetCacheValueHelper[*string](re, rules.BrandVariable)
	if err != nil {
		logger.Log.Error(ctx, err, "Rule Id - filterByVehicleBrand, error fetching brandValue")
		return qe, err
	}

	qe.Where(
		orm.Eq(orm.Column(*vehicleTable, "brand"), brandValue),
	)

	if len(vehicleSuggestionIds) > 0 {
		qe.Where(orm.In(orm.Column(*vehicleTable, "id"), ruleEngine.TypeCastToInterfaceSlice(vehicleSuggestionIds)))
	}
	return qe, nil
}

func VehicleBrandFilter() *ruleEngine.DbRule {
	return ruleEngine.CreateDbRule("vehicle_brand_filter").
		AddTask(filterByVehicleBrand)
}
