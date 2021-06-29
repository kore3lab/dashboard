package api

import (
	"database/sql"
	"fmt"
	"github.com/kore3lab/dashboard-metrics-scraper/pkg/config"
	"net/http"

	"github.com/gorilla/mux"
	dashboardProvider "github.com/kore3lab/dashboard-metrics-scraper/pkg/api/dashboard"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// Manager provides a handler for all api calls
func Manager(r *mux.Router, db *sql.DB) {
	dashboardRouter := r.PathPrefix("/api/v1").Subrouter() // customized by kore-board
	dashboardProvider.DashboardRouter(dashboardRouter, db)
	dashboardProvider.DashboardExpandRouter(dashboardRouter, db) // customized by kore-board
	r.HandleFunc("/api/kubeconfig", LoadConfig).Methods("PUT")
	r.PathPrefix("/").HandlerFunc(DefaultHandler)
}

func LoadConfig(w http.ResponseWriter, _ *http.Request) {
	config.Setup()
	msg := fmt.Sprint("Kubeconfig updata successful")
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Errorf("Error cannot reload kubeconfig: %v", err)
	}
}

// DefaultHandler provides a handler for all http calls
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("URL: %s", r.URL)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Errorf("Error cannot write response: %v", err)
	}
}
