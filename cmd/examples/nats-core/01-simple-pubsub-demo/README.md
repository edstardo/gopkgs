
## Simple PubSub Demo using Core NATS

This demos a simple pubsub using available nats subscribe methods in the NATS.go library. <br />

## NATS Cluster Setup
For local NATS cluster setup see [Guide to Setup a Local NATS Cluster](../../../../cmd/setup/local-nats-cluster/README.md). <br />
## Dependencies
Download dependencies
```
$ go mod tidy
```

## Run
Steps to run this demo:
- setup a local nats cluster, see [Setup](#setup) 
- cd into the demo folder
- build
```
$ go build -o demo
```
- run
```
$ ./demo ./
```

