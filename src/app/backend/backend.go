// main.go

package main

import (
	flag "github.com/spf13/pflag"

	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/acornsoftlab/dashboard/router"
)

func main() {
	var kubeconfig *string

	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	flag.Parse()

	config.Setup(*kubeconfig)
	router.CreateUrlMappings()
	router.Router.Run(":3001")

}
