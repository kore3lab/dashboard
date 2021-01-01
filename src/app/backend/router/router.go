package router

import (
	"net/http"
	"strings"

	"github.com/acornsoftlab/dashboard/docs"
	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/router/apis/_raw"
	api "github.com/acornsoftlab/dashboard/router/apis/clusters"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var Router *gin.Engine

func CreateUrlMappings() {

	// swagger docs
	docs.SwaggerInfo.Title = "acornsoft-dashboard API"
	docs.SwaggerInfo.Description = "mulit-cluster kubernetes dashboard api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "github.com/acornsoftlab"
	docs.SwaggerInfo.BasePath = "/swegger"

	// gin
	Router = gin.Default()
	Router.Use(cors())
	Router.Use(route())
	// Router.Use(gin.Logger())
	// Router.Use(Recovery())

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // restful-api docs
	Router.GET("/healthy", healthy)                                           // healthy
	Router.GET("/api/clusters", api.ListContexts)

	clustersAPI := Router.Group("/api/clusters/:CLUSTER")
	{
		clustersAPI.GET("/topology", api.Topology)
		clustersAPI.GET("/topology/namespaces/:NAMESPACE", api.Topology)
		clustersAPI.GET("/dashboard", api.Dashboard)
	}

	// RAW-API POST/PUT (apply, patch)
	Router.POST("/raw/clusters/:CLUSTER", _raw.ApplyRaw)
	Router.PUT("/raw/clusters/:CLUSTER", _raw.ApplyRaw)

	// RAW-API Core
	//      non-Namespaced
	//          /api/v1/namespaces/acornsoft-dashboard
	//          /api/v1/nodes/apps-113
	//      Namespaced
	//          /api/v1/namespaces/default/services/kubernetes
	//          /api/v1/namespaces/default/serviceaccounts/default
	rawAPI := Router.Group("/raw/clusters/:CLUSTER/api/:VERSION")
	{
		rawAPI.GET("/:A", _raw.GetRaw)                          // "/:RESOURCE"                             > core apiGroup - list
		rawAPI.GET("/:A/:B", _raw.GetRaw)                       // "/:RESOURCE/:NAME"                       > core apiGroup - get
		rawAPI.DELETE("/:A/:B", _raw.DeleteRaw)                 // "/:RESOURCE/:NAME"                       > core apiGroup - delete
		rawAPI.PATCH("/:A/:B", _raw.PatchRaw)                   // "/:RESOURCE/:NAME"                       > core apiGroup - patch
		rawAPI.GET("/:A/:B/:RESOURCE", _raw.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"       > namespaced core apiGroup - list
		rawAPI.GET("/:A/:B/:RESOURCE/:NAME", _raw.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - get
		rawAPI.DELETE("/:A/:B/:RESOURCE/:NAME", _raw.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - delete
		rawAPI.PATCH("/:A/:B/:RESOURCE/:NAME", _raw.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - patch
	}

	// RAW-API Grouped
	//      non-Namespaced
	//          /apis/metrics.k8s.io/v1beta1/nodes/apps-115
	//      Namespaced
	//          /apis/apps/v1/namespaces/kube-system/deployments/nginx
	//          /apis/rbac.authorization.k8s.io/v1/namespaces/default/rolebindings/clusterrolebinding-2g782
	rawAPIs := Router.Group("/raw/clusters/:CLUSTER/apis/:GROUP/:VERSION")
	{
		rawAPIs.GET("/:A", _raw.GetRaw)                          // "/:RESOURCE"                                > apiGroup - list
		rawAPIs.GET("/:A/:B", _raw.GetRaw)                       // "/:RESOURCE/:NAME"                          > apiGroup - get
		rawAPIs.DELETE("/:A/:B", _raw.DeleteRaw)                 // "/:RESOURCE/:NAME"                          > apiGroup - delete
		rawAPIs.PATCH("/:A/:B", _raw.PatchRaw)                   // "/:RESOURCE/:NAME"                          > apiGroup - patch
		rawAPIs.GET("/:A/:B/:RESOURCE", _raw.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"          > namespaced apiGroup - list
		rawAPIs.GET("/:A/:B/:RESOURCE/:NAME", _raw.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - get
		rawAPIs.DELETE("/:A/:B/:RESOURCE/:NAME", _raw.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - delete
		rawAPIs.PATCH("/:A/:B/:RESOURCE/:NAME", _raw.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - patch
	}

}

/**
  RAW-API  URL Route resolved handler
*/
func route() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/raw/clusters") {
			if c.Param("RESOURCE") == "" {
				c.Params = append(c.Params,
					gin.Param{Key: "RESOURCE", Value: c.Param("A")},
					gin.Param{Key: "NAME", Value: c.Param("B")})
			} else if c.Param("A") == "namespaces" {
				c.Params = append(c.Params, gin.Param{Key: "NAMESPACE", Value: c.Param("B")})
			}
		}
	}
}

func cors() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func healthy(c *gin.Context) {
	g := app.Gin{C: c}
	g.Send(http.StatusOK, "healthy")
}
