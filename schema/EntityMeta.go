package schema

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

const (
	BOOLEXP string = "BoolExp"
	ORDERBY string = "OrderBy"
)

const (
	Entity_NORMAL    string = "Normal"
	Entity_ENUM      string = "Enum"
	Entity_INTERFACE string = "Interface"
)

type EntityMeta struct {
	Uuid        string       `json:"uuid"`
	Name        string       `json:"name"`
	TableName   string       `json:"tableName"`
	EntityType  string       `json:"entityType"`
	Columns     []ColumnMeta `json:"columns"`
	Eventable   bool         `json:"eventable"`
	Description string       `json:"description"`
}

func (entity *EntityMeta) createQueryFields() graphql.Fields {
	fields := graphql.Fields{}
	for _, column := range entity.Columns {
		fields[column.Name] = &graphql.Field{
			Type: column.toOutputType(),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fmt.Println(p.Context.Value("data"))
				return "world", nil
			},
		}
	}
	return fields
}

func (entity *EntityMeta) toOutputType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   entity.Name,
			Fields: entity.createQueryFields(),
		},
	)
}

func (entity *EntityMeta) toWhereExp() *graphql.InputObject {
	andExp := graphql.InputObjectFieldConfig{}
	notExp := graphql.InputObjectFieldConfig{}
	orExp := graphql.InputObjectFieldConfig{}

	fields := graphql.InputObjectConfigFieldMap{
		"and": &andExp,
		"not": &notExp,
		"or":  &orExp,
	}

	boolExp := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:   entity.Name + BOOLEXP,
			Fields: fields,
		},
	)
	andExp.Type = boolExp
	notExp.Type = boolExp
	orExp.Type = boolExp

	for _, column := range entity.Columns {
		columnExp := column.ToExp()

		if columnExp != nil {
			fields[column.Name] = columnExp
		}
	}

	return boolExp
}

func (entity *EntityMeta) toOrderBy() *graphql.InputObject {
	fields := graphql.InputObjectConfigFieldMap{}

	orderByExp := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:   entity.Name + ORDERBY,
			Fields: fields,
		},
	)

	for _, column := range entity.Columns {
		columnOrderBy := column.ToOrderBy()

		if columnOrderBy != nil {
			fields[column.Name] = &graphql.InputObjectFieldConfig{Type: columnOrderBy}
		}
	}

	return orderByExp
}
