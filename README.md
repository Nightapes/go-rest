# go-rest

Generate openapi definition from code.

## Usage

```go
import "github.com/Nightapes/go-rest/pkg/openapi"
```

```go
package main 

import "github.com/Nightapes/go-rest/pkg/openapi"

var GetUser = &openapi.Get{
	Summary:        "Get User",
	Description:    "Get User with given ID",
	OperationID:    "GetMyTest",
	Tags:           []string{"UserService"},
	Authentication: map[string][]string{"mybasic": nil, "mybearer": {"users:read"}},
	Response: map[string]openapi.MethodResponse{
		"200": {
			Description: "The response with userID",
			Value: &User{
				UserID: "exampleID",
			},
		},
	},
	Headers: []openapi.Parameter{{Description: "My custom header", Name: "test-header", Required: false, Type: openapi.INTEGER}},
	Path: openapi.NewPathBuilder().
		Add("users").
		AddParameter("userId", openapi.STRING, "UserID").
		WithQueryParameter("filter", openapi.STRING, "Filter stuff", false),
	HandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
		user := &User{UserID: "userID"}
		resp, _ := json.Marshal(user)
		writer.Write(resp)
	},
}

func main() {
	
api := openapi.NewOpenAPI()
api.Title = "MyAPI"
api.Get(GetUser)
...
}
```

See `./example` for complete setup 

## Authentication

Middleware to get authentication is available for chi `router`


```go
authMiddleware := api.ChiAuthMiddleware(func(authName string, scopes []string, r *http.Request) bool {
    log.Printf("Auth check %s %s", authName, scopes)
    if authName == "mybearer" {
        return false
    }
    return true
})

r := chi.NewRouter()
for _, handleConfig := range api.GetHandleFunc() {
    log.Printf("Add func %s %s", handleConfig.Method, handleConfig.Path)
    r.With(authMiddleware).MethodFunc(handleConfig.Method, handleConfig.Path, handleConfig.HandlerFunc)
}
```
