package schema

import "fmt"

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
