# Shortly: Event-Driven URL Shortener Microservices

**Shortly** is a portfolio project designed to showcase a robust, scalable, and event-driven microservices architecture built with **Go**. The system enables shortening long URLs and performing high-speed redirections through distributed caching and efficient inter-service communication.

## 🏗️ System Architecture

The project implements a decoupled architecture with the following components:

- **KrakenD (API Gateway):** A single entry point that manages routing, rate limiting, and service aggregation.
- **Shortener Service:** Responsible for generating unique short codes, persisting data in PostgreSQL, and publishing events.
- **Redirect Service:** Optimized for ultra-fast redirections. It uses a distributed caching strategy (Valkey) and synchronizes via events or direct gRPC lookups.
- **RabbitMQ:** Message broker for asynchronous communication and event propagation (pre-populating cache upon URL creation).
- **Valkey (Cache):** High-performance in-memory storage to minimize redirection latency.
- **PostgreSQL:** Source of truth for long-term data persistence.

### Data Flow

1. **Creation:** User sends a URL to the Gateway -> `shortener_svc` generates a code -> Persists to Postgres -> Publishes a "created" event to RabbitMQ.
2. **Synchronization:** `redirect_svc` consumes the event from RabbitMQ and stores the mapping in Valkey (Cache).
3. **Redirection:** User requests a code -> `redirect_svc` checks Valkey -> If a Cache Miss occurs, it queries `shortener_svc` via **gRPC** -> Updates cache and redirects the user.

---

## 🛠️ Tech Stack & Tools

- **Language:** Go (v1.25) utilizing Workspaces.
- **Communication:** gRPC (Protocol Buffers) & REST.
- **Messaging:** RabbitMQ 4 (Event-Driven Architecture).
- **Databases:** PostgreSQL 17 & Valkey (Redis-compatible).
- **API Gateway:** KrakenD.
- **Infrastructure:** Podman Compose / Docker Compose.
- **Design Patterns:** Clean Architecture, Repository Pattern, Cache-aside, Pub/Sub.

---

## 🚀 Getting Started

### Prerequisites
- Go 1.25+
- Docker / Podman.
- **Minikube** (for Kubernetes deployment).
- **kubectl** (Kubernetes CLI).

### Installation (Docker Compose)
1. Clone the repository.
2. Configure environment variables:
   ```bash
   cp .env.example .env
   ```
3. Start the infrastructure and services:
   ```bash
   make up
   ```

### ☸️ Kubernetes Deployment (Minikube)
For a production-like local environment, you can deploy the entire stack to Kubernetes:

1. **Start Minikube:**
   ```bash
   minikube start
   ```

2. **Build images directly in Minikube's Docker daemon:**
   ```bash
   eval $(minikube docker-env)
   docker build -t shortener_svc:latest ./src/shortener_svc
   docker build -t redirect_svc:latest ./src/redirect_svc
   ```

3. **Deploy all manifests:**
   ```bash
   kubectl apply -f k8s/
   ```

4. **Access the Gateway:**
   ```bash
   # Get the URL for the KrakenD service
   minikube service krakend --url
   ```
   *The gateway is exposed on NodePort 30080 by default.*

### 📦 Helm Deployment (Automated)
The most professional way to deploy this project is using **Helm**. It allows you to manage the entire stack as a single package.

1. **Install the Chart:**
   ```bash
   helm install shortly ./charts/shortly
   ```

2. **Verify the installation:**
   ```bash
   helm list
   kubectl get pods
   ```

3. **Upgrade configuration (optional):**
   If you change something in `values.yaml`, simply run:
   ```bash
   helm upgrade shortly ./charts/shortly
   ```

4. **Uninstall:**
   ```bash
   helm uninstall shortly
   ```

---

## 📖 API Documentation

The system exposes its services through the port configured in `KRAKEND_GATEWAY_PORT` (default is `8080`).

### 1. Create a Short URL
**Endpoint:** `POST /shortly/create`

**Request:**
```json
{
  "raw_url": "https://www.google.com"
}
```

**Response:**
```json
{
  "id": 1,
  "raw_url": "https://www.google.com",
  "short_code": "xY8z2A",
  "created_at": "2026-03-13T10:00:00Z"
}
```

### 2. Redirection
**Endpoint:** `GET /shortly/{short_code}`

**Example:** `GET http://localhost:8080/shortly/xY8z2A`
*Result: 302 Redirect to the original URL.*

---

## 📐 Design Decisions

- **gRPC vs REST:** Internal communication uses gRPC for its efficiency and strong typing, while the Gateway exposes REST for compatibility with web/mobile clients.
- **Valkey (Redis fork):** Chosen for its performance in read-intensive scenarios and open-source compatibility.
- **Event-Driven Cache:** The cache is pre-populated via RabbitMQ events to ensure redirections are near-instant immediately after creation, without waiting for the first manual lookup.

---

## 👨‍💻 Author
[Your Name/FreyreCorona]
- [LinkedIn](https://www.linkedin.com/in/your-profile)
- [GitHub](https://github.com/FreyreCorona)
