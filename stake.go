package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	restclient "teamstaking.com/go-jormungandr-monitor/pkg"
)

type StakeCollector struct {
	LiveStake      *prometheus.GaugeVec
	RewardsStakers *prometheus.GaugeVec
	RewardsPool    *prometheus.GaugeVec
	PoolTaxFixed   *prometheus.GaugeVec
	PoolTaxRatio   *prometheus.GaugeVec
}

func NewStakeCollector() *StakeCollector {
	return &StakeCollector{
		LiveStake: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_pool_live_stake_lovelace",
				Help: "The live stake of pool",
			}, []string{}),
		RewardsStakers: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_reward_stakers_lovelace",
				Help: "The rewards for stakers",
			}, []string{}),
		RewardsPool: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_reward_pool_lovelace",
				Help: "The reward for the pool",
			}, []string{}),
		PoolTaxFixed: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_poolTaxFixed",
				Help: "The fixed tax for the pool",
			}, []string{}),
		PoolTaxRatio: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jormungandr_poolTaxRatio",
				Help: "The ratio tax for the pool",
			}, []string{}),
	}
}

func (sc StakeCollector) Describe(ch chan<- *prometheus.Desc) {
	sc.LiveStake.Describe(ch)
	sc.RewardsStakers.Describe(ch)
	sc.RewardsPool.Describe(ch)
	sc.PoolTaxFixed.Describe(ch)
	sc.PoolTaxRatio.Describe(ch)
}

func (sc StakeCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Info("Collecting stake metrics...")
	stakePoolId := os.Getenv("GJM_STAKE_POOL_ID")

	if len(stakePoolId) > 0 {
		stake, err := restclient.GetStake(stakePoolId)
		if err != nil {
			logrus.Error(fmt.Sprintf("Error getting stake: %s", err))
			return
		} else {
			sc.LiveStake.WithLabelValues().Set(stake.TotalStake)
			sc.LiveStake.Collect(ch)

			sc.RewardsPool.WithLabelValues().Set(stake.Rewards.ValueTaxed)
			sc.RewardsPool.Collect(ch)

			sc.RewardsStakers.WithLabelValues().Set(stake.Rewards.ValueForStakers)
			sc.RewardsStakers.Collect(ch)

			sc.PoolTaxFixed.WithLabelValues().Set(stake.Tax.Fixed)
			sc.PoolTaxFixed.Collect(ch)

			sc.PoolTaxRatio.WithLabelValues().Set(stake.Tax.Ratio.Numerator / stake.Tax.Ratio.Denominator)
			sc.PoolTaxRatio.Collect(ch)
		}
	}
}
