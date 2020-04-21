// main.go

package main

import (
	"github.com/acornsoftlab/kore3/pkg/config"
	"github.com/acornsoftlab/kore3/router"
)

func main() {

	config.Setup()
	router.CreateUrlMappings()
	router.Router.Run(":3001")

}
