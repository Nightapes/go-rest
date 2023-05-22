package openapi

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

func (a *API) ChiAuthMiddleware(isAllowed func(authName string, scopes []string, r *http.Request) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rctx := chi.RouteContext(r.Context())
			routePattern := strings.Join(rctx.RoutePatterns, "")
			log.Printf("Check auth with path %s", routePattern)

			isAuthenticated := false
			if config, ok := a.handlers[routePattern]; ok {
				if _, methodFound := config[rctx.RouteMethod]; !methodFound {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				authTypes := config[rctx.RouteMethod].AuthTypes
				for name, scopes := range authTypes {
					if isAllowed(name, scopes, r) {
						isAuthenticated = true
						break
					}
				}
				if !isAuthenticated && len(authTypes) != 0 {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
