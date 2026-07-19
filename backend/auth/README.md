                              CLIENT
                                 │
                                 │ HTTP Request
                                 ▼
                        ┌──────────────────┐
                        │   Chi Router     │
                        └────────┬─────────┘
                                 │
                                 ▼
                      Authentication Middleware
                                 │
                    Validate Access Token (JWT)
                                 │
                                 ▼
                      Authorization Middleware
                         RequireRole("admin")
                                 │
                                 ▼
                          HTTP Handlers
                                 │
                 Parse Request / Validate Input
                                 │
                                 ▼
                           Auth Service
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        │ Business Rules     │ JWT Manager        │ Password Hashing
        │                    │                    │
        └────────────────────┼────────────────────┘
                             │
                             ▼
                        Repository Layer
                             │
                      SQL Queries (Postgres)
                             │
                             ▼
                        PostgreSQL Database
