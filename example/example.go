package main

import (
	// This is the path to Gococ
	"fmt"

	"github.com/Khulnasoft-lab/gococ"
)

func main() {

	// Build a new Gococ app.
	registry, router, context := gococ.Gococ()

	// Fill the registry.
	registry.AddRoutes(
		gococ.Route{
			Name: "TEST",
			Help: "A test route",
			Does: gococ.Tasks{
				gococ.Cmd{
					Name: "hi",
					Fn:   HelloWorld,
				},
			},
		},
	)

	// Execute the route.
	router.HandleRequest("TEST", context, false)
}

func HelloWorld(cxt gococ.Context, params *gococ.Params) (interface{}, gococ.Interrupt) {
	fmt.Println("Hello World")
	return true, nil
}
