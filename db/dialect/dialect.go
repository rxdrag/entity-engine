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

	BuildQuerySQLBody(argEntity *graph.ArgEntity, fields []*graph.Attribute) string
	BuildWhereSQL(argEntity *graph.ArgEntity, fields []*graph.Attribute, where map[string]interface{}) (string, []interface{})
	BuildOrderBySQL(argEntity *graph.ArgEntity, fields []*graph.Attribute, orderBy map[string]interface{}) string
	BuildQuerySQL(tableName string, fields []*graph.Attribute, args map[string]interface{}) (string, []interface{})

	BuildInsertSQL(fields []*data.Field, table *table.Table) string
	BuildUpdateSQL(id uint64, fields []*data.Field, table *table.Table) string

	BuildQueryByIdsSQL(entity *graph.Entity, idCounts int) string
	BuildClearAssociationSQL(ownerId uint64, tableName string, ownerFieldName string) string
	BuildQueryAssociatedInstancesSQL(node graph.Noder,
		ownerId uint64,
		povitTableName string,
		ownerFieldName string,
		typeFieldName string,
	) string
	BuildBatchAssociationSQL(
		tableName string,
		fields []*graph.Attribute,
		ids []uint64,
		povitTableName string,
		ownerFieldName string,
		typeFieldName string,
	) string
	BuildDeleteSQL(id uint64, tableName string) string
	BuildSoftDeleteSQL(id uint64, tableName string) string

	BuildQueryPovitSQL(povit *data.AssociationPovit) string
	BuildInsertPovitSQL(povit *data.AssociationPovit) string
	BuildDeletePovitSQL(povit *data.AssociationPovit) string
}

func GetSQLBuilder() SQLBuilder {
	var builder MySQLBuilder
	return &builder
}
