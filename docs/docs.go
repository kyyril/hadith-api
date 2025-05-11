// Updated docs/swagger.go
package docs

import (
	"fmt"
	"github.com/hadith-api/config"
	"github.com/swaggo/swag"
)

// The base template for the swagger documentation
var docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "API Hadis dengan terjemahan dari 9 perawi",
        "title": "Hadith API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "1.0"
    },
    "host": "%s",
    "basePath": "/api/v1",
    "paths": {
        "/hadis": {
            "get": {
                "description": "Returns all hadiths with pagination and optional search filtering",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "hadiths"
                ],
                "summary": "Get all hadiths",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number for pagination (default: 1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Items per page for pagination (default: 10)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search query to filter hadiths by ID (translation)",
                        "name": "q",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PaginatedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/hadis/{slug}": {
            "get": {
                "description": "Returns all hadiths from a specific narrator with optional pagination and filtering",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "hadiths"
                ],
                "summary": "Get hadiths by narrator",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Narrator slug (e.g., muslim, bukhari)",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Items per page for pagination",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search query to filter hadiths",
                        "name": "q",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PaginatedResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/hadis/{slug}/{number}": {
            "get": {
                "description": "Returns a specific hadith from a narrator by its number",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "hadiths"
                ],
                "summary": "Get hadith by narrator and number",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Narrator slug (e.g., muslim, bukhari)",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Hadith number",
                        "name": "number",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HadithResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/narrators": {
            "get": {
                "description": "Returns a list of all available hadith narrators",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "narrators"
                ],
                "summary": "Get list of available narrators",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HadithResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.HadithResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.Pagination": {
            "type": "object",
            "properties": {
                "current_page": {
                    "type": "integer"
                },
                "per_page": {
                    "type": "integer"
                },
                "total_items": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                }
            }
        },
        "models.PaginatedResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "pagination": {
                    "$ref": "#/definitions/models.Pagination"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

func init() {
	// Get the configuration
	cfg := config.GetConfig()

	// Extract the host part from the baseURL (remove protocol)
	host := cfg.BaseURL

	// Remove protocol part if present
	if host[:7] == "http://" {
		host = host[7:]
	} else if host[:8] == "https://" {
		host = host[8:]
	}

	// Format the doc with the dynamic host
	formattedDoc := fmt.Sprintf(docTemplate, host)

	// Register the swagger spec
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  formattedDoc,
	})
}
