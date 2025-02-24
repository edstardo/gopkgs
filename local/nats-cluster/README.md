
## Setup a NATS Cluster using Docker Compose
Here's a step by step guide on how to setup a local NATS cluster composed of three nodes using Docker Compose <br />
Refer to [Synadia's Youtube video](https://www.youtube.com/watch?v=srARy0m9SdI&ab_channel=Synadia) and NATS Documentation on [Clustering](https://docs.nats.io/running-a-nats-service/configuration/clustering) <br />
Note: Only do steps 1-3 if first time setting up the cluster context.
### 1. Create new context named "nats_local_cluster"
```
$ nats context save nats_local_cluster
```
### 2. Edit "nats_local_cluster" context
```
$ nats context edit nats_local_cluster
```
### 3. Set description, user and password
```
description: "Context used to connect to local cluster"
url:         "nats://127.0.0.1:4222"
user:        "admin"
password:    "password"
```
### 4. Use "nats_local_cluster" context
```
$ nats context select nats_local_cluster
```
### 5. Check that "nats_local_cluster" context is selected
```
$ nats context ls
```
### 6. Start cluster using Docker Compose
```
$ docker-compose up -d
```

## Monitoring and Testing of a Local NATS Cluster
### 1. Monitor cluster using watch
You should see a total of 3 nodes under cluster "nats-cluster": nats-1, nats-2, nats-3. <br />
with a total of 1 connection, which the connection currently used by the watch command.
```
$ watch nats server ls --sort=name --reverse
```
### 2. Test cluster using NATS bench
This command starts 20 subscribers to subject "hello" that are randomly assigned to each node in the cluster.
```
$ nats bench --server=nats://localhost:4222,nats://localhost:4223,nats://localhost:4224 --sub=20 hello
```
In the watch command above, you should see that there are now 21 total active connections to the cluster.
### 3. Test publishing to cluster
This command publishes to the subject "hello" indefinitely with a 10ms interval.
```
$ nats pub hello hi --count=-1 --sleep=10ms
```
In the watch command above, you should see that there are now 22 total active connections to the cluster. <br />
In the bench terminal, you should see the 20 subscribers receiving messages.
### 4. Kill one of the cluster nodes
```
$ docker kill nats-3
```
The cluster should be able to figure out that node nats-3 is not responding so cluster should be able to reroute all connections to other healty nodes which are nats-1 and nat-2. <br />
In the watch command above, after a few seconds, the total numbers of connections to the cluster should equal 22, distributed between nats-1 and nats-2.