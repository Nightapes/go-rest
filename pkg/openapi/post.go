package openapi

import (
	"net/http"
)

type RequestBodies struct {
	Description string
	Required    bool
	Bodies      map[string]interface{}
}

type FileUpload string

const (
	FileUploadBinary FileUpload = "binary"
	FileUploadBase64 FileUpload = "base64"
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
	// Deprecated: Please migrate to RequestBodies
	RequestBody   interface{}
	RequestBodies *RequestBodies
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

func (m *Post) GetRequestBodies() *RequestBodies {
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

func (m *Post) GetTags() []string {
	return m.Tags
}

func (m *Post) GetHandlerFunc() http.HandlerFunc {
	return m.HandlerFunc
}
