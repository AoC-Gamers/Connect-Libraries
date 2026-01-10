package openapi

import (
	"strings"

	"github.com/AoC-Gamers/connect-libraries/swagger/detector"
	"github.com/AoC-Gamers/connect-libraries/swagger/schema"
)

// generatePaths converts detected routes to OpenAPI paths format
func generatePaths(routes []detector.Route, schemas *schema.Registry) map[string]interface{} {
	paths := make(map[string]interface{})

	for _, route := range routes {
		if _, exists := paths[route.Path]; !exists {
			paths[route.Path] = make(map[string]interface{})
		}

		operation := buildOperation(route, schemas)
		method := strings.ToLower(route.Method)
		paths[route.Path].(map[string]interface{})[method] = operation
	}

	return paths
}

// buildOperation construye la operación completa para un endpoint
func buildOperation(route detector.Route, schemas *schema.Registry) map[string]interface{} {
	operation := map[string]interface{}{
		"summary":     route.Summary,
		"description": route.Description,
		"tags":        route.Tags,
	}

	// Agregar parámetros
	addParameters(operation, route, schemas)

	// Agregar request body si está registrado
	if requestBody, exists := schemas.GetRequestBody(route.Method, route.Path); exists {
		operation["requestBody"] = schema.ConvertRequestBodyToSwagger(requestBody)
	}

	// Agregar responses
	operation["responses"] = buildResponses(route, schemas)

	// Agregar security
	addSecurity(operation, route)

	return operation
}

// addParameters agrega parámetros (path y query) a la operación
func addParameters(operation map[string]interface{}, route detector.Route, schemas *schema.Registry) {
	// AUTO-DETECTAR path params desde la ruta
	autoPathParams := schema.ExtractPathParamsFromRoute(route.Path)

	// OBTENER params registrados manualmente
	manualQueryParams := schemas.GetQueryParams(route.Method, route.Path)
	manualPathParams := schemas.GetPathParams(route.Method, route.Path)

	// COMBINAR parámetros (manual override auto)
	var allParams []schema.ParamSchema
	if len(manualPathParams) > 0 {
		allParams = manualPathParams
	} else {
		allParams = autoPathParams
	}
	allParams = append(allParams, manualQueryParams...)

	// AGREGAR parámetros a la operación
	if len(allParams) > 0 {
		operation["parameters"] = schema.ConvertToSwaggerParams(allParams)
	}
}

// buildResponses construye las responses combinando defaults y registradas
func buildResponses(route detector.Route, schemas *schema.Registry) map[string]interface{} {
	responses := map[string]interface{}{
		"400": map[string]interface{}{"description": "Bad request"},
		"401": map[string]interface{}{"description": "Unauthorized"},
		"500": map[string]interface{}{"description": "Internal server error"},
	}

	// OBTENER responses registradas manualmente
	registeredResponses := schemas.GetResponses(route.Method, route.Path)

	// Agregar responses registradas (override defaults)
	if len(registeredResponses) > 0 {
		for code, response := range schema.ConvertResponsesToSwagger(registeredResponses) {
			responses[code] = response
		}
	} else {
		// Si no hay responses registradas, usar default 200
		responses["200"] = map[string]interface{}{
			"description": "Successful response",
		}
	}

	return responses
}

// addSecurity agrega seguridad a la operación si es necesaria
func addSecurity(operation map[string]interface{}, route detector.Route) {
	if len(route.Security) > 0 {
		security := make([]map[string][]string, 0)
		for _, sec := range route.Security {
			security = append(security, map[string][]string{sec: {}})
		}
		operation["security"] = security
	}
}
