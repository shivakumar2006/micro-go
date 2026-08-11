# 🚗 FleetOps

### Distributed Vehicle Commerce Platform

FleetOps is a **Go-based distributed vehicle commerce platform** built with a microservices architecture, API Gateway, Redis caching, PostgreSQL, Stripe payments, and Kafka-based asynchronous processing.

The project is focused on understanding and implementing real-world backend engineering concepts such as **service-to-service communication, resilience patterns, caching, asynchronous event processing, payment workflows, and inventory management**.

> 🚧 **Status:** Actively under development
> **Current milestone:** Core microservices + Kafka-based payment and inventory workflow

---

## 🏗️ Architecture

```text
                              ┌──────────────────────┐
                              │      React App       │
                              │      Frontend        │
                              └──────────┬───────────┘
                                         │
                                         ▼
                              ┌──────────────────────┐
                              │     API Gateway      │
                              │                      │
                              │ • Reverse Proxy      │
                              │ • Request Routing    │
                              │ • Rate Limiting      │
                              │ • Circuit Breaker    │
                              │ • JWT-aware Flow     │
                              └──────────┬───────────┘
                                         │
              ┌──────────────────────────┼──────────────────────────┐
              │                          │                          │
              ▼                          ▼                          ▼
      ┌──────────────┐           ┌──────────────┐           ┌──────────────┐
      │ Auth Service │           │Vehicle Service│          │ Cart Service │
      │              │           │              │           │              │
      │ • JWT        │           │ • Vehicles   │           │ • Cart CRUD  │
      │ • RBAC       │           │ • Pricing    │           │ • Redis      │
      │ • Auth       │           │ • Stock      │           │ • Vehicle    │
      └──────────────┘           └──────┬───────┘           │   Client     │
                                        │                   └──────┬───────┘
                                        │                          │
                                        │                    ┌─────▼─────┐
                                        │                    │   Redis   │
                                        │                    └───────────┘
                                        │
                                ┌───────▼────────┐
                                │  Order Service │
                                │                │
                                │ • Orders       │
                                │ • Order Items  │
                                │ • Order Status │
                                └───────┬────────┘
                                        ▲
                                        │
                                        │ Get Order + Items
                                        │
                         ┌──────────────┴──────────────┐
                         │                             │
                         │                             │
                  ┌──────┴───────┐              ┌──────▼────────┐
                  │   Payment    │              │   Inventory   │
                  │   Service    │              │    Service    │
                  │              │              │               │
                  │ • Stripe     │              │ • Kafka       │
                  │ • Payments   │              │   Consumer    │
                  └──────┬───────┘              │ • Stock Flow  │
                         │                      └──────┬────────┘
                         │                             │
                  Payment Success                      │
                         │                             │
                         ▼                             │
                  ┌──────────────┐                     │
                  │    Kafka     │─────────────────────┘
                  │              │
                  │ Payment      │
                  │ Success Event│
                  └──────────────┘

Inventory Service communicates directly with:

        Inventory ──────────────► Order Service
                    Get Order / Items

        Inventory ──────────────► Vehicle Service
                    Decrease Stock
```

---

# 🧩 Services

FleetOps currently consists of **6 backend services**, an API Gateway, and a React frontend.

| Component             | Responsibility                                                         |
| --------------------- | ---------------------------------------------------------------------- |
| **Auth Service**      | Authentication, JWT, user identity, and role-based access              |
| **Vehicle Service**   | Vehicle management, pricing, details, availability, and stock          |
| **Cart Service**      | User cart management, Redis caching, and Vehicle Service communication |
| **Order Service**     | Order and order-item management                                        |
| **Payment Service**   | Payment processing through Stripe                                      |
| **Inventory Service** | Event-driven inventory processing using Kafka                          |
| **API Gateway**       | Central entry point, routing, rate limiting, and resilience            |
| **React Frontend**    | Customer and admin-facing application                                  |

---

# 🔐 Authentication

FleetOps uses **JWT-based authentication** with access and refresh tokens.

### Authentication Flow

```text
Client
   │
   ▼
Auth Service
   │
   ├── Access Token
   └── Refresh Token
```

