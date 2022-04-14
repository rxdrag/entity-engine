package dialect

import (
	"rxdrag.com/entity-engine/model"
	"rxdrag.com/entity-engine/model/data"
	"rxdrag.com/entity-engine/model/graph"
	"rxdrag.com/entity-engine/model/table"
)

const (
	MySQL = "mysql"
)

type SQLBuilder interface {
	BuildBoolExp(where map[string]interface{}) (string, []interface{})
	BuildFieldExp(fieldName string, fieldArgs map[string]interface{}) (string, []interface{})

	BuildCreateTableSQL(table *table.Table) string
	BuildDeleteTableSQL(table *table.Table) string
	BuildColumnSQL(column *table.Column) string
	BuildModifyTableAtoms(diff *model.TableDiff) []model.ModifyAtom
	ColumnTypeSQL(column *table.Column) string

	BuildQuerySQL(entity graph.Node, args map[string]interface{}) (string, []interface{})

	BuildInsertSQL(fields []*data.Field, table *table.Table) (string, []interface{})
	BuildUpdateSQL(id uint64, fields []*data.Field, table *table.Table) (string, []interface{})
}

func GetSQLBuilder() SQLBuilder {
	var builder MySQLBuilder
	return &builder
}
