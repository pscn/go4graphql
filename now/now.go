package notyet

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/pscn/go4graphql/api"
	"github.com/pscn/go4graphql/graph"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	h := handler.GraphQL(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: api.NewResolver(true)}),
	)
	h(w, r)
}
