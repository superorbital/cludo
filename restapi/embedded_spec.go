// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Cludo - Cloud Sudo",
    "title": "Cludo API",
    "contact": {
      "name": "SuperOrbital",
      "url": "http://superorbital.io/",
      "email": "info@superorbital.io"
    },
    "version": "1.0.0"
  },
  "basePath": "/v1",
  "paths": {
    "/environment": {
      "post": {
        "description": "Generate a temporary environment (set of environment variables)",
        "tags": [
          "environment"
        ],
        "summary": "Generate a temporary environment",
        "operationId": "generate-environment",
        "parameters": [
          {
            "description": "Temporary Environment Request definition",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/models.EnvironmentRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/models.EnvironmentResponse"
            }
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/role": {
      "get": {
        "description": "List all roles available to current user",
        "tags": [
          "role"
        ],
        "summary": "List all roles",
        "operationId": "list-roles",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "models.EnvironmentBundle": {
      "type": "object",
      "additionalProperties": {
        "type": "string"
      },
      "example": {
        "ENV_VAR_1": "Hello!",
        "ENV_VAR_2": "Bonjour!"
      }
    },
    "models.EnvironmentRequest": {
      "type": "object",
      "properties": {
        "roleID": {
          "description": "The id of the role to apply.",
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "models.EnvironmentResponse": {
      "type": "object",
      "properties": {
        "bundle": {
          "$ref": "#/definitions/models.EnvironmentBundle"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Cludo - Cloud Sudo",
    "title": "Cludo API",
    "contact": {
      "name": "SuperOrbital",
      "url": "http://superorbital.io/",
      "email": "info@superorbital.io"
    },
    "version": "1.0.0"
  },
  "basePath": "/v1",
  "paths": {
    "/environment": {
      "post": {
        "description": "Generate a temporary environment (set of environment variables)",
        "tags": [
          "environment"
        ],
        "summary": "Generate a temporary environment",
        "operationId": "generate-environment",
        "parameters": [
          {
            "description": "Temporary Environment Request definition",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/models.EnvironmentRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/models.EnvironmentResponse"
            }
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/role": {
      "get": {
        "description": "List all roles available to current user",
        "tags": [
          "role"
        ],
        "summary": "List all roles",
        "operationId": "list-roles",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "models.EnvironmentBundle": {
      "type": "object",
      "additionalProperties": {
        "type": "string"
      },
      "example": {
        "ENV_VAR_1": "Hello!",
        "ENV_VAR_2": "Bonjour!"
      }
    },
    "models.EnvironmentRequest": {
      "type": "object",
      "properties": {
        "roleID": {
          "description": "The id of the role to apply.",
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "models.EnvironmentResponse": {
      "type": "object",
      "properties": {
        "bundle": {
          "$ref": "#/definitions/models.EnvironmentBundle"
        }
      }
    }
  }
}`))
}