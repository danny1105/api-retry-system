# API Retry System

## Overview

This project implements a **message delivery system** between two microservices using a queue-based architecture and a REST endpoint.

The system ensures:

* **No message loss**
* **At-least-once delivery**
* **Strict ordering (FIFO)**
* **Infinite retry on failure**

---

## Architecture

```
Producer → Kafka (Redpanda) → Microservice-1 → Microservice-2
```

### Components

* **Producer**

  * Publishes messages to Kafka topic

* **Microservice-1 (Consumer + Retry Engine)**

  * Consumes messages sequentially
  * Calls Microservice-2 via POST API
  * Retries indefinitely every 10 seconds on failure
  * Commits Kafka offset only after successful delivery

* **Microservice-2 (API Service)**

  * Exposes POST `/events`
  * Processes incoming messages
  * Ensures **idempotency** using message ID

* **Redpanda (Kafka-compatible broker)**

  * Provides FIFO message queue

---

## Key Design Decisions

### 1. Sequential Processing (Ordering Guarantee)

Messages are processed one at a time to preserve strict ordering.

> Tradeoff: Throughput is reduced, but ordering is guaranteed.

---

### 2. Blocking Retry Mechanism

If a message fails:

* Microservice-1 retries **indefinitely every 10 seconds**
* Next message is NOT consumed until current succeeds

> Ensures no message loss and maintains ordering.

---

### 3. Offset Commit After Success

Kafka offsets are committed **only after successful API delivery**.

> Guarantees at-least-once delivery and prevents message loss.

---

### 4. Idempotent API (Microservice-2)

Since retries may cause duplicate requests:

* Each message contains a unique `id`
* Processed IDs are tracked (in-memory for this demo)
* Duplicate requests are safely ignored

---

## Prerequisites

* Docker
* Docker Compose

---

## How to Run

```bash
docker-compose up --build
```

---

## Producing Messages

In a new terminal:

```bash
docker-compose exec producer ./app
```

---

## Failure Simulation (Important)

1. Start the system
2. Stop API service:

   ```bash
   docker-compose stop microservice-2
   ```
3. Produce messages:

   ```bash
   docker-compose exec producer ./app
   ```
4. Observe retries in logs
5. Restart API:

   ```bash
   docker-compose start microservice-2
   ```

Messages will be delivered successfully after recovery.

---

## Logs You Should See

```text
Received message: 1
Retry in 10s for message: 1
Delivered message: 1
Offset committed for message: 1
```

---

## Guarantees

* **At-least-once delivery**
* **No message loss**
* **Ordering preserved**
* **Eventual consistency**

---

## Notes

* Idempotency is implemented using an in-memory store for simplicity
* In production, this should be replaced with a persistent datastore

---

## Tech Stack

* Go (Golang)
* Kafka (via Redpanda)
* Docker & Docker Compose

---

## Summary

This system demonstrates a reliable, fault-tolerant message delivery pipeline with strong consistency guarantees and clear tradeoffs.

---
