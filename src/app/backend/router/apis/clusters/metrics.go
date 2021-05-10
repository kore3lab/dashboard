package apis

import (
	"fmt"
	"net/http"

	resty "github.com/go-resty/resty/v2"

	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/acornsoftlab/dashboard/pkg/lang"
	"github.com/gin-gonic/gin"
)

// Get node metrics
func GetNodeMetrics(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Value.DefaultContext)

	// invoke metrics-scraper api
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("%s/api/v1/clusters/%s/nodes/%s/metrics/%s", config.Value.MetricsScraperUrl, cluster, c.Param("NODE"), c.Param("METRICS")))
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to get scrapping metrics", err)
	} else {
		g.C.Data(http.StatusOK, "application/json", resp.Body())
	}

}

// Get pod metrics
func GetPodMetrics(c *gin.Context) {
	g := app.Gin{C: c}

	cluster := lang.NVL(g.C.Param("CLUSTER"), config.Value.DefaultContext)

	// invoke metrics-scraper api
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("%s/api/v1/clusters/%s/namespaces/%s/pods/%s/metrics/%s", config.Value.MetricsScraperUrl, cluster, c.Param("NAMESPACE"), c.Param("POD"), c.Param("METRICS")))
	if err != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to get scrapping metrics", err)
	} else {
		g.C.Data(http.StatusOK, "application/json", resp.Body())
	}

}
