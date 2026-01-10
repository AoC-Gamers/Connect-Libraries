package swagger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ParamLocation especifica dónde se encuentra el parámetro
type ParamLocation string

const (
	ParamInPath   ParamLocation = "path"
	ParamInQuery  ParamLocation = "query"
	ParamInHeader ParamLocation = "header"
)

// ParamSchema define la estructura de un parámetro
type ParamSchema struct {
	Name        string        `json:"name"`
	In          ParamLocation `json:"in"`
	Type        string        `json:"type"`
	Format      string        `json:"format,omitempty"`
	Required    bool          `json:"required"`
	Description string        `json:"description,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	Example     interface{}   `json:"example,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
}

// RequestBodySchema define el schema de un request body
type RequestBodySchema struct {
	Description string                 `json:"description,omitempty"`
	Required    bool                   `json:"required"`
	Content     map[string]interface{} `json:"content"`
}

// ResponseSchema define el schema de una respuesta
type ResponseSchema struct {
	Description string                 `json:"description"`
	Content     map[string]interface{} `json:"content,omitempty"`
}

// SchemaRegistry mantiene schemas registrados manualmente
type SchemaRegistry struct {
	queryParams map[string][]ParamSchema          // key: "METHOD:PATH"
	pathParams  map[string][]ParamSchema          // key: "METHOD:PATH"
	requestBody map[string]RequestBodySchema      // key: "METHOD:PATH"
	responses   map[string]map[int]ResponseSchema // key: "METHOD:PATH", inner key: statusCode
}

// NewSchemaRegistry crea un nuevo registro de schemas
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		queryParams: make(map[string][]ParamSchema),
		pathParams:  make(map[string][]ParamSchema),
		requestBody: make(map[string]RequestBodySchema),
		responses:   make(map[string]map[int]ResponseSchema),
	}
}

// makeKey genera la clave para el mapa de schemas
func (sr *SchemaRegistry) makeKey(method, path string) string {
	return fmt.Sprintf("%s:%s", strings.ToUpper(method), path)
}

// RegisterQueryParams registra parámetros de query usando reflection en un struct
func (sr *SchemaRegistry) RegisterQueryParams(path, method string, structType interface{}) {
	params := extractParamsFromStruct(structType, ParamInQuery)
	key := sr.makeKey(method, path)
	sr.queryParams[key] = params
}

// RegisterPathParams registra parámetros de path manualmente
func (sr *SchemaRegistry) RegisterPathParams(path, method string, params []ParamSchema) {
	key := sr.makeKey(method, path)
	sr.pathParams[key] = params
}

// RegisterRequestBody registra el schema de request body usando reflection
func (sr *SchemaRegistry) RegisterRequestBody(path, method string, structType interface{}, description string, required bool) {
	schema := buildSchemaFromStruct(structType)

	key := sr.makeKey(method, path)
	sr.requestBody[key] = RequestBodySchema{
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
func (sr *SchemaRegistry) RegisterResponse(path, method string, statusCode int, structType interface{}, description string) {
	schema := buildSchemaFromStruct(structType)

	key := sr.makeKey(method, path)
	if sr.responses[key] == nil {
		sr.responses[key] = make(map[int]ResponseSchema)
	}

	sr.responses[key][statusCode] = ResponseSchema{
		Description: description,
		Content: map[string]interface{}{
			"application/json": map[string]interface{}{
				"schema": schema,
			},
		},
	}
}

// GetQueryParams obtiene los query params registrados para un endpoint
func (sr *SchemaRegistry) GetQueryParams(method, path string) []ParamSchema {
	key := sr.makeKey(method, path)
	return sr.queryParams[key]
}

// GetPathParams obtiene los path params registrados para un endpoint
func (sr *SchemaRegistry) GetPathParams(method, path string) []ParamSchema {
	key := sr.makeKey(method, path)
	return sr.pathParams[key]
}

// GetRequestBody obtiene el request body registrado para un endpoint
func (sr *SchemaRegistry) GetRequestBody(method, path string) (RequestBodySchema, bool) {
	key := sr.makeKey(method, path)
	body, exists := sr.requestBody[key]
	return body, exists
}

// GetResponses obtiene las respuestas registradas para un endpoint
func (sr *SchemaRegistry) GetResponses(method, path string) map[int]ResponseSchema {
	key := sr.makeKey(method, path)
	return sr.responses[key]
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

// parseJSONTag parsea el tag JSON y retorna el nombre y si es omitempty
func parseJSONTag(tag string) (name string, omitempty bool) {
	parts := strings.Split(tag, ",")
	name = parts[0]

	for _, part := range parts[1:] {
		if part == "omitempty" {
			omitempty = true
			break
		}
	}

	return name, omitempty
}

// buildSchemaFromStruct construye un schema JSON desde un struct usando reflection
func buildSchemaFromStruct(v interface{}) map[string]interface{} {
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

// ConvertToSwaggerParams convierte ParamSchema a formato Swagger/OpenAPI
func ConvertToSwaggerParams(params []ParamSchema) []map[string]interface{} {
	result := make([]map[string]interface{}, len(params))

	for i, param := range params {
		p := map[string]interface{}{
			"name":     param.Name,
			"in":       string(param.In),
			"required": param.Required,
			"schema": map[string]interface{}{
				"type": param.Type,
			},
		}

		if param.Format != "" {
			p["schema"].(map[string]interface{})["format"] = param.Format
		}

		if param.Description != "" {
			p["description"] = param.Description
		}

		if param.Default != nil {
			p["schema"].(map[string]interface{})["default"] = param.Default
		}

		if param.Example != nil {
			p["example"] = param.Example
		}

		if len(param.Enum) > 0 {
			p["schema"].(map[string]interface{})["enum"] = param.Enum
		}

		result[i] = p
	}

	return result
}

// ConvertRequestBodyToSwagger convierte RequestBodySchema a formato Swagger/OpenAPI
func ConvertRequestBodyToSwagger(body RequestBodySchema) map[string]interface{} {
	result := map[string]interface{}{
		"required": body.Required,
		"content":  body.Content,
	}

	if body.Description != "" {
		result["description"] = body.Description
	}

	return result
}

// ConvertResponsesToSwagger convierte ResponseSchema a formato Swagger/OpenAPI
func ConvertResponsesToSwagger(responses map[int]ResponseSchema) map[string]interface{} {
	result := make(map[string]interface{})

	for code, response := range responses {
		codeStr := fmt.Sprintf("%d", code)
		respMap := map[string]interface{}{
			"description": response.Description,
		}

		if response.Content != nil {
			respMap["content"] = response.Content
		}

		result[codeStr] = respMap
	}

	return result
}

// MarshalSchemaToJSON convierte un schema a JSON (útil para debugging)
func MarshalSchemaToJSON(schema interface{}) (string, error) {
	bytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
