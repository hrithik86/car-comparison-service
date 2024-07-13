package priority

import (
	"car-comparison-service/db/model"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules"
	"car-comparison-service/utils"
	"context"
	"sort"
)

func pricePriority(ctx context.Context, re *ruleEngine.RuleEngineExecutor) error {
	vehicleSuggestionPtr, err := ruleEngine.GetCacheValueHelper[[]model.VehicleSuggestionResult](re, rules.VehicleSuggestions)
	if err != nil {
		return err
	}

	result := utils.GetValFromPtr(vehicleSuggestionPtr)
	var priceList []int64
	priceToRankMap := make(map[int64]int64)
	for index := range result {
		if _, ok := priceToRankMap[result[index].Price]; !ok {
			priceList = append(priceList, result[index].Price)
			priceToRankMap[result[index].Price] = -1
		}
	}
	priceList = utils.GetUniqueValuesFromArray(priceList)

	// Sort in ascending order
	sort.Slice(priceList, func(i, j int) bool {
		return priceList[i] < priceList[j]
	})

	var currRank int64 = 1
	for index := range priceList {
		priceToRankMap[priceList[index]] = currRank
		currRank = currRank + 1
	}

	var multiplicationFactor int64 = 10
	for currRank > multiplicationFactor {
		multiplicationFactor = multiplicationFactor * 10
	}

	for index := range result {
		priority := int64(99999)
		if val, ok := priceToRankMap[result[index].Price]; ok {
			priority = val
		}
		result[index].Rank = result[index].Rank*multiplicationFactor + priority
	}

	return nil
}

func PricePriority() *ruleEngine.Rule {
	rule := &ruleEngine.Rule{
		RuleId:      "price_priority",
		ExecuteRule: pricePriority,
	}
	return rule
}
