// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/ronilsonalves/go-wallet-watcher/blob/main/LICENSE.md",
        "contact": {
            "name": "Ronilson Alves",
            "url": "https://www.linkedin.com/in/ronilsonalves"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/ronilsonalves/go-wallet-watcher/blob/main/LICENSE.md"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/eth/wallets/{address}": {
            "get": {
                "description": "Get wallet info from an address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallets"
                ],
                "summary": "Get a wallet info and balance from a wallet address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Wallet address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/wallet.DTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.errResponse"
                        }
                    }
                }
            }
        },
        "/eth/wallets/{address}/transactions": {
            "get": {
                "description": "Get wallet info from an address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallets"
                ],
                "summary": "Retrieves up to 10000 transactions by given adrress in a paggeable response",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Wallet address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Items per page",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/web.PageableResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.errResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "wallet.DTO": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "balance": {
                    "type": "number"
                }
            }
        },
        "web.PageableResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "items": {
                    "type": "string"
                },
                "page": {
                    "type": "string"
                }
            }
        },
        "web.errResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "status-code": {
                    "type": "integer"
                },
                "time-stamp": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Wallet Watcher API",
	Description:      "This API handle query to check crypto wallet info such as balance, transactions...",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
