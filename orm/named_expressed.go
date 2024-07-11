package orm

import "gorm.io/gorm/clause"

type namedExpression struct {
	argExp  map[string]interface{}
	argList []interface{}
	sql     string
}

func (ne *namedExpression) AddParamArgs(value interface{}) *namedExpression {
	ne.argList = append(ne.argList, value)
	return ne
}

func (ne *namedExpression) ToExpression() clause.Expr {
	return clause.Expr{
		SQL:  ne.sql,
		Vars: ne.argList,
	}
}

func ExpressionBuilder(sql string) *namedExpression {
	return &namedExpression{
		sql:     sql,
		argExp:  map[string]interface{}{},
		argList: []interface{}{},
	}
}
