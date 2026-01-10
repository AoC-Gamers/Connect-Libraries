package schema

import (
	"reflect"
	"strings"
	"time"
)

// goTypeToSwaggerType convierte un tipo de Go a tipo Swagger/OpenAPI
func goTypeToSwaggerType(t reflect.Type) string {
	// Manejar punteros
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Array, reflect.Slice:
		return "array"
	case reflect.Map, reflect.Struct:
		// Casos especiales para structs comunes
		if t == reflect.TypeOf(time.Time{}) {
			return "string"
		}
		return "object"
	default:
		return "string"
	}
}

// goTypeToSwaggerFormat determina el formato Swagger desde un tipo Go
func goTypeToSwaggerFormat(t reflect.Type) string {
	// Manejar punteros
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// time.Time
	if t == reflect.TypeOf(time.Time{}) {
		return "date-time"
	}

	switch t.Kind() {
	case reflect.Int32:
		return "int32"
	case reflect.Int64:
		return "int64"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.Uint8: // byte
		return "byte"
	default:
		return ""
	}
}

// ExtractPathParamsFromRoute detecta automáticamente path params desde la ruta
// Ejemplo: "/users/{id}/teams/{teamId}" → ["id", "teamId"]
func ExtractPathParamsFromRoute(path string) []ParamSchema {
	params := []ParamSchema{}

	// Buscar patrones {param}
	start := -1
	for i, ch := range path {
		if ch == '{' {
			start = i + 1
		} else if ch == '}' && start != -1 {
			paramName := path[start:i]

			// Determinar tipo basado en nombre
			paramType := inferParamType(paramName)

			params = append(params, ParamSchema{
				Name:        paramName,
				In:          ParamInPath,
				Type:        paramType,
				Required:    true,
				Description: generateParamDescription(paramName),
			})

			start = -1
		}
	}

	return params
}

// inferParamType intenta inferir el tipo desde el nombre del parámetro
func inferParamType(paramName string) string {
	lower := strings.ToLower(paramName)

	// IDs son strings por defecto (pueden ser UUIDs)
	if strings.HasSuffix(lower, "id") {
		return "string"
	}

	// Números comunes
	if strings.Contains(lower, "page") || strings.Contains(lower, "limit") ||
		strings.Contains(lower, "count") || strings.Contains(lower, "size") {
		return "integer"
	}

	// Default a string
	return "string"
}

// generateParamDescription genera una descripción básica desde el nombre
func generateParamDescription(paramName string) string {
	// Convertir camelCase o snake_case a palabras
	words := splitParamName(paramName)

	// Capitalizar primera palabra
	if len(words) > 0 {
		words[0] = strings.Title(words[0])
	}

	description := strings.Join(words, " ")

	// Agregar contexto basado en sufijos comunes
	if strings.HasSuffix(strings.ToLower(paramName), "id") {
		description += " identifier"
	}

	return description
}

// splitParamName divide un nombre de parámetro en palabras
func splitParamName(name string) []string {
	// Manejar snake_case
	if strings.Contains(name, "_") {
		return strings.Split(name, "_")
	}

	// Manejar camelCase
	words := []string{}
	word := ""

	for i, ch := range name {
		if i > 0 && ch >= 'A' && ch <= 'Z' {
			if word != "" {
				words = append(words, word)
			}
			word = string(ch)
		} else {
			word += string(ch)
		}
	}

	if word != "" {
		words = append(words, word)
	}

	return words
}

// MergeParams combina múltiples listas de parámetros eliminando duplicados
func MergeParams(paramLists ...[]ParamSchema) []ParamSchema {
	seen := make(map[string]bool)
	result := []ParamSchema{}

	for _, params := range paramLists {
		for _, param := range params {
			key := string(param.In) + ":" + param.Name
			if !seen[key] {
				seen[key] = true
				result = append(result, param)
			}
		}
	}

	return result
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

	// Manejar arrays/slices
	if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
		elemType := field.Type.Elem()
		fieldSchema["items"] = map[string]interface{}{
			"type": goTypeToSwaggerType(elemType),
		}
		if format := goTypeToSwaggerFormat(elemType); format != "" {
			fieldSchema["items"].(map[string]interface{})["format"] = format
		}
	}

	return fieldSchema
}

// ExtractParamsFromStruct extrae parámetros desde struct tags usando reflection
func ExtractParamsFromStruct(v interface{}, location ParamLocation) []ParamSchema {
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

		// Extraer descripción desde tag
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
