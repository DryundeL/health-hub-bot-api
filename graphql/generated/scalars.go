package generated

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	graphqlScalars "github.com/health-hub-bot-api/internal/infrastructure/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// unmarshalInputDate - метод для unmarshal Date скаляра
func (ec *executionContext) unmarshalInputDate(ctx context.Context, v interface{}) (time.Time, error) {
	return graphqlScalars.UnmarshalDate(v)
}

// _Date - метод для marshal Date скаляра
func (ec *executionContext) _Date(ctx context.Context, sel ast.SelectionSet, v *time.Time) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return graphqlScalars.MarshalDate(*v)
}

// unmarshalInputTime - метод для unmarshal Time скаляра
func (ec *executionContext) unmarshalInputTime(ctx context.Context, v interface{}) (time.Time, error) {
	return graphqlScalars.UnmarshalTime(v)
}

// _Time - метод для marshal Time скаляра
func (ec *executionContext) _Time(ctx context.Context, sel ast.SelectionSet, v *time.Time) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return graphqlScalars.MarshalTime(*v)
}

