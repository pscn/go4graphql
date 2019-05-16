package notyet

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
)

func PlayGround(w http.ResponseWriter, r *http.Request) {
	h := handler.Playground("GraphQL playground", "/query")
	h(w, r)
}
