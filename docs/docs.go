// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/announcements": {
            "post": {
                "description": "Create an announcement and save it to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Announcements"
                ],
                "summary": "Create an announcement",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Announcement object",
                        "name": "announcement",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Announcement"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Announcement created successfully",
                        "schema": {
                            "$ref": "#/definitions/utils.CreateAnnouncementSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Could not parse request body",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Authorization token is required or invalid",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "Authenticate a user and return a JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User login credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.LoginData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User logged in successfully with JWT token",
                        "schema": {
                            "$ref": "#/definitions/utils.LoginSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request - could not parse the request or generate token",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized - invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error - server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/signup": {
            "post": {
                "description": "Create a new user in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Sign up a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created successfully",
                        "schema": {
                            "$ref": "#/definitions/utils.UserSuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict - user already exists",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{email}": {
            "get": {
                "description": "Get a user by their email address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Retrieve user by email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email",
                        "name": "email",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/utils.UserSuccessResponse"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Announcement": {
            "type": "object",
            "properties": {
                "create_date": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "owner_id": {
                    "type": "integer"
                },
                "start_date": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/models.Status"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "models.Status": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3,
                4
            ],
            "x-enum-comments": {
                "Accepted": "1",
                "Active": "3",
                "Deactivated": "4",
                "Declined": "2",
                "Pending": "0"
            },
            "x-enum-varnames": [
                "Pending",
                "Accepted",
                "Declined",
                "Active",
                "Deactivated"
            ]
        },
        "models.User": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                }
            }
        },
        "utils.CreateAnnouncementSuccessResponse": {
            "type": "object",
            "properties": {
                "announcement": {
                    "$ref": "#/definitions/models.Announcement"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "utils.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "description": "The error message",
                    "type": "string"
                }
            }
        },
        "utils.LoginData": {
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
        "utils.LoginSuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "utils.UserSuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "description": "The success message",
                    "type": "string"
                },
                "user": {
                    "description": "The user data",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.User"
                        }
                    ]
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "AnnounceIT API",
	Description:      "AnnounceIT is a solution for broadcasting agencies to receive and manage announcements effectively.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
