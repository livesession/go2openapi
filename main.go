package main

import (
	"log"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"restflix/test/app"
)

// TODO: support methods and functions
// TODO: support recursion search
func main() {
	go func() {
		app.Init() // TODO:
	}()
	time.Sleep(time.Second)

	//operation := swag.NewOperation(nil)

	openapi := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			//Version: "", TODO: currently ls does not support versioning
			Title: "LiveSession Internal Web API",
		},
		Servers: openapi3.Servers{
			{
				ExtensionProps: openapi3.ExtensionProps{},
				URL:            "https://api.livesession.io",
				Description:    "Production",
				Variables:      nil,
			},
			{
				ExtensionProps: openapi3.ExtensionProps{},
				URL:            "https://api-labs.livesession.io/",
				Description:    "Labs",
				Variables:      nil,
			},
			{
				ExtensionProps: openapi3.ExtensionProps{},
				URL:            "http://localhost:3001",
				Description:    "Local",
				Variables:      nil,
			},
		},
		Paths:        make(openapi3.Paths),
		Security:     nil,
		Tags:         nil,
		ExternalDocs: nil,
		Components: openapi3.Components{
			RequestBodies: make(openapi3.RequestBodies),
		},
	}

	var validationIdentifier = func() []string { // TODO: more identifiers
		return []string{"baseController", "ValidateBody"}
	}

	// TODO: more router strategies
	irisRouterStrategy(app.App(), openapi, validationIdentifier) // 1. iris strategy

	data, _ := openapi.MarshalJSON() // 2. parse open api struct into json

	f, err := os.Create("./openapi.json") // 3. write into file
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(string(data))

	return
}
