package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	_ "github.com/mattn/go-sqlite3"

	sideapi "github.com/kore3lab/dashboard-metrics-scraper/pkg/api"
	sidedb "github.com/kore3lab/dashboard-metrics-scraper/pkg/database"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/gorilla/mux"

	"github.com/fsnotify/fsnotify"
	"github.com/kore3lab/dashboard-metrics-scraper/pkg/config"
	"github.com/spf13/viper"
)

func main() {
	var kubeconfig *string
	var dbFile *string
	var metricResolution *time.Duration
	var metricDuration *time.Duration
	var logLevel *string
	var metricNamespace *[]string

	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	dbFile = flag.String("db-file", "/tmp/metrics.db", "What file to use as a SQLite3 database.")
	metricResolution = flag.Duration("metric-resolution", 1*time.Minute, "The resolution at which metrics-scraper will poll metrics.")
	metricDuration = flag.Duration("metric-duration", 15*time.Minute, "The duration after which metrics are purged from the database.")
	logLevel = flag.String("log-level", "debug", "The log level")
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
	//      customized by kore-board
	config.SetKubeconfig(*kubeconfig)
	config.Setup()
	//ConfigMap update delay problem quick fix..
	WatchConfig()

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
			// customized by kore-board
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

	// customized by kore-board
	conf, err := config.KubeConfigs(cluster)
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	client, err := metricsclient.NewForConfig(conf)
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
	err = sidedb.UpdateDatabase(db, cluster, nodeMetrics, podMetrics) // customized by kore-board
	if err != nil {
		log.Errorf("Error updating database: %s", err)
		return err
	}

	// Delete rows outside of the metricDuration time
	err = sidedb.CullDatabase(db, cluster, metricDuration)
	if err != nil {
		log.Errorf("Error culling database: %s on %s", err, cluster) // customized by kore-board
		return err
	}

	log.Infof("Database updated: %s cluster, %d nodes, %d pods", cluster, len(nodeMetrics.Items), len(podMetrics.Items)) // customized by kore-board
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

func WatchConfig() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigType("yaml")
	v.SetConfigFile(config.Value.ConfigLoadingRules.GetExplicitFile())

	// monitor the changes in the config file
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.WithField("file", e.Name).Info("Config file changed")
		config.Setup()
	})
}
