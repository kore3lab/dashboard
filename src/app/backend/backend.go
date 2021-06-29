// main.go

package main

import (
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/router"
)

func main() {

	config.Setup()
	router.CreateUrlMappings()
	router.Router.Run(":3001")

}
