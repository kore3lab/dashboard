package router

import (
	"net/http"

	"github.com/acornsoftlab/dashboard/docs"
	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/acornsoftlab/dashboard/pkg/lang"
	"github.com/acornsoftlab/dashboard/router/apis/_raw"
	api "github.com/acornsoftlab/dashboard/router/apis/clusters"
	"github.com/acornsoftlab/dashboard/router/apis/termapi"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var Router *gin.Engine

func CreateUrlMappings() {

	// swagger docs
	docs.SwaggerInfo.Title = "kore-board API"
	docs.SwaggerInfo.Description = "mulit-cluster kubernetes dashboard api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "github.com/acornsoftlab"
	docs.SwaggerInfo.BasePath = "/swegger"

	// gin
	Router = gin.Default()
	Router.Use(cors()) // cors
	authenticate := config.Authenticator.HandlerFunc
	// Router.Use(gin.Logger())
	// Router.Use(Recovery())

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // restful-api docs
	Router.GET("/healthy", healthy)                                           // healthy

	// Authentication API
	authAPI := Router.Group("/api/auth")
	{
		authAPI.POST("/login", api.Login)                                // authentication & issuance  access/refresh token
		authAPI.GET("/login", api.GetAuth)                               // get auth-config
		authAPI.GET("/logout", api.Logout)                               // logout
		authAPI.GET("/user", authenticate(), api.GetUser)                // user info
		authAPI.POST("/token/refresh", authenticate(), api.RefreshToken) // authentication & get access token
	}

	// contexts API
	contextsAPI := Router.Group("/api/contexts", authenticate())
	{
		contextsAPI.GET("", api.ListContexts)                     // list contexts
		contextsAPI.POST("", api.CreateContexts)                  // create contexts (kubeconfig yaml)
		contextsAPI.GET("/:CLUSTER", api.GetContext)              // get a context (context meta-data : resources & namespaces)
		contextsAPI.GET("/:CLUSTER/config", api.GetContextConfig) // get a context (server, user, context)
		contextsAPI.POST("/:CLUSTER", api.AddContext)             // add a context (server, user, context)
		contextsAPI.DELETE("/:CLUSTER", api.DeleteContext)        // delete a context
	}

	// clusters API
	clustersAPI := Router.Group("/api/clusters/:CLUSTER", authenticate())
	{
		clustersAPI.GET("/nodes/:NODE/metrics/:METRICS", api.GetNodeMetrics)                    // get node metrics
		clustersAPI.GET("/namespaces/:NAMESPACE/pods/:POD/metrics/:METRICS", api.GetPodMetrics) // get pod metrics
		clustersAPI.GET("/topology", api.Topology)                                              // get cluster topology graph
		clustersAPI.GET("/topology/namespaces/:NAMESPACE", api.Topology)                        // get namespace topology graph
		clustersAPI.GET("/dashboard", api.Dashboard)                                            // get dashboard
		//terminal API
		clustersAPI.GET("/terminal", termapi.ProcCluster)
		clustersAPI.GET("/namespaces/:NAMESPACE/pods/:POD/terminal", termapi.ProcPod)
		clustersAPI.GET("/namespaces/:NAMESPACE/pods/:POD/containers/:CONTAINER/terminal", termapi.ProcContainer)
	}
	//for terminal websocket connect
	Router.GET("/api/terminal/ws", termapi.GenerateHandleWS)

	//// ROOT API
	//API := Router.Group("/api", authenticate())
	//{
	//	API.GET("/topology", api.Topology)
	//	API.GET("/topology/namespaces/:NAMESPACE", api.Topology)
	//	API.GET("/dashboard", api.Dashboard)
	//}

	// RAW-API > POST/PUT (apply, patch)
	Router.POST("/raw/clusters/:CLUSTER", authenticate(), _raw.ApplyRaw)
	Router.PUT("/raw/clusters/:CLUSTER", authenticate(), _raw.ApplyRaw)
	Router.POST("/raw", authenticate(), _raw.ApplyRaw)
	Router.PUT("/raw", authenticate(), _raw.ApplyRaw)

	// RAW-API > API-Group List
	Router.GET("/raw/clusters/:CLUSTER/apis/", authenticate(), _raw.GetAPIGroupList)
	Router.GET("/raw/apis/", authenticate(), _raw.GetAPIGroupList)

	// RAW-API Core
	//      non-Namespaced
	//          /api/v1/namespaces/kore
	//          /api/v1/nodes/apps-113
	//      Namespaced
	//          /api/v1/namespaces/default/services/kubernetes
	//          /api/v1/namespaces/default/serviceaccounts/default
	Router.GET("/raw/clusters/:CLUSTER/api/", authenticate(), _raw.GetRaw) // Core APIVersions
	rawAPI := Router.Group("/raw/clusters/:CLUSTER/api/:VERSION", authenticate(), route())
	{
		rawAPI.GET("", _raw.GetRaw)                             // ""                                       > core apiGroup - APIResourceLis
		rawAPI.GET("/:A", _raw.GetRaw)                          // "/:RESOURCE"                             > core apiGroup - list
		rawAPI.GET("/:A/:B", _raw.GetRaw)                       // "/:RESOURCE/:NAME"                       > core apiGroup - get
		rawAPI.DELETE("/:A/:B", _raw.DeleteRaw)                 // "/:RESOURCE/:NAME"                       > core apiGroup - delete
		rawAPI.PATCH("/:A/:B", _raw.PatchRaw)                   // "/:RESOURCE/:NAME"                       > core apiGroup - patch
		rawAPI.GET("/:A/:B/:RESOURCE", _raw.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"       > namespaced core apiGroup - list
		rawAPI.GET("/:A/:B/:RESOURCE/:NAME", _raw.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - get
		rawAPI.DELETE("/:A/:B/:RESOURCE/:NAME", _raw.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - delete
		rawAPI.PATCH("/:A/:B/:RESOURCE/:NAME", _raw.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - patch
	}
	//Router.GET("/raw/api/", authenticate(), _raw.GetRaw) // Core APIVersions
	//rawAPI_d := Router.Group("/raw/api/:VERSION", authenticate(), route())
	//{
	//	rawAPI_d.GET("", _raw.GetRaw)                             // ""                                       > core apiGroup - APIResourceList
	//	rawAPI_d.GET("/:A", _raw.GetRaw)                          // "/:RESOURCE"                             > core apiGroup - list
	//	rawAPI_d.GET("/:A/:B", _raw.GetRaw)                       // "/:RESOURCE/:NAME"                       > core apiGroup - get
	//	rawAPI_d.DELETE("/:A/:B", _raw.DeleteRaw)                 // "/:RESOURCE/:NAME"                       > core apiGroup - delete
	//	rawAPI_d.PATCH("/:A/:B", _raw.PatchRaw)                   // "/:RESOURCE/:NAME"                       > core apiGroup - patch
	//	rawAPI_d.GET("/:A/:B/:RESOURCE", _raw.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"       > namespaced core apiGroup - list
	//	rawAPI_d.GET("/:A/:B/:RESOURCE/:NAME", _raw.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - get
	//	rawAPI_d.DELETE("/:A/:B/:RESOURCE/:NAME", _raw.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - delete
	//	rawAPI_d.PATCH("/:A/:B/:RESOURCE/:NAME", _raw.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - patch
	//}

	// RAW-API Grouped
	//      non-Namespaced
	//          /apis/metrics.k8s.io/v1beta1/nodes/apps-115
	//      Namespaced
	//          /apis/apps/v1/namespaces/kube-system/deployments/nginx
	//          /apis/rbac.authorization.k8s.io/v1/namespaces/default/rolebindings/clusterrolebinding-2g782
	Router.GET("/raw/clusters/:CLUSTER/apis/:GROUP", authenticate(), _raw.GetRaw) // APIGroup
	rawAPIs := Router.Group("/raw/clusters/:CLUSTER/apis/:GROUP/:VERSION", authenticate(), route())
	{
		rawAPIs.GET("", _raw.GetRaw)                             // ""                                          > apiGroup - APIResourceList
		rawAPIs.GET("/:A", _raw.GetRaw)                          // "/:RESOURCE"                                > apiGroup - list
		rawAPIs.GET("/:A/:B", _raw.GetRaw)                       // "/:RESOURCE/:NAME"                          > apiGroup - get
		rawAPIs.DELETE("/:A/:B", _raw.DeleteRaw)                 // "/:RESOURCE/:NAME"                          > apiGroup - delete
		rawAPIs.PATCH("/:A/:B", _raw.PatchRaw)                   // "/:RESOURCE/:NAME"                          > apiGroup - patch
		rawAPIs.GET("/:A/:B/:RESOURCE", _raw.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"          > namespaced apiGroup - list
		rawAPIs.GET("/:A/:B/:RESOURCE/:NAME", _raw.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - get
		rawAPIs.DELETE("/:A/:B/:RESOURCE/:NAME", _raw.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - delete
		rawAPIs.PATCH("/:A/:B/:RESOURCE/:NAME", _raw.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - patch
	}
	//Router.GET("/raw/apis/:GROUP", authenticate(), _raw.GetRaw) // APIGroup
	//rawAPIs_d := Router.Group("/raw/apis/:GROUP/:VERSION", authenticate(), route())
	//{
	//	rawAPIs_d.GET("", _raw.GetRaw)                             // ""                                          > apiGroup - APIResourceList
	//	rawAPIs_d.GET("/:A", _raw.GetRaw)                          // "/:RESOURCE"                                > apiGroup - list
	//	rawAPIs_d.GET("/:A/:B", _raw.GetRaw)                       // "/:RESOURCE/:NAME"                          > apiGroup - get
	//	rawAPIs_d.DELETE("/:A/:B", _raw.DeleteRaw)                 // "/:RESOURCE/:NAME"                          > apiGroup - delete
	//	rawAPIs_d.PATCH("/:A/:B", _raw.PatchRaw)                   // "/:RESOURCE/:NAME"                          > apiGroup - patch
	//	rawAPIs_d.GET("/:A/:B/:RESOURCE", _raw.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"          > namespaced apiGroup - list
	//	rawAPIs_d.GET("/:A/:B/:RESOURCE/:NAME", _raw.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - get
	//	rawAPIs_d.DELETE("/:A/:B/:RESOURCE/:NAME", _raw.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - delete
	//	rawAPIs_d.PATCH("/:A/:B/:RESOURCE/:NAME", _raw.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - patch
	//}

}

/**
  RAW-API  URL Route resolved handler
*/
func route() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("RESOURCE") == "" {
			c.Params = append(c.Params,
				gin.Param{Key: "RESOURCE", Value: c.Param("A")},
				gin.Param{Key: "NAME", Value: c.Param("B")})
		} else if c.Param("A") == "namespaces" {
			c.Params = append(c.Params, gin.Param{Key: "NAMESPACE", Value: c.Param("B")})
		}
	}
}

func cors() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", lang.NVL(c.Request.Header.Get("Origin"), "*"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func healthy(c *gin.Context) {
	g := app.Gin{C: c}
	g.Send(http.StatusOK, "healthy")
}
