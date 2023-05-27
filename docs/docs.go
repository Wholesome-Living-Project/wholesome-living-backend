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
        "/finance": {
            "get": {
                "description": "fetch a single investment session.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "finance"
                ],
                "summary": "Get a single investment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "investment ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start time",
                        "name": "startTime",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "end time",
                        "name": "endTime",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/finance.getInvestmentResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new investment.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "finance"
                ],
                "summary": "Create a investment.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "investment to create",
                        "name": "investment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/finance.createInvestmentRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/finance.createInvestmentResponse"
                        }
                    }
                }
            }
        },
        "/meditation": {
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
                        "name": "userId",
                        "in": "header"
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
            },
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
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header"
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
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header",
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
        "/settings": {
            "get": {
                "description": "fetch settings for a user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "settings"
                ],
                "summary": "Get settings for a user.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Plugin name",
                        "name": "plugin",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/settings.getInvestmentResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates settings for a user.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "settings"
                ],
                "summary": "Create onboarding in backend, set settings.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "onboarding to create",
                        "name": "settings",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/settings.createSettingsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/settings.createInvestmentResponse"
                        }
                    }
                }
            }
        },
        "/settings/finance": {
            "post": {
                "description": "Creates settings for a user for onr Plugin.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "settings"
                ],
                "summary": "Create settings for the finance plugin.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "onboarding to create",
                        "name": "settings",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/settings.FinanceSettings"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/settings.createInvestmentResponse"
                        }
                    }
                }
            }
        },
        "/settings/meditation": {
            "post": {
                "description": "Creates settings for a user",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "settings"
                ],
                "summary": "Create settings for the meditation Plugin.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "onboarding to create",
                        "name": "settings",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/settings.MeditationSettings"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/settings.createInvestmentResponse"
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
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "header"
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
        "finance.createInvestmentRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "investmentTime": {
                    "type": "integer"
                }
            }
        },
        "finance.createInvestmentResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "finance.getInvestmentResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "investmentTime": {
                    "type": "integer"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "meditation.createMeditationRequest": {
            "type": "object",
            "properties": {
                "endTime": {
                    "type": "integer"
                },
                "meditationTime": {
                    "type": "integer"
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
                    "type": "integer"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "settings.FinanceSettings": {
            "type": "object",
            "properties": {
                "amountNotifications": {
                    "type": "integer"
                },
                "investmentGoal": {
                    "description": "The user's investment goal.",
                    "type": "integer"
                },
                "investmentTimeGoal": {
                    "description": "The user's investment time goal.",
                    "type": "integer"
                },
                "notifications": {
                    "type": "boolean"
                },
                "periodNotifications": {
                    "$ref": "#/definitions/settings.enumNotificationPeriod"
                },
                "strategy": {
                    "$ref": "#/definitions/settings.enumStrategy"
                },
                "strategyAmount": {
                    "type": "integer"
                }
            }
        },
        "settings.MeditationSettings": {
            "type": "object",
            "properties": {
                "amountNotifications": {
                    "type": "integer"
                },
                "meditationTime": {
                    "description": "The user's meditation time goal.",
                    "type": "integer"
                },
                "notifications": {
                    "type": "boolean"
                },
                "periodNotifications": {
                    "$ref": "#/definitions/settings.enumNotificationPeriod"
                }
            }
        },
        "settings.createInvestmentResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "settings.createSettingsRequest": {
            "type": "object",
            "properties": {
                "enabledPlugins": {
                    "description": "A list with the Plugins that the user has enabled.",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "finance": {
                    "description": "The user's settings for the finance plugin.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/settings.FinanceSettings"
                        }
                    ]
                },
                "meditation": {
                    "description": "The user's settings for the meditation plugin.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/settings.MeditationSettings"
                        }
                    ]
                }
            }
        },
        "settings.enumNotificationPeriod": {
            "type": "object",
            "properties": {
                "day": {
                    "type": "boolean"
                },
                "month": {
                    "type": "boolean"
                },
                "week": {
                    "type": "boolean"
                }
            }
        },
        "settings.enumStrategy": {
            "type": "object",
            "properties": {
                "percent": {
                    "type": "boolean"
                },
                "plus": {
                    "type": "boolean"
                },
                "round": {
                    "type": "boolean"
                }
            }
        },
        "settings.getInvestmentResponse": {
            "type": "object",
            "properties": {
                "enabledPlugins": {
                    "description": "A list with the Plugins that the user has enabled.",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "finance": {
                    "description": "The user's settings for the finance plugin.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/settings.FinanceSettings"
                        }
                    ]
                },
                "meditation": {
                    "description": "The user's settings for the meditation plugin.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/settings.MeditationSettings"
                        }
                    ]
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
