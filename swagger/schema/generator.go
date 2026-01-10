package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Registry mantiene schemas registrados manualmente
type Registry struct {
	queryParams map[string][]ParamSchema          // key: "METHOD:PATH"
	pathParams  map[string][]ParamSchema          // key: "METHOD:PATH"
	requestBody map[string]RequestBodySchema      // key: "METHOD:PATH"
	responses   map[string]map[int]ResponseSchema // key: "METHOD:PATH", inner key: statusCode
}

// NewRegistry crea un nuevo registro de schemas
func NewRegistry() *Registry {
	return &Registry{
		queryParams: make(map[string][]ParamSchema),
		pathParams:  make(map[string][]ParamSchema),
		requestBody: make(map[string]RequestBodySchema),
		responses:   make(map[string]map[int]ResponseSchema),
	}
}

// makeKey genera la clave para el mapa de schemas
func (r *Registry) makeKey(method, path string) string {
	return fmt.Sprintf("%s:%s", strings.ToUpper(method), path)
}

// RegisterQueryParams registra parámetros de query usando reflection en un struct
func (r *Registry) RegisterQueryParams(path, method string, structType interface{}) {
	params := extractParamsFromStruct(structType, ParamInQuery)
	key := r.makeKey(method, path)
	r.queryParams[key] = params
}

// RegisterPathParams registra parámetros de path manualmente
func (r *Registry) RegisterPathParams(path, method string, params []ParamSchema) {
	key := r.makeKey(method, path)
	r.pathParams[key] = params
}

// RegisterRequestBody registra el schema de request body usando reflection
func (r *Registry) RegisterRequestBody(path, method string, structType interface{}, description string, required bool) {
	schema := BuildSchemaFromStruct(structType)

	key := r.makeKey(method, path)
	r.requestBody[key] = RequestBodySchema{
		Description: description,
		Required:    required,
		Content: map[string]interface{}{
			"application/json": map[string]interface{}{
				"schema": schema,
			},
		},
	}
}

// RegisterResponse registra el schema de una respuesta
func (r *Registry) RegisterResponse(path, method string, statusCode int, structType interface{}, description string) {
	schema := BuildSchemaFromStruct(structType)

	key := r.makeKey(method, path)
	if r.responses[key] == nil {
		r.responses[key] = make(map[int]ResponseSchema)
	}

	r.responses[key][statusCode] = ResponseSchema{
		Description: description,
		Content: map[string]interface{}{
			"application/json": map[string]interface{}{
				"schema": schema,
			},
		},
	}
}

// GetQueryParams obtiene los query params registrados para un endpoint
func (r *Registry) GetQueryParams(method, path string) []ParamSchema {
	key := r.makeKey(method, path)
	return r.queryParams[key]
}

// GetPathParams obtiene los path params registrados para un endpoint
func (r *Registry) GetPathParams(method, path string) []ParamSchema {
	key := r.makeKey(method, path)
	return r.pathParams[key]
}

// GetRequestBody obtiene el request body registrado para un endpoint
func (r *Registry) GetRequestBody(method, path string) (RequestBodySchema, bool) {
	key := r.makeKey(method, path)
	body, exists := r.requestBody[key]
	return body, exists
}

// GetResponses obtiene las respuestas registradas para un endpoint
func (r *Registry) GetResponses(method, path string) map[int]ResponseSchema {
	key := r.makeKey(method, path)
	return r.responses[key]
}

// extractParamsFromStruct extrae parámetros desde struct tags usando reflection
func extractParamsFromStruct(v interface{}, location ParamLocation) []ParamSchema {
	t := reflect.TypeOf(v)
	if t == nil {
		return []ParamSchema{}
	}

	// Si es un puntero, obtener el tipo subyacente
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return []ParamSchema{}
	}

	params := []ParamSchema{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Ignorar campos no exportados
		if !field.IsExported() {
			continue
		}

		// Obtener nombre desde tag json
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		paramName, omitempty := parseJSONTag(jsonTag)

		param := ParamSchema{
			Name:     paramName,
			In:       location,
			Type:     goTypeToSwaggerType(field.Type),
			Format:   goTypeToSwaggerFormat(field.Type),
			Required: !omitempty,
		}

		// Extraer descripción desde tag comment o validate
		if desc := field.Tag.Get("description"); desc != "" {
			param.Description = desc
		}

		// Extraer default desde tag
		if def := field.Tag.Get("default"); def != "" {
			param.Default = def
		}

		// Extraer ejemplo desde tag
		if example := field.Tag.Get("example"); example != "" {
			param.Example = example
		}

		params = append(params, param)
	}

	return params
}

// BuildSchemaFromStruct construye un schema JSON desde un struct usando reflection
func BuildSchemaFromStruct(v interface{}) map[string]interface{} {
	t := reflect.TypeOf(v)
	if t == nil {
		return map[string]interface{}{"type": "object"}
	}

	// Si es un puntero, obtener el tipo subyacente
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		// Para tipos primitivos
		return map[string]interface{}{
			"type": goTypeToSwaggerType(t),
		}
	}

	properties, required := buildStructProperties(t)

	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// buildStructProperties construye las propiedades y campos requeridos de un struct
func buildStructProperties(t reflect.Type) (map[string]interface{}, []string) {
	properties := make(map[string]interface{})
	required := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		fieldName, omitempty := parseJSONTag(jsonTag)
		fieldSchema := buildFieldSchema(field)

		properties[fieldName] = fieldSchema

		if !omitempty {
			required = append(required, fieldName)
		}
	}

	return properties, required
}

// buildFieldSchema construye el schema de un campo individual
func buildFieldSchema(field reflect.StructField) map[string]interface{} {
	fieldSchema := map[string]interface{}{
		"type": goTypeToSwaggerType(field.Type),
	}

	if format := goTypeToSwaggerFormat(field.Type); format != "" {
		fieldSchema["format"] = format
	}

	if desc := field.Tag.Get("description"); desc != "" {
		fieldSchema["description"] = desc
	}

	if example := field.Tag.Get("example"); example != "" {
		fieldSchema["example"] = example
	}

	return fieldSchema
}

// MarshalToJSON convierte un schema a JSON (útil para debugging)
func MarshalToJSON(schema interface{}) (string, error) {
	bytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
