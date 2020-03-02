### Overview

Provides a go implementation to monitor the Cardano Jormungandr Rust Node. Used to drive some of the data at https://teamstaking.com

May need to set GOOS and GOARCH values if compiling on system not running on (https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)

### Installing / Running

1. Build from source (go build) or download executable in releases
2. Configure and run bash script included in files folder
```
# port to run monitor on
export GJM_MONITOR_PORT=8000
# url of jormungandr rest api
export GJM_BASE_REST_URL=http://127.0.0.1:3101
# your stake pool id, can be left blank for passive nodes
export GJM_STAKE_POOL_ID=
```
3. Make sure `netstat` is installed as that is used for the server connection metrics

### Examples Metrics

curl http://localhost:8000/metrics

    # HELP jormungandr_buildInfo Build info
    # TYPE jormungandr_buildInfo gauge
    jormungandr_buildInfo{state="Running",version="0.8.9-46df05c7"} 1
    # HELP jormungandr_connections Jormungandr connections
    # TYPE jormungandr_connections gauge
    jormungandr_connections 908
    # HELP jormungandr_epoch The current epoch
    # TYPE jormungandr_epoch gauge
    jormungandr_epoch 74
    # HELP jormungandr_lastBlockEpoch Epoch of last block
    # TYPE jormungandr_lastBlockEpoch gauge
    jormungandr_lastBlockEpoch 74
    # HELP jormungandr_lastBlockHeight Last block height from the status endpoint
    # TYPE jormungandr_lastBlockHeight gauge
    jormungandr_lastBlockHeight 243982
    # HELP jormungandr_lastBlockSlot Slot of last block
    # TYPE jormungandr_lastBlockSlot gauge
    jormungandr_lastBlockSlot 39662
    # HELP jormungandr_peerAvailableCnt Peers available
    # TYPE jormungandr_peerAvailableCnt gauge
    jormungandr_peerAvailableCnt 3584
    # HELP jormungandr_peerQuarantinedCnt Peers quarantined
    # TYPE jormungandr_peerQuarantinedCnt gauge
    jormungandr_peerQuarantinedCnt 6584
    # HELP jormungandr_peerTotalCnt Peer total
    # TYPE jormungandr_peerTotalCnt gauge
    jormungandr_peerTotalCnt 10255
    # HELP jormungandr_poolTaxFixed The fixed tax for the pool
    # TYPE jormungandr_poolTaxFixed gauge
    jormungandr_poolTaxFixed 0
    # HELP jormungandr_poolTaxRatio The ratio tax for the pool
    # TYPE jormungandr_poolTaxRatio gauge
    jormungandr_poolTaxRatio 0.05
    # HELP jormungandr_pool_live_stake_lovelace The live stake of pool
    # TYPE jormungandr_pool_live_stake_lovelace gauge
    jormungandr_pool_live_stake_lovelace 0
    # HELP jormungandr_reward_pool_lovelace The reward for the pool
    # TYPE jormungandr_reward_pool_lovelace gauge
    jormungandr_reward_pool_lovelace 0
    # HELP jormungandr_reward_stakers_lovelace The rewards for stakers
    # TYPE jormungandr_reward_stakers_lovelace gauge
    jormungandr_reward_stakers_lovelace 0
    # HELP jormungandr_server_connections The server level connections
    # TYPE jormungandr_server_connections gauge
    jormungandr_server_connections{state="CLOSING"} 1
    jormungandr_server_connections{state="ESTABLISHED"} 589
    jormungandr_server_connections{state="FIN_WAIT1"} 1
    jormungandr_server_connections{state="FIN_WAIT2"} 1
    jormungandr_server_connections{state="LAST_ACK"} 1
    jormungandr_server_connections{state="SYN_RECV"} 2
    jormungandr_server_connections{state="SYN_SENT"} 17
    jormungandr_server_connections{state="TIME_WAIT"} 61
    # HELP jormungandr_total_staked_lovelace Total staked in the current epoch
    # TYPE jormungandr_total_staked_lovelace gauge
    jormungandr_total_staked_lovelace 1.1685566218546152e+16
    # HELP jormungandr_txRecvCnt Transactions received
    # TYPE jormungandr_txRecvCnt gauge
    jormungandr_txRecvCnt 21
    # HELP jormungandr_uptime Uptime from the status endpoint
    # TYPE jormungandr_uptime gauge
    jormungandr_uptime 3398

#### Offline Metric Example

Build Info will return offline with a value of 0 if the stats rest call fails

```
jormungandr_buildInfo{state="Offline",version="Offline"} 0
```


### REST Calls

