package openapi

import (
	"net/http"
)

type MethodResponse struct {
	Description string
	Value       interface{}
}

type Get struct {
	Summary        string
	Description    string
	OperationID    string
	Tags           []string
	Authentication map[string][]string
	Response       map[string]MethodResponse
	Path           *PathBuilder
	Headers        []Parameter
	Extensions
	http.HandlerFunc
}

func (m *Get) GetSummary() string {
	return m.Summary
}

func (m *Get) GetOperationID() string {
	return m.OperationID
}

func (m *Get) GetDescription() string {
	return m.Description
}

func (m *Get) GetExtensions() map[string]interface{} {
	return m.Extensions
}

func (m *Get) GetResponse(s string) (string, interface{}) {
	if r, ok := m.Response[s]; ok {
		return r.Description, r.Value
	}
	return "", nil
}

func (m *Get) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Get) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Get) GetRequestBodies() *RequestBodies {
	return nil
}

func (m *Get) GetTags() []string {
	return m.Tags
}

func (m *Get) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
