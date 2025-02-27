
# NATS TODO
## NATS Learning References
- [NATS by Example](https://natsbyexample.com/)
- [NATS Blog](https://nats.io/blog/)
- [NATS Documentation](https://docs.nats.io/)

## Setup
- ✅ nats cluster using docker / docker-compose
- ⬜️ setup nats super cluster using kubernetes
- ⬜️ nats authorization / authentication

## TODOS
### Core NATS
- ✅ Simple PubSub Demo
- ✅ Chat Application
- ✅ Simple Request-Reply Demo
- ✅ Scatter Gather with Load Balancing
- NATS Micro
    - refer to [Synadia's video](https://www.youtube.com/watch?v=AiUazlrtgyU&t=14s&ab_channel=Synadia)
- Distributed Logger (Easy-Medium)
    - Description: A logging system where multiple services send logs to a central NATS subject. A subscriber writes logs to a file or stdout.
    - Concepts Used: Structured logging, message filtering.
    - Example: Services log messages to logs.error, logs.info, and a logger writes them to logs.txt.
- Task Queue (Medium)
    - Description: A simple job queue where publishers submit tasks, and workers consume and process them.
    - Concepts Used: Load distribution using queue subscribers.
    - Example: A PDF processing system where tasks (pdf.process) are consumed by multiple worker instances.
- Metrics Aggregator (Medium)
    - Description: Services publish metrics (CPU, memory, request count) to NATS, and a metrics collector aggregates and stores them.
    - Concepts Used: Streaming data, message aggregation.
    - Example: Services send metrics.cpu.usage, and a collector stores them in Prometheus.
- Distributed Key-Value Store (Medium-Hard)
    - Description: Implement a key-value store where clients publish set and get requests, and a NATS subscriber manages state.
    - Concepts Used: Stateful processing, simple in-memory database.
    - Example: SET key1 value1 is sent to store.set.key1, and GET key1 fetches from store.get.key1.
- Remote Procedure Call (RPC) System (Hard)
    - Description: Implement a microservices RPC framework using NATS Request-Reply. Services expose functions via NATS subjects.
    - Concepts Used: Synchronous messaging, request-response pattern.
    - Example: rpc.math.add receives { "a": 5, "b": 10 } and responds with 15.
- Real-Time Stock Price Updater (Hard)
    - Description: A system where stock prices are published in real-time, and traders subscribe to updates.
    - Concepts Used: High-frequency data streaming, multi-client pub/sub.
    - Example: stocks.AAPL streams prices every second, and multiple subscribers receive real-time updates.
- Fault-Tolerant Service Registry (Very Hard)
    - Description: Build a lightweight service discovery system where microservices register themselves and other services query their availability.
    - Concepts Used: Dynamic service registration, health checks.
    - Example: Services register with services.register, and clients query services.lookup.serviceA.
- Distributed Load Balancer (Very Hard)
    - Description: Implement a load balancer where requests are published to NATS, and backend services consume requests based on availability.
    - Concepts Used: Dynamic workload distribution, custom load balancing.
    - Example: A requests.api queue where backend workers dynamically scale and pick requests.

### NATS Jetstresm demo todos:
- Persistent Event Logger (Easy)
    - Description: A system where services log events (errors, warnings, info) to a JetStream stream, ensuring they are stored and replayable.
    - Concepts Used: Streams, durable message storage, basic consumer subscription.
    - Example: Services publish logs to LOGS, and a subscriber writes them to a file/database.
- At-Least-Once Message Queue (Easy)
    - Description: A task queue where messages persist even if the worker crashes, ensuring reliable processing.
    - Concepts Used: Work queues, durable consumers, message acknowledgment.
    - Example: A tasks.process queue where workers pick jobs and acknowledge after processing.
- Delayed Message Processing (Medium)
    - Description: A system where messages are published but only consumed after a delay (e.g., scheduled notifications).
    - Concepts Used: Message retention, time-based pull consumers.
    - Example: A notifications.email queue that delays sending based on a timestamp.
- Replayable Event Store (Medium)
    - Description: A system where past events can be replayed for debugging or reprocessing failures.
    - Concepts Used: Stream replay, consumer offsets.
    - Example: A PAYMENTS stream where a service can replay failed transactions.
- User Activity Tracker (Medium-Hard)
    - Description: A system where user activity events (logins, clicks, purchases) are stored and can be queried later.
    - Concepts Used: Retained streams, partitioning by user ID.
    - Example: Store USER_ACTIVITY events and replay them for analytics.
- Dead Letter Queue (Hard)
    - Description: A system where failed messages (after multiple retries) are moved to a separate stream for later inspection.
    - Concepts Used: Message retries, filtering, stream forwarding.
    - Example: Messages in ORDERS failing 3 times move to ORDERS_DLQ for manual review.
- Multi-Region Replicated Messaging (Hard)
    - Description: A JetStream setup where messages are automatically replicated across multiple data centers.
    - Concepts Used: Stream mirroring, cross-region durability.
    - Example: A PAYMENTS stream replicated between us-east and eu-west.
- Distributed Workflow Orchestrator (Very Hard)
    - Description: A system where workflows are modeled as sequences of steps, with retries and timeouts.
    - Concepts Used: Message sequencing, durable consumers, scheduled tasks.
    - Example: Order processing with ORDER_RECEIVED -> PAYMENT_PROCESSED -> ORDER_SHIPPED.
- Event-Sourced Microservices (Very Hard)
    - Description: A microservices system where every state change is stored as an event, enabling auditability and CQRS.
    - Concepts Used: Event sourcing, multiple streams for state reconstruction.
    - Example: A USER_EVENTS stream tracks all account changes, allowing rebuilding the user state.
- High-Throughput Market Data System (Extremely Hard)
    - Description: A real-time market data streaming service where stock price updates are efficiently stored and served to thousands of clients.
    - Concepts Used: High-frequency ingestion, multiple consumers, snapshot/restore.
    - Example: STOCKS.AAPL publishes prices, and different traders subscribe with varying replay needs.

