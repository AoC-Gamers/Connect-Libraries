package schema

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
