package provider

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gorilla/mux"

	"github.com/kubernetes-sigs/dashboard-metrics-scraper/pkg/config"
)

// DashboardRouter defines the usable API routes
func DashboardRouter(r *mux.Router, db *sql.DB) {
	r.Path("/dashboard/nodes/{Name}/metrics/{MetricName}/{Whatever}").HandlerFunc(nodeHandler(db))
	r.Path("/dashboard/namespaces/{Namespace}/pod-list/{Name}/metrics/{MetricName}/{Whatever}").HandlerFunc(podHandler(db))
	// customized by kore-board
	r.Path("/clusters/{Cluster}/nodes/{Name}/metrics/{MetricName}").HandlerFunc(nodeHandler(db))
	r.Path("/clusters/{Cluster}/namespaces/{Namespace}/pods/{Name}/metrics/{MetricName}").HandlerFunc(podHandler(db))
	r.Path("/nodes/{Name}/metrics/{MetricName}").HandlerFunc(nodeHandler(db))
	r.Path("/namespaces/{Namespace}/pods/{Name}/metrics/{MetricName}").HandlerFunc(podHandler(db))
	// --END
	// r.PathPrefix("/").HandlerFunc(defaultHandler)    // by kore-board
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("%v - URL: %s", time.Now(), r.URL)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Errorf("Error cannot write response: %v", err)
	}
}

func nodeHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// customized by kore-board
		clsuter := vars["Cluster"]
		if clsuter == "" {
			clsuter = config.Value.DefaultContext
		}
		resp, err := getNodeMetrics(db, clsuter, vars["MetricName"], ResourceSelector{
			Namespace:    "",
			ResourceName: vars["Name"],
		})
		// --END

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(fmt.Sprintf("Node Metrics Error - %v", err.Error())))
			if err != nil {
				log.Errorf("Error cannot write response: %v", err)
			}
		}

		j, err := json.Marshal(resp)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(fmt.Sprintf("JSON Error - %v", err.Error())))
			if err != nil {
				log.Errorf("Error cannot write response: %v", err)
			}
		}

		_, err = w.Write(j)
		if err != nil {
			log.Errorf("Error cannot write response: %v", err)
		}
	}

	return fn
}

func podHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// customized by kore-board
		clsuter := vars["Cluster"]
		if clsuter == "" {
			clsuter = config.Value.DefaultContext
		}
		resp, err := getPodMetrics(db, clsuter, vars["MetricName"], ResourceSelector{
			Namespace:    vars["Namespace"],
			ResourceName: vars["Name"],
		})
		// --END

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(fmt.Sprintf("Pod Metrics Error - %v", err.Error())))
			if err != nil {
				log.Errorf("Error cannot write response: %v", err)
			}
		}

		j, err := json.Marshal(resp)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(fmt.Sprintf("JSON Error - %v", err.Error())))
			if err != nil {
				log.Errorf("Error cannot write response: %v", err)
			}
		}

		_, err = w.Write(j)
		if err != nil {
			log.Errorf("Error cannot write response: %v", err)
		}
	}

	return fn
}

func getRows(db *sql.DB, cluster string, table string, metricName string, selector ResourceSelector) (*sql.Rows, error) { // customized by kore-board
	var query string
	var values []interface{}
	var args []string
	orderBy := []string{"name", "time"}
	if metricName == "cpu" {
		query = fmt.Sprintf("select sum(cpu), name, uid, time from %s where cluster='%s' and ", table, cluster) // customized by kore-board
	} else {
		//default to metricName == "memory/usage"
		// metricName = "memory"
		query = fmt.Sprintf("select sum(memory), name, uid, time from %s where cluster='%s' and ", table, cluster) // customized by kore-board
	}

	if table == "pods" {
		orderBy = []string{"namespace", "name", "time"}
		args = append(args, "namespace=?")
		if selector.Namespace != "" {
			values = append(values, selector.Namespace)
		} else {
			values = append(values, "default")
		}
	}

	if selector.ResourceName != "" {
		if strings.ContainsAny(selector.ResourceName, ",") {
			subargs := []string{}
			for _, v := range strings.Split(selector.ResourceName, ",") {
				subargs = append(subargs, "?")
				values = append(values, v)
			}
			args = append(args, " name in ("+strings.Join(subargs, ",")+")")
		} else {
			values = append(values, selector.ResourceName)
			args = append(args, " name = ?")
		}
	}
	if selector.UID != "" {
		args = append(args, " uid = ?")
		values = append(values, selector.UID)
	}

	query = fmt.Sprintf(query+strings.Join(args, " and ")+" group by name, time order by %v;", strings.Join(orderBy, ", ")) // customized by kore-board

	return db.Query(query, values...)
}

/*
	getPodMetrics: With a database connection and a resource selector
	Queries SQLite and returns a list of metrics.
*/
func getPodMetrics(db *sql.DB, cluster string, metricName string, selector ResourceSelector) (SidecarMetricResultList, error) { // customized by kore-board
	rows, err := getRows(db, cluster, "pods", metricName, selector) // customized by kore-board
	if err != nil {
		log.Errorf("Error getting pod metrics: %v", err)
		return SidecarMetricResultList{}, err
	}

	defer rows.Close()

	resultList := make(map[string]SidecarMetric)

	for rows.Next() {
		var metricValue string
		var pod string
		var metricTime string
		var uid string
		var newMetric MetricPoint
		err = rows.Scan(&metricValue, &pod, &uid, &metricTime)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, metricTime)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		v, err := strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		newMetric = MetricPoint{
			Timestamp: t,
			Value:     v,
		}

		if _, ok := resultList[pod]; ok {
			metricThing := resultList[pod]
			metricThing.AddMetricPoint(newMetric)
			resultList[pod] = metricThing
		} else {
			resultList[pod] = SidecarMetric{
				MetricName:   metricName,
				MetricPoints: []MetricPoint{newMetric},
				DataPoints:   []DataPoint{},
				UIDs: []types.UID{
					types.UID(pod),
				},
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return SidecarMetricResultList{}, err
	}

	result := SidecarMetricResultList{}
	for _, v := range resultList {
		result.Items = append(result.Items, v)
	}

	return result, nil
}

/*
	getNodeMetrics: With a database connection and a resource selector
	Queries SQLite and returns a list of metrics.
*/
func getNodeMetrics(db *sql.DB, cluster string, metricName string, selector ResourceSelector) (SidecarMetricResultList, error) { // customized by kore-board
	resultList := make(map[string]SidecarMetric)
	rows, err := getRows(db, cluster, "nodes", metricName, selector) // customized by kore-board

	if err != nil {
		log.Errorf("Error getting node metrics: %v", err)
		return SidecarMetricResultList{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var metricValue string
		var node string
		var metricTime string
		var uid string
		var newMetric MetricPoint
		err = rows.Scan(&metricValue, &node, &uid, &metricTime)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, metricTime)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		v, err := strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		newMetric = MetricPoint{
			Timestamp: t,
			Value:     v,
		}

		if _, ok := resultList[node]; ok {
			metricThing := resultList[node]
			metricThing.AddMetricPoint(newMetric)
			resultList[node] = metricThing
		} else {
			resultList[node] = SidecarMetric{
				MetricName:   metricName,
				MetricPoints: []MetricPoint{newMetric},
				DataPoints:   []DataPoint{},
				UIDs: []types.UID{
					types.UID(node),
				},
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return SidecarMetricResultList{}, err
	}

	result := SidecarMetricResultList{}
	for _, v := range resultList {
		result.Items = append(result.Items, v)
	}

	return result, nil
}
