package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"

	restclient "teamstaking.com/go-jormungandr-monitor/pkg"
)

type StatsCollector struct {
	LastBlockHeight      *prometheus.GaugeVec
	Uptime               *prometheus.GaugeVec
	LastBlockEpoch       *prometheus.GaugeVec
	LastBlockSlot        *prometheus.GaugeVec
	TxRecvCount          *prometheus.GaugeVec
	BuildInfo            *prometheus.GaugeVec
	PeerAvailableCount   *prometheus.GaugeVec
	PeerQuarantinedCount *prometheus.GaugeVec
	PeerTotalCount       *prometheus.GaugeVec
}

func NewStatsCollector() *StatsCollector {
	return &StatsCollector{
		LastBlockHeight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_lastBlockHeight",
				Help: "Last block height from the status endpoint",
			}, []string{}),
		Uptime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_uptime",
				Help: "Uptime from the status endpoint",
			}, []string{}),
		LastBlockEpoch: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_lastBlockEpoch",
				Help: "Epoch of last block",
			}, []string{}),
		LastBlockSlot: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_lastBlockSlot",
				Help: "Slot of last block",
			}, []string{}),
		TxRecvCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_txRecvCnt",
				Help: "Transactions received",
			}, []string{}),
		PeerAvailableCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_peerAvailableCnt",
				Help: "Peers available",
			}, []string{}),
		PeerQuarantinedCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_peerQuarantinedCnt",
				Help: "Peers quarantined",
			}, []string{}),
		PeerTotalCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_peerTotalCnt",
				Help: "Peer total",
			}, []string{}),
		BuildInfo: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_buildInfo",
				Help: "Build info",
			},
			[]string{"version", "state"}),
	}
}

func (sc StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	sc.LastBlockHeight.Describe(ch)
	sc.Uptime.Describe(ch)
	sc.LastBlockEpoch.Describe(ch)
	sc.LastBlockSlot.Describe(ch)
	sc.TxRecvCount.Describe(ch)
	sc.BuildInfo.Describe(ch)
	sc.PeerAvailableCount.Describe(ch)
	sc.PeerQuarantinedCount.Describe(ch)
	sc.PeerTotalCount.Describe(ch)
}

func (sc StatsCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("Collecting stats metrics...")

	sc.BuildInfo.Reset()

	// get stats info
	stats, err := restclient.GetStats()
	if err != nil {
		logrus.Error(fmt.Sprintf("Error getting stats info: %s", err))
		sc.BuildInfo.WithLabelValues("Offline", "Offline").Set(0)
		sc.BuildInfo.Collect(ch)
		return
	} else {
		versionSplit := strings.Split(*stats.Version, " ")
		sc.BuildInfo.WithLabelValues(versionSplit[1], *stats.State).Set(1)
		sc.BuildInfo.Collect(ch)

		if stats.LastBlockDate != nil {
			lastBlockFloat, err := strconv.ParseFloat(*stats.LastBlockHeight, 64)
			if err != nil {
				logrus.Error(fmt.Sprintf("Error getting last block height %s", err))
			} else {
				sc.LastBlockHeight.WithLabelValues().Set(lastBlockFloat)
				sc.LastBlockHeight.Collect(ch)
			}
		}

		sc.Uptime.WithLabelValues().Set(stats.Uptime)
		sc.Uptime.Collect(ch)

		sc.TxRecvCount.WithLabelValues().Set(stats.TxRecvCnt)
		sc.TxRecvCount.Collect(ch)

		sc.PeerAvailableCount.WithLabelValues().Set(stats.PeerAvailable)
		sc.PeerAvailableCount.Collect(ch)

		sc.PeerQuarantinedCount.WithLabelValues().Set(stats.PeerQuarantined)
		sc.PeerQuarantinedCount.Collect(ch)

		sc.PeerTotalCount.WithLabelValues().Set(stats.PeerTotal)
		sc.PeerTotalCount.Collect(ch)

		if stats.LastBlockDate != nil {
			var splits = strings.Split(*stats.LastBlockDate, ".")
			epoch, err := strconv.ParseFloat(splits[0], 64)
			slot, err := strconv.ParseFloat(splits[1], 64)
			if err != nil {
				logrus.Error(fmt.Sprintf("Error getting last block date %s", err))
			} else {
				sc.LastBlockEpoch.WithLabelValues().Set(epoch)
				sc.LastBlockSlot.WithLabelValues().Set(slot)
				sc.LastBlockEpoch.Collect(ch)
				sc.LastBlockSlot.Collect(ch)
			}
		}
	}

}
