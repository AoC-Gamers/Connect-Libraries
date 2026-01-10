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

// getJSONFieldName obtiene el nombre del campo JSON desde los tags
func getJSONFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return ""
	}

	parts := strings.Split(tag, ",")
	return parts[0]
}

// isRequiredField determina si un campo es requerido basado en tags
func isRequiredField(field reflect.StructField) bool {
	jsonTag := field.Tag.Get("json")

	// Si tiene omitempty, no es requerido
	if strings.Contains(jsonTag, "omitempty") {
		return false
	}

	// Si tiene tag binding con required, es requerido
	bindingTag := field.Tag.Get("binding")
	if strings.Contains(bindingTag, "required") {
		return true
	}

	// Si es un puntero, generalmente es opcional
	if field.Type.Kind() == reflect.Ptr {
		return false
	}

	// Por defecto, si no tiene omitempty, considerarlo requerido
	return jsonTag != "" && !strings.Contains(jsonTag, "omitempty")
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
