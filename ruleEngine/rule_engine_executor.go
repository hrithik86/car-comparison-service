package ruleEngine

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-set"
)

type ExecuteFunc func(re *RuleEngineExecutor, result interface{}) error

type IRule interface {
	GetRuleId() string
	Execute(ctx context.Context, executor *RuleEngineExecutor) error
}

type RuleEngineExecutor struct {
	ruleSet     *set.Set[string]
	store       map[string]interface{}
	rules       []IRule
	Ctx         context.Context
	fetchResult ExecuteFunc
}

func (re *RuleEngineExecutor) SetValue(key string, value interface{}) {
	re.store[key] = value
}

func (re *RuleEngineExecutor) HasKey(key string) bool {
	_, ok := re.store[key]
	return ok
}

func (re *RuleEngineExecutor) GetValue(key string) interface{} {
	val, _ := re.store[key]
	return val
}

func (re *RuleEngineExecutor) HasRule(ruleId string) bool {
	return re.ruleSet.Contains(ruleId)
}

func (re *RuleEngineExecutor) AddRule(rule IRule) {
	re.rules = append(re.rules, rule)
}

func (re *RuleEngineExecutor) ClearRules() {
	re.rules = []IRule{}
	re.fetchResult = nil
}

func (re *RuleEngineExecutor) ClearQuery() {
	re.SetValue("query", nil)
}

func (re *RuleEngineExecutor) SetGetterResult(executeFunc ExecuteFunc) {
	re.fetchResult = executeFunc
}

func (re *RuleEngineExecutor) Execute(ctx context.Context, result interface{}) error {
	for _, rule := range re.rules {
		if error := rule.Execute(ctx, re); error != nil {
			return fmt.Errorf("execute failed for %s : %s", rule.GetRuleId(), error.Error())
		}
		re.ruleSet.Insert(rule.GetRuleId())
	}
	if re.fetchResult != nil {
		return re.fetchResult(re, result)
	}
	return nil
}

func Init(ctx context.Context) *RuleEngineExecutor {
	return &RuleEngineExecutor{
		ruleSet: set.New[string](0),
		store:   map[string]interface{}{},
		Ctx:     ctx,
	}
}
