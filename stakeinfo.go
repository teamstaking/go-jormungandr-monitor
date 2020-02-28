package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	restclient "teamstaking.com/go-jormungandr-monitor/pkg"
)

type StakeInfoCollector struct {
	Epoch           *prometheus.GaugeVec
	TotalStake      *prometheus.GaugeVec
	ControlledStake *prometheus.GaugeVec
}

func NewStakeInfoCollector() *StakeInfoCollector {
	return &StakeInfoCollector{
		Epoch: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_epoch",
				Help: "The current epoch",
			}, []string{}),
		TotalStake: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_total_staked_lovelace",
				Help: "Total staked in the current epoch",
			}, []string{}),
		ControlledStake: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_pool_controlled_stake_lovelace",
				Help: "Controlled staked in the current epoch",
			}, []string{}),
	}
}

func (sc StakeInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	sc.Epoch.Describe(ch)
	sc.TotalStake.Describe(ch)
	sc.ControlledStake.Describe(ch)
}

func (sc StakeInfoCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Info("Collecting stake info metrics...")
	stakePoolId := os.Getenv("GJM_STAKE_POOL_ID")

	// get stake info
	stakeInfo, err := restclient.GetStakeInfo()
	if err != nil {
		logrus.Error(fmt.Sprintf("Error getting stake info: %s", err))
		return
	} else {
		sc.Epoch.WithLabelValues().Set(stakeInfo.Epoch)
		sc.Epoch.Collect(ch)

		totalStake := float64(0)
		poolInfo := stakeInfo.Stake.Pools
		for _, pool := range poolInfo {
			//add up total stake
			totalStake += pool[1].(float64)

			if len(stakePoolId) > 0 {
				//get controlled stake for pool
				poolId := pool[0]
				if poolId == stakePoolId {
					sc.ControlledStake.WithLabelValues().Set(pool[1].(float64))
					sc.ControlledStake.Collect(ch)
				}
			}
		}

		sc.TotalStake.WithLabelValues().Set(totalStake)
		sc.TotalStake.Collect(ch)
	}
}
