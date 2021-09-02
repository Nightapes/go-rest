package openapi

import (
	"fmt"
	"strings"
)

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

func (p *PathBuilder) WithQueryParameter(name string, t TYPE, description string, required bool) *PathBuilder {
	p.queryParameters = append(p.queryParameters, Parameter{
		Name:        name,
		Description: description,
		Type:        t,
		Required:    required,
	})
	return p
}
