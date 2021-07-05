package main

import (
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	api := openapi.NewOpenAPI()
	api.Title = "MyAPI"
	api.Description = "This is an example for go-rest"
	api.Version = "1.2.0"
	api.AddServer("http://localhost:8080", "my test server")
	api.AddTag("UserService", "All request for user management")
	err := api.WithBasicAuth("mybasic")
	checkErr(err)
	err = api.WithBearerAuth("mybearer", "bearer", "JWT")
	checkErr(err)
	api.DefaultResponse(&openapi.MethodResponse{
		Description: "DefaultError",
		Value:       &APIError{},
	})
	err = api.Get(MyGet)
	checkErr(err)
	authMiddleware := api.ChiAuthMiddleware(func(authName string, scopes []string, r *http.Request) bool {
		log.Printf("Auth check %s %s", authName, scopes)
		return authName != "mybearer"
	})

	r := chi.NewRouter()
	for _, handleConfig := range api.GetHandleFunc() {
		log.Printf("Add func %s %s", handleConfig.Method, handleConfig.Path)
		r.With(authMiddleware).MethodFunc(handleConfig.Method, handleConfig.Path, handleConfig.HandlerFunc)
	}

	h, err := api.OpenAPIHandlerFunc()
	checkErr(err)
	r.MethodFunc(http.MethodGet, "/openapi", h)

	log.Fatal(http.ListenAndServe(":8080", r))

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
