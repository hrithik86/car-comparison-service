package ruleEngine

import (
	"car-comparison-service/logger"
	"car-comparison-service/orm"
	"context"
	"fmt"
)

type DbRuleExecuteFunc func(ctx context.Context, qe *orm.QueryEngine, re *RuleEngineExecutor) (*orm.QueryEngine, error)

type DbRule struct {
	IRule
	ruleId string
	tasks  []DbRuleExecuteFunc
}

func (rule *DbRule) GetRuleId() string {
	return rule.ruleId
}

func (rule *DbRule) AddTask(task DbRuleExecuteFunc) *DbRule {
	rule.tasks = append(rule.tasks, task)
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
