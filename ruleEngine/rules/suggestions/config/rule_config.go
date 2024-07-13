package config

import (
	"car-comparison-service/logger"
	"car-comparison-service/ruleEngine"
	"car-comparison-service/ruleEngine/rules/suggestions/filter"
	"car-comparison-service/ruleEngine/rules/suggestions/priority"
	"encoding/json"
	"os"
	"sort"
)

const (
	RuleSuggestionsConfigFile = "ruleEngine/rules/suggestions/config/suggestions_config.json"
)

type Rule struct {
	RuleId         *string `json:"ruleId,omitempty"`
	IsEnabled      *bool   `json:"isEnabled,omitempty"`
	SequenceNumber *int32  `json:"sequenceNumber,omitempty"`
}

type RuleType string

const (
	FILTER         RuleType = "Filter"
	PRIORITY_RULES RuleType = "PriorityRules"
)

type RuleDetails struct {
	RuleFunc ruleEngine.IRule
	Type     RuleType
}

var RuleIdToRuleMap = func() map[string]RuleDetails {
	m := make(map[string]RuleDetails)

	m["vehicle_type_filter"] = RuleDetails{RuleFunc: filter.VehicleTypeFilter(), Type: FILTER}
	m["vehicle_brand_filter"] = RuleDetails{RuleFunc: filter.VehicleBrandFilter(), Type: FILTER}
	m["vehicle_fuel_type_filter"] = RuleDetails{RuleFunc: filter.VehicleFuelTypeFilter(), Type: FILTER}
	m["price_priority"] = RuleDetails{RuleFunc: priority.PricePriority(), Type: PRIORITY_RULES}
	m["manufacturing_year_priority"] = RuleDetails{RuleFunc: priority.ManufacturingYearPriority(), Type: PRIORITY_RULES}
	return m
}()

func loadRuleConfig(fileName string) []Rule {
	config, err := os.ReadFile(fileName)
	if err != nil {
		logger.Log.Error("Error reading rule config file", err)
	}

	var defaultConfig []Rule
	err = json.Unmarshal(config, &defaultConfig)

	if err != nil {
		logger.Log.Error("Error reading rule config file", err)
	}
	return defaultConfig
}

type RuleConfig struct {
	FilterRules   []Rule
	PriorityRules []Rule
}

func GetRuleConfig(filename string) RuleConfig {
	rules := loadRuleConfig(filename)
	var filterRules []Rule
	var priorityRules []Rule
	for _, rule := range rules {
		if *rule.IsEnabled {
			if _, ok := RuleIdToRuleMap[*rule.RuleId]; ok {
				ruleDetails := RuleIdToRuleMap[*rule.RuleId]
				if ruleDetails.Type == FILTER {
					filterRules = append(filterRules, rule)
				} else {
					priorityRules = append(priorityRules, rule)
				}
			}
		}
	}

	sort.Slice(priorityRules, func(i, j int) bool {
		return *priorityRules[i].SequenceNumber < *priorityRules[j].SequenceNumber
	})

	sort.Slice(filterRules, func(i, j int) bool {
		return *filterRules[i].SequenceNumber < *filterRules[j].SequenceNumber
	})

	Config := RuleConfig{
		FilterRules:   filterRules,
		PriorityRules: priorityRules,
	}

	return Config
}
