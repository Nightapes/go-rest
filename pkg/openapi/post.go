package openapi

import (
	"net/http"
)

type Post struct {
	Summary        string
	Description    string
	OperationID    string
	Tags           []string
	Authentication map[string][]string
	Response       map[string]MethodResponse
	Path           *PathBuilder
	Headers        []Parameter
	RequestBody    interface{}
	http.HandlerFunc
}

func (m *Post) GetSummary() string {
	return m.Summary
}

func (m *Post) GetOperationID() string {
	return m.OperationID
}

func (m *Post) GetDescription() string {
	return m.Description
}

func (m *Post) GetResponse(s string) (string, interface{}) {
	if r, ok := m.Response[s]; ok {
		return r.Description, r.Value
	}
	return "", nil
}

func (m *Post) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Post) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Post) GetRequestBody() interface{} {
	return m.RequestBody
}

func (m *Post) GetTags() []string {
	return m.Tags
}

func (m *Post) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
