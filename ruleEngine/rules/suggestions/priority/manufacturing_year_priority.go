package priority

import (
	"car-comparison-service/db/model"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules"
	"car-comparison-service/utils"
	"context"
	"sort"
)

func manufacturingYearPriority(ctx context.Context, re *ruleEngine.RuleEngineExecutor) error {
	vehicleSuggestionPtr, err := ruleEngine.GetCacheValueHelper[[]model.VehicleSuggestionResult](re, rules.VehicleSuggestions)
	if err != nil {
		return err
	}

	result := utils.GetValFromPtr(vehicleSuggestionPtr)
	var manufacturingYearList []int
	manufacturingYearToRankMap := make(map[int]int64)
	for index := range result {
		if _, ok := manufacturingYearToRankMap[result[index].ManufacturingYear]; !ok {
			manufacturingYearList = append(manufacturingYearList, result[index].ManufacturingYear)
			manufacturingYearToRankMap[result[index].ManufacturingYear] = -1
		}
	}
	manufacturingYearList = utils.GetUniqueValuesFromArray(manufacturingYearList)

	// Sort in descending order
	sort.Slice(manufacturingYearList, func(i, j int) bool {
		return manufacturingYearList[i] > manufacturingYearList[j]
	})

	var currRank int64 = 1
	for index := range manufacturingYearList {
		manufacturingYearToRankMap[manufacturingYearList[index]] = currRank
		currRank = currRank + 1
	}

	var multiplicationFactor int64 = 10
	for currRank > multiplicationFactor {
		multiplicationFactor = multiplicationFactor * 10
	}

	for index := range result {
		priority := int64(99999)
		if val, ok := manufacturingYearToRankMap[result[index].ManufacturingYear]; ok {
			priority = val
		}
		result[index].Rank = result[index].Rank*multiplicationFactor + priority
	}

	return nil
}

func ManufacturingYearPriority() *ruleEngine.Rule {
	rule := &ruleEngine.Rule{
		RuleId:      "manufacturing_year_priority",
		ExecuteRule: manufacturingYearPriority,
	}
	return rule
}
