package openapi

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/jsoninfo"
)

type OpenAPI struct {
	OpenAPI      string `json:"openapi" yaml:"openapi" validate:"required"`
	Info         `json:"info" yaml:"info" validate:"required"`
	Servers      []*Server            `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        map[string]*PathItem `json:"paths" yaml:"paths"  validate:"required"`
	Components   *Components          `json:"components" yaml:"components"`
	Tags         []*Tag               `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocs        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type Tag struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Info struct {
	Title          string `json:"title" yaml:"title" validate:"required"`
	Description    string `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	//Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	//License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version string `json:"version" yaml:"version" validate:"required"` // Required
}

type Components struct {
	SecuritySchemes map[string]interface{} `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Schemas         map[string]interface{} `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

type PathItem struct {
	Summary     string     `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Connect     *Operation `json:"connect,omitempty" yaml:"connect,omitempty"`
	Delete      *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Head        *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	Options     *Operation `json:"options,omitempty" yaml:"options,omitempty"`
	Patch       *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	Post        *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Put         *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Trace       *Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
}

type RefValue struct {
	Ref   string      `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

type SchemaRef struct {
	Ref   string      `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

func (value *SchemaRef) MarshalJSON() ([]byte, error) {
	return jsoninfo.MarshalRef(value.Ref, value.Value)
}

func (value *SchemaRef) UnmarshalJSON(data []byte) error {
	return jsoninfo.UnmarshalRef(data, &value.Ref, &value.Value)
}

type SchemaRefs []*SchemaRef
type Schemas map[string]*SchemaRef

type Response struct {
	Description string                `json:"description" yaml:"description" validate:"required"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

type RequestBody struct {
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

type MediaType struct {
	Schema  interface{} `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example interface{} `json:"example,omitempty" yaml:"example,omitempty"`
}

type Operation struct {
	Tags        []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary     string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Deprecated  bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Responses   map[string]*Response   `json:"responses" yaml:"responses" validate:"required"`
	Security    []map[string][]string  `json:"security,omitempty" yaml:"security,omitempty"`
	RequestBody *RequestBody           `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Parameters  []GenericParameter     `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	extensions  map[string]interface{} `json:"-" yaml:"-"`
}

func (u *Operation) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(&struct {
		Tags        []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
		Summary     string                `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string                `json:"description,omitempty" yaml:"description,omitempty"`
		OperationID string                `json:"operationId,omitempty" yaml:"operationId,omitempty"`
		Deprecated  bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
		Responses   map[string]*Response  `json:"responses" yaml:"responses" validate:"required"`
		Security    []map[string][]string `json:"security,omitempty" yaml:"security,omitempty"`
		RequestBody *RequestBody          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
		Parameters  []GenericParameter    `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	}{
		Tags:        u.Tags,
		Summary:     u.Summary,
		Description: u.Description,
		OperationID: u.Description,
		Deprecated:  u.Deprecated,
		Responses:   u.Responses,
		Security:    u.Security,
		RequestBody: u.RequestBody,
		Parameters:  u.Parameters,
	})

	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)

	for i, extension := range u.extensions {
		m[i] = extension
	}

	return json.Marshal(m)
}

type GenericParameter struct {
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
	In          string      `json:"in,omitempty" yaml:"in,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Deprecated  bool        `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Required    bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Example     interface{} `json:"example,omitempty" yaml:"example,omitempty"`
	Schema      *SchemaRef  `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type Server struct {
	URL         string `json:"url" yaml:"url" validate:"required,url"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	//Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

type BasicAuth struct {
	Type   string `json:"type,omitempty" yaml:"type,omitempty" validate:"required"`
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty" validate:"required"`
}

type BearerAuth struct {
	Type         string `json:"type,omitempty" yaml:"type,omitempty" validate:"required"`
	Scheme       string `json:"scheme,omitempty" yaml:"scheme,omitempty" validate:"required"`
	BearerFormat string `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty" validate:"required"`
}

type InTypes string

const (
	HEADER InTypes = "header"
	QUERY  InTypes = "query"
	COOKIE InTypes = "cookie"
)

type ApiKeyAuth struct {
	Type string  `json:"type,omitempty" yaml:"type,omitempty" validate:"required"`
	In   InTypes `json:"in,omitempty" yaml:"in,omitempty" validate:"required"`
	Name string  `json:"name,omitempty" yaml:"name,omitempty" validate:"required"`
}

type OpenIDConnectAuth struct {
	Type             string `json:"type,omitempty" yaml:"type,omitempty" validate:"required"`
	OpenIdConnectUrl string `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty" validate:"required"`
}

type Extensions map[string]interface{}
