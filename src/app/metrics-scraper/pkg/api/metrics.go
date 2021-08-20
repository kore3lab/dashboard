package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kore3lab/dashboard-metrics-scraper/pkg/config"
	sidedb "github.com/kore3lab/dashboard-metrics-scraper/pkg/database"
	log "github.com/sirupsen/logrus"
)

// DashboardRouter defines the usable API routes
func MetricsRouter(router *mux.Router, db *sql.DB) {

	router.Path("/clusters/{CLUSTER}").HandlerFunc(clusterHandler(db))
	router.Path("/clusters/{CLUSTER}/nodes/{NAME}").HandlerFunc(nodeHandler(db))
	router.Path("/clusters/{CLUSTER}/namespaces/{NAMESPACE}/pods/{NAME}").HandlerFunc(podHandler(db))
	router.Path("/clusters/{CLUSTER}/namespaces/{NAMESPACE}/pods/{NAME}/{OP}").HandlerFunc(podHandler(db))
}

func clusterHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cluster := vars["CLUSTER"]
		if cluster == "" {
			cluster = config.Value.DefaultContext
		}

		app := App{Writer: w}
		metrics, err := sidedb.Select(db, "nodes", cluster, "", "", "SUM")
		if err != nil {
			app.Send(nil)
		} else {
			app.Send(metrics)
		}

	}

}

func nodeHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cluster := vars["CLUSTER"]
		if cluster == "" {
			cluster = config.Value.DefaultContext
		}

		app := App{Writer: w}
		metrics, err := sidedb.Select(db, "nodes", cluster, "", vars["NAME"], "SUM")
		if err != nil {
			app.Send(nil)
		} else {
			app.Send(metrics)
		}

	}

}

func podHandler(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// customized by kore-board
		cluster := vars["CLUSTER"]
		if cluster == "" {
			cluster = config.Value.DefaultContext
		}

		app := App{Writer: w}
		metrics, err := sidedb.Select(db, "pods", cluster, vars["NAMESPACE"], vars["NAME"], vars["OP"])
		if err != nil {
			app.Send(nil)
		} else {
			app.Send(metrics)
		}

	}

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("%v - URL: %s", time.Now(), r.URL)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Errorf("Error cannot write response: %v", err)
	}
}
