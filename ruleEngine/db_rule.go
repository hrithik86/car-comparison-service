package ruleEngine

import (
	"car-comparison-service/logger"
	"car-comparison-service/orm"
	"context"
	"fmt"
)

type ConditionType string

const orCondition ConditionType = "OR"
const andCondition ConditionType = "AND"

type DbRuleExecuteFunc func(ctx context.Context, qe *orm.QueryEngine, re *RuleEngineExecutor) (*orm.QueryEngine, error)

type Condition struct {
	conditionFunc RulePreconditionFunc
	conditionType ConditionType
}

type DbRule struct {
	IRule
	ruleId        string
	preconditions []Condition
	preload       RulePreloadFunc
	tasks         []DbRuleExecuteFunc
}

func (rule *DbRule) GetRuleId() string {
	return rule.ruleId
}

func (rule *DbRule) Precondition(ctx context.Context, re *RuleEngineExecutor) bool {
	result := true
	for i, precondition := range rule.preconditions {
		if i == 0 {
			result = precondition.conditionFunc(ctx, re)
			continue
		}
		if precondition.conditionType == orCondition {
			result = result || precondition.conditionFunc(ctx, re)
		}
		if precondition.conditionType == andCondition {
			result = result && precondition.conditionFunc(ctx, re)
		}
	}
	return result
}

func (rule *DbRule) AddOrCondition(conditionFunc RulePreconditionFunc) *DbRule {
	rule.preconditions = append(rule.preconditions, Condition{
		conditionFunc: conditionFunc,
		conditionType: orCondition,
	})
	return rule
}

func (rule *DbRule) AddAndCondition(conditionFunc RulePreconditionFunc) *DbRule {
	rule.preconditions = append(rule.preconditions, Condition{
		conditionFunc: conditionFunc,
		conditionType: andCondition,
	})
	return rule
}

func (rule *DbRule) AddTask(task DbRuleExecuteFunc) *DbRule {
	rule.tasks = append(rule.tasks, task)
	return rule
}

func (rule *DbRule) Preload(ctx context.Context, re *RuleEngineExecutor) error {
	if rule.preload == nil {
		return nil
	}
	return rule.preload(ctx, re)
}

func (rule *DbRule) SetPreload(preloadFunc RulePreloadFunc) *DbRule {
	rule.preload = preloadFunc
	return rule
}

func (rule *DbRule) Execute(ctx context.Context, executor *RuleEngineExecutor) error {
	query, _ := GetCacheValueHelper[orm.QueryEngine](executor, "query")
	db, err := GetDbWithError(executor)
	if err != nil {
		return fmt.Errorf("db Object not set")
	}
	if query == nil {
		query = orm.Query(db)
	}
	for _, task := range rule.tasks {
		query, err = task(ctx, query, executor)
		if err != nil {
			logger.Get(ctx).Errorf("error excecuting rule %s with error %v", rule.GetRuleId(), err.Error())
			return err
		}
	}
	executor.SetValue("query", *query)
	return err
}

func CreateDbRule(ruleId string) *DbRule {
	return &DbRule{
		ruleId: ruleId,
	}
}
