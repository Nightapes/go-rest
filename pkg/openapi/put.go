package openapi

import (
	"net/http"
)

type Put struct {
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

func (m *Put) GetSummary() string {
	return m.Summary
}

func (m *Put) GetOperationID() string {
	return m.OperationID
}

func (m *Put) GetDescription() string {
	return m.Description
}

func (m *Put) GetResponse(s string) (string, interface{}) {
	if r, ok := m.Response[s]; ok {
		return r.Description, r.Value
	}
	return "", nil
}

func (m *Put) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Put) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Put) GetRequestBody() interface{} {
	return m.RequestBody
}

func (m *Put) GetTags() []string {
	return m.Tags
}

func (m *Put) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
