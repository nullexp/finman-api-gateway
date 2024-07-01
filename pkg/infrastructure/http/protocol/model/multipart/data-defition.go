package multipart

import (
	"context"

	http "github.com/nullexp/finman-gateway-service/pkg/infrastructure/http/protocol"
)

type DataDefinition struct {
	Name     string
	Object   http.Verifier
	Optional bool
	Single   bool
}

func (f *DataDefinition) GetPartName() string {
	return f.Name
}

func (f *DataDefinition) GetSupportedTypes() []string {
	return []string{}
}

func (f *DataDefinition) GetObject() interface{} {
	return f.Object
}

func (f *DataDefinition) IsSingle() bool {
	return f.Single
}

func (f *DataDefinition) IsOptional() bool {
	return f.Optional
}

func (f *DataDefinition) Validate(context.Context) error {
	return f.Object.Validate(context.Background())
}

const UnknownData = "unknown data"