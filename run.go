package restflix

import (
	"log"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kataras/iris/v12"
)

// TODO: custom open api
func initOpenAPI() *openapi3.T {
	return &openapi3.T{
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
}

func save(openapi *openapi3.T, path string) {
	if path == "" {
		path = "./openapi.json"
	}
	data, _ := openapi.MarshalJSON()

	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(string(data))
	if err != nil {
		log.Fatal(err)
	}
}

type Options struct {
	SearchIdentifiers      []*SearchIdentifier
	StructsMappingRootPath string
	SavePath               string
	GoModName              string
	IgnoreRoutes           []string
	Require                map[string]bool
	iris                   *iris.Application
}

func (o *Options) WithIris(app *iris.Application) *Options {
	o.iris = app

	return o
}

func Init(options *Options) {
	// wait for router registration
	go func() {
		time.Sleep(time.Second * 2)

		openapi := initOpenAPI()

		if options.iris != nil {
			if err := irisRouterStrategy(options, openapi); err != nil {
				panic(err)
			}
		}

		save(openapi, options.SavePath)

		log.Println("finished openapi with iris strategy")
	}()
}