JWT claims contain information such as:

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "admin",
  "token_type": "access"
}
```

### Roles

- `customer`
- `admin`

The authenticated user identity is propagated through the request flow and used for authorization.

---

# 🚪 API Gateway

The API Gateway acts as the **single entry point** for the frontend.

```text
Frontend
    │
    ▼
API Gateway
    │
    ├── Auth Service
    ├── Vehicle Service
    ├── Cart Service
    ├── Order Service
    ├── Payment Service
    └── Inventory Service
```

### Responsibilities

- Reverse proxy
- Request routing
- Authentication-aware request flow
- Rate limiting
- Circuit breaker
- Service isolation

The gateway prevents clients from directly depending on individual internal service endpoints.

---

# 🛡️ Resilience

FleetOps implements resilience patterns for distributed service communication.

## Retry

Transient failures are handled using retry mechanisms with exponential backoff.

Example:

```text
Request
   │
   ▼
Failure
   │
   ▼
Retry #1
   │
   ▼
Failure
   │
   ▼
Retry #2
   │
   ▼
Failure
   │
   ▼
Retry #3
```

Retries are intended for transient failures rather than blindly retrying every error.

---

## Circuit Breaker

Circuit breakers prevent an unhealthy downstream service from continuously receiving requests.

```text
             ┌──────────────┐
             │    CLOSED    │
             └──────┬───────┘
                    │
              Failures exceed
                  threshold
                    │
                    ▼
             ┌──────────────┐
             │     OPEN     │
             └──────┬───────┘
                    │
                 Cooldown
                    │
                    ▼
             ┌──────────────┐
             │  HALF-OPEN   │
             └──────┬───────┘
                    │
             ┌──────┴──────┐
             │             │
          Success        Failure
             │             │
             ▼             ▼
          CLOSED         OPEN
```

This helps isolate failures between services.

---

# 🚘 Vehicle Service

Vehicle Service owns vehicle-related information and inventory state.

### Responsibilities

- Vehicle creation
- Vehicle retrieval
- Vehicle updates
- Vehicle deletion
- Vehicle pricing
- Vehicle details
- Vehicle availability
- Stock management

Vehicle data is persisted in PostgreSQL.

The service is consumed internally by other services, including:

```text
Cart Service
     │
     └──────► Vehicle Service
              Get Vehicle Details


Inventory Service
     │
     └──────► Vehicle Service
              Decrease Vehicle Stock
```

---

# 🛒 Cart Service

Cart Service manages individual user shopping carts.

### Responsibilities

- Add vehicle to cart
- Update cart quantity
- Remove cart items
- Delete cart
- Calculate cart total
- Count cart items
- Retrieve user cart
- Communicate with Vehicle Service
- Redis caching
- Cache invalidation

### Add to Cart Flow

```text
Client
  │
  ▼
API Gateway
  │
  ▼
Cart Service
  │
  ▼
Vehicle Service
  │
  ├── Vehicle exists?
  ├── Price
  └── Stock available?
  │
  ▼
Cart Database
  │
  ▼
Invalidate Redis Cache
```

The Cart Service does not blindly trust client-provided vehicle information. Vehicle information is retrieved from the Vehicle Service before adding an item to the cart.

---

# ⚡ Redis Caching

Redis is used to cache frequently accessed cart information.

Current cache keys include:

```text
cart:{user_id}
cart:total:{user_id}
cart:count:{user_id}
```

### Cache-Aside Flow

```text
GET Cart
   │
   ▼
Redis
   │
   ├──────── HIT ────────► Return Cached Data
   │
   └──────── MISS
          │
          ▼
      PostgreSQL
          │
          ▼
      Store in Redis
          │
          ▼
      Return Response
```

### Cache Invalidation

Cart mutations invalidate the relevant cache entries.

```text
Add / Update / Delete Cart
          │
          ▼
    PostgreSQL Update
          │
          ▼
     Redis Invalidate
```

This prevents stale cart data from being served after mutations.

---

# 📦 Order Service

Order Service is responsible for managing orders and their associated items.

An order contains information such as:

- Order ID
- User ID
- Order status
- Total amount
- Order items
- Vehicle IDs
- Quantities
- Prices

The Order Service acts as the source of truth for order information.

---

# 💳 Payment Service

Payment processing is implemented using **Stripe**.

### Payment Flow

```text
Customer
   │
   ▼
Payment Service
   │
   ▼
Stripe
   │
   ▼
