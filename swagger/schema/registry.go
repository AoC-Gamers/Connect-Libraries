package schema

import (
	"fmt"
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
