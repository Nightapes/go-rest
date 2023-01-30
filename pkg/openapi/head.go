package openapi

import (
	"net/http"
)

type Head struct {
	Summary        string
	Description    string
	OperationID    string
	Tags           []string
	Authentication map[string][]string
	Path           *PathBuilder
	Headers        []Parameter
	Extensions
	http.HandlerFunc
}

func (m *Head) GetSummary() string {
	return m.Summary
}

func (m *Head) GetOperationID() string {
	return m.OperationID
}

func (m *Head) GetDescription() string {
	return m.Description
}

func (m *Head) GetExtensions() map[string]interface{} {
	return m.Extensions
}

func (m *Head) GetResponse(s string) (string, interface{}) {
	return "", nil
}

func (m *Head) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Head) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Head) GetRequestBodies() *RequestBodies {
	return nil
}

func (m *Head) GetTags() []string {
	return m.Tags
}

func (m *Head) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
