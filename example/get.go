package main

import (
	"encoding/json"
	"github.com/Nightapes/go-rest/pkg/openapi"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type User struct {
	UserID string `json:"userID"`
}

type User2 struct {
	UserID string `json:"userID"`
}

type UserList = []string

var MyGet = &openapi.Get{
	Summary:        "Get User",
	Description:    "Get User with given ID",
	OperationID:    "GetMyTest",
	Tags:           []string{"UserService"},
	Authentication: map[string][]string{"mybasic": nil, "mybearer": {"users:read"}},
	Extensions:     map[string]interface{}{"x-custom": "test", "x-custom1": map[string]string{"test": "test"}},
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
		AddParameterList("ids", openapi.STRING, "DDDD", nil).
		WithQueryParameter("filter", openapi.STRING, "Filter stuff", false),
	HandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
		userID := chi.URLParam(request, "userID")
		user := &User{UserID: userID}
		resp, _ := json.Marshal(user)
		_, _ = writer.Write(resp)
	},
}
