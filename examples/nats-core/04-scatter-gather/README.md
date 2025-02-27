
## Scatter Gather With Load Balancing Using NATS Request-Reply
This demos a simple scatter application with load balancing and request-reply. <br />

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
The program is set with 10 publishers and 10 subscribers and will run for 10sec or until terminated.

### Output
```
INFO[0011] responder[6] received[18861] processed[18861]
INFO[0011] responder[1] received[19200] processed[19200]
INFO[0011] responder[4] received[19239] processed[19239]
INFO[0011] responder[8] received[19076] processed[19076]
INFO[0011] responder[3] received[19085] processed[19085]
INFO[0011] responder[2] received[19123] processed[19123]
INFO[0011] responder[7] received[18899] processed[18899]
INFO[0011] responder[0] received[18933] processed[18933]
INFO[0011] responder[9] received[19070] processed[19070]
INFO[0011] responder[5] received[19295] processed[19295]
INFO[0012] publisher[1] sent[19049] responded[19049]
INFO[0012] publisher[10] sent[19077] responded[19077]
INFO[0012] publisher[7] sent[19076] responded[19076]
INFO[0012] publisher[2] sent[19072] responded[19072]
INFO[0012] publisher[5] sent[19090] responded[19090]
INFO[0012] publisher[8] sent[19083] responded[19083]
INFO[0012] publisher[3] sent[19087] responded[19087]
INFO[0012] publisher[9] sent[19088] responded[19088]
INFO[0012] publisher[6] sent[19077] responded[19077]
INFO[0012] publisher[4] sent[19082] responded[19082]
INFO[0015] num of total requests received by all responders: 190781
INFO[0015] num of total responses received by all publishers: 190781
```
In the output, you can see how many messages are published and received by each publisher and subscriber.
You can also see that all messages are load balanced among all subscribers and total requests equals total responses.