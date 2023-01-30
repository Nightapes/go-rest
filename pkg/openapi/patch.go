package openapi

import (
	"net/http"
)

type Patch struct {
	Summary        string
	Description    string
	OperationID    string
	Tags           []string
	Authentication map[string][]string
	Response       map[string]MethodResponse
	Path           *PathBuilder
	Headers        []Parameter
	// Deprecated: Please migrate to RequestBodies
	RequestBody   interface{}
	RequestBodies *RequestBodies
	Extensions
	http.HandlerFunc
}

func (m *Patch) GetSummary() string {
	return m.Summary
}

func (m *Patch) GetOperationID() string {
	return m.OperationID
}

func (m *Patch) GetDescription() string {
	return m.Description
}

func (m *Patch) GetExtensions() map[string]interface{} {
	return m.Extensions
}

func (m *Patch) GetResponse(s string) (string, interface{}) {
	if r, ok := m.Response[s]; ok {
		return r.Description, r.Value
	}
	return "", nil
}

func (m *Patch) GetAuthentication(key string) (bool, []string) {
	if a, ok := m.Authentication[key]; ok {
		return true, a
	}
	return false, nil
}

func (m *Patch) GetHeaders() []Parameter {
	return m.Headers
}

func (m *Patch) GetRequestBodies() *RequestBodies {
	if m.RequestBodies != nil {
		return m.RequestBodies
	}

	if m.RequestBody != nil {
		return &RequestBodies{
			Required: true,
			Bodies: map[string]interface{}{
				"application/json": m.RequestBody,
			},
		}
	}

	return nil
}

func (m *Patch) GetTags() []string {
	return m.Tags
}

func (m *Patch) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
