package provider

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"github.com/acornsoftlab/dashboard-metrics-scraper/pkg/config"
)

type ClusterMetrics struct {
	CPU     SidecarMetric `json:"cpu"`
	Memory  SidecarMetric `json:"memory"`
	Storage SidecarMetric `json:"storage"`
}

func DashboardExpandRouter(r *mux.Router, db *sql.DB) {
	r.Path("/config").HandlerFunc(patchHandler(db)).Methods("PATCH")
	r.Path("/clusters/{Cluster}").HandlerFunc(clusterHandler(db)).Methods("GET")
}

func patchHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		config.Setup()
	}

}

func clusterHandler(db *sql.DB) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		clsuter := vars["Cluster"]
		if clsuter == "" {
			clsuter = config.Value.DefaultContext
		}

		resp, err := getClusterMetrics(db, clsuter)
		if err != nil {
			sendError(w, fmt.Sprintf("Cluster Metrics Error - %v", err.Error()))
		} else {
			send(w, resp)
		}

	}

	return fn
}

func getClusterMetrics(db *sql.DB, cluster string) (metrics map[string]SidecarMetric, err error) {

	metrics = make(map[string]SidecarMetric)
	cMetrics := SidecarMetric{}
	mMetrics := SidecarMetric{}
	sMetrics := SidecarMetric{}

	query := fmt.Sprintf("SELECT time, SUM(cpu), SUM(memory), SUM(storage) FROM nodes WHERE cluster='%s' GROUP BY time", cluster)
	rows, err := db.Query(query)
	if err != nil {
		log.Errorf("Error getting pod metrics: %v", err)
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var cpuValue string
		var memoryValue string
		var storageValue string
		var metricTime string

		err = rows.Scan(&metricTime, &cpuValue, &memoryValue, &storageValue)
		if err != nil {
			return metrics, err
		}

		layout := "2006-01-02T15:04:05Z"
		tm, err := time.Parse(layout, metricTime)
		if err != nil {
			return metrics, err
		}

		cpu, err := strconv.ParseUint(cpuValue, 10, 64)
		if err != nil {
			return metrics, err
		}
		memory, err := strconv.ParseUint(memoryValue, 10, 64)
		if err != nil {
			return metrics, err
		}
		storage, err := strconv.ParseUint(storageValue, 10, 64)
		if err != nil {
			return metrics, err
		}
		cMetrics.MetricPoints = append(cMetrics.MetricPoints, MetricPoint{Timestamp: tm, Value: cpu})
		mMetrics.MetricPoints = append(mMetrics.MetricPoints, MetricPoint{Timestamp: tm, Value: memory})
		sMetrics.MetricPoints = append(sMetrics.MetricPoints, MetricPoint{Timestamp: tm, Value: storage})
		cMetrics.DataPoints = append(cMetrics.DataPoints, DataPoint{X: tm.UnixNano() / int64(time.Second), Y: int64(cpu)})
		mMetrics.DataPoints = append(mMetrics.DataPoints, DataPoint{X: tm.UnixNano() / int64(time.Second), Y: int64(memory)})
		sMetrics.DataPoints = append(sMetrics.DataPoints, DataPoint{X: tm.UnixNano() / int64(time.Second), Y: int64(storage)})

	}
	err = rows.Err()
	if err != nil {
		return metrics, err
	}
	cMetrics.MetricName = "cpu/usage"
	mMetrics.MetricName = "memory/usage"
	sMetrics.MetricName = "usage/usage"

	metrics["cpu"] = cMetrics
	metrics["memory"] = mMetrics
	metrics["storage"] = sMetrics

	return metrics, nil
}
func sendError(w http.ResponseWriter, msg string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, fmt.Sprintf("{\"message\": \"%s\"}", msg))
	log.Errorf(msg)

}

func send(w http.ResponseWriter, data interface{}) {

	json, err := json.Marshal(data)
	if err != nil {
		sendError(w, fmt.Sprintf("JSON Error - %v", err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(json)
		if err != nil {
			sendError(w, fmt.Sprintf("Error cannot write response: %v", err))
		}
	}

}
