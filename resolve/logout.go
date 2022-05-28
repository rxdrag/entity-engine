package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/authentication"
	"rxdrag.com/entify/common"
	"rxdrag.com/entify/utils"
)

func Logout(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()
	token := common.ParseContextValues(p).Token
	if token != "" {
		authentication.Logout(token)
	}
	return true, nil
}
