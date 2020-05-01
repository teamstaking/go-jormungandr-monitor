package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ConnectionsCollector struct {
	Connections            *prometheus.GaugeVec
	JormungandrConnections *prometheus.GaugeVec
}

func NewConnectionsCollector() *ConnectionsCollector {
	return &ConnectionsCollector{
		Connections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_server_connections",
				Help: "The server level connections",
			}, []string{"state"}),
		JormungandrConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_connections",
				Help: "Jormungandr connections",
			}, []string{}),
	}
}

func (sc ConnectionsCollector) Describe(ch chan<- *prometheus.Desc) {
	sc.Connections.Describe(ch)
	sc.JormungandrConnections.Describe(ch)
}

func (cc ConnectionsCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("Collecting connections metrics...")
	cc.getNetstats(ch)
}

func (cc ConnectionsCollector) getNetstats(ch chan<- prometheus.Metric) {
	cc.Connections.Reset()

	netstat := "netstat -tn | tail -n +3 | awk \"{ print \\$6 }\" | sort | uniq -c | sort -n"
	cmd := exec.Command("bash", "-c", netstat)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(fmt.Sprintf("  Error running command: %s", err))
		return
	}

	netstatOutput := strings.Split(string(out), "\n")
	for _, output := range netstatOutput {
		outputTrim := strings.Trim(output, " ")
		if outputTrim != "" {
			outputSplit := strings.Split(outputTrim, " ")
			count, _ := strconv.ParseFloat(outputSplit[0], 64)

			cc.Connections.WithLabelValues(outputSplit[1]).Set(count)
		}
	}

	cc.Connections.Collect(ch)

	connections, err := jormungandrClient.GetConnections()
	if err != nil {
		logrus.Error(fmt.Sprintf("Error getting connection info: %s", err))
		return
	} else {
		cc.JormungandrConnections.WithLabelValues().Set(float64(len(connections)))
		cc.JormungandrConnections.Collect(ch)
	}
}
