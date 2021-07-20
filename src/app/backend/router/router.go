package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kore3lab/dashboard/docs"
	"github.com/kore3lab/dashboard/pkg/app"
	"github.com/kore3lab/dashboard/pkg/config"
	"github.com/kore3lab/dashboard/pkg/lang"
	"github.com/kore3lab/dashboard/router/apis"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var Router *gin.Engine

func CreateUrlMappings() {

	// swagger docs
	docs.SwaggerInfo.Title = "kore-board API"
	docs.SwaggerInfo.Description = "mulit-cluster kubernetes dashboard api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "github.com/kore3lab"
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
		authAPI.POST("/login", apis.Login)                                // authentication & issuance  access/refresh token
		authAPI.GET("/login", apis.GetAuth)                               // get auth-config
		authAPI.GET("/logout", apis.Logout)                               // logout
		authAPI.GET("/user", authenticate(), apis.GetUser)                // user info
		authAPI.POST("/token/refresh", authenticate(), apis.RefreshToken) // authentication & get access token
	}

	// contexts API
	contextsAPI := Router.Group("/api/contexts", authenticate())
	{
		contextsAPI.GET("", apis.ListContexts)                     // list contexts
		contextsAPI.POST("", apis.CreateContexts)                  // create contexts (kubeconfig yaml)
		contextsAPI.GET("/:CLUSTER", apis.GetContext)              // get a context (context meta-data : resources & namespaces)
		contextsAPI.GET("/:CLUSTER/config", apis.GetContextConfig) // get a context (server, user, context)
		contextsAPI.POST("/:CLUSTER", apis.AddContext)             // add a context (server, user, context)
		contextsAPI.DELETE("/:CLUSTER", apis.DeleteContext)        // delete a context
	}

	// custom API
	clustersAPI := Router.Group("/api/clusters/:CLUSTER", authenticate())
	{
		clustersAPI.GET("/nodes/:NAME/metrics", apis.GetNodeMetrics)                               // get node metrics
		clustersAPI.GET("/namespaces/:NAMESPACE/:RESOURCE/:NAME/metrics", apis.GetMetrics)         // get metrics (pod, deployment, statefulset, daemonset, replicaset)
		clustersAPI.GET("/nodes/:NAME/pods", apis.GetNodePodListWithMetrics)                       // get node pod list
		clustersAPI.GET("/namespaces/:NAMESPACE/:RESOURCE/:NAME/pods", apis.GetPodListWithMetrics) // get pod list (deployment, statefulset, daemonset, replicaset)
		clustersAPI.GET("/topology", apis.Topology)                                                // get cluster topology graph
		clustersAPI.GET("/topology/namespaces/:NAMESPACE", apis.Topology)                          // get namespace topology graph
		clustersAPI.GET("/dashboard", apis.Dashboard)                                              // get dashboard
	}

	// RAW-API > POST/PUT (apply, patch)
	Router.POST("/raw/clusters/:CLUSTER", authenticate(), apis.ApplyRaw)
	Router.PUT("/raw/clusters/:CLUSTER", authenticate(), apis.ApplyRaw)
	Router.POST("/raw", authenticate(), apis.ApplyRaw)
	Router.PUT("/raw", authenticate(), apis.ApplyRaw)

	// RAW-API > API-Group List
	Router.GET("/raw/clusters/:CLUSTER/apis/", authenticate(), apis.GetAPIGroupList)
	Router.GET("/raw/apis/", authenticate(), apis.GetAPIGroupList)

	// RAW-API Core
	//      non-Namespaced
	//          /api/v1/namespaces/kore
	//          /api/v1/nodes/apps-113
	//      Namespaced
	//          /api/v1/namespaces/default/services/kubernetes
	//          /api/v1/namespaces/default/serviceaccounts/default
	Router.GET("/raw/clusters/:CLUSTER/api/", authenticate(), apis.GetRaw) // Core APIVersions
	rawAPI := Router.Group("/raw/clusters/:CLUSTER/api/:VERSION", authenticate(), route())
	{
		rawAPI.GET("", apis.GetRaw)                             // ""                                       > core apiGroup - APIResourceLis
		rawAPI.GET("/:A", apis.GetRaw)                          // "/:RESOURCE"                             > core apiGroup - list
		rawAPI.GET("/:A/:B", apis.GetRaw)                       // "/:RESOURCE/:NAME"                       > core apiGroup - get
		rawAPI.DELETE("/:A/:B", apis.DeleteRaw)                 // "/:RESOURCE/:NAME"                       > core apiGroup - delete
		rawAPI.PATCH("/:A/:B", apis.PatchRaw)                   // "/:RESOURCE/:NAME"                       > core apiGroup - patch
		rawAPI.GET("/:A/:B/:RESOURCE", apis.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"       > namespaced core apiGroup - list
		rawAPI.GET("/:A/:B/:RESOURCE/:NAME", apis.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - get
		rawAPI.DELETE("/:A/:B/:RESOURCE/:NAME", apis.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - delete
		rawAPI.PATCH("/:A/:B/:RESOURCE/:NAME", apis.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME" > namespaced core apiGroup - patch
	}

	// RAW-API Grouped
	//      non-Namespaced
	//          /apis/metrics.k8s.io/v1beta1/nodes/apps-115
	//      Namespaced
	//          /apis/apps/v1/namespaces/kube-system/deployments/nginx
	//          /apis/rbac.authorization.k8s.io/v1/namespaces/default/rolebindings/clusterrolebinding-2g782
	Router.GET("/raw/clusters/:CLUSTER/apis/:GROUP", authenticate(), apis.GetRaw) // APIGroup
	rawAPIs := Router.Group("/raw/clusters/:CLUSTER/apis/:GROUP/:VERSION", authenticate(), route())
	{
		rawAPIs.GET("", apis.GetRaw)                             // ""                                          > apiGroup - APIResourceList
		rawAPIs.GET("/:A", apis.GetRaw)                          // "/:RESOURCE"                                > apiGroup - list
		rawAPIs.GET("/:A/:B", apis.GetRaw)                       // "/:RESOURCE/:NAME"                          > apiGroup - get
		rawAPIs.DELETE("/:A/:B", apis.DeleteRaw)                 // "/:RESOURCE/:NAME"                          > apiGroup - delete
		rawAPIs.PATCH("/:A/:B", apis.PatchRaw)                   // "/:RESOURCE/:NAME"                          > apiGroup - patch
		rawAPIs.GET("/:A/:B/:RESOURCE", apis.GetRaw)             // "/namespaces/:NAMESPACE/:RESOURCE"          > namespaced apiGroup - list
		rawAPIs.GET("/:A/:B/:RESOURCE/:NAME", apis.GetRaw)       // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - get
		rawAPIs.DELETE("/:A/:B/:RESOURCE/:NAME", apis.DeleteRaw) // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - delete
		rawAPIs.PATCH("/:A/:B/:RESOURCE/:NAME", apis.PatchRaw)   // "/namespaces/:NAMESPACE/:RESOURCE/:NAME"    > namespaced apiGroup - patch
	}

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
