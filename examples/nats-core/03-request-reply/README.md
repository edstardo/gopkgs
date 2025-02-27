
## NATS Request-Reply Simple Demo
This demos a simple application using NATS Request-Reply. <br />
Check its [Documentation](https://docs.nats.io/nats-concepts/core-nats/reqreply) and [Semantics](https://docs.nats.io/using-nats/developer/sending/request_reply). <br />


### NATS Cluster Setup
-----
For local NATS cluster setup see [Guide to Setup a Local NATS Cluster](../../../local/nats-cluster/README.md). <br />

### Dependencies
----
Download dependencies
```
$ go mod tidy
```

### Run
-----
```
$ go run main.go
```

### Description
-----
The program will start to show all subscribers receiving the message being published. <br />
To terminate the program, just hit CTRL+C. Next you'll set the Stats for each subscribers like this:
```
subscriber[1] closed: 1000000 reveived = 100.00 percent
```
The Request-Reply mechanism ensures that all messages are published and replied to by all subscribers. <bt />

This also demonstrates NATS ability to allow multiple responders. A requester utilizes the first reponse it receives and NATS effeciently discards the others.
