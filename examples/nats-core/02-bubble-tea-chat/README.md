
## Chat using Bubble Tea
This is a simple terminal-based chat application using Core Nats and [Bubble Tea](https://github.com/charmbracelet/bubbletea)

### NATS Cluster Setup
-----
For local NATS cluster setup see [Guide to Setup a Local NATS Cluster](../../../local/nats-cluster/README.md). <br />

### Dependencies
-----
Download dependencies
```
$ go mod tidy
```

### Run
-----
CD into this project directory and run the following command in multiple separate terminals to simulate a live chat using Core NATS. 
```
$ go run ./ ./
```
Every instance of the application connects to a local nats cluster, subscribes and publishes to subject "chat.bubble" creating a mesh of clients forming a chat room.

### Test client instances
-----
Make that a few instances of the application are running then on a separate terminal run the following command:
```
$nats pub  "chat.bubble" '{"sender_id":"admin","message":"message-{{Count}}"}'  --count=-1 --sleep=500ms --no-context
```
This command publishes a message to the "chat.bubble" subject every 500ms indefinitely until closed. <br />
All active instances should receive message with sender "admin".
