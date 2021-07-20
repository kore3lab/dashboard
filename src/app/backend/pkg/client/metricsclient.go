package client

import (
	"fmt"
	"strings"

	resty "github.com/go-resty/resty/v2"
)

const (
	METRICS_RESOURCE_CLUSTER = 1
	METRICS_RESOURCE_NODE    = 2
	METRICS_RESOURCE_POD     = 3
)

// type resourceVerber struct {
type CumulativeMetricsClient struct {
	metricsScraperUrl string
	context           string
}

type CumulativeMetricsResourceSelector struct {
	Node      string
	Namespace string
	Pods      []string
	Function  string
}

func (self *CumulativeMetricsResourceSelector) getUrl() string {
	if len(self.Pods) > 0 {
		return fmt.Sprintf("/namespaces/%s/pods/%s", self.Namespace, strings.Join(self.Pods, ","))
	} else if self.Node != "" {
		return fmt.Sprintf("/nodes/%s", self.Node)
	} else {
		return ""
	}
}

// RestfulClient 리턴
func NewCumulativeMetricsClient(metricsScraperUrl string, ctx string) *CumulativeMetricsClient {
	return &CumulativeMetricsClient{metricsScraperUrl: metricsScraperUrl, context: ctx}
}

func (self *CumulativeMetricsClient) Get(selector CumulativeMetricsResourceSelector) ([]interface{}, error) {

	result := []interface{}{}

	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("%s/api/v1/clusters/%s%s", self.metricsScraperUrl, self.context, selector.getUrl()))
	if err != nil {
		return nil, err
	}

	return result, nil
}
