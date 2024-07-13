package suggestions

import (
	"car-comparison-service/db/model"
	"car-comparison-service/logger"
	"car-comparison-service/orm"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules"
	"context"
	"github.com/google/uuid"
)

func initSelect(ctx context.Context, qe *orm.QueryEngine, re *ruleEngine.RuleEngineExecutor) (*orm.QueryEngine, error) {
	vehicleTable := qe.GetTable(orm.TableWithAlias(model.TableNameVehicle, "v1"))

	re.SetValue(rules.VehicleTable, vehicleTable)
	qe.Select(orm.RawColumn(vehicleTable, "*")).From(vehicleTable)

	vehicleId, err := ruleEngine.GetCacheValueHelper[uuid.UUID](re, rules.VehicleId)
	if err != nil {
		logger.Log.Error(ctx, err, "Rule Id - filterByVehicleBrand, error fetching brandValue")
		return qe, err
	}

	qe.Where(orm.Neq(orm.Column(vehicleTable, "id"), *vehicleId))

	return qe, nil
}

func InitSelectQuery() *ruleEngine.DbRule {
	return ruleEngine.CreateDbRule("vehicle_suggestions_init").
		AddTask(initSelect)
}
