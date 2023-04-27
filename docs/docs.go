// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Wholesome Living"
        },
        "license": {
            "name": "MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/meditation": {
            "post": {
                "description": "Creates a new meditation.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "meditation"
                ],
                "summary": "Create meditation.",
                "parameters": [
                    {
                        "description": "Meditation to create",
                        "name": "meditation",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/meditation.createMeditationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/meditation.createMeditationResponse"
                        }
                    }
                }
            }
        },
        "/meditation/getAll/{userID}": {
            "get": {
                "description": "fetch all meditation sessions of a user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "meditation"
                ],
                "summary": "Get all meditation session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "endTime": {
                                        "type": "integer"
                                    },
                                    "id": {
                                        "type": "string"
                                    },
                                    "meditationTime": {
                                        "type": "integer"
                                    },
                                    "userId": {
                                        "type": "string"
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "/meditation/{id}": {
            "get": {
                "description": "fetch a single meditation session.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "meditation"
                ],
                "summary": "Get a meditation session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Meditation ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/meditation.getMeditationResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "fetch every user available.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get all users.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/user.UserDB"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "update a user by id.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update a user.",
                "parameters": [
                    {
                        "description": "User to update",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.updateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.UserDB"
                        }
                    }
                }
            },
            "post": {
                "description": "creates one user.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create one user.",
                "parameters": [
                    {
                        "description": "User to create",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.createUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.createUserResponse"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "fetch a user by id.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.UserDB"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "meditation.createMeditationRequest": {
            "type": "object",
            "properties": {
                "meditationTime": {
                    "type": "integer"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "meditation.createMeditationResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "meditation.getMeditationResponse": {
            "type": "object",
            "properties": {
                "endTime": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "meditationTime": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "user.UserDB": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "integer"
                },
                "dateOfBirth": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "plugins": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.pluginType"
                    }
                }
            }
        },
        "user.createUserRequest": {
            "type": "object",
            "properties": {
                "dateOfBirth": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                }
            }
        },
        "user.createUserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "user.pluginType": {
            "type": "string",
            "enum": [
                "meditation",
                "workout"
            ],
            "x-enum-varnames": [
                "PluginTypeMeditation",
                "PluginTypeWorkout"
            ]
        },
        "user.updateUserRequest": {
            "type": "object",
            "properties": {
                "dateOfBirth": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "plugins": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.pluginType"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Wholesome Living Backend",
	Description:      "A backend for Wholesome Living written in Golang backend API using Fiber and MongoDB",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
