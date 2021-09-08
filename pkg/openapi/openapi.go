package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/Nightapes/go-rest/pkg/jsonschema"
	"github.com/go-playground/validator/v10"

	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type GetDescription interface {
	GetDescription() string
}

type Parameter struct {
	Description string
	Name        string
	Required    bool
	Type        TYPE
}

type PathDesc interface {
	GetSummary() string
	GetOperationID() string
	GetDescription
	GetResponse(string) (string, interface{})
	GetAuthentication(key string) (bool, []string)
	GetHeaders() []Parameter
	GetRequestBody() interface{}
	GetTags() []string
	GetHandlerFunc() http.HandlerFunc
}

type HandlerConfig struct {
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
	authTypes   map[string][]string
}

type API struct {
	OpenAPI
	validate        *validator.Validate
	jsonschemaRefl  *jsonschema.Reflector
	defaultResponse *MethodResponse
	handlers        map[string]map[string]HandlerConfig
}

func NewOpenAPI() *API {
	return &API{
		validate: validator.New(),
		handlers: map[string]map[string]HandlerConfig{},
		jsonschemaRefl: &jsonschema.Reflector{
			TypeMapper: func(r reflect.Type) *jsonschema.Type {
				modelType := reflect.TypeOf((*Enum)(nil)).Elem()

				if r.Implements(modelType) {
					enumPointer := reflect.New(r)
					enumValue := enumPointer.Elem()
					enumInterface := enumValue.Interface()
					enumObject := enumInterface.(Enum)
					if len(enumObject.GetValues()) == 0 {
						log.Printf("Could not set enum values, no values found in type %s", r.Kind().String())
						return nil
					}

					definition := &jsonschema.Type{
						Enum: enumObject.GetValues(),
						Type: reflect.TypeOf(enumObject.GetValues()[0]).Kind().String(),
					}

					return &jsonschema.Type{

						Ref: "#/definitions/" + r.Name(),
						Definitions: map[string]*jsonschema.Type{
							r.Name(): definition,
						},
					}
				}

				return nil
			}},
		OpenAPI: OpenAPI{
			Tags:    []*Tag{},
			OpenAPI: "3.0.0",
			Paths:   map[string]*PathItem{},
			Components: &Components{
				SecuritySchemes: map[string]interface{}{},
				Schemas:         map[string]interface{}{},
			},
		}}
}

func (a *API) AddExternalDocs(url, externalDescription string) {
	a.ExternalDocs = &ExternalDocs{
		Description: externalDescription,
		URL:         url,
	}
}

func (a *API) AddTagExternalDocs(name, description string, url, externalDescription string) {
	a.Tags = append(a.Tags, &Tag{
		Name:        name,
		Description: description,
		ExternalDocs: &ExternalDocs{
			Description: externalDescription,
			URL:         url,
		},
	})
}

func (a *API) AddTag(name, description string) {
	a.Tags = append(a.Tags, &Tag{
		Name:        name,
		Description: description,
	})
}

func (a *API) Get(option *Get) error {
	return a.AddPath(http.MethodGet, option.Path, option)
}
func (a *API) Options(option *Options) error {
	return a.AddPath(http.MethodOptions, option.Path, option)
}
func (a *API) Post(option *Post) error {
	return a.AddPath(http.MethodPost, option.Path, option)
}
func (a *API) Put(option *Put) error {
	return a.AddPath(http.MethodPut, option.Path, option)
}
func (a *API) Patch(option *Patch) error {
	return a.AddPath(http.MethodPatch, option.Path, option)
}
func (a *API) Delete(option *Delete) error {
	return a.AddPath(http.MethodDelete, option.Path, option)
}
func (a *API) Head(option *Patch) error {
	return a.AddPath(http.MethodHead, option.Path, option)
}
func (a *API) Trace(option *Trace) error {
	return a.AddPath(http.MethodTrace, option.Path, option)
}

func (a *API) AddPath(method string, pathBuilder *PathBuilder, desc PathDesc) error {
	path := pathBuilder.path
	if _, ok := a.OpenAPI.Paths[path]; !ok {
		a.OpenAPI.Paths[path] = &PathItem{
			Summary:     desc.GetSummary(),
			Description: desc.GetDescription(),
		}
	}

	ops, handlerConfig := a.toPath(desc, pathBuilder)
	handlerConfig.Method = method
	if a.handlers[path] == nil {
		a.handlers[path] = map[string]HandlerConfig{}
	}
	a.handlers[path][method] = *handlerConfig
	switch method {
	case http.MethodGet:
		a.OpenAPI.Paths[path].Get = ops
	case http.MethodPost:
		a.OpenAPI.Paths[path].Post = ops
	case http.MethodPut:
		a.OpenAPI.Paths[path].Put = ops
	case http.MethodDelete:
		a.OpenAPI.Paths[path].Delete = ops
	case http.MethodHead:
		a.OpenAPI.Paths[path].Head = ops
	case http.MethodPatch:
		a.OpenAPI.Paths[path].Patch = ops
	case http.MethodTrace:
		a.OpenAPI.Paths[path].Trace = ops
	case http.MethodConnect:
		a.OpenAPI.Paths[path].Trace = ops
	}

	return nil
}

