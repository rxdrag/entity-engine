package schema

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entity-engine/config"
	"rxdrag.com/entity-engine/utils"
)

const (
	BOOLEXP     string = "BoolExp"
	ORDERBY     string = "OrderBy"
	DISTINCTEXP string = "DistinctExp"
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
	EnumValues  []byte       `json:"enumValues"`
}

//where表达式缓存，query跟mutation都用
var whereExpMap = make(map[string]*graphql.InputObject)

//类型缓存， query用
var outputTypeMap = make(map[string]*graphql.Output)

var distinctOnEnumMap = make(map[string]*graphql.Enum)

var orderByMap = make(map[string]*graphql.InputObject)

func (entity *EntityMeta) createQueryFields() graphql.Fields {
	fields := graphql.Fields{}
	for _, column := range entity.Columns {
		fields[column.Name] = &graphql.Field{
			Type: column.toType(),
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	fmt.Println(p.Context.Value("data"))
			// 	return "world", nil
			// },
		}
	}
	return fields
}

func (entity *EntityMeta) toOutputType() graphql.Output {
	if outputTypeMap[entity.Name] != nil {
		return *outputTypeMap[entity.Name]
	}
	var returnValue graphql.Output

	if entity.EntityType == Entity_ENUM {
		enumValues := make(map[string]interface{})
		json.Unmarshal(entity.EnumValues, &enumValues)
		enumValueConfigMap := graphql.EnumValueConfigMap{}
		for enumName, enumValue := range enumValues {
			var value, ok = enumValue.(string)
			if !ok {
				value = enumValue.(map[string]string)["value"]
			}
			enumValueConfigMap[enumName] = &graphql.EnumValueConfig{
				Value: value,
			}
		}
		returnValue = graphql.NewEnum(
			graphql.EnumConfig{
				Name:   entity.Name,
				Values: enumValueConfigMap,
			},
		)
	} else {
		returnValue = graphql.NewObject(
			graphql.ObjectConfig{
				Name:   entity.Name,
				Fields: entity.createQueryFields(),
			},
		)
	}
	outputTypeMap[entity.Name] = &returnValue
	return returnValue
}

func (entity *EntityMeta) toWhereExp() *graphql.InputObject {
	expName := entity.Name + BOOLEXP
	if whereExpMap[expName] != nil {
		return whereExpMap[expName]
	}

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
			Name:   expName,
			Fields: fields,
		},
	)
	andExp.Type = &graphql.List{
		OfType: &graphql.NonNull{
			OfType: boolExp,
		},
	}
	notExp.Type = boolExp
	orExp.Type = &graphql.List{
		OfType: &graphql.NonNull{
			OfType: boolExp,
		},
	}

	for _, column := range entity.Columns {
		columnExp := column.ToExp()

		if columnExp != nil {
			fields[column.Name] = columnExp
		}
	}
	whereExpMap[expName] = boolExp
	return boolExp
}

func (entity *EntityMeta) toOrderBy() *graphql.InputObject {
	if orderByMap[entity.Name] != nil {
		return orderByMap[entity.Name]
	}
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

	orderByMap[entity.Name] = orderByExp
	return orderByExp
}

func (entity *EntityMeta) toDistinctOnEnum() *graphql.Enum {
	if distinctOnEnumMap[entity.Name] != nil {
		return distinctOnEnumMap[entity.Name]
	}
	enumValueConfigMap := graphql.EnumValueConfigMap{}
	for _, column := range entity.Columns {
		enumValueConfigMap[column.Name] = &graphql.EnumValueConfig{
			Value: column.Name,
		}
	}

	entEnum := graphql.NewEnum(
		graphql.EnumConfig{
			Name:   entity.Name + DISTINCTEXP,
			Values: enumValueConfigMap,
		},
	)
	distinctOnEnumMap[entity.Name] = entEnum
	return entEnum
}

func (entity *EntityMeta) getTableName() string {
	if (*entity).TableName != "" {
		return (*entity).TableName
	}
	return utils.SnakeString((*entity).Name)
}

func (entity *EntityMeta) QueryResolve() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		db, err := sql.Open("mysql", config.MYSQL_CONFIG)
		defer db.Close()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		queryStr := "select * from %s"

		queryStr = fmt.Sprintf(queryStr, entity.getTableName())
		//err = db.Select(&instances, queryStr)
		rows, err := db.Query(queryStr)
		columns, err := rows.Columns()
		var instances []utils.JSON
		for rows.Next() {
			row := make(map[string]interface{})
			values := make([]interface{}, len(columns))
			for i, columnName := range columns {
				if columnName == "content" {
					var value utils.JSON
					values[i] = &value
				} else {
					var value string
					values[i] = &value
				}

			}
			err = rows.Scan(values...)
			for i, value := range values {
				row[columns[i]] = value
			}
			//err = rows.MapScan(row)
			// for i, encoded := range row {
			// 	switch encoded.(type) {
			// 	case byte:
			// 		row[i] = encoded.(byte)
			// 		break
			// 	case []byte:
			// 		row[i] = string(encoded.([]byte))
			// 		break
			// 	case time.Time:
			// 		row[i] = encoded
			// 		// if val.IsZero() {
			// 		// 	ret[columns[i]] = nil
			// 		// } else {
			// 		// 	ret[columns[i]] = val.Format("2006-01-02 15:04:05")
			// 		// }
			// 		break
			// 	default:
			// 		row[i] = encoded
			// 	}
			// }
			instances = append(instances, row)
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fmt.Println("Resolve entity:" + entity.Name)
		fmt.Println(p.Args)
		fmt.Println(p.Context.Value("data"))
		return instances, nil
	}
}
