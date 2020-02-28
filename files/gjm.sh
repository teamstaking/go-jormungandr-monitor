#!/bin/bash
export GJM_MONITOR_PORT=8000
export GJM_BASE_REST_URL=http://127.0.0.1:3101
export GJM_STAKE_POOL_ID=

./go-jormungandr-monitor
