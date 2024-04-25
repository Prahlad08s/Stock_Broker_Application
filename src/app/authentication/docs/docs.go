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
        "/validateOTP": {
            "post": {
                "description": "Validates the OTP for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OTP"
                ],
                "summary": "Validate OTP",
                "parameters": [
                    {
                        "description": "OTP Request",
                        "name": "otpRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.OTPValidationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OTP validated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "OTP is expired or invalid",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.OTPValidationRequest": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "otp": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}