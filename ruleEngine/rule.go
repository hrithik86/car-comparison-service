package ruleEngine

import "context"

type Rule struct {
	IRule
	RuleId      string
	Conditions  []RulePreconditionFunc
	Preloads    []RulePreloadFunc
	ExecuteRule func(ctx context.Context, re *RuleEngineExecutor) error
}

func (rule *Rule) GetRuleId() string {
	return rule.RuleId
}

func (rule *Rule) Precondition(ctx context.Context, re *RuleEngineExecutor) bool {
	result := true
	for _, condition := range rule.Conditions {
		result = result && condition(ctx, re)
	}
	return result
}

func (rule *Rule) Preload(ctx context.Context, re *RuleEngineExecutor) error {
	for _, preload := range rule.Preloads {
		err := preload(ctx, re)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rule *Rule) Execute(ctx context.Context, re *RuleEngineExecutor) error {
	return rule.ExecuteRule(ctx, re)
}

var _ IRule = (*Rule)(nil)