Payment Success
   │
   ▼
Kafka Producer
```

When payment succeeds, the Payment Service publishes a payment-success event to Kafka.

This decouples payment processing from downstream inventory processing.

---

# 📨 Kafka Event-Driven Architecture

Kafka is currently used for asynchronous communication between the Payment and Inventory workflows.

### Payment Success Event

```text
Payment Service
      │
      │ Payment Success
      ▼
Kafka Producer
      │
      ▼
┌──────────────────────────┐
│          Kafka            │
│                          │
│    PaymentSuccess Event  │
└────────────┬─────────────┘
             │
             ▼
     Inventory Consumer
```

The Payment Service does not need to synchronously execute the entire inventory workflow after payment succeeds.

Instead, it publishes an event and allows the Inventory Service to process the event asynchronously.

---

# 📦 Inventory Service

Inventory Service is an **event-driven consumer** responsible for processing inventory changes after successful payments.

It has direct communication with **both Order Service and Vehicle Service**.

### Inventory Flow

```text
                     Kafka
                       │
               PaymentSuccess
                       │
                       ▼
              ┌─────────────────┐
              │ Inventory       │
              │ Service         │
              └───────┬─────────┘
                      / \
                     /   \
                    /     \
                   ▼       ▼
          Order Service   Vehicle Service
                │               │
                │               │
          Get Order +       Decrease
          Order Items        Stock
```

### Processing Steps

1. Inventory Service consumes the `PaymentSuccess` event.
2. It extracts the `order_id`.
3. It communicates with Order Service.
4. Order Service returns the order and its items.
5. Inventory Service extracts:
   - `vehicle_id`
   - `quantity`

6. Inventory Service communicates directly with Vehicle Service.
7. Vehicle Service decreases the corresponding vehicle stock.

### Important Architecture Boundary

Inventory does **not** get vehicle IDs directly from the payment event.

Instead:

```text
PaymentSuccess
      │
      ▼
Inventory
      │
      ▼
Order Service
      │
      ▼
Order Items
      │
      ├── vehicle_id
      └── quantity
      │
      ▼
Inventory
      │
      ▼
Vehicle Service
      │
      ▼
Decrease Stock
```

This keeps order information owned by the Order Service while vehicle inventory remains owned by the Vehicle Service.

---

# 🔄 Synchronous vs Asynchronous Communication

FleetOps currently uses both communication models.

## Synchronous

Used when a service needs an immediate response.

Example:

```text
Cart Service
     │
     │ HTTP
     ▼
Vehicle Service
     │
     ▼
Vehicle Details
```

Another example:

```text
Inventory Service
     │
     │ HTTP
     ▼
Order Service
     │
     ▼
Order + Items
```

## Asynchronous

Used for event-driven workflows.

```text
Payment Service
     │
     │ PaymentSuccess
     ▼
Kafka
     │
     ▼
Inventory Service
```

This combination allows FleetOps to use synchronous communication where immediate data is required and asynchronous communication where decoupling is beneficial.

---

# 🗄️ Data Layer

FleetOps uses **PostgreSQL** for persistent application data.

The services are organized around their own domain responsibilities rather than treating the entire system as a single monolithic application.

Redis acts as the caching layer, while Kafka provides asynchronous event transport.

---

# 🧱 Technology Stack

## Backend

- Go
- REST APIs
- PostgreSQL
- Redis
- Apache Kafka

## Frontend

- React
- Redux Toolkit
- RTK Query
- Tailwind CSS

## Payments

- Stripe

## Infrastructure

- Docker
- Docker Compose

## Architecture

- Microservices
- API Gateway
- Event-Driven Architecture
- Synchronous HTTP communication
- Asynchronous Kafka communication
- Retry with exponential backoff
- Circuit Breaker
- Rate Limiting
- Redis Cache

---

# 📊 Architecture Overview

```text
                         ┌──────────────────┐
                         │     Frontend     │
                         │      React       │
                         └────────┬─────────┘
                                  │
                                  ▼
                         ┌──────────────────┐
                         │   API Gateway    │
                         └────────┬─────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
        ▼                         ▼                         ▼
   ┌─────────┐              ┌──────────┐              ┌─────────┐
   │  Auth   │              │ Vehicle  │◄─────────────│  Cart   │
   │ Service │              │ Service  │              │ Service │
   └─────────┘              └────┬─────┘              └────┬────┘
                                  │                         │
                                  │                    ┌────▼────┐
                                  │                    │  Redis  │
                                  │                    └─────────┘
                                  │
                            ┌─────▼─────┐
                            │   Order   │
                            │  Service  │
                            └─────▲─────┘
                                  │
                                  │ Get Order + Items
                                  │
                            ┌─────┴──────┐
                            │ Inventory  │
                            │  Service   │
                            └─────▲───┬──┘
                                  │   │
                         Consume   │   │ Decrease Stock
                                  │   │
                             ┌────┴┐  │
                             │Kafka│  │
                             └──▲──┘  │
                                │     │
                         Payment│     │
                         Success│     │
                                │     │
                         ┌──────┴─┐   │
                         │Payment │   │
                         │Service │   │
                         └───┬────┘   │
                             │        │
                           Stripe     │
                                      │
                                      ▼
                               Vehicle Service
