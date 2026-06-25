# webhook-event-processor

A Redpanda consumer service written in Go that processes inventory webhook events and persists them to PostgreSQL — built as a study project with production-grade patterns in mind.

## Overview

`webhook-event-processor` is one part of a two-service pipeline. An upstream webhook receiver (separate service) receives HTTP events from external systems and publishes them to a Redpanda topic. This service sits downstream: it consumes those messages, validates the payload, creates or updates an inventory record, and appends an audit entry to the `inventory_tracks` table.

The separation is intentional — receiving webhooks and processing their events are distinct responsibilities with different scaling and failure characteristics. Keeping them in separate services makes each one easier to reason about, deploy, and scale independently.

The project was built for learning purposes, with a deliberate focus on patterns commonly found in production Go services: dependency injection via context, separation of concerns across packages, graceful shutdown with signal handling, and a clean worker dispatch layer.

## Architecture & Design Decisions

### Flow

```
External system (HTTP POST)
        │
        ▼
  webhook-receiver                  ← upstream service (separate project)
        │  publishes to topic
        ▼
Redpanda (event.on_demand)
        │
        ▼
  Consumer.Poll()                   ← blocks until message arrives or signal received
        │
        ▼
Worker.ProcessWebhookEvent()        ← validates payload, dispatches to services
        │
        ├──▶ Service.Inventory.Create()        ──▶ PostgreSQL (inventory)
        │
        └──▶ Service.InventoryTracks.Create()  ──▶ PostgreSQL (inventory_tracks)
```

### Key Decisions

**Redpanda over a generic message broker**
Redpanda is Kafka-compatible but ships as a single binary with no ZooKeeper dependency, making it significantly easier to run locally and in simple environments. For a study project, this lowers the operational burden without sacrificing the Kafka API surface.

**`twmb/franz-go` as the Kafka client**
Pure Go, no CGO, no JVM. It has a clean API and is the de facto choice for high-performance Kafka clients in Go. `confluent-kafka-go` was explicitly avoided because it wraps a C library, introducing external build dependencies.

**Dependency injection via `contexts/`**
Rather than using global state or `init()` functions, dependencies (DB, cache, logger) are wired at startup into a `Context` struct and passed explicitly. This makes the data flow traceable and components independently testable.

**Separated `worker`, `service`, and `redpanda` packages**
The worker layer handles orchestration (validate → dispatch), the service layer owns business rules (what to persist), and the redpanda package owns all queue concerns. Each layer can evolve independently.

**Graceful shutdown via `signal.NotifyContext`**
The consumer's poll loop receives a context that is cancelled on `SIGINT`/`SIGTERM`. This ensures in-flight messages are not dropped when the process is stopped.

## Tech Stack

| Technology | Role | Why |
|---|---|---|
| Go 1.25 | Language | Strong concurrency primitives, fast compile, suited for long-running services |
| Redpanda | Message queue | Kafka-compatible, single binary, no ZooKeeper |
| `twmb/franz-go` | Kafka client | Pure Go, no CGO, clean and performant API |
| PostgreSQL | Primary database | Relational model fits inventory with audit trail requirements |
| Redis / Dragonfly | Cache | Fast lookups, reduces DB reads for hot records |
| Cobra | CLI | Standard Go CLI framework; allows clean multi-command structure |
| logrus | Logging | Structured logging with fields, compatible with log aggregators |
| godotenv | Config | Simple `.env`-based configuration for local development |
| Docker Compose | Local infrastructure | Spins up Redpanda, PostgreSQL, Dragonfly, and Redpanda Console with a single command |

## Project Structure

```
webhook-event-processor/
├── cmd/
│   └── event_processor.go    # CLI command: starts the consumer
├── common/
│   ├── helper.go              # Utility functions (env vars, UUID generation)
│   └── system_constants.go   # Topic name and consumer group ID
├── contexts/
│   ├── context.go             # Root context — wires all dependencies together
│   └── logs.go                # Logger initialization
├── internal/
│   ├── cache/
│   │   └── dragonfly.go       # Redis/Dragonfly client setup
│   ├── db/
│   │   ├── db.go              # PostgreSQL connection
│   │   ├── cache.go           # DB-level cache abstraction
│   │   ├── db_config.go       # Connection config from environment variables
│   │   └── migrations/        # SQL migration files
│   ├── redpanda/
│   │   ├── client.go          # Raw kgo client initialization
│   │   ├── consumer.go        # Poll loop and message deserialization
│   │   ├── event.go           # Event struct and Kafka record mapping
│   │   ├── producer.go        # Synchronous producer
│   │   └── redpanda.go        # Public facade (PublishTopic)
│   ├── service/
│   │   ├── inventory.go       # Inventory persistence logic
│   │   └── inventory_tracks.go# Audit trail persistence
│   └── worker/
│       └── process_event.go   # Validation and dispatch orchestration
└── main.go
```

## Quick Start

### Environment variables

| Variable | Description | Example |
|---|---|---|
| `BROKERS` | Comma-separated Redpanda broker addresses | `localhost:9092` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `secret` |
| `DB_NAME` | Database name | `webhook_db` |
| `SSL_MODE` | PostgreSQL SSL mode | `disable` |

### Running

```bash
# Copy and fill in your environment variables
cp .env.example .env

# Start infrastructure (Redpanda, PostgreSQL, Dragonfly, Redpanda Console)
docker compose -f scripts/docker/docker-compose.yml up -d

# Run database migrations
psql -U $DB_USER -d $DB_NAME -f internal/db/migrations/000001_create_base_tables.up.sql

# Start the consumer
go run main.go event-consumer
```

---

> Built as a study project, intentionally structured to reflect production patterns: clean package separation, explicit dependency injection, graceful shutdown, and a layered worker/service model.
