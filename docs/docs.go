// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/login": {
            "post": {
                "description": "accepts json sent by the user as input and authorize it",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "authenticates the user",
                "parameters": [
                    {
                        "description": "account info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Authentication was successful"
                    },
                    "400": {
                        "description": "error: Invalid to insert token"
                    },
                    "403": {
                        "description": "error: Invalid email or password"
                    },
                    "422": {
                        "description": "error: Invalid password size"
                    },
                    "500": {
                        "description": "error: Invalid to create token"
                    }
                }
            }
        },
        "/api/auth/refresh": {
            "post": {
                "description": "accept json and refresh user refresh and access tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "refresh user's tokens",
                "parameters": [
                    {
                        "description": "account info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: RefreshTokens was successful"
                    },
                    "401": {
                        "description": "error: Invalid to get refresh token from cookie"
                    },
                    "500": {
                        "description": "error: Invalid to create token"
                    }
                }
            }
        },
        "/api/auth/registration": {
            "post": {
                "description": "accepts json sent by the user as input and registers it",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "registration"
                ],
                "summary": "registers a user",
                "parameters": [
                    {
                        "description": "account info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Registration was successful"
                    },
                    "400": {
                        "description": "error: Failed to read body"
                    },
                    "409": {
                        "description": "error: email or channel already been use"
                    },
                    "422": {
                        "description": "error: Failed create password, because it exceeds the character limit or backwards"
                    },
                    "500": {
                        "description": "error: Failed to hash password. Please, try again later"
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "example"
                ],
                "summary": "ping example",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "plain"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "models.RegisterRequest": {
            "type": "object",
            "properties": {
                "channel_name": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Video-hosting API",
	Description:      "API Server for video-hosting application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
