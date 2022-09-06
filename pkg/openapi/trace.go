package openapi

import (
	"net/http"
)

type Trace struct {
	Summary        string
	Description    string
	OperationID    string
	Tags           []string
	Authentication map[string][]string
	Path           *PathBuilder
	Headers        []Parameter
	http.HandlerFunc
}

func (m *Trace) GetSummary() string {
	return m.Summary
}

func (m *Trace) GetOperationID() string {
	return m.OperationID
}

func (m *Trace) GetDescription() string {
	return m.Description
}

func (m *Trace) GetResponse(s string) (string, interface{}) {
	return "", nil
}

func (m *Trace) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Trace) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Trace) GetRequestBodies() *RequestBodies {
	return nil
}

func (m *Trace) GetTags() []string {
	return m.Tags
}

func (m *Trace) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
