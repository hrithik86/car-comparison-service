package orm

import (
	"car-comparison-service/db/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type QueryEngine struct {
	db *gorm.DB
	tableMaster
	groupColumnMap map[clause.Column]bool
	FromClause     *clause.From
	WhereClause    *clause.Where
	SelectClause   *clause.Select
}

func (qe *QueryEngine) GetTable(table clause.Table) clause.Table {
	qe.tableMaster.addMapping(table)
	tableName := qe.tableMaster.getTableName(table)
	return *tableName
}

func (qe *QueryEngine) getFromClause() *clause.From {
	if qe.FromClause == nil {
		qe.FromClause = &clause.From{}
	}
	return qe.FromClause
}

func (qe *QueryEngine) getSelectClause() *clause.Select {
	if qe.SelectClause == nil {
		qe.SelectClause = &clause.Select{}
	}
	return qe.SelectClause
}

func (qe *QueryEngine) getWhereClause() *clause.Where {
	if qe.WhereClause == nil {
		qe.WhereClause = &clause.Where{}
	}
	return qe.WhereClause
}

func (qe *QueryEngine) setCommaExpression(clauseInterface interface{}, exprs ...clause.Expression) {
	var commaExpression clause.CommaExpression
	switch clauseInterface.(type) {
	case *clause.Select:
		if qe.getSelectClause().Expression == nil {
			qe.getSelectClause().Expression = clause.CommaExpression{}
		}
		commaExpression = qe.getSelectClause().Expression.(clause.CommaExpression)
		for _, exp := range exprs {
			commaExpression.Exprs = append(commaExpression.Exprs, exp)
		}
		qe.getSelectClause().Expression = commaExpression
	}
}

func (qe *QueryEngine) From(tables ...clause.Table) *QueryEngine {
	for _, table := range tables {
		qe.getFromClause().Tables = append(qe.getFromClause().Tables, table)
	}
	return qe
}

func (qe *QueryEngine) Where(exprs ...clause.Expression) *QueryEngine {
	for _, exp := range exprs {
		qe.getWhereClause().Exprs = append(qe.getWhereClause().Exprs, exp)
	}
	return qe
}

func (qe *QueryEngine) Select(columns ...clause.Column) *QueryEngine {
	var columnExpr []clause.Expression
	for _, col := range columns {
		columnExpr = append(columnExpr, ExpressionBuilder("?").
			AddParamArgs(col).
			ToExpression(),
		)
	}
	qe.setCommaExpression(qe.getSelectClause(), columnExpr...)
	return qe
}

func (qe *QueryEngine) BuildQuery() []clause.Expression {
	var expr []clause.Expression
	val := reflect.Value{}
	val = reflect.ValueOf(qe).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Type().Field(i).IsExported() || val.Field(i).IsNil() {
			continue
		}
		field := val.Field(i).Interface()
		if clauseObject, ok := field.(clause.Expression); ok {
			expr = append(expr, clauseObject)
		}
	}
	return expr
}

func (qe *QueryEngine) GetQuery() *gorm.DB {
	return qe.db.Clauses(qe.BuildQuery()...)
}

func (qe *QueryEngine) Execute(result interface{}) error {
	res := qe.GetQuery().Scan(result)
	err := utils.ValidateResultSuccess(res)
	return err
}

func Query(db *gorm.DB) *QueryEngine {
	qe := &QueryEngine{
		db: db,
		tableMaster: tableMaster{
			countMapping: map[string]int{},
			tableMapping: map[string]map[string]clause.Table{},
		},
		groupColumnMap: map[clause.Column]bool{},
	}
	return qe
}
