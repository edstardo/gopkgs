
## Chat using Bubble Tea
This is a simple terminal-based chat application using Core Nats and [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## NATS Cluster Setup
For local NATS cluster setup see [Guide to Setup a Local NATS Cluster](../../../local/nats-cluster/README.md). <br />
## Dependencies
Download dependencies
```
$ go mod tidy
```

## Run
CD into this project directory and run the following command in multiple separate terminals to simulate a live chat using Core NATS. 
```
$ go run ./ ./
```
Every instances of the application connects to a local nats cluster, subscribes and publishes to subject "chat.bubble" creating a mesh of clients forming a chat room.

