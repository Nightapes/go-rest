package openapi

import (
	"net/http"
)

type Options struct {
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

func (m *Options) GetSummary() string {
	return m.Summary
}

func (m *Options) GetOperationID() string {
	return m.OperationID
}

func (m *Options) GetDescription() string {
	return m.Description
}

func (m *Options) GetResponse(s string) (string, interface{}) {
	if r, ok := m.Response[s]; ok {
		return r.Description, r.Value
	}
	return "", nil
}

func (m *Options) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Options) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Options) GetRequestBodies() *RequestBodies {
	return nil
}

func (m *Options) GetTags() []string {
	return m.Tags
}

func (m *Options) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
