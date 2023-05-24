package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	"github.com/pothulapati/tailscale-talk/pkg/restapi"
	"github.com/pothulapati/tailscale-talk/pkg/restapi/operations"
	tspkg "github.com/pothulapati/tailscale-talk/pkg/tailscale"
	"tailscale.com/tsnet"
)

var (
	tsKey struct {
		Key string `json:"key"`
	}
)

// This file was generated by the swagger tool.
// Make sure not to overwrite this file after you generated it because all your edits would be lost!

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewTodoListAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "A Todo list application"
	parser.LongDescription = "From the todo list tutorial on goswagger.io"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	authKey, err := tspkg.GetTodoAuthKeyFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Tailscale Server
	s := &tsnet.Server{
		Hostname: "todo-server",
		AuthKey:  authKey,
	}

	defer s.Close()

	ln, err := s.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	server.ConfigureAPI()

	// wrap OpenAI server into ln
	http.Serve(ln, server.GetHandler())
}
