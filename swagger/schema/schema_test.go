package schema

import "testing"

type sampleQuery struct {
	UserID int    `json:"user_id" description:"User identifier"`
	Search string `json:"search,omitempty" default:"all" example:"john"`
	Skip   string `json:"-"`
}

type sampleBody struct {
	ID    string   `json:"id" description:"Identifier"`
	Items []int64  `json:"items,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

func TestPathParamsAndMerge(t *testing.T) {
	params := ExtractPathParamsFromRoute("/users/{id}/teams/{teamId}")
	if len(params) != 2 {
		t.Fatalf("expected 2 path params, got %d", len(params))
	}
	if params[0].Name != "id" || params[0].In != ParamInPath || !params[0].Required {
		t.Fatal("unexpected first path param")
	}

	merged := MergeParams(params, []ParamSchema{{Name: "id", In: ParamInPath}, {Name: "page", In: ParamInQuery}})
	if len(merged) != 3 {
		t.Fatalf("expected deduplicated merge size 3, got %d", len(merged))
	}
}

func TestBuildSchemaFromStructAndParamExtraction(t *testing.T) {
	schemaMap := BuildSchemaFromStruct(sampleBody{})
	properties, ok := schemaMap["properties"].(map[string]interface{})
	if !ok || len(properties) == 0 {
		t.Fatal("expected properties in schema")
	}

	idField := properties["id"].(map[string]interface{})
	if idField["type"] != "string" {
		t.Fatalf("expected id type string, got %v", idField["type"])
	}
	itemsField := properties["items"].(map[string]interface{})
	if itemsField["type"] != "array" {
		t.Fatalf("expected items type array, got %v", itemsField["type"])
	}

	params := ExtractParamsFromStruct(sampleQuery{}, ParamInQuery)
	if len(params) != 2 {
		t.Fatalf("expected 2 params from struct tags, got %d", len(params))
	}
	if params[0].Name != "user_id" || !params[0].Required {
		t.Fatal("expected required user_id param")
	}
	if params[1].Name != "search" || params[1].Required {
		t.Fatal("expected optional search param")
	}
}

func TestConverters(t *testing.T) {
	params := []ParamSchema{{
		Name:        "status",
		In:          ParamInQuery,
		Type:        "string",
		Required:    false,
		Description: "Current status",
		Default:     "active",
		Example:     "inactive",
		Enum:        []interface{}{"active", "inactive"},
	}}
	swaggerParams := ConvertToSwaggerParams(params)
	if len(swaggerParams) != 1 {
		t.Fatalf("expected 1 swagger param, got %d", len(swaggerParams))
	}
	if swaggerParams[0]["name"] != "status" {
		t.Fatal("expected param name status")
	}

	body := ConvertRequestBodyToSwagger(RequestBodySchema{
		Description: "Body",
		Required:    true,
		Content:     map[string]interface{}{"application/json": map[string]interface{}{}},
	})
	if body["required"] != true {
		t.Fatal("expected request body required true")
	}

	responses := ConvertResponsesToSwagger(map[int]ResponseSchema{
		200: {Description: "OK"},
		404: {Description: "Not found"},
	})
	if len(responses) != 2 || responses["200"] == nil || responses["404"] == nil {
		t.Fatal("expected converted response entries for 200 and 404")
	}
}
