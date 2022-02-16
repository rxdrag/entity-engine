package schema

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entity-engine/model"
	"rxdrag.com/entity-engine/utils"
)

const (
	DISTINCTEXP       string = "DistinctExp"
	INT_COMPARISONEXP string = "IntComparisonExp"
	BOOLEXP           string = "BoolExp"
)

func createFieldFtype(column *model.ColumnMeta) graphql.Output {
	switch column.Type {
	case model.COLUMN_NUMBER:
		return graphql.Int
	case model.COLUMN_BOOLEAN:
		return graphql.Boolean
	case model.COLUMN_STRING:
		return graphql.String
	case model.COLUMN_TEXT:
		return graphql.String
	case model.COLUMN_MEDIUM_TEXT:
		return graphql.String
	case model.COLUMN_LONG_TEXT:
		return graphql.String
	case model.COLUMN_DATE:
		return graphql.DateTime
	case model.COLUMN_SIMPLE_JSON:
		return graphql.NewScalar(graphql.ScalarConfig{Name: "JSON"})
	case model.COLUMN_SIMPLE_ARRAY:
		return graphql.NewScalar(graphql.ScalarConfig{Name: "JSON"})
	case model.COLUMN_JSON_ARRAY:
		return graphql.NewScalar(graphql.ScalarConfig{Name: "JSON"})
	case model.COLUMN_ENUM:
		return graphql.EnumValueType

	}

	panic("No column type")
}

func CreateEntityFields(entity *model.EntityMeta) *graphql.Fields {
	fields := &graphql.Fields{}
	for _, column := range entity.Columns {
		(*fields)[column.Name] = &graphql.Field{
			Type: createFieldFtype(&column),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fmt.Println(p.Context.Value("data"))
				return "world", nil
			},
		}
	}
	return fields
}

func AppendEntityToQueryFields(entity *model.EntityMeta, feilds *graphql.Fields) {
	metaType := graphql.NewObject(graphql.ObjectConfig{Name: entity.Name, Fields: *CreateEntityFields(entity)})
	metaDistinctType := graphql.NewEnum(graphql.EnumConfig{
		Name: entity.Name + DISTINCTEXP,
		Values: graphql.EnumValueConfigMap{
			"name": &graphql.EnumValueConfig{
				Value: "name",
			},
		},
	})

	andExp := graphql.InputObjectFieldConfig{}
	notExp := graphql.InputObjectFieldConfig{}
	orExp := graphql.InputObjectFieldConfig{}
	intComparisonExp := graphql.InputObjectFieldConfig{
		Type: graphql.NewInputObject(
			graphql.InputObjectConfig{
				Name: INT_COMPARISONEXP,
				Fields: graphql.InputObjectConfigFieldMap{
					"eq": &graphql.InputObjectFieldConfig{
						Type: graphql.Int,
					},
					"gt": &graphql.InputObjectFieldConfig{
						Type: graphql.Int,
					},
					"gte": &graphql.InputObjectFieldConfig{
						Type: graphql.Int,
					},
					"in": &graphql.InputObjectFieldConfig{
						Type: graphql.NewList(graphql.Int),
					},
				},
			},
		),
	}

	metaBoolExp := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: entity.Name + BOOLEXP,
			Fields: graphql.InputObjectConfigFieldMap{
				"and": &andExp,
				"not": &notExp,
				"or":  &orExp,
				"id":  &intComparisonExp,
			},
		},
	)
	andExp.Type = metaBoolExp
	notExp.Type = metaBoolExp
	orExp.Type = metaBoolExp

	(*feilds)[utils.FirstLower(entity.Name)] = &graphql.Field{
		Type: metaType,
		Args: graphql.FieldConfigArgument{
			"distinctOn": &graphql.ArgumentConfig{
				Type: metaDistinctType,
			},
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"orderBy": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"where": &graphql.ArgumentConfig{
				Type: metaBoolExp,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			fmt.Println(p.Context.Value("data"))
			return "world", nil
		},
	}
}
