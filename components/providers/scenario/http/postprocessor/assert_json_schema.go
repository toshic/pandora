package postprocessor

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
	"strings"
)

type AssertJsonSchema struct {
	SchemaPath string `config:"schema_path"`
}

func (a AssertJsonSchema) Process(resp *http.Response, _ io.Reader) (map[string]any, error) {
	schema := gojsonschema.NewReferenceLoader(a.SchemaPath)
	respJson, _ := gojsonschema.NewReaderLoader(resp.Body)
	result, err := gojsonschema.Validate(schema, respJson)

	if err != nil {
		return nil, fmt.Errorf("assert/json-schema check error: %w", err)
	}

	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}

		return nil, fmt.Errorf("assert/json-schema validation errors: %s", strings.Join(errors, ", "))
	}

	return nil, nil
}

func (a AssertJsonSchema) Validate() error {

	if a.SchemaPath == "" {
		return fmt.Errorf("assert/json-schema validation json-schema not found in path %s", a.SchemaPath)
	}

	return nil
}
