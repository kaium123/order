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
        "license": {
            "name": "Apache 2.0"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/healthz": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/todos": {
            "get": {
                "tags": [
                    "todos"
                ],
                "summary": "Find all todos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "Data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.Todo"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todos"
                ],
                "summary": "Create a new todo",
                "parameters": [
                    {
                        "description": "json",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.ResponseError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.Todo"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    }
                }
            }
        },
        "/todos/:id": {
            "get": {
                "tags": [
                    "todos"
                ],
                "summary": "Find a todo",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "Data": {
                                            "$ref": "#/definitions/model.Todo"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "todos"
                ],
                "summary": "Update a todo",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UpdateRequestBody"
                        }
                    },
                    {
                        "type": "integer",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "Data": {
                                            "$ref": "#/definitions/model.Todo"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "todos"
                ],
                "summary": "Delete a todo",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.ResponseData": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Data is the utils data."
                }
            }
        },
        "handler.ResponseError": {
            "type": "object",
            "properties": {
                "errors": {
                    "description": "Errors is the utils errors.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.Error"
                    }
                }
            }
        },
        "model.CreateRequest": {
            "type": "object",
            "required": [
                "priority",
                "task"
            ],
            "properties": {
                "priority": {
                    "type": "string"
                },
                "task": {
                    "type": "string"
                }
            }
        },
        "model.Status": {
            "type": "string",
            "enum": [
                "created",
                "processing",
                "done"
            ],
            "x-enum-varnames": [
                "Created",
                "Processing",
                "Done"
            ]
        },
        "model.Todo": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "priority": {
                    "$ref": "#/definitions/model.TodoPriority"
                },
                "status": {
                    "$ref": "#/definitions/model.Status"
                },
                "task": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "model.TodoPriority": {
            "type": "string",
            "enum": [
                "low",
                "medium",
                "high"
            ],
            "x-enum-varnames": [
                "TP_Low",
                "TP_Medium",
                "TP_High"
            ]
        },
        "model.UpdateRequestBody": {
            "type": "object",
            "properties": {
                "status": {
                    "$ref": "#/definitions/model.Status"
                },
                "task": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "fullstack-examination-2024 API",
	Description:      "This is a server for fullstack-examination-2024.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