func (a *API) toPath(desc PathDesc, pathBuilder *PathBuilder) (*Operation, *HandlerConfig) {
	ops := &Operation{
		Tags:        desc.GetTags(),
		Summary:     desc.GetSummary(),
		Description: desc.GetDescription(),
		Responses:   map[string]*Response{},
		Parameters:  []GenericParameter{},
		OperationID: desc.GetOperationID(),
		Deprecated:  false,
		Security:    []map[string][]string{},
	}

	for _, parameter := range pathBuilder.parameters {
		ops.Parameters = append(ops.Parameters, GenericParameter{
			Name:        parameter.Name,
			Description: parameter.Description,
			In:          "path",
			Deprecated:  false,
			Required:    parameter.Required,
			Schema:      parameter.Type.toSchema(),
		})
	}

	for _, parameter := range pathBuilder.queryParameters {
		ops.Parameters = append(ops.Parameters, GenericParameter{
			Name:        parameter.Name,
			Description: parameter.Description,
			In:          "query",
			Deprecated:  false,
			Required:    parameter.Required,
			Schema:      parameter.Type.toSchema(),
		})
	}

	for _, parameter := range desc.GetHeaders() {
		ops.Parameters = append(ops.Parameters, GenericParameter{
			Name:        parameter.Name,
			Description: parameter.Description,
			In:          "header",
			Deprecated:  false,
			Required:    parameter.Required,
			Schema:      parameter.Type.toSchema(),
		})
	}

	handlerConfig := &HandlerConfig{
		Path:        pathBuilder.path,
		HandlerFunc: desc.GetHandlerFunc(),
		authTypes:   map[string][]string{},
	}

	for s := range a.OpenAPI.Components.SecuritySchemes {
		if ok, scopes := desc.GetAuthentication(s); ok {
			if scopes == nil {
				scopes = make([]string, 0)
			}
			sec := map[string][]string{
				s: scopes,
			}
			ops.Security = append(ops.Security, sec)
			handlerConfig.authTypes[s] = scopes
		}
	}

	if desc.GetRequestBody() != nil {
		body := a.jsonschemaRefl.Reflect(desc.GetRequestBody())
		for s, t := range body.Definitions {
			a.OpenAPI.Components.Schemas[s] = t
		}
		ops.RequestBody = &RequestBody{
			Required: true,
			Content: map[string]*MediaType{
				"application/json": {Schema: &SchemaRef{Ref: body.Ref}, Example: desc.GetRequestBody()},
			},
		}
	}

	if a.defaultResponse != nil {
		a.handleResponse(a.defaultResponse.Description, a.defaultResponse.Value, "default", ops)
	}

	respDescDefault, respDefault := desc.GetResponse("default")
	a.handleResponse(respDescDefault, respDefault, "default", ops)
	for i := 200; i < 561; i++ {
		respDesc, resp := desc.GetResponse(strconv.Itoa(i))
		a.handleResponse(respDesc, resp, strconv.Itoa(i), ops)
	}
	return ops, handlerConfig
}

func (a *API) handleResponse(respDesc string, resp interface{}, i string, ops *Operation) {
	if resp != nil {
		schema := a.jsonschemaRefl.Reflect(resp)
		//TODO allow rename of objects
		/*for s, _ := range schema.Definitions {
				fmt.Printf("Change key %s with name %s", s, name)
				schema.Definitions[name] = schema.Definitions[s]
				schema.Ref = strings.Replace(schema.Ref, s, name, 1)
				delete(schema.Definitions, s)
				break
		}*/

		for s, t := range schema.Definitions {
			a.handleEnumInProperties(t)
			a.handleEnumInArrays(t)
			a.OpenAPI.Components.Schemas[s] = t
		}

		ops.Responses[i] = &Response{
			Content: map[string]*MediaType{
				"application/json": {Schema: &SchemaRef{Value: schema.Type}, Example: resp},
			},
			Description: respDesc,
		}
		return
	}

	if respDesc != "" {
		ops.Responses[i] = &Response{
			Description: respDesc,
		}
	}

}

func (a *API) handleEnumInArrays(t *jsonschema.Type) {
	items := t.Items
	if items == nil {
		return
	}
	if len(items.Definitions) > 0 {
		for defKey, defValue := range items.Definitions {
			a.OpenAPI.Components.Schemas[defKey] = defValue
		}
	}
	t.Items.Definitions = nil
}

