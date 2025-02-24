# go-availability-checker
A micro-service written in golang used to test the availability and latency of an endpoint within a GKE cluster. The metrics are exposed in prometheus metric syntax.


## Getting started

export URL="https://google.com"
export FREQUENCY="10"

## Containerize
docker build -t gac:1.0.0 .