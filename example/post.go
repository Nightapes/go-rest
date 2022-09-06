package main

import (
	"encoding/json"
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var example = []string{}

var MyPost = &openapi.Post{
	Summary:        "Post User",
	Description:    "Post User with given ID",
	OperationID:    "PostMyTest",
	Tags:           []string{"UserService"},
	Authentication: map[string][]string{"mybasic": nil, "mybearer": {"users:read"}},
	RequestBodies: &openapi.RequestBodies{
		Description: "Test",
		Required:    true,
		Bodies: map[string]interface{}{
			"application/json": &example,
			"application/xml":  openapi.FileUploadBinary,
		},
	},
	Response: map[string]openapi.MethodResponse{
		"200": {
			Description: "The response with userID",
			Value:       UserList{"test"},
		},
		"201": {
			Description: "The response with userID",
			Value:       &[]User{{"test"}},
		},
		"202": {
			Description: "Map Test",
			Value: map[string]User2{
				"test": {
					UserID: "myID",
				},
			},
		},
		"204": {
			Description: "The response with userID",
		},
	},
	Headers: []openapi.Parameter{{Description: "My custom header", Name: "test-header", Required: false, Type: openapi.INTEGER}},
	Path: openapi.NewPathBuilder().
		Add("users").
		AddParameter("userId", openapi.TYPEENUM([]string{"aws"}), "UserID").
		WithQueryParameter("filter", openapi.STRING, "Filter stuff", false),
	HandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
		userID := chi.URLParam(request, "userID")
		user := &User{UserID: userID}
		resp, _ := json.Marshal(user)
		_, _ = writer.Write(resp)
	},
}
