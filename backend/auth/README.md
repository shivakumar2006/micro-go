## Overall authentication service structure

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

## Request Lifecycle

Client

│

▼

HTTP Request

│

▼

Router

│

▼

JWT Authentication Middleware

│

▼

RequireRole Middleware (optional)

│

▼

Handler

│

▼

Service

│

▼

Repository

│

▼

PostgreSQL

│

▼

Repository

│

▼

Service

│

▼

Handler

│

▼

JSON Response

## Register Flow

Client

│

▼

POST /register

│

▼

Validate Request

│

▼

Check Email Exists

│

▼

Hash Password

│

▼

Normalize Role

│

▼

Get Role ID

│

▼

Create User

│

▼

Generate JWT Pair

│

▼

Hash Refresh Token

│

▼

Store Refresh Token

│

▼

Return

Access Token
Refresh Token
User

## Login Flow

Client

│

▼

POST /login

│

▼

Find User By Email

│

▼

Compare Password

│

▼

Generate Access Token

│

▼

Generate Refresh Token

│

▼

Hash Refresh Token

│

▼

Store Session

│

▼

Return Tokens

## Refresh Flow

Client

│

▼

POST /refresh

│

▼

Validate Refresh Token

│

▼

Find Session

│

▼

Check Expiry

│

▼

Find User

│

▼

Delete Old Refresh Token

│

▼

Generate New Tokens

│

▼

Store New Refresh Token

│

▼

Return New Token Pair

## authentication flow

Client

│

Authorization Header

│

▼

Auth Middleware

│

▼

Extract Bearer Token

│

▼

Validate JWT

│

▼

Claims

(UserID
Email
Role)

│

▼

Store Claims In Context

│

▼

Next Handler

## RBAC (Role Based Access Control)

Request

│

▼

Auth Middleware

│

▼

JWT Claims

│

▼

Role

│

▼

RequireRole("admin")

      │

┌────┴─────┐

│ │

Yes No

│ │

▼ ▼

Handler 403 Forbidden

## Token Architecture

                 Login/Register

                       │

                       ▼

             Generate Token Pair

           ┌────────────┴─────────────┐

           │                          │

           ▼                          ▼

     Access Token              Refresh Token

(15 Minutes) (7 Days)

           │                          │

           ▼                          ▼

Sent To Client Hash & Store In DB

                                      │

                                      ▼

                               user_sessions
