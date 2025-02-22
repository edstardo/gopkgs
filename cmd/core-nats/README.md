
## Core NATS Examples
Here's a list of mini projects using pure Core NATS
- [Simple Publisher Subcriber](./01-simple-pubsub-demo/README.md)

## Setup a NATS Cluster using Docker Compose
Here's a step by step guide on how to setup a local NATS cluster composed of three nodes using Docker Compose <br />
Refer to [Synadia's Youtube video](https://www.youtube.com/watch?v=srARy0m9SdI&ab_channel=Synadia) and NATS Documentation on [Clustering](https://docs.nats.io/running-a-nats-service/configuration/clustering) <br />
Note: Only do steps 1-3 if first time setting up the cluster

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

### 6. Monitor cluster
```
$ watch nats server ls --sort=name
```

### 7. Start cluster using Docker Compose
```
$ docker-compose up -d
```

## Monitoring and Testing of a Local NATS Cluster
