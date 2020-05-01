package jormungandrrestclent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type StakeDto struct {
	Rewards    StakeRewardsDto `json:"rewards"`
	Tax        StakeTaxDto     `json:"tax"`
	TotalStake float64         `json:"total_stake"`
}

type StakeRewardsDto struct {
	Epoch           float64 `json:"epoch"`
	ValueForStakers float64 `json:"value_for_stakers"`
	ValueTaxed      float64 `json:"value_taxed"`
}

type StakeTaxDto struct {
	Fixed float64          `json:"fixed"`
	Ratio StakeTaxRatioDto `json:"ratio"`
}

type StakeTaxRatioDto struct {
	Numerator   float64 `json:"numerator"`
	Denominator float64 `json:"denominator"`
}

type Client struct {
	httpClient http.Client
	baseUrl    string
}

func NewClient() Client {
	baseUrl := os.Getenv("GJM_BASE_REST_URL")
	httpClient := &http.Client{
		Timeout: 1 * time.Second,
	}

	client := Client{
		httpClient: *httpClient,
		baseUrl:    baseUrl,
	}

	return client
}

func (client *Client) GetStake(stakePoolId string) (*StakeDto, error) {
	resp, err := client.httpClient.Get(fmt.Sprintf("%s/api/v0/stake_pool/%s", client.baseUrl, stakePoolId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stake := StakeDto{}
	err = json.Unmarshal(bytes, &stake)
	if err != nil {
		return nil, err
	}

	return &stake, nil
}

type StatsDto struct {
	LastBlockHeight *string `json:"lastBlockHeight"`
	Uptime          float64 `json:"uptime"`
	TxRecvCnt       float64 `json:"txRecvCnt"`
	LastBlockDate   *string `json:"lastBlockDate"`
	Version         *string `json:"version"`
	State           *string `json:"state"`
	PeerAvailable   float64 `json:"peerAvailableCnt"`
	PeerQuarantined float64 `json:"peerQuarantinedCnt"`
	PeerTotal       float64 `json:"peerTotalCnt"`
}

func (client *Client) GetStats() (*StatsDto, error) {
	resp, err := client.httpClient.Get(fmt.Sprintf("%s/api/v0/node/stats", client.baseUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stats := StatsDto{}
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

type StakeInfoDto struct {
	Epoch float64           `json:"epoch"`
	Stake StakeInfoStakeDto `json:"stake"`
}

type StakeInfoStakeDto struct {
	Pools [][]interface{} `json:"pools"`
}

func (client *Client) GetStakeInfo() (*StakeInfoDto, error) {
	resp, err := client.httpClient.Get(fmt.Sprintf("%s/api/v0/stake", client.baseUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stake := StakeInfoDto{}
	err = json.Unmarshal(bytes, &stake)
	if err != nil {
		return nil, err
	}

	return &stake, nil
}

type ConnectionsDto struct {
	NodeId *string `json:"nodeId"`
}

func (client *Client) GetConnections() ([]ConnectionsDto, error) {
	resp, err := client.httpClient.Get(fmt.Sprintf("%s/api/v0/network/stats", client.baseUrl))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	connections := []ConnectionsDto{}
	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return nil, err
	}

	return connections, nil
}
