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

	clustersAPI := Router.Group("/api/clusters")
	{
		clustersAPI.GET("/", api.ListContexts)
		clustersAPI.GET("/:CLUSTER/topology", api.Topology)
		clustersAPI.GET("/:CLUSTER/topology/namespaces/:NAMESPACE", api.Topology)
		clustersAPI.GET("/:CLUSTER/dashboard", api.Dashboard)
	}
	rawAPI := Router.Group("/api/_raw")
	{
		rawAPI.POST("/:CLUSTER", _raw.ApplyRaw)
		rawAPI.PUT("/:CLUSTER", _raw.ApplyRaw)
		rawAPI.GET("/:CLUSTER/:GROUP/:VERSION/*PATH", _raw.GetRaw)
		rawAPI.DELETE("/:CLUSTER/:GROUP/:VERSION/*PATH", _raw.DeleteRaw)
	}

}

/**
  Route URL "PATH" resolved handler
      "/:CLUSTER/:GROUP/:VERSION/:RESOURCETYPE/:NAME"
      "/:CLUSTER/:GROUP/:VERSION/:RESOURCETYPE"
      "/:CLUSTER/:GROUP/:VERSION/namespaces/:NAMESPACE/:RESOURCETYPE/:NAME"
      "/:CLUSTER/:GROUP/:VERSION/namespaces/:NAMESPACE/:RESOURCETYPE"
*/
func route() gin.HandlerFunc {

	return func(c *gin.Context) {

		if strings.HasPrefix(c.Request.RequestURI, "/api/_raw") && c.Param("PATH") != "" {

			path := c.Param("PATH")
			arr := strings.Split(path, "/")
			if len(arr) >= 2 {
				namespace := ""
				kind := ""
				name := ""
				if arr[1] == "namespaces" && len(arr) > 2 {
					namespace = arr[2]
					kind = arr[3]
					if len(arr) > 4 {
						name = arr[4]
					}
				} else {
					kind = arr[1]
					if len(arr) > 2 {
						name = arr[2]
					}
				}

				c.Params = append(c.Params,
					gin.Param{Key: "NAMESPACE", Value: namespace},
					gin.Param{Key: "RESOURCETYPE", Value: kind},
					gin.Param{Key: "NAME", Value: name})
				c.Next()

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
