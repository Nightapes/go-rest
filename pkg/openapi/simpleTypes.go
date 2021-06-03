package openapi

type Enum interface {
	GetValues() []interface{}
}

type DataType struct {
	Type   string
	Format string
}

type TYPE DataType

var (
	INTEGER = TYPE{
		Type: "integer",
	}
	//INT32 signed 32 bits
	INT32 = TYPE{
		Type:   "integer",
		Format: "int32",
	}
	//INT64 signed 64 bits (a.k.a long)
	INT64 = TYPE{
		Type:   "integer",
		Format: "int64",
	}
	FLOAT = TYPE{
		Type:   "number",
		Format: "float",
	}
	DOUBLE = TYPE{
		Type:   "number",
		Format: "double",
	}
	STRING = TYPE{
		Type: "string",
	}
	BOOLEAN = TYPE{
		Type: "boolean",
	}
	//BINARY any sequence of octets
	BINARY = TYPE{
		Type:   "string",
		Format: "binary",
	}
	//BYTE base64 encoded characters
	BYTE = TYPE{
		Type:   "string",
		Format: "byte",
	}
	//DATE As defined by full-date - RFC3339
	DATE = TYPE{
		Type:   "string",
		Format: "date",
	}
	//DATETIME As defined by date-time - RFC3339
	DATETIME = TYPE{
		Type:   "string",
		Format: "date-time",
	}
	//PASSWORD A hint to UIs to obscure input.
	PASSWORD = TYPE{
		Type:   "string",
		Format: "password",
	}
)

type Schema struct {
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
}

func (t *TYPE) toSchema() *SchemaRef {
	return &SchemaRef{
		Value: &Schema{
			Type:   t.Type,
			Format: t.Format,
		},
	}
}