Swagger Doc: https://raw.githubusercontent.com/input-output-hk/jormungandr/master/doc/openapi.yaml

#### Stake Pool Info

Url: {base_url}/api/v0/stake_pool/{pool_id}
Metrics:
* jormungandr_pool_live_stake_lovelace
* jormungandr_reward_stakers_lovelace
* jormungandr_reward_pool_lovelace

```
{
  "tax": {
    "fixed": 5,
    "ratio": {
      "numerator": 1
      "denominator": 10000,
    }
    "max": 100,
  },
  "rewards": {
    "epoch": 42,
    "value_taxed": 2901,
    "value_for_stakers": 2028
  },
  "total_stake": 2000000000000,
  "kesPublicKey": "kes25519-12-pk1q7susucqwje0lpetqzjgzncgcrjzx7e2guh900qszdjskkeyqpusf3p39r",
  "vrfPublicKey": "vrf_pk1rcm4qm3q9dtwq22x9a4avnan7a3k987zvepuxwekzj3uyu6a8v0s6sdy0l"
}
```

#### Stake Info

Url: {base_url}/api/v0/stake
Metrics:
* jormungandr_total_staked_lovelace
* jormungandr_epoch
* jormungandr_pool_controlled_stake_lovelace

```
{
  "epoch": 0,
  "stake": {
    "dangling": 0,
    "pools": [
      [
        "d882fc32c4b4b901cb29dfb4162e070d7650e937abb7bc2947d3a7d48b6c86a6",
        2000000000000
      ]
    ],
    "unassigned": 0
  }
}
```

#### Status

Url: {base_url}/api/v0/node/stats
Metrics:
* jormungandr_lastBlockHeight
* jormungandr_uptime
* jormungandr_lastBlockSlot
* jormungandr_lastBlockEpoch
* jormungandr_txRecvCnt
* jormungandr_buildInfo
* jormungandr_peerAvailableCnt
* jormungandr_peerQuarantinedCnt
* jormungandr_peerTotalCnt

```
{
  "version": "jormungandr 0.8.7-364cd84",
  "state": "Bootstrapping"
}
```

```
{
  "blockRecvCnt": 1102,
  "lastBlockContentSize": 484,
  "lastBlockDate": "20.29",
  "lastBlockFees": 534,
  "lastBlockHash": "b9597b45a402451540e6aabb58f2ee4d65c67953b338e04c52c00aa0886bd1f0",
  "lastBlockHeight": 202901,
  "lastBlockSum": 51604,
  "lastBlockTime": "2019-08-12T11:20:52.316544007+00:00",
  "lastBlockTx": 2,
  "peerAvailableCnt": 7437,
  "peerQuarantinedCnt": 18875,
  "peerTotalCnt": 26393,
  "state": "Running",
  "txRecvCnt": 5440,
  "uptime": 20032,
  "version": "jormungandr 0.8.0-rc4-67477249"
}
```

#### Network

Url: {base_url}/api/v0/network/stats
Metrics:
  * jormungandr_connections
  
```
[
  {
    "addr": "3.124.55.91:3000",
    "nodeId": "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
    "establishedAt": "2019-10-14T06:24:12.010231281+00:00",
    "lastBlockReceived": "2019-10-14T07:54:32.014432772+00:00",
    "lastFragmentReceived": "2019-10-14T07:54:33.014432831+00:00",
    "lastGossipReceived": "2019-10-14T07:54:34.014432887+00:00"
  },
  {
    "nodeId": "02f1e1d1c1b1a191817161514131211101f0e0d0c0b0a0908070605040302010"
  }
]
```

#### Settings

```
{
  "block0Hash": "8d94ecfcc9a566f492e6335858db645691f628b012bed4ac2b1338b5690355a7",
  "block0Time": "2019-07-09T12:32:51+00:00",
  "blockContentMaxSize": 102400,
  "consensusVersion": "bft",
  "currSlotStartTime": "2019-07-18T22:01:17+00:00",
  "fees": {
    "certificate": 4,
    "per_certificate_fees": {
      "certificate_pool_registration": 5,
      "certificate_stake_delegation": 3,
      "certificate_owner_stake_delegation": 4,
    }
    "coefficient": 1,
    "constant": 2
  },
  "rewardParams": {
      "compoundingRatio": {
          "denominator": 1024,
          "numerator": 1
      },
      "compoundingType": "Linear",
      "epochRate": 100,
      "epochStart": 0,
      "initialValue": 10000
  },
  "slotDuration": 10,
  "slotsPerEpoch": 60,
  "treasuryTax": {
    "fixed": 5,
    "ratio": {
      "numerator": 1
      "denominator": 10000,
    }
    "max": 100,
  }
}
```
