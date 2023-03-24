// Code generated by swaggo/swag. DO NOT EDIT.

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
            "email": "artawower@protonmail.com"
        },
        "license": {
            "name": "GPL 3.0",
            "url": "https://www.gnu.org/licenses/gpl-3.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/api-tokens": {
            "get": {
                "description": "Return all available API tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Get API tokens",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-array_models_APIToken-any"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/auth/github/callback": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Callback for github OAuth",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/auth/github/login": {
            "get": {
                "description": "Entrypoint for login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "provider",
                        "name": "provider",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-handlers_OAuthRedirectData-any"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/auth/token": {
            "post": {
                "description": "Create API token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Create API token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-models_APIToken-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete API token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Delete API token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/auth/verify": {
            "get": {
                "description": "Return found user by provided token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-models_PublicUser-any"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/notes/": {
            "get": {
                "description": "Get all notes with optional filter",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Get notes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search text",
                        "name": "searchText",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit for pagination",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-array_models_Note-models_Pagination"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            },
            "post": {
                "description": "Create note",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Create note",
                "parameters": [
                    {
                        "description": "Note model",
                        "name": "note",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Note"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/notes/bulk-upsert": {
            "put": {
                "description": "Bulk update or insert notes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Upsert notes",
                "parameters": [
                    {
                        "description": "Notes list",
                        "name": "notes",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Note"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/notes/graph": {
            "get": {
                "description": "Return graph model with links between connected notes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Get notes graph",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-models_NoteGraph-any"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/notes/{id}": {
            "get": {
                "description": "get note by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Get note",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Note ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-models_Note-any"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        },
        "/tags": {
            "get": {
                "description": "Return list of al registered tags",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tags"
                ],
                "summary": "Get tags",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-array_string-any"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpError-any"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.HttpError-any": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "handlers.HttpResponse-array_models_APIToken-any": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.APIToken"
                    }
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-array_models_Note-models_Pagination": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Note"
                    }
                },
                "meta": {
                    "$ref": "#/definitions/models.Pagination"
                }
            }
        },
        "handlers.HttpResponse-array_string-any": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-handlers_OAuthRedirectData-any": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/handlers.OAuthRedirectData"
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-models_APIToken-any": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.APIToken"
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-models_Note-any": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.Note"
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-models_NoteGraph-any": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.NoteGraph"
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-models_PublicUser-any": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.PublicUser"
                },
                "meta": {}
            }
        },
        "handlers.OAuthRedirectData": {
            "type": "object",
            "properties": {
                "redirectUrl": {
                    "type": "string"
                }
            }
        },
        "models.APIToken": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "permission": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.GraphNoteLink": {
            "type": "object",
            "properties": {
                "source": {
                    "type": "string"
                },
                "target": {
                    "type": "string"
                }
            }
        },
        "models.GraphNoteNode": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "weight": {
                    "type": "integer"
                }
            }
        },
        "models.Note": {
            "type": "object",
            "properties": {
                "authorId": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "likes": {
                    "type": "integer"
                },
                "meta": {
                    "$ref": "#/definitions/models.NoteMeta"
                },
                "updatedAt": {
                    "type": "string"
                },
                "views": {
                    "type": "integer"
                }
            }
        },
        "models.NoteGraph": {
            "type": "object",
            "properties": {
                "links": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GraphNoteLink"
                    }
                },
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.GraphNoteNode"
                    }
                }
            }
        },
        "models.NoteHeading": {
            "type": "object",
            "properties": {
                "level": {
                    "type": "integer"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "models.NoteLink": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "models.NoteMeta": {
            "type": "object",
            "properties": {
                "category": {
                    "$ref": "#/definitions/models.category"
                },
                "description": {
                    "type": "string"
                },
                "externalLinks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.NoteLink"
                    }
                },
                "fileTags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "headings": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.NoteHeading"
                    }
                },
                "images": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "linkedArticles": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.NoteLink"
                    }
                },
                "previewImg": {
                    "type": "string"
                },
                "published": {
                    "type": "boolean"
                },
                "startup": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Pagination": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "models.PublicUser": {
            "type": "object",
            "properties": {
                "avatarUrl": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "nickName": {
                    "type": "string"
                },
                "profileUrl": {
                    "type": "string"
                }
            }
        },
        "models.category": {
            "type": "string",
            "enum": [
                "article",
                "book",
                "schedule"
            ],
            "x-enum-varnames": [
                "CategoryArticle",
                "CategoryBook",
                "CategorySchedule"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Second Brain API",
	Description:      "List of methods for work with second brain.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
