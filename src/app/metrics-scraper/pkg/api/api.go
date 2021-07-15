package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kore3lab/dashboard-metrics-scraper/pkg/config"
)

// Manager provides a handler for all api calls
func Manager(r *mux.Router, db *sql.DB) {

	router := r.PathPrefix("/api/v1").Subrouter()
	MetricsRouter(router, db)

	r.HandleFunc("/api/kubeconfig", LoadConfig).Methods("PUT")
	r.PathPrefix("/").HandlerFunc(DefaultHandler)

}

func LoadConfig(w http.ResponseWriter, _ *http.Request) {
	config.Setup()

	app := App{Writer: w}
	app.SendMessage(http.StatusOK, "Kubeconfig updata successful")
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	app := App{Writer: w}
	app.SendMessage(http.StatusOK, fmt.Sprintf("URL: %s", r.URL))

}
