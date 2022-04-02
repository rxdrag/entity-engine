package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entity-engine/consts"
	"rxdrag.com/entity-engine/model/graph"
	"rxdrag.com/entity-engine/utils"
)

var Cache TypeCache

type TypeCache struct {
	ObjectTypeMap        map[string]*graphql.Object
	EnumTypeMap          map[string]*graphql.Enum
	InterfaceTypeMap     map[string]*graphql.Interface
	UpdateInputMap       map[string]*graphql.InputObject
	SaveInputMap         map[string]*graphql.InputObject
	HasManyInputMap      map[string]*graphql.InputObject
	HasOneInputMap       map[string]*graphql.InputObject
	WhereExpMap          map[string]*graphql.InputObject
	DistinctOnEnumMap    map[string]*graphql.Enum
	OrderByMap           map[string]*graphql.InputObject
	EnumComparisonExpMap map[string]*graphql.InputObjectFieldConfig
	MutationResponseMap  map[string]*graphql.Output
	AggregateMap         map[string]*graphql.Output
}

var NodeInterfaceType = graphql.NewInterface(
	graphql.InterfaceConfig{
		Name: utils.FirstUpper(consts.NODE),
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
		},
		Description: "Node interface",
	},
)

func (c *TypeCache) MakeCache() {
	c.clearCache()
	c.makeEnums(Model.graph.Enums)
	c.makeOutputInterfaces(Model.graph.Interfaces)
	c.makeOutputObjects(Model.graph.Entities)
	c.makeRelations()
	c.makeArgs()
	c.makeInputs()
}

// func (c *TypeCache) OutputInterfaceType(entity *model.Entity) graphql.Type {
// 	return c.InterfaceTypeMap[entity.Name]
// }

func (c *TypeCache) OutputType(name string) graphql.Type {
	intf := c.InterfaceTypeMap[name]
	if intf != nil {
		return intf
	}
	obj := c.ObjectTypeMap[name]
	if obj == nil {
		panic("Can not find output type of " + name)
	}
	return obj
}

func (c *TypeCache) EnumType(name string) graphql.Type {
	return c.EnumTypeMap[name]
}

func (c *TypeCache) WhereExp(name string) *graphql.InputObject {
	return c.WhereExpMap[name]
}

func (c *TypeCache) OrderByExp(name string) *graphql.InputObject {
	return c.OrderByMap[name]
}

func (c *TypeCache) DistinctOnEnum(name string) *graphql.Enum {
	return c.DistinctOnEnumMap[name]
}

func (c *TypeCache) SaveInput(name string) *graphql.InputObject {
	return c.SaveInputMap[name]
}

func (c *TypeCache) UpdateInput(name string) *graphql.InputObject {
	return c.UpdateInputMap[name]
}
func (c *TypeCache) HasManyInput(name string) *graphql.InputObject {
	return c.HasManyInputMap[name]
}
func (c *TypeCache) HasOneInput(name string) *graphql.InputObject {
	return c.HasOneInputMap[name]
}

func (c *TypeCache) MutationResponse(name string) *graphql.Output {
	return c.MutationResponseMap[name]
}

func (c *TypeCache) mapInterfaces(entities []*graph.Entity) []*graphql.Interface {
	interfaces := []*graphql.Interface{NodeInterfaceType}
	for i := range entities {
		interfaces = append(interfaces, c.InterfaceTypeMap[entities[i].Name()])
	}

	return interfaces
}

func (c *TypeCache) clearCache() {
	c.ObjectTypeMap = make(map[string]*graphql.Object)
	c.EnumTypeMap = make(map[string]*graphql.Enum)
	c.InterfaceTypeMap = make(map[string]*graphql.Interface)
	c.UpdateInputMap = make(map[string]*graphql.InputObject)
	c.SaveInputMap = make(map[string]*graphql.InputObject)
	c.HasManyInputMap = make(map[string]*graphql.InputObject)
	c.HasOneInputMap = make(map[string]*graphql.InputObject)
	c.WhereExpMap = make(map[string]*graphql.InputObject)
	c.DistinctOnEnumMap = make(map[string]*graphql.Enum)
	c.OrderByMap = make(map[string]*graphql.InputObject)
	c.EnumComparisonExpMap = make(map[string]*graphql.InputObjectFieldConfig)
	c.MutationResponseMap = make(map[string]*graphql.Output)
	c.AggregateMap = make(map[string]*graphql.Output)
}