```

---

# 🧪 Current Engineering Practices

FleetOps currently implements:

- JWT authentication
- Access and refresh tokens
- Role-based authorization
- API Gateway
- Reverse proxy routing
- Rate limiting
- Circuit breaker
- Retry with exponential backoff
- Redis caching
- Cache invalidation
- PostgreSQL persistence
- Synchronous inter-service HTTP communication
- Kafka producer/consumer architecture
- Stripe payment processing
- Event-driven inventory processing
- Order → Inventory → Vehicle workflow
- Service responsibility boundaries

---

# 🚧 Roadmap

The following capabilities are planned for future milestones.

## Reliability

- [ ] Outbox Pattern
- [ ] Idempotent event processing
- [ ] Saga / compensating transactions
- [ ] Dead Letter Queue (DLQ)
- [ ] Improved inventory concurrency handling
- [ ] Better event retry and recovery mechanisms

## Observability

- [ ] Prometheus
- [ ] Grafana
- [ ] Jaeger / OpenTelemetry
- [ ] Distributed tracing
- [ ] Service-level metrics
- [ ] Kafka consumer metrics

## Infrastructure

- [ ] Production Docker configuration
- [ ] Kubernetes deployment
- [ ] CI/CD pipeline
- [ ] Container registry
- [ ] Cloud deployment

## Frontend

- [ ] Complete payment UI flow
- [ ] Customer checkout experience
- [ ] Order history
- [ ] Order status tracking
- [ ] Final UI/UX polish

---

# 🎯 Goals

FleetOps is being built as a practical exploration of modern backend and distributed-system engineering.

The project focuses on understanding:

- How microservices communicate
- When to use synchronous vs asynchronous communication
- How Redis improves read performance
- How Kafka decouples services
- How payment workflows trigger downstream processes
- How inventory should be coordinated with orders
- How resilience patterns protect distributed services
- How service boundaries affect system design
- How failures propagate through distributed systems

The goal is not simply to create multiple services, but to understand the **engineering problems introduced when a system becomes distributed**.

---

# 📌 Current Status

### Completed

- [x] Auth Service
- [x] JWT Authentication
- [x] Role-Based Authorization
- [x] API Gateway
- [x] Vehicle Service
- [x] Cart Service
- [x] Redis Caching
- [x] Order Service
- [x] Stripe Payment Service
- [x] Inventory Service
- [x] Kafka Producer
- [x] Kafka Consumer
- [x] Payment → Kafka → Inventory workflow
- [x] Inventory → Order Service communication
- [x] Inventory → Vehicle Service communication
- [x] Vehicle stock update workflow
- [x] React Frontend
- [x] Retry / Circuit Breaker / Rate Limiting

### In Progress

- [ ] Outbox Pattern
- [ ] Idempotency
- [ ] Saga
- [ ] DLQ
- [ ] Observability
- [ ] Kubernetes
- [ ] CI/CD
- [ ] Cloud Deployment
- [ ] Final frontend integration

---

# 👨‍💻 Author

**Shiva**

Backend-focused developer building with **Go, microservices, distributed systems, event-driven architecture, Redis, Kafka, and cloud-native technologies.**

---

> **FleetOps is built to learn by building — from REST APIs and microservices to asynchronous event processing and production-oriented distributed system design.**
