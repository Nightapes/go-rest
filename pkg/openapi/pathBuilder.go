package openapi

import (
	"fmt"
	"strings"
)

type ParameterListOption struct {
	Explode bool
	Style   Style
}

type PathBuilder struct {
	path            string
	parameters      []Parameter
	queryParameters []Parameter
}

func NewPathBuilder() *PathBuilder {
	return &PathBuilder{
		parameters:      []Parameter{},
		queryParameters: []Parameter{},
	}
}

func GetPathFromBuilder(p *PathBuilder) string {
	return p.path
}

func (p *PathBuilder) Add(segment string) *PathBuilder {
	p.path = p.path + "/" + strings.TrimLeft(segment, "/")
	return p
}

func (p *PathBuilder) AddParameter(name string, t TYPE, description string) *PathBuilder {
	p.Add(fmt.Sprintf("{%s}", name))
	p.parameters = append(p.parameters, Parameter{
		Name:        name,
		Description: description,
		Type:        t,
		Required:    true,
	})
	return p
}

func (p *PathBuilder) AddParameterList(name string, t TYPE, description string, option *ParameterListOption) *PathBuilder {
	p.Add(fmt.Sprintf("{%s}", name))
	parameter := Parameter{
		Name:        name,
		Description: description,
		Type:        t,
		IsArray:     true,
		Required:    true,
	}
	if option != nil {
		parameter.Style = option.Style
		parameter.Explode = option.Explode
	}
	p.parameters = append(p.parameters, parameter)

	return p
}

func (p *PathBuilder) WithQueryParameter(name string, t TYPE, description string, required bool) *PathBuilder {
	p.queryParameters = append(p.queryParameters, Parameter{
		Name:        name,
		Description: description,
		Type:        t,
		Required:    required,
	})
	return p
}

func (p *PathBuilder) WithQueryParameterList(name string, t TYPE, description string, required bool, option *ParameterListOption) *PathBuilder {
	parameter := Parameter{
		Name:        name,
		Description: description,
		Type:        t,
		Required:    required,
		IsArray:     true,
	}
	if option != nil {
		parameter.Style = option.Style
		parameter.Explode = option.Explode
	}
	p.queryParameters = append(p.queryParameters, parameter)
	return p
}