func (a *API) handleEnumInProperties(t *jsonschema.Type) {
	for _, propKey := range t.Properties.Keys() {
		prop, _ := t.Properties.Get(propKey)
		parsed, _ := prop.(*jsonschema.Type)
		a.handleEnumInArrays(parsed)
		if len(parsed.Definitions) > 0 {
			for defKey, defValue := range parsed.Definitions {
				a.OpenAPI.Components.Schemas[defKey] = defValue
			}
		}
		parsed.Definitions = nil
		t.Properties.Set(propKey, prop)
	}
}

func (a *API) AddServer(url, description string) {
	if a.OpenAPI.Servers == nil {
		a.OpenAPI.Servers = make([]*Server, 0)
	}
	a.OpenAPI.Servers = append(a.OpenAPI.Servers, &Server{
		URL:         url,
		Description: description,
	})
}

func (a *API) WithBasicAuth(key string) error {
	if _, ok := a.OpenAPI.Components.SecuritySchemes[key]; ok {
		return fmt.Errorf("a authentication with key %s already exists", key)
	}
	a.OpenAPI.Components.SecuritySchemes[key] = &BasicAuth{
		Type:   "http",
		Scheme: "basic",
	}
	return nil
}
func (a *API) WithBearerAuth(key, scheme, bearerFormat string) error {
	if _, ok := a.OpenAPI.Components.SecuritySchemes[key]; ok {
		return fmt.Errorf("a authentication with key %s already exists", key)
	}
	a.OpenAPI.Components.SecuritySchemes[key] = &BearerAuth{
		Type:         "http",
		Scheme:       scheme,
		BearerFormat: bearerFormat,
	}
	return nil
}

func (a *API) WithApiKey(key string, in InTypes, name string) error {
	if _, ok := a.OpenAPI.Components.SecuritySchemes[key]; ok {
		return fmt.Errorf("a authentication with key %s already exists", key)
	}
	a.OpenAPI.Components.SecuritySchemes[key] = &ApiKeyAuth{
		Type: "apiKey",
		In:   in,
		Name: name,
	}
	return nil
}

func (a *API) WithOpenIDConnect(key string, openIdConnectUrl string) error {
	if _, ok := a.OpenAPI.Components.SecuritySchemes[key]; ok {
		return fmt.Errorf("a authentication with key %s already exists", key)
	}
	a.OpenAPI.Components.SecuritySchemes[key] = &OpenIDConnectAuth{
		Type:             "openIdConnect",
		OpenIdConnectUrl: openIdConnectUrl,
	}
	return nil
}

func (a *API) Print() {
	for key, path := range a.OpenAPI.Paths {
		fmt.Printf("My path %s \n", key)

		if path.Get != nil {
			fmt.Printf("My get path %s \n", path.Get.Description)
		}
		if path.Post != nil {
			fmt.Printf("My post path %s \n", path.Post.Description)
		}
	}
}

func (a *API) ToJSON() ([]byte, error) {
	err := a.Validate()
	if err != nil {
		return nil, err
	}

	j, err := json.MarshalIndent(a.OpenAPI, "", "    ")
	if err != nil {
		return nil, err
	}
	return a.patchOpenApi(j), nil
}

func (a *API) ToYAML() ([]byte, error) {
	err := a.Validate()
	if err != nil {
		return nil, err
	}

	j, err := a.ToJSON()
	if err != nil {
		return nil, err
	}

	jsonObj := &yaml.MapSlice{}
	err = yaml.Unmarshal(j, jsonObj)
	if err != nil {
		return nil, err
	}

	// Marshal this object into YAML.
	y, err := yaml.Marshal(jsonObj)
	if err != nil {
		return nil, err
	}

	return y, nil
}

func (a *API) Validate() error {
	return a.validate.Struct(a.OpenAPI)
}

func (a *API) GetHandleFunc() []HandlerConfig {
	c := make([]HandlerConfig, 0)

	for _, path := range a.handlers {
		for _, config := range path {
			c = append(c, config)
		}
	}
	return c
}

func (a *API) OpenAPIHandlerFunc() (http.HandlerFunc, error) {
	yamlApi, err := a.ToYAML()
	if err != nil {
		return nil, err
	}
	jsonApi, err := a.ToJSON()
	if err != nil {
		return nil, err
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Accept") == "application/x-yaml" {
			writer.Header().Set("Content-Type", "application/x-yaml")
			_, _ = writer.Write(yamlApi)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(jsonApi)

	}, nil
}

func (a *API) DefaultResponse(response *MethodResponse) {
	a.defaultResponse = response
}

func (a *API) patchOpenApi(openapi []byte) []byte {
	// Replace json definitions with open api schemas
	patch := strings.ReplaceAll(string(openapi), "#/definitions", "#/components/schemas")
	// Remove unneeded json draft version
	patch = strings.ReplaceAll(patch, `"$schema": "http://json-schema.org/draft-04/schema#",`, "")
	return []byte(patch)
}
