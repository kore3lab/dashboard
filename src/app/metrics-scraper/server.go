package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	_ "github.com/mattn/go-sqlite3"

	sideapi "github.com/kubernetes-sigs/dashboard-metrics-scraper/pkg/api"
	sidedb "github.com/kubernetes-sigs/dashboard-metrics-scraper/pkg/database"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/kubernetes-sigs/dashboard-metrics-scraper/pkg/config"
)

func main() {
	var kubeconfig *string
	var dbFile *string
	var metricResolution *time.Duration
	var metricDuration *time.Duration
	var logLevel *string
	var metricNamespace *[]string

	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	dbFile = flag.String("db-file", "metrics.db", "What file to use as a SQLite3 database.")
	metricResolution = flag.Duration("metric-resolution", 1*time.Minute, "The resolution at which metrics-scraper will poll metrics.")
	metricDuration = flag.Duration("metric-duration", 15*time.Minute, "The duration after which metrics are purged from the database.")
	logLevel = flag.String("log-level", "info", "The log level")
	// When running in a scoped namespace, disable Node lookup and only capture metrics for the given namespace(s)
	metricNamespace = flag.StringSliceP("namespace", "n", []string{getEnv("POD_NAMESPACE", "")}, "The namespace to use for all metric calls. When provided, skip node metrics. (defaults to cluster level metrics)")

	err := flag.Set("logtostderr", "true")
	if err != nil {
		log.Errorf("Error cannot set logtostderr: %v", err)
	}
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(level)
	}

	// configuration
	//      customized by acornsoft-dashboard
	config.Setup(*kubeconfig)

	// Create the db "connection"
	db, err := sql.Open("sqlite3", *dbFile)
	if err != nil {
		log.Fatalf("Unable to open Sqlite database: %s", err)
	}
	defer db.Close()

	// Populate tables
	err = sidedb.CreateDatabase(db)
	if err != nil {
		log.Fatalf("Unable to initialize database tables: %s", err)
	}

	go func() {
		r := mux.NewRouter()

		sideapi.Manager(r, db)
		// Bind to a port and pass our router in
		log.Fatal(http.ListenAndServe(":8000", handlers.CombinedLoggingHandler(os.Stdout, r)))
	}()

	// Start the machine. Scrape every metricResolution
	ticker := time.NewTicker(*metricResolution)
	quit := make(chan struct{})

	for {
		select {
		case <-quit:
			ticker.Stop()
			return

		case <-ticker.C:
			// customized by acornsoft-dashboard
			for _, ctx := range config.Value.Contexts {
				err = update(ctx, db, metricDuration, metricNamespace)
			}
			// --END
		}
	}
}

/**
* Update the Node and Pod metrics in the provided DB
 */
func update(cluster string, db *sql.DB, metricDuration *time.Duration, metricNamespace *[]string) error {
	nodeMetrics := &v1beta1.NodeMetricsList{}
	podMetrics := &v1beta1.PodMetricsList{}
	var err error

	// customized by acornsoft-dashboard
	client, err := metricsclient.NewForConfig(config.Value.KubeConfigs[cluster])
	if err != nil {
		log.Fatalf("Unable to generate a clientset: %s", err)
		return err
	}
	// --END

	// If no namespace is provided, make a call to the Node
	if len(*metricNamespace) == 1 && (*metricNamespace)[0] == "" {
		// List node metrics across the cluster
		nodeMetrics, err = client.MetricsV1beta1().NodeMetricses().List(v1.ListOptions{})
		if err != nil {
			log.Errorf("Error scraping node metrics: %s", err)
			return err
		}
	}

	// List pod metrics across the cluster, or for a given namespace
	for _, namespace := range *metricNamespace {
		pod, err := client.MetricsV1beta1().PodMetricses(namespace).List(v1.ListOptions{})
		if err != nil {
			log.Errorf("Error scraping '%s' for pod metrics: %s", namespace, err)
			return err
		}
		podMetrics.TypeMeta = pod.TypeMeta
		podMetrics.ListMeta = pod.ListMeta
		podMetrics.Items = append(podMetrics.Items, pod.Items...)
	}

	// Insert scrapes into DB
	err = sidedb.UpdateDatabase(db, cluster, nodeMetrics, podMetrics) // customized by acornsoft-dashboard
	if err != nil {
		log.Errorf("Error updating database: %s", err)
		return err
	}

	// Delete rows outside of the metricDuration time
	err = sidedb.CullDatabase(db, cluster, metricDuration)
	if err != nil {
		log.Errorf("Error culling database: %s on %s", err, cluster) // customized by acornsoft-dashboard
		return err
	}

	log.Infof("Database updated: %s cluster, %d nodes, %d pods", cluster, len(nodeMetrics.Items), len(podMetrics.Items)) // customized by acornsoft-dashboard
	return nil
}

/**
* Lookup the environment variable provided and set to default value if variable isn't found
 */
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}
