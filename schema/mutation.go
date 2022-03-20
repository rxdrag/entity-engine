package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entity-engine/consts"
	"rxdrag.com/entity-engine/meta"
	"rxdrag.com/entity-engine/resolve"
	"rxdrag.com/entity-engine/utils"
)

func AppendToMutationFields(entity *meta.Entity, feilds *graphql.Fields) {
	//如果是枚举
	if entity.EntityType == meta.ENTITY_ENUM {
		return
	}

	name := utils.FirstUpper(entity.Name)

	(*feilds)[consts.DELETE+name] = &graphql.Field{
		Type: *Cache.MutationResponse(entity),
		Args: graphql.FieldConfigArgument{
			consts.ARG_WHERE: &graphql.ArgumentConfig{
				Type: Cache.WhereExp(entity),
			},
		},
		//Resolve: entity.QueryResolve(),
	}
	(*feilds)[consts.DELETE+name+consts.BY_ID] = &graphql.Field{
		Type: Cache.OutputType(entity),
		Args: graphql.FieldConfigArgument{
			consts.ID: &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		//Resolve: entity.QueryResolve(),
	}
	(*feilds)[consts.SAVE+name] = &graphql.Field{
		Type: Cache.OutputType(entity),
		Args: graphql.FieldConfigArgument{
			consts.ARG_OBJECTS: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{
					OfType: &graphql.List{
						OfType: &graphql.NonNull{
							OfType: *Cache.SaveInput(entity),
						},
					},
				},
			},
		},
	}
	//Resolve: entity.QueryResolve(),
	(*feilds)[consts.SAVE_ONE+name] = &graphql.Field{
		Type: Cache.OutputType(entity),
		Args: graphql.FieldConfigArgument{
			consts.ARG_OBJECT: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{
					OfType: *Cache.SaveInput(entity),
				},
			},
		},
		Resolve: resolve.PostOneResolveFn(entity),
	}

	(*feilds)[consts.UPDATE+name] = &graphql.Field{
		Type: *Cache.MutationResponse(entity),
		Args: graphql.FieldConfigArgument{
			consts.ARG_OBJECTS: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{
					OfType: &graphql.List{
						OfType: &graphql.NonNull{
							OfType: *Cache.UpdateInput(entity),
						},
					},
				},
			},
			consts.ARG_WHERE: &graphql.ArgumentConfig{
				Type: Cache.WhereExp(entity),
			},
		},
		//Resolve: entity.QueryResolve(),
	}
}
