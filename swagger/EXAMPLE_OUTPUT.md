# Ejemplo de Output OpenAPI con Parámetros

Este documento muestra cómo se ve el JSON generado por la biblioteca swagger v1.1.0 con el sistema de parámetros y schemas.

## Endpoint con Path Parameter (Auto-Detectado)

### Definición
```go
// Ruta Chi
r.Get("/users/{id}", GetUserHandler)

// Schema registrado
registry.RegisterResponse("/users/{id}", "GET", 200,
    models.User{},
    "User details")
```

### Output JSON
```json
{
  "paths": {
    "/users/{id}": {
      "get": {
        "tags": ["Users"],
        "summary": "Get user by ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "User identifier"
          }
        ],
        "responses": {
          "200": {
            "description": "User details",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            }
          },
          "400": {
            "description": "Bad Request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal Server Error"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    }
  }
}
```

---

## Endpoint con Query Parameters

### Definición
```go
// Query params struct
type UserFilters struct {
    Page     int    `json:"page" binding:"min=1" description:"Page number"`
    Limit    int    `json:"limit" binding:"min=1,max=100" description:"Items per page"`
    Status   string `json:"status" description:"User status filter"`
    Sort     string `json:"sort" description:"Sort field (username|created_at)"`
}

// Ruta Chi
r.Get("/users", ListUsersHandler)

// Schema registrado
registry.RegisterQueryParams("/users", "GET", UserFilters{})
registry.RegisterResponse("/users", "GET", 200,
    models.UserListResponse{},
    "List of users")
```

### Output JSON
```json
{
  "paths": {
    "/users": {
      "get": {
        "tags": ["Users"],
        "summary": "List users",
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "minimum": 1
            },
            "description": "Page number"
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "minimum": 1,
              "maximum": 100
            },
            "description": "Items per page"
          },
          {
            "name": "status",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "User status filter"
          },
          {
            "name": "sort",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "Sort field (username|created_at)"
          }
        ],
        "responses": {
          "200": {
            "description": "List of users",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/UserListResponse"
                }
              }
            }
          }
        }
      }
    }
  }
}
```

---

## Endpoint con Request Body

### Definición
```go
// Request struct
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20" description:"Unique username"`
    Email    string `json:"email" binding:"required,email" description:"Email address"`
    Password string `json:"password" binding:"required,min=8" description:"User password"`
    Bio      string `json:"bio,omitempty" description:"User biography"`
}

// Ruta Chi
r.Post("/users", CreateUserHandler)

// Schema registrado
registry.RegisterRequestBody("/users", "POST",
    CreateUserRequest{},
    "User creation data",
    true)
registry.RegisterResponse("/users", "POST", 201,
    models.User{},
    "User created successfully")
```

### Output JSON
```json
{
  "paths": {
    "/users": {
      "post": {
        "tags": ["Users"],
        "summary": "Create user",
        "requestBody": {
          "required": true,
          "description": "User creation data",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "username": {
                    "type": "string",
                    "description": "Unique username",
                    "minLength": 3,
                    "maxLength": 20
                  },
                  "email": {
                    "type": "string",
                    "description": "Email address",
                    "format": "email"
                  },
                  "password": {
                    "type": "string",
                    "description": "User password",
                    "minLength": 8
                  },
                  "bio": {
                    "type": "string",
                    "description": "User biography"
                  }
                },
                "required": ["username", "email", "password"]
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "User created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            }
          },
          "400": {
            "description": "Bad Request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal Server Error"
          }
        }
      }
    }
  }
}
```

---

## Endpoint Completo (Path + Body + Responses)

### Definición
```go
// Update request
type UpdateUserRequest struct {
    Username string `json:"username,omitempty" binding:"min=3,max=20" description:"New username"`
    Bio      string `json:"bio,omitempty" description:"New biography"`
}

// Ruta Chi
r.Patch("/users/{id}", UpdateUserHandler)

// Schemas registrados
registry.RegisterRequestBody("/users/{id}", "PATCH",
    UpdateUserRequest{},
    "User update data",
    false) // not required (body opcional)

