// Package swagger provides automatic OpenAPI/Swagger documentation generation for Go Chi routers.
//
// This package offers a simple and powerful way to automatically generate OpenAPI 3.0
// specifications from your Chi router definitions. It automatically detects endpoints,
// infers security requirements, and allows manual schema registration for complete API documentation.
//
// # Key Features
//
//   - Automatic endpoint detection from Chi routers
//   - Security scheme detection (JWT Bearer, API Key)
//   - Tag-based endpoint organization
//   - Automatic path parameter detection
//   - Manual schema registration for request/response bodies
//   - Complete OpenAPI 3.0 specification generation
//
// # Quick Start
//
// Create a configuration and detector:
//
//	cfg := swagger.DefaultConfig().
//		WithServiceInfo("My API", "1.0.0").
//		WithDescription("A wonderful API").
//		AddTagRule("/users", "Users").
//		AddTagRule("/posts", "Posts")
//
//	detector := swagger.NewDetector(cfg)
//
// Scan your Chi router:
//
//	err := detector.ScanRouter(router)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Generate OpenAPI specification:
//
//	spec, err := detector.ExportSwaggerSpec()
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Manual Schema Registration
//
// Register query parameters:
//
//	type ListParams struct {
//		Page  int    `json:"page" description:"Page number"`
//		Limit int    `json:"limit" description:"Items per page"`
//	}
//
//	registry := detector.GetSchemaRegistry()
//	registry.RegisterQueryParams("/users", "GET", ListParams{})
//
// Register request body:
//
//	type CreateUserRequest struct {
//		Name  string `json:"name" binding:"required"`
//		Email string `json:"email" binding:"required,email"`
//	}
//
//	registry.RegisterRequestBody("/users", "POST", CreateUserRequest{}, "User creation payload", true)
//
// Register response:
//
//	type UserResponse struct {
//		ID    string `json:"id"`
//		Name  string `json:"name"`
//		Email string `json:"email"`
//	}
//
//	registry.RegisterResponse("/users/{id}", "GET", 200, UserResponse{}, "User details")
//
// # Package Organization
//
// The package is organized into the following subpackages:
//
//   - config: Configuration structures and builders
//   - detector: Endpoint detection and route scanning
//   - schema: Schema registry and type reflection
//   - openapi: OpenAPI specification generation
//
// All main types are re-exported from the root package for backward compatibility.
package swagger
