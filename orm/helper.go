package orm

import "gorm.io/gorm/clause"

func Table(tableName string) clause.Table {
	return clause.Table{
		Name: tableName,
	}
}

func TableWithAlias(tableName string, alias string) clause.Table {
	tableClause := Table(tableName)
	tableClause.Alias = alias
	return tableClause
}

func Column(table clause.Table, column string) clause.Column {
	return clause.Column{
		Table: table.Alias,
		Name:  column,
	}
}

func RawColumn(table clause.Table, column string) clause.Column {
	return clause.Column{
		Table: table.Alias,
		Name:  column,
		Raw:   true,
	}
}

func Eq(column interface{}, value interface{}) clause.Expression {
	return clause.Eq{Column: column, Value: value}
}

func In(column interface{}, values []interface{}) clause.Expression {
	return clause.IN{Column: column, Values: values}
}

func Neq(column interface{}, value interface{}) clause.Expression {
	return clause.Neq{Column: column, Value: value}
}
