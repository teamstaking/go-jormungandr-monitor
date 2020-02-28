package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logrus.Info("Starting go-jormungandr-monitor...")

	stakeCollector := NewStakeCollector()
	stakeInfoCollector := NewStakeInfoCollector()
	statCollector := NewStatsCollector()
	connectionsCollector := NewConnectionsCollector()

	var registry = prometheus.NewRegistry()
	registry.MustRegister(
		stakeCollector,
		statCollector,
		connectionsCollector,
		stakeInfoCollector,
	)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	port := os.Getenv("GJM_MONITOR_PORT")
	http.ListenAndServe(":"+port, nil)
}
