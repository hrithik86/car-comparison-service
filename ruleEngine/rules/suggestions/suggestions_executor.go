package suggestions

import (
	"car-comparison-service/db/model"
	serviceErrors "car-comparison-service/errors"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules"
	"car-comparison-service/ruleEngine/rules/suggestions/config"
	"context"
	"errors"
	"gorm.io/gorm"
	"sort"
)

func ExecuteRules(ctx context.Context, db *gorm.DB, vehicle *model.Vehicle) ([]model.VehicleSuggestionResult, error) {
	ruleConfig := config.GetRuleConfig(config.RuleSuggestionsConfigFile)
	db = db.Session(&gorm.Session{})

	re := ruleEngine.Init(ctx)
	ruleEngine.SetDbForEngine(re, db)

	// Set variables used for suggestions in rule engine
	re.SetValue(rules.BrandVariable, vehicle.Brand)
	re.SetValue(rules.ModelVariable, vehicle.Model)
	re.SetValue(rules.PriceVariable, vehicle.Price)
	re.SetValue(rules.VehicleTypeVariable, vehicle.Type)
	re.SetValue(rules.VehicleSuggestions, make([]model.VehicleSuggestionResult, 0, 1))

	err := applyFilterRules(ctx, ruleConfig.FilterRules, re)
	if err != nil {
		return nil, err
	}

	applyPriorityRules(ruleConfig.PriorityRules, re)
	err = re.Execute(ctx, nil)
	if err != nil {
		return nil, err
	}

	result, err := ruleEngine.GetCacheValueHelper[[]model.VehicleSuggestionResult](re, rules.VehicleSuggestions)
	if err != nil {
		return nil, err
	}
	vehicleSuggestionResults := ruleEngine.GetValFromPtr(result)

	// Sort in ascending order
	sort.Slice(vehicleSuggestionResults, func(i, j int) bool {
		return vehicleSuggestionResults[i].Rank < vehicleSuggestionResults[j].Rank
	})
	return vehicleSuggestionResults, nil
}

func applyFilterRules(ctx context.Context, configRules []config.Rule, re *ruleEngine.RuleEngineExecutor) error {
	for _, rule := range configRules {
		re.AddRule(InitSelectQuery())
		if _, ok := config.RuleIdToRuleMap[*rule.RuleId]; ok {
			re.AddRule(config.RuleIdToRuleMap[*rule.RuleId].RuleFunc)
		}
		re.SetGetterResult(ruleEngine.DbExecuteFunc)
		var result []model.VehicleSuggestionResult
		err := re.Execute(ctx, &result)
		if err != nil {
			if !errors.Is(err, serviceErrors.RECORD_NOT_FOUND) {
				return err
			}
		}
		re.SetValue(rules.VehicleSuggestions, result)
		if len(result) < rules.MaxDesiredVehicleSuggestionCount {
			return nil
		}
		re.ClearRules()
		re.ClearQuery()
	}
	return nil
}

func applyPriorityRules(rules []config.Rule, re *ruleEngine.RuleEngineExecutor) {
	for _, rule := range rules {
		if _, ok := config.RuleIdToRuleMap[*rule.RuleId]; ok {
			re.AddRule(config.RuleIdToRuleMap[*rule.RuleId].RuleFunc)
		}
	}
}
