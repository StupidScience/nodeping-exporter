# NodePing exporter
[![Build Status](https://travis-ci.org/StupidScience/nodeping-exporter.svg?branch=master)](https://travis-ci.org/StupidScience/nodeping-exporter)
[![Coverage Status](https://coveralls.io/repos/github/StupidScience/nodeping-exporter/badge.svg)](https://coveralls.io/github/StupidScience/nodeping-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/StupidScience/nodeping-exporter)](https://goreportcard.com/report/github.com/StupidScience/nodeping-exporter)

Get checks statistics from [NodePing service API](https://nodeping.com/docs-api-overview.html) and expose it in prometheus format.

## Configuration

Exporter configurates via environment variables:

|Env var|Description|
|---|---|
|NODEPING_TOKEN|Authorization token for access to NodePing API|

Exporter listen on tcp-port **9503**. Metrics available on `/metrics` path.

## Exposed metrics

|Metric|Labels|Descrpition|
|---|---|---|
|nodeping_check_up|label, target, type|Current check status: 1 - for up, 0 - for down|
|nodeping_check_duration_seconds|label, target, type|Last check duration in seconds|
