# ⚡ Smart Energy Management Platform

A **microservices-based** platform built with **Golang** to manage and process real-time electricity consumption data from IoT loggers.
This project is designed as a **proof of work** to showcase skills in backend architecture, gRPC communication, event-driven processing, caching, monitoring, and containerized deployments.

---

## 🚀 Features

- **Microservices Architecture** – Independent services for ingestion, billing, analytics, etc.
- **gRPC Communication** – Fast and type-safe inter-service communication.
- **RabbitMQ** – Message broker for asynchronous processing.
- **Redis** – High-speed cache for recent readings and bills.
- **MongoDB** – Persistent storage for readings and billing data.
- **Prometheus + Grafana** – Observability stack for metrics collection and visualization.
- **Docker & Docker Compose** – Local containerized development environment.
- **Mostly Go Standard Library** – Minimal third-party dependencies for routing, concurrency, and networking.

---

## 📂 Folder Structure

```bash
smart-energy-platform/
│
├── api-gateway/             # gRPC/HTTP gateway
│   ├── main.go
│   ├── proto/
│   ├── handlers/
│   ├── config/
│   └── Dockerfile
│
├── ingestion-service/       # Receives logger data & pushes to RabbitMQ
│   ├── main.go
│   ├── proto/
│   ├── rabbitmq/
│   ├── config/
│   └── Dockerfile
│
├── billing-service/         # Consumes from RabbitMQ, stores bills in MongoDB
│   ├── main.go
│   ├── proto/
│   ├── mongo/
│   ├── redis/
│   └── Dockerfile
│
├── analytics-service/       # Exposes metrics to Prometheus
│   ├── main.go
│   ├── prometheus/
│   ├── mongo/
│   ├── redis/
│   └── Dockerfile
│
├── common/                  # Shared packages
│   ├── models/
│   ├── logger/
│   ├── utils/
│   └── config/
│
├── deployments/
│   ├── docker-compose.yml
│   ├── prometheus.yml
│   ├── grafana-provisioning/
│   └── README.md
│
├── scripts/
│   ├── build.sh
│   ├── run.sh
│   └── test.sh
│
├── README.md
└── Makefile
```

---

## 🔄 Service Interaction Flow

1. **Logger → API Gateway**
   Device sends gRPC `SendReading` request.
2. **API Gateway → Ingestion Service**
   Forwards request via gRPC.
3. **Ingestion Service → RabbitMQ**
   Publishes message `{ deviceID, timestamp, usage_kWh }`.
4. **Billing Service**
   Subscribes to RabbitMQ, processes usage, stores in MongoDB, updates Redis cache.
5. **Analytics Service**
   Aggregates usage stats, exposes `/metrics` for Prometheus.
6. **Prometheus → Grafana**
   Metrics collected and visualized in dashboards.

---

## 🛠️ Technology Stack

| Tool         | Purpose |
|--------------|---------|
| **Go stdlib** | HTTP server, gRPC, concurrency, JSON encoding |
| **gRPC**     | Service-to-service communication |
| **RabbitMQ** | Event queue for ingestion → processing |
| **Redis**    | Fast in-memory caching |
| **MongoDB**  | Persistent storage |
| **Prometheus** | Metrics collection |
| **Grafana**  | Dashboard visualization |
| **Docker**   | Containerization |
| **Docker Compose** | Local orchestration |

---

## 🗺️ Roadmap

### Phase 1 – Base Setup
- Initialize repo and services.
- Write `docker-compose.yml` for RabbitMQ, Redis, MongoDB, Prometheus, Grafana.
- Define `.proto` files for gRPC services.

### Phase 2 – Core Services
- Implement **API Gateway** with gRPC endpoints.
- Implement **Ingestion Service** to push messages to RabbitMQ.
- Implement **Billing Service** to consume RabbitMQ & write to MongoDB.

### Phase 3 – Observability
- Implement **Analytics Service** exposing Prometheus metrics.
- Configure Grafana dashboards.

### Phase 4 – Extras
- Add Redis caching in Billing & Analytics.
- Add authentication to API Gateway.
- Write integration tests.

---

## 📦 Deployment

**Start the environment:**
```bash
docker-compose up -d
```

View services:
-	RabbitMQ Management: http://localhost:15672
-	Prometheus: http://localhost:9090
-	Grafana: http://localhost:3000
-	MongoDB: mongodb://localhost:27017
-	Redis: redis://localhost:6379

---

## 📜 License

This project is open-source and available under the MIT License.
