package ruleEngine

import "context"

type Rule struct {
	IRule
	RuleId      string
	ExecuteRule func(ctx context.Context, re *RuleEngineExecutor) error
}

func (rule *Rule) GetRuleId() string {
	return rule.RuleId
}

func (rule *Rule) Execute(ctx context.Context, re *RuleEngineExecutor) error {
	return rule.ExecuteRule(ctx, re)
}

var _ IRule = (*Rule)(nil)
