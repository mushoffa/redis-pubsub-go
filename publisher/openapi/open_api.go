package openapi

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func NewOpenAPI3() openapi3.T {
	swagger := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Publisher API",
			Description: "REST APIs used for interacting with the Publisher Service",
			Version:     "0.0.0",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			// Contact: &openapi3.Contact{
			// 	URL: "https://github.com/MarioCarrion/todo-api-microservice-example",
			// },
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Local development",
				URL:         "http://localhost:9001",
			},
		},
	}

	swagger.Components.Schemas = openapi3.Schemas{
		"Publisher": openapi3.NewSchemaRef("",
			openapi3.NewObjectSchema().
				WithProperty("topic", openapi3.NewStringSchema()).
				WithProperty("data", openapi3.NewStringSchema())),
	}

	swagger.Components.RequestBodies = openapi3.RequestBodies{
		"PublisherRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for creating a task.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("topic", openapi3.NewStringSchema().
						WithMinLength(1).
						WithDefault("publisher.")).
					WithProperty("data", openapi3.NewStringSchema().
						WithMinLength(1)),
				),
		},
	}

	swagger.Components.Responses = openapi3.Responses{
		"ErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response when errors happen.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("error", openapi3.NewStringSchema()))),
		},
		"PublisherResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("OK", openapi3.NewStringSchema()))),
		},
	}

	swagger.Paths = openapi3.Paths{
		"/publish": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "Publish Topic",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/PublisherRequest",
				},
				Responses: openapi3.Responses{
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/PublisherResponse",
					},
				},
			},
		},
	}

	return swagger
}

func RegisterOpenAPI(router *gin.Engine) {
	swagger := NewOpenAPI3()

	router.GET("/openapi3.json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, &swagger)
	})

	router.GET("/openapi3.yaml", func(ctx *gin.Context) {
		data, _ := yaml.Marshal(&swagger)

		ctx.YAML(http.StatusOK, data)
	})
}
