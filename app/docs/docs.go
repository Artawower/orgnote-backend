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
            }
        },
        "/auth/token/{tokenId}": {
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
                "parameters": [
                    {
                        "type": "string",
                        "description": "token id",
                        "name": "tokenId",
                        "in": "path",
                        "required": true
                    }
                ],
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
        "/auth/{provider}/callback": {
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
                "summary": "Callback for OAuth",
                "parameters": [
                    {
                        "type": "string",
                        "description": "provider",
                        "name": "provider",
                        "in": "path",
                        "required": true
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
        "/auth/{provider}/login": {
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
                "summary": "OAuth Login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "provider",
                        "name": "provider",
                        "in": "path",
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
        "/files/upload": {
            "post": {
                "description": "Upload files.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Upload files",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "files",
                        "name": "files",
                        "in": "formData",
                        "required": true
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
        "/notes": {
            "delete": {
                "description": "Mark notes as deleted by provided list of ids",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Delete notes",
                "parameters": [
                    {
                        "description": "List of ids of deleted notes",
                        "name": "ids",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-any-any"
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
                        "type": "integer",
                        "x-order": "1",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "x-order": "2",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "x-order": "3",
                        "description": "User id of which notes to load",
                        "name": "userId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "x-order": "4",
                        "name": "searchText",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "x-order": "5",
                        "description": "Load all my own notes (user will be used from provided token)",
                        "name": "my",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "x-order": "6",
                        "name": "from",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-array_models_PublicNote-models_Pagination"
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
                            "$ref": "#/definitions/handlers.CreatingNote"
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
                        "description": "List of crated notes",
                        "name": "notes",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.CreatingNote"
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
        "/notes/sync": {
            "post": {
                "description": "Synchronize notes with specific timestamp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notes"
                ],
                "summary": "Synchronize notes",
                "parameters": [
                    {
                        "description": "Sync notes request",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.SyncNotesRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.HttpResponse-array_models_PublicNote-any"
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
                            "$ref": "#/definitions/handlers.HttpResponse-models_PublicNote-any"
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
        "handlers.CreatingNote": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "filePath": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "meta": {
                    "$ref": "#/definitions/models.NoteMeta"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "handlers.HttpError-any": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "handlers.HttpResponse-any-any": {
            "type": "object",
            "properties": {
                "data": {},
                "meta": {}
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
        "handlers.HttpResponse-array_models_PublicNote-any": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PublicNote"
                    }
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-array_models_PublicNote-models_Pagination": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PublicNote"
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
        "handlers.HttpResponse-models_NoteGraph-any": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.NoteGraph"
                },
                "meta": {}
            }
        },
        "handlers.HttpResponse-models_PublicNote-any": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.PublicNote"
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
        "handlers.SyncNotesRequest": {
            "type": "object",
            "properties": {
                "notes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.CreatingNote"
                    }
                },
                "timestamp": {
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
        "models.PublicNote": {
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/models.PublicUser"
                },
                "content": {
                    "type": "string"
                },
                "filePath": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "meta": {
                    "$ref": "#/definitions/models.NoteMeta"
                },
                "updatedAt": {
                    "type": "string"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