registry.RegisterResponse("/users/{id}", "PATCH", 200,
    models.User{},
    "User updated successfully")

registry.RegisterResponse("/users/{id}", "PATCH", 404,
    models.ErrorResponse{},
    "User not found")
```

### Output JSON
```json
{
  "paths": {
    "/users/{id}": {
      "patch": {
        "tags": ["Users"],
        "summary": "Update user",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "User identifier"
          }
        ],
        "requestBody": {
          "required": false,
          "description": "User update data",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "username": {
                    "type": "string",
                    "description": "New username",
                    "minLength": 3,
                    "maxLength": 20
                  },
                  "bio": {
                    "type": "string",
                    "description": "New biography"
                  }
                },
                "required": []
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "User updated successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/User"
                }
              }
            }
          },
          "404": {
            "description": "User not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "400": {
            "description": "Bad Request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal Server Error"
          }
        },
        "security": [
          {
            "BearerAuth": []
          }
        ]
      }
    }
  }
}
```

---

## Components/Schemas Generados

Los schemas de modelos se generan automáticamente en la sección `components/schemas`:

```json
{
  "components": {
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string",
            "description": "User ID"
          },
          "username": {
            "type": "string",
            "description": "Username"
          },
          "email": {
            "type": "string",
            "description": "Email address"
          },
          "bio": {
            "type": "string",
            "description": "User biography"
          },
          "created_at": {
            "type": "string",
            "format": "date-time",
            "description": "Creation timestamp"
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {
            "type": "string",
            "description": "Error message"
          },
          "code": {
            "type": "string",
            "description": "Error code"
          }
        }
      }
    },
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  }
}
```

---

## Ejemplo Real: Connect-RT Presence Endpoint

### Código
```go
// PATCH /rt/presence
registry.RegisterRequestBody("/rt/presence", "PATCH",
    models.UpdatePresenceRequest{
        Status models.UserPlatformStatus `json:"status" binding:"required,min=1,max=5"`
    },
    "Update user presence status",
    true)

registry.RegisterResponse("/rt/presence", "PATCH", 200,
    models.PresenceResponse{
        Presence models.UserPresence `json:"presence"`
    },
    "Presence updated successfully")
```

### JSON Generado
```json
{
  "paths": {
    "/rt/presence": {
      "patch": {
        "tags": ["Presence"],
        "summary": "Update presence",
        "requestBody": {
          "required": true,
          "description": "Update user presence status",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "status": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 5
                  }
                },
                "required": ["status"]
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Presence updated successfully",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/PresenceResponse"
                }
              }
            }
          },
          "400": {"description": "Bad Request"},
          "401": {"description": "Unauthorized"},
          "500": {"description": "Internal Server Error"}
        },
        "security": [{"BearerAuth": []}]
      }
    }
  }
}
```

---

## Comparación v1.0.0 vs v1.1.0

### v1.0.0 (Sin Parámetros)
```json
{
  "paths": {
    "/users/{id}": {
      "get": {
        "tags": ["Users"],
        "summary": "Get user by ID",
        "responses": {
          "200": {"description": "OK"},
          "400": {"description": "Bad Request"}
        }
      }
    }
  }
}
```

### v1.1.0 (Con Parámetros)
```json
{
  "paths": {
    "/users/{id}": {
      "get": {
        "tags": ["Users"],
        "summary": "Get user by ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {"type": "string"},
            "description": "User identifier"
          }
        ],
        "responses": {
          "200": {
            "description": "User details",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/User"}
              }
            }
          },
          "400": {"description": "Bad Request"},
          "401": {"description": "Unauthorized"},
          "404": {"description": "Not Found"},
          "500": {"description": "Internal Server Error"}
        }
      }
    }
  }
}
```

**Beneficios:**
✅ Path parameters documentados automáticamente  
✅ Response schemas con tipos estructurados  
✅ Swagger UI permite probar endpoints con valores reales  
✅ Validaciones visibles (min, max, required, format)  
✅ Descripciones contextuales en cada campo  
