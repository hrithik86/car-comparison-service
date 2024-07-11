package suggestions

import (
	"car-comparison-service/db/model"
	"car-comparison-service/orm"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules"
	"context"
)

func initSelect(ctx context.Context, qe *orm.QueryEngine, re *ruleEngine.RuleEngineExecutor) (*orm.QueryEngine, error) {
	vehicleTable := qe.GetTable(orm.TableWithAlias(model.TableNameVehicle, "v1"))

	re.SetValue(rules.VehicleTable, vehicleTable)
	qe.Select(orm.RawColumn(vehicleTable, "*")).From(vehicleTable)
	return qe, nil
}

func InitSelectQuery() *ruleEngine.DbRule {
	return ruleEngine.CreateDbRule("vehicle_suggestions_init").
		AddTask(initSelect)
}
