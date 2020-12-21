// main.go

package main

import (
	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/acornsoftlab/dashboard/router"
)

func main() {

	config.Setup()
	router.CreateUrlMappings()
	router.Router.Run(":3001")

}
