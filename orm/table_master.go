package orm

import (
	"fmt"
	"gorm.io/gorm/clause"
)

type tableMaster struct {
	tableMapping map[string]map[string]clause.Table
	countMapping map[string]int
}

func (tm *tableMaster) addMapping(table clause.Table) {
	if tm.getTableName(table) != nil && len(tm.getTableName(table).Name) > 0 {
		return
	}
	_, ok := tm.tableMapping[table.Name]
	if !ok {
		tm.tableMapping[table.Name] = map[string]clause.Table{}
	}
	tm.countMapping[table.Name] += 1
	tm.tableMapping[table.Name][table.Alias] = clause.Table{
		Name:  table.Name,
		Alias: fmt.Sprintf("%s%d", table.Name, tm.countMapping[table.Name]),
	}
}

func (tm *tableMaster) getTableName(table clause.Table) *clause.Table {
	tableMap, ok := tm.tableMapping[table.Name]
	if !ok {
		return nil
	}
	queryAlias, ok := tableMap[table.Alias]
	if !ok {
		return nil
	}
	return &queryAlias
}
