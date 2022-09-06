package openapi

import (
	"net/http"
)

type Delete struct {
	Summary        string
	Description    string
	OperationID    string
	Tags           []string
	Authentication map[string][]string
	Response       map[string]MethodResponse
	Path           *PathBuilder
	Headers        []Parameter
	http.HandlerFunc
}

func (m *Delete) GetSummary() string {
	return m.Summary
}

func (m *Delete) GetOperationID() string {
	return m.OperationID
}

func (m *Delete) GetDescription() string {
	return m.Description
}

func (m *Delete) GetResponse(s string) (string, interface{}) {
	if r, ok := m.Response[s]; ok {
		return r.Description, r.Value
	}
	return "", nil
}

func (m *Delete) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Delete) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Delete) GetRequestBodies() *RequestBodies {
	return nil
}

func (m *Delete) GetTags() []string {
	return m.Tags
}

func (m *Delete) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
