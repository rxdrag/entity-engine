package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entity-engine/model"
)

func (c *TypeCache) makeRelations() {
	for i := range model.GlobalModel.Graph.Interfaces {
		intf := model.GlobalModel.Graph.Interfaces[i]
		objectType := c.ObjectTypeMap[intf.Name()]
		for _, assocition := range intf.Associations() {
			objectType.AddFieldConfig(assocition.Name(), &graphql.Field{
				Name:        assocition.Name(),
				Type:        c.OutputType(assocition.TypeClass().Name()),
				Description: assocition.Description(),
			})
		}
	}
	for i := range model.GlobalModel.Graph.Entities {
		entity := model.GlobalModel.Graph.Entities[i]
		objectType := c.ObjectTypeMap[entity.Name()]
		for _, assocition := range entity.Associations() {
			objectType.AddFieldConfig(assocition.Name(), &graphql.Field{
				Name:        assocition.Name(),
				Type:        c.OutputType(assocition.TypeClass().Name()),
				Description: assocition.Description(),
			})
		}
	}
}
