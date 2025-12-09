<div align="center">
  <img src="logo.png" alt="GoProjectSkeleton logo" height="400">

  *A foundation for clean architecture projects in Go*

  <!-- Badges -->
  <p>
    <img src="https://img.shields.io/badge/Go-1.25-blue?logo=go&logoColor=white" alt="Go Version">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
    <img src="https://img.shields.io/badge/Architecture-Clean%20Architecture-blueviolet" alt="Architecture">
    <img src="https://img.shields.io/badge/Pattern-Hexagonal%20Architecture-orange" alt="Pattern">
    <img src="https://img.shields.io/badge/Status-Production%20Ready-success" alt="Status">
    <img src="https://img.shields.io/badge/Serverless-AWS%20%7C%20Azure-lightgrey" alt="Serverless">
  </p>

  <p>
    <img src="https://img.shields.io/badge/GORM-1.25.12-blue" alt="GORM">
    <img src="https://img.shields.io/badge/Gin-1.10.0-cyan" alt="Gin">
    <img src="https://img.shields.io/badge/PostgreSQL-Supported-blue" alt="PostgreSQL">
    <img src="https://img.shields.io/badge/Redis-Supported-red" alt="Redis">
    <img src="https://img.shields.io/badge/JWT-Supported-yellow" alt="JWT">
    <img src="https://img.shields.io/badge/Swagger-Documented-green" alt="Swagger">
  </p>
</div>

## 🚀 Quick Start

```bash
# 1. Clone the repository
git clone <repository-url>
cd GoProjectSkeleton

# 2. Configure environment variables
cp dev.env.example dev.env
# Edit dev.env with your configurations

# 3. Start services with Docker
docker network create goprojectskeleton
docker volume create goprojectskeleton-db-data
docker-compose -f docker/docker-compose.dev.yml up -d

# 4. Run the application
go run src/infrastructure/server/cmd/main.go

# 5. Access Swagger documentation (independent service)
# Swagger runs as a separate service on port 8081
# http://localhost:8081/docs/
```

## 📋 Table of Contents

1. [Introduction](#introduction)
2. [Complete Development Environment](#-complete-development-environment)
3. [Project Architecture](#project-architecture)
4. [Scalability and Serverless](#scalability-and-serverless)
5. [Complete Request Flow](#complete-request-flow)
6. [Virtues and Benefits](#virtues-and-benefits)
7. [Project Structure - Layer by Layer](#project-structure---layer-by-layer)
8. [Exhaustive Review by Folders](#exhaustive-review-by-folders)
9. [Technologies and Dependencies](#technologies-and-dependencies)
10. [Configuration and Setup](#configuration-and-setup)
11. [Business Modules](#business-modules)
12. [API and Endpoints](#api-and-endpoints)
13. [Database and Persistence](#database-and-persistence)
14. [Authentication and Security](#authentication-and-security)
15. [Testing](#testing)
16. [Docker and Deployment](#docker-and-deployment)
17. [GitHub Actions Deployment](#github-actions-deployment)
18. [Development Guide](#development-guide)

---

## Introduction

**GoProjectSkeleton** is a professional Go project skeleton that implements **Clean Architecture** and **Hexagonal Architecture** (Ports & Adapters). This project provides a solid, scalable, and maintainable foundation for building enterprise applications in Go.

### Project Philosophy

The core philosophy of **GoProjectSkeleton** is that the **domain** and **application logic** (business logic) must be completely independent of **infrastructure** (frameworks, libraries, databases, etc.). This enables:

- **Interchangeability**: Change any infrastructure component without affecting the business
- **Testability**: Test business logic without external dependencies
- **Scalability**: Migrate to serverless or microservices architectures without rewriting code
- **Maintainability**: Clean and well-structured code

### Main Features

#### 🏗️ Architecture and Design
- ✅ **Clean Architecture** - Clear separation of responsibilities in layers (Domain, Application, Infrastructure)
- ✅ **Hexagonal Architecture** - Complete decoupling of external dependencies through Ports & Adapters
- ✅ **Dependency Inversion** - Inner layers do not depend on outer layers
- ✅ **SOLID Principles** - Design principles applied consistently
- ✅ **Repository Pattern** - Data access abstraction
- ✅ **Use Case Pattern** - Business logic encapsulation
- ✅ **Factory Pattern** - Instance creation with dependency injection

#### 🔐 Authentication and Security
- ✅ **Complete JWT Authentication** - Access tokens and refresh tokens with flexible configuration
- ✅ **OTP (One-Time Password)** - Two-factor authentication with temporary codes
- ✅ **Secure Password System** - Bcrypt hashing, password reset with tokens
- ✅ **Guards and Authorization** - Access control based on roles and permissions
- ✅ **Multi-layer Validation** - Validation in DTOs, use cases, and repositories
- ✅ **CORS Configured** - Security for web applications

#### 👥 User Management
- ✅ **Complete CRUD** - Create, read, update, and delete users
- ✅ **Role Management** - Role system with priorities
- ✅ **User States** - Pending, Active, Inactive, Suspended, Deleted
- ✅ **Account Activation** - Activation system through tokens
- ✅ **Pagination and Filtering** - Efficient queries with Query Payload
- ✅ **Smart Cache** - List caching with Redis

#### 🌐 Internationalization and Communication
- ✅ **Internationalization (i18n)** - Multi-language support (Spanish, English, extensible)
- ✅ **Email System** - HTML templates with dynamic rendering
- ✅ **Transactional Emails** - Registration, password reset, OTP
- ✅ **Locale per Request** - Each request can have its own language

#### 📚 Documentation and Testing
- ✅ **Independent Swagger Service** - Swagger runs as a separate, deployable service (port 8081)
- ✅ **Docker Registry Ready** - Swagger can be built and deployed independently
- ✅ **Complete Testing** - Unit, integration, and E2E tests
- ✅ **Complete Mocks** - Repository and provider mocks for testing
- ✅ **Bruno Collection** - Ready-to-use collection for E2E testing

#### 🛠️ Development Environment
- ✅ **Complete IDE Configuration** - Pre-configured VS Code/IDE settings for debugging
- ✅ **Hot Reload with Air** - Automatic code reloading on file changes
- ✅ **Remote Debugging** - Delve debugger configured for Docker and local development
- ✅ **Development Docker Setup** - Complete development environment with all tools pre-installed
- ✅ **Independent Swagger Service** - Swagger runs in separate container with hot reload
- ✅ **Development Tools Included** - Air, Delve, Swag pre-configured and ready to use
- ✅ **Email Testing** - Mailpit integrated for email testing in development
- ✅ **Redis Management UI** - Redis Commander for cache inspection

#### 🐳 DevOps and Deployment
- ✅ **Complete Docker** - Multi-service for development, test, and E2E
- ✅ **Serverless Ready** - Support for AWS Lambda and Azure Functions
- ✅ **Terraform** - Infrastructure as code for AWS and Azure
- ✅ **Secrets Management** - Integration with AWS Secrets Manager and Azure Key Vault
- ✅ **Hot Reload** - Efficient development with automatic reloading

#### ⚡ Performance and Scalability
- ✅ **Redis Cache** - Performance optimization with configurable TTL
- ✅ **Connection Pooling** - Database connection reuse
- ✅ **Pipes System (DAG)** - Sequential use case orchestration
- ✅ **Parallelization** - Concurrent execution of use cases with goroutines
- ✅ **Stateless Design** - Ready for horizontal scalability
- ✅ **Tree Shaking** - Automatic binary optimization in serverless

---

## 🛠️ Complete Development Environment

**GoProjectSkeleton** comes with a **fully configured development environment** that allows developers to start coding immediately without spending time on setup. Everything is pre-configured and ready to use.

### Development Environment Features

#### ✅ Pre-configured IDE Settings

The project includes complete IDE configuration files in the `IDE/` directory:

- **`launch.json`**: VS Code debugging configurations
  - **Attach to Docker Go (Delve)**: Remote debugging configuration for Docker containers
  - **Launch Go Program**: Local debugging with automatic Swagger generation
  - Pre-configured environment variables for development
  - Automatic path mapping for Docker debugging

- **`tasks.json`**: Build tasks
  - Automatic Swagger documentation generation
  - Pre-configured build commands

**Usage:**
```bash
# Copy IDE configuration to your workspace
cp -r IDE/.vscode .  # For VS Code
# Or use the IDE/ folder directly
```

#### ✅ Hot Reload with Air

The project uses **Air** for automatic code reloading during development:

```toml
# .air.toml configuration
[build]
cmd = "swag init && go build -gcflags 'all=-N -l' -o /app/tmp/main"
bin = "/app/tmp/main"
full_bin = "dlv exec /app/tmp/main --headless --listen=:40000"
include_ext = ["go", "tpl", "tmpl", "html"]
```

**Features:**
- ✅ Automatic rebuild on file changes
- ✅ Integrated with Delve debugger
- ✅ Watches Go files, templates, and HTML
- ✅ Colored and timestamped logs
- ✅ Polling mode for better file system compatibility

**How it works:**
1. Air watches for file changes
2. Automatically regenerates Swagger documentation
3. Rebuilds the application with debug symbols (`-N -l` flags)
4. Restarts the application with Delve debugger attached
5. Ready for remote debugging on port 40000

#### ✅ Remote Debugging with Delve

**Delve (dlv)** is pre-configured for both local and remote debugging:

**Docker Remote Debugging:**
```json
{
  "name": "Attach to Docker Go (Delve)",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "port": 40000,
  "host": "127.0.0.1"
}
```

**Features:**
- ✅ Debug running application in Docker
- ✅ Set breakpoints in VS Code
- ✅ Inspect variables and call stack
- ✅ Step through code execution
- ✅ Path mapping for seamless debugging

**Debugging Workflow:**
1. Start Docker development environment: `docker-compose -f docker/docker-compose.dev.yml up -d`
2. Application starts with Delve on port 40000
3. Attach debugger from VS Code using "Attach to Docker Go (Delve)"
4. Set breakpoints and debug as if running locally

#### ✅ Complete Docker Development Setup

The `docker-compose.dev.yml` includes everything needed for development:

**Services:**
- **Application**: Go application with hot reload and debugging
- **PostgreSQL**: Database on port 5436
- **Redis**: Cache server
- **Mailpit**: Email testing (port 8025 for UI, 1025 for SMTP)
- **Redis Commander**: Redis management UI (port 8081)

**Development Features:**
```yaml
services:
  goprojectskeleton:
    ports:
      - "8080:8080"    # Application
      - "40000:40000"  # Delve debugger
    volumes:
      - ../src:/app/src              # Live code sync
      - ../dev.env:/app/.env:ro       # Environment variables
      - ../.air.toml:/app/.air.toml   # Air configuration
    command: ["air", "-c", ".air.toml"]  # Hot reload
```

**Pre-installed Tools in Docker:**
- ✅ **Air** (v1.61.7) - Hot reload
- ✅ **Delve** (dlv) - Debugger
- ✅ **Swag** - Swagger generator
- ✅ All Go dependencies

#### ✅ Independent Swagger Documentation Service

Swagger documentation is now a **completely independent service** that runs separately from the main application:

**Architecture:**
- **Independent Server**: Swagger runs on its own HTTP server (port 8081)
- **Separate Deployment**: Can be deployed independently using `dockerfile.swagger`
- **Docker Registry Ready**: Can be pushed to any Docker registry and deployed separately
- **No Dependencies**: Completely decoupled from the main application

**Development:**
- Swagger service runs in a separate Docker container
- Automatically regenerates documentation on code changes
- Available at `http://localhost:8081/docs/`
- Configured via environment variables (see `dev.env.example`)

**Production Deployment:**
- Build: `docker build -f docker/dockerfile.swagger -t your-registry/swagger:latest .`
- Push: `docker push your-registry/swagger:latest`
- Run: `docker run -p 8081:8081 -e API_HOST=your-api-host your-registry/swagger:latest`

**Configuration:**
- `SWAGGER_PORT`: Port for Swagger server (default: 8081)
- `API_HOST`: Host of the main API (e.g., `api.example.com:8080`)
- `API_TITLE`: API title in Swagger UI
- `API_VERSION`: API version
- `API_DESCRIPTION`: API description
- `API_BASE_PATH`: Base path for API endpoints

#### ✅ Development Tools Integration

**Email Testing with Mailpit:**
- Web UI: `http://localhost:8025`
- SMTP: `localhost:1025`
- Test all transactional emails without real SMTP server

**Redis Management:**
- Redis Commander: `http://localhost:8081`
- Visual cache inspection
- Key browsing and editing

**Database Access:**
- PostgreSQL: `localhost:5436`
- Direct database access for debugging
- Persistent data with Docker volumes

### Quick Start Development

**1. Start Development Environment:**
```bash
# Create network and volume
docker network create goprojectskeleton
docker volume create goprojectskeleton-db-data

# Start all services
docker-compose -f docker/docker-compose.dev.yml up -d
```

**2. Configure IDE:**
```bash
# Copy IDE configuration (if using VS Code)
cp -r IDE/.vscode .
```

**3. Start Debugging:**
- Open VS Code
- Go to Run and Debug (F5)
- Select "Attach to Docker Go (Delve)"
- Set breakpoints and start debugging!

**4. Development Workflow:**
- Edit code in `src/` directory
- Air automatically detects changes
- Application rebuilds and restarts
- Debugger reattaches automatically
- See changes immediately

### Development Environment Architecture

```mermaid
graph TB
    subgraph Developer["👨‍💻 Developer"]
        IDE[VS Code/IDE<br/>Breakpoints & Debugging]
        Code[Source Code<br/>src/]
    end

    subgraph Docker["🐳 Docker Container"]
        Air[Air<br/>Hot Reload Watcher]
        Delve[Delve Debugger<br/>Port 40000]
        App[Go Application<br/>Port 8080]
    end

    subgraph Services["🔧 Development Services"]
        DB[(PostgreSQL<br/>Port 5436)]
        Redis[(Redis)]
        Mailpit[Mailpit<br/>Port 8025/1025]
        RedisUI[Redis Commander<br/>Port 8081]
    end

    Code -->|File Changes| Air
    Air -->|Rebuild & Restart| Delve
    Delve -->|Debug Session| App
    IDE -->|Attach Debugger| Delve
    App --> DB
    App --> Redis
    App -->|Send Emails| Mailpit

    style Air fill:#e3f2fd
    style Delve fill:#fff9c4
    style App fill:#c8e6c9
```

### Development Configuration Files

| File | Purpose |
|------|---------|
| `IDE/launch.json` | VS Code debugging configurations |
| `IDE/tasks.json` | Build tasks (Swagger generation) |
| `.air.toml` | Air hot reload configuration |
| `docker/dockerfile.dev` | Development Docker image |
| `docker/docker-compose.dev.yml` | Complete development stack |
| `dev.env.example` | Development environment variables template |

### Benefits for Developers

1. **Zero Setup Time**
   - Everything pre-configured
   - Start coding immediately
   - No manual tool installation needed

2. **Productive Debugging**
   - Remote debugging in Docker
   - Breakpoints work seamlessly
   - Variable inspection and call stack

3. **Fast Development Cycle**
   - Hot reload on file changes
   - No manual restarts needed
   - See changes instantly

4. **Complete Tooling**
   - Email testing without SMTP
   - Redis inspection UI
   - Database access ready

5. **Team Consistency**
   - Same environment for all developers
   - No "works on my machine" issues
   - Docker ensures consistency

---

## ☁️ Cloud and Serverless Capabilities

**GoProjectSkeleton** is designed to work in both traditional environments and modern serverless architectures.

### 🚀 AWS Lambda

The project includes complete support for **AWS Lambda** with:

- ✅ **Automatic Function Generation** - Generation system from `functions.json`
- ✅ **Independent Modules** - Each Lambda function has its own Go module
- ✅ **Optimized Tree Shaking** - 5-15 MB binaries vs 50+ MB without optimization
- ✅ **AWS Secrets Manager** - Automatic secret loading from Secrets Manager
- ✅ **Lambda Adapter** - Adapter for API Gateway events
- ✅ **Terraform** - Infrastructure as code ready to deploy
- ✅ **Optimized Compilation** - Specific flags for Lambda (`lambda.norpc`, `-ldflags="-s -w"`)

**Lambda functions structure:**
```
src/infrastructure/clouds/aws/
├── functions.json          # Function definitions
├── init.go                 # AWS initialization
├── lambda_adapter.go       # Lambda adapter
├── secrets_manager.go      # Secrets management
├── terraform/              # Terraform infrastructure
└── functions/              # Function generator
```

### 🔷 Azure Functions

Complete support for **Azure Functions** with:

- ✅ **HTTP Adapter** - Adapter for Azure Functions HTTP triggers
- ✅ **Azure Key Vault** - Integration with Azure Key Vault for secrets
- ✅ **Terraform** - Infrastructure as code for Azure
- ✅ **Independent Modules** - Each function has its own module
- ✅ **Automatic Generation** - Generation system from `functions.json`

**Azure functions structure:**
```
src/infrastructure/clouds/azure/
├── functions.json          # Function definitions
├── init.go                 # Azure initialization
├── http_adapter.go         # HTTP adapter
├── vault.go                # Key Vault integration
├── terraform/              # Terraform infrastructure
└── functions/              # Function generator
```

### 📊 Architecture Comparison

| Feature | Traditional Monolith | AWS Lambda | Azure Functions |
|---------|---------------------|------------|-----------------|
| **Initialization** | Once at startup | Per function | Per function |
| **Scalability** | Manual/Horizontal | Automatic | Automatic |
| **Cost** | Fixed | Pay per use | Pay per use |
| **Cold Start** | N/A | ~100-500ms | ~200-800ms |
| **Binary Size** | ~50 MB | ~5-15 MB | ~5-15 MB |
| **Secrets Management** | Environment variables | Secrets Manager | Key Vault |
| **Deployment** | Docker/VM | ZIP to Lambda | ZIP to Functions |

---

## Project Architecture

### Architecture Overview

The project implements **Clean Architecture** with three main layers:

```mermaid
graph TB
    subgraph Infrastructure["🔧 Infrastructure Layer"]
        API["API Layer<br/>(Gin HTTP)"]
        DB["Database<br/>(GORM/Redis)"]
        Providers["Providers<br/>(JWT/Email)"]
        Handlers["Handlers<br/>(HTTP)"]
        Repos["Repositories<br/>(GORM)"]
        Middlewares["Middlewares<br/>(Auth/CORS)"]
    end

    subgraph Application["💼 Application Layer"]
        UseCases["Use Cases<br/>(Business Logic)"]
        Services["Services<br/>(Email/OTP)"]
        Contracts["Contracts<br/>(Interfaces)"]
        Pipes["Pipes<br/>(DAG)"]
        DTOs["DTOs<br/>(Transfer)"]
        Errors["Errors<br/>(Handling)"]
    end

    subgraph Domain["🎯 Domain Layer"]
        Models["Models<br/>(User/Role)"]
        Utils["Utils<br/>(Query)"]
        Entities["Entities<br/>(Business)"]
    end

    Infrastructure --> Application
    Application --> Domain

    style Infrastructure fill:#e1f5ff
    style Application fill:#fff4e1
    style Domain fill:#e8f5e9
```

### Hexagonal Architecture Diagram

```mermaid
graph LR
    subgraph External["🌐 External World"]
        HTTP["HTTP Clients"]
        DB_EXT["PostgreSQL"]
        REDIS_EXT["Redis"]
        EMAIL_EXT["SMTP"]
    end

    subgraph Adapters["🔌 Adapters (Infrastructure)"]
        GinAdapter["Gin Adapter<br/>(HTTP)"]
        GORMAdapter["GORM Adapter<br/>(Database)"]
        RedisAdapter["Redis Adapter<br/>(Cache)"]
        SMTPAdapter["SMTP Adapter<br/>(Email)"]
    end

    subgraph Ports["🔌 Ports (Contracts)"]
        HTTPPort["HTTP Port<br/>(Handlers)"]
        RepoPort["Repository Port<br/>(Interfaces)"]
        CachePort["Cache Port<br/>(Interface)"]
        EmailPort["Email Port<br/>(Interface)"]
    end

    subgraph ApplicationCore["💼 Application Core"]
        UseCase["Use Cases"]
        Services["Services"]
    end

    subgraph DomainCore["🎯 Domain Core"]
        Entities["Entities"]
        ValueObjects["Value Objects"]
    end

    HTTP --> GinAdapter
    DB_EXT --> GORMAdapter
    REDIS_EXT --> RedisAdapter
    EMAIL_EXT --> SMTPAdapter

    GinAdapter --> HTTPPort
    GORMAdapter --> RepoPort
    RedisAdapter --> CachePort
    SMTPAdapter --> EmailPort

    HTTPPort --> UseCase
    RepoPort --> UseCase
    CachePort --> UseCase
    EmailPort --> Services

    UseCase --> Services
    Services --> Entities
    UseCase --> Entities

    style External fill:#ffebee
    style Adapters fill:#e3f2fd
    style Ports fill:#fff9c4
    style ApplicationCore fill:#e8f5e9
    style DomainCore fill:#f3e5f5
```

### Design Principles

#### 1. Dependency Inversion

Inner layers (Domain and Application) **never** depend on outer layers (Infrastructure). Instead, they define **interfaces (contracts)** that infrastructure implements.

**Example:**
```go
// Domain/Application defines the interface
type IUserRepository interface {
    Create(input UserCreate) (*User, error)
    GetByID(id uint) (*User, error)
}

// Infrastructure implements the interface
type UserRepository struct {
    DB *gorm.DB
}

func (r *UserRepository) Create(input UserCreate) (*User, error) {
    // Implementation with GORM
}
```

#### 2. Separation of Concerns

Each layer has a single, well-defined responsibility:

- **Domain**: Pure business entities, no dependencies
- **Application**: Business logic and use cases
- **Infrastructure**: Technical implementations (HTTP, DB, etc.)

#### 3. Testability

Thanks to interfaces, we can test business logic without needing real databases or services:

```go
// In tests, we use mocks
mockRepo := &MockUserRepository{}
useCase := NewCreateUserUseCase(logger, mockRepo)
```

#### 4. Extensibility

Adding new features is simple and doesn't affect existing code:

1. Create model in `domain/models/`
2. Define interface in `application/contracts/`
3. Implement use case in `application/modules/`
4. Create repository in `infrastructure/repositories/`
5. Add handler in `infrastructure/handlers/`

### Implemented Design Patterns

#### 1. Repository Pattern

Abstracts data access, allowing database changes without affecting business logic.

```go
type IUserRepository interface {
    IRepositoryBase[UserCreate, UserUpdate, User, UserInDB]
    CreateWithPassword(input UserAndPasswordCreate) (*User, error)
    GetUserWithRole(id uint) (*UserWithRole, error)
    GetByEmailOrPhone(emailOrPhone string) (*User, error)
}
```

#### 2. Use Case Pattern

Encapsulates business logic in reusable and testable use cases.

```go
type CreateUserUseCase struct {
    log  ILoggerProvider
    repo IUserRepository
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, locale LocaleTypeEnum, input UserCreate) *UseCaseResult[User] {
    // Business logic here
}
```

#### 3. DAG (Directed Acyclic Graph) Pattern

Orchestrates multiple use cases sequentially or in parallel.

```go
// Sequential execution
dag := NewDag(NewStep(createUserUC), locale, ctx)
dag = Then(dag, NewStep(sendEmailUC))
result := dag.Execute(input)

// Parallel execution
parallelDag := NewUseCaseParallelDag[Input, Output]()
parallelDag.Usecases = []BaseUseCase{uc1, uc2, uc3}
result := parallelDag.Execute(ctx, locale, input)
```

#### 4. Factory Pattern

Creates use case instances with injected dependencies.

```go
func NewCreateUserUseCase(
    log ILoggerProvider,
    repo IUserRepository,
) *CreateUserUseCase {
    return &CreateUserUseCase{
        log:  log,
        repo: repo,
    }
}
```

#### 5. Strategy Pattern

Allows swapping implementations through interfaces (providers).

```go
// Interface
type IHashProvider interface {
    Hash(password string) (string, error)
    Compare(hashed, plain string) bool
}

// Interchangeable implementations
type BcryptHashProvider struct {}
type Argon2HashProvider struct {}
```

---

## Scalability and Serverless

### Scalability Capabilities

**GoProjectSkeleton** is designed to scale both **horizontally** and **vertically**, and can easily migrate to **serverless** or **serverless monolith** architectures.

### Horizontal Scalability

#### Features that Facilitate Horizontal Scalability

1. **Stateless Design**
   - The application does not maintain state in memory
   - Each request is independent
   - Perfect for load balancers

2. **Distributed Cache (Redis)**
   - Cache shared between instances
   - No local session dependencies

3. **External Database**
   - PostgreSQL independent of the application
   - Multiple instances can connect

4. **No Session State**
   - JWT-based authentication (stateless)
   - No sticky sessions required

#### Implementation for Horizontal Scalability

```go
// Each handler is stateless
func CreateUser(ctx HandlerContext) {
    // No shared state
    // Each request is independent
    // Can run on any instance
}
```

### Vertical Scalability

#### Implemented Optimizations

1. **Native Goroutines**
   - Efficient Go concurrency
   - Multiple requests processed simultaneously

2. **Connection Pooling**
   - Database connection reuse
   - Optimized GORM configuration

3. **Smart Cache**
   - Reduces database queries
   - Configurable TTL

4. **Query Optimization**
   - Efficient filters and sorting
   - Pagination for memory control

### Migration to Serverless Monolith

**GoProjectSkeleton** can easily migrate to a **serverless monolith** architecture (like AWS Lambda, Google Cloud Functions, Azure Functions) thanks to:

#### 1. Decoupled Architecture

```go
// Business logic does not depend on Gin
type CreateUserUseCase struct {
    // No HTTP references
    // Only interfaces
}
```

#### 2. Modular Initialization

```go
// container.go - Separate initialization
func Initialize() {
    // Configuration
    // Database
    // Providers
    // Services
}
```

#### 3. Independent Handlers

Handlers can easily adapt to different frameworks or serverless environments:

```go
// Current handler (Gin)
func CreateUser(ctx HandlerContext) { ... }

// Lambda adaptation
func CreateUserLambda(ctx context.Context, event APIGatewayEvent) (Response, error) {
    // Same logic, different wrapper
}
```

#### 4. No Global State

- No global state variables
- Everything is passed through context or dependency injection

#### Steps for Serverless Migration

1. **Extract Business Logic**
   ```go
   // Already done - Use Cases are independent
   ```

2. **Create Serverless Adapter**
   ```go
   // lambda/handlers/user.go
   func CreateUserHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
       // Initialize infrastructure
       infrastructure.Initialize()

       // Adapt event to DTO
       var userCreate dtos.UserCreate
       json.Unmarshal([]byte(event.Body), &userCreate)

       // Execute use case (same code)
       ucResult := usecases_user.NewCreateUserUseCase(...).Execute(...)

       // Adapt response
       return adaptResponse(ucResult), nil
   }
   ```

3. **Configure Environment Variables**
   - Use the same configuration system
   - Adapt to cloud provider environment variables

4. **Maintain Same Logic**
   - Use Cases don't change
   - Repositories don't change
   - Only the HTTP infrastructure layer changes

### Production Scalability

#### Horizontal Scalability Diagram

```mermaid
graph TB
    subgraph Clients["Clients"]
        Web[Web App]
        Mobile[Mobile App]
        API_Client[API Clients]
    end

    subgraph LB["Load Balancer"]
        Nginx[Nginx/HAProxy<br/>Health Checks]
    end

    subgraph AppInstances["Application Instances"]
        App1[App Instance 1<br/>Stateless]
        App2[App Instance 2<br/>Stateless]
        App3[App Instance N<br/>Stateless]
    end

    subgraph DataLayer["Data Layer"]
        RedisCluster[(Redis Cluster<br/>Distributed Cache)]
        PGPrimary[(PostgreSQL<br/>Primary)]
        PGReplica1[(PostgreSQL<br/>Read Replica 1)]
        PGReplica2[(PostgreSQL<br/>Read Replica 2)]
    end

    Web --> Nginx
    Mobile --> Nginx
    API_Client --> Nginx

    Nginx --> App1
    Nginx --> App2
    Nginx --> App3

    App1 --> RedisCluster
    App2 --> RedisCluster
    App3 --> RedisCluster

    App1 -->|Write| PGPrimary
    App2 -->|Write| PGPrimary
    App3 -->|Write| PGPrimary

    App1 -->|Read| PGReplica1
    App2 -->|Read| PGReplica2
    App3 -->|Read| PGReplica1

    PGPrimary -->|Replication| PGReplica1
    PGPrimary -->|Replication| PGReplica2

    style App1 fill:#e3f2fd
    style App2 fill:#e3f2fd
    style App3 fill:#e3f2fd
    style RedisCluster fill:#ffcdd2
    style PGPrimary fill:#c8e6c9
    style PGReplica1 fill:#fff9c4
    style PGReplica2 fill:#fff9c4
```

#### Recommended Configuration

1. **Load Balancer**
   - Distribute traffic among multiple instances
   - Health checks at `/api/health-check`

2. **Database**
   - Read replicas for queries
   - Optimized connection pooling
   - Indexes on frequent columns

3. **Cache**
   - Redis cluster for high availability
   - Cache strategies (cache-aside, write-through)

4. **Monitoring**
   - Performance metrics
   - Structured logging
   - Error alerts

---

## AWS Serverless Deployment and Initialization

### Go Modules Architecture for Serverless

**GoProjectSkeleton** uses a **multiple Go modules** architecture to optimize deployment on AWS Lambda. Each serverless function has its own independent module that only includes the necessary dependencies to run.

#### Module Structure

```mermaid
graph TB
    subgraph RootModule["Root Module: goprojectskeleton"]
        Domain[Domain Models]
        Application[Application Layer]
        Infrastructure[Infrastructure Core]
    end

    subgraph AWSModule["AWS Module: goprojectskeleton/aws"]
        AWSInit[AWS Init]
        LambdaAdapter[Lambda Adapter]
        SecretsManager[Secrets Manager]
    end

    subgraph FunctionModules["Function Modules"]
        HealthCheck[health-check<br/>goprojectskeleton/functions/aws/status/health_check]
        AuthLogin[auth-login<br/>goprojectskeleton/functions/aws/auth/login]
        UserGet[user-get<br/>goprojectskeleton/functions/aws/user/get]
    end

    RootModule --> AWSModule
    AWSModule --> FunctionModules
    RootModule --> FunctionModules

    style RootModule fill:#e8f5e9
    style AWSModule fill:#e3f2fd
    style FunctionModules fill:#fff9c4
```

#### Dependency Resolution

Each Lambda function has its own `go.mod` that uses **replace directives** to point to local modules:

```go
// go.mod de una función Lambda
module goprojectskeleton/functions/aws/status/health_check

go 1.25

require (
    github.com/aws/aws-lambda-go v1.47.0
    goprojectskeleton v0.0.0
    goprojectskeleton/aws v0.0.0
)

replace goprojectskeleton => ../../../../../../../..
replace goprojectskeleton/aws => ../../../..
```

**Advantages of this architecture:**

1. **Binary Optimization**: Go only compiles code that is actually used
2. **Minimal Dependencies**: Each function only includes what's necessary
3. **Automatic Tree Shaking**: Go removes unused code at compile time
4. **Independent Modules**: Each function can evolve independently

### Generation and Compilation Process

#### 1. Function Generation

Lambda functions are automatically generated from `functions.json`:

```json
{
  "name": "health-check",
  "path": "status/health_check",
  "handler": "GetHealthCheck",
  "route": "health-check",
  "method": "get",
  "authLevel": "anonymous"
}
```

**Generation process:**

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Gen as Generator
    participant Template as Templates
    participant Func as Lambda Function

    Dev->>Gen: GenerateFunctions(functions.json)
    Gen->>Gen: Reads functions.json
    loop For each function
        Gen->>Template: Loads main.go.tmpl
        Gen->>Template: Loads go.mod.tmpl
        Gen->>Func: Generates main.go
        Gen->>Func: Generates go.mod
        Gen->>Func: Creates directory structure
    end
    Gen-->>Dev: ✅ Functions generated
```

**Generated structure:**

```
tmp/
└── status/
    └── health_check/
        ├── main.go          # Generated Lambda handler
        ├── go.mod           # Independent module
        └── bin/             # Compilation directory
            ├── bootstrap    # Compiled binary
            └── src/         # Copied templates
```

#### 2. Optimized Compilation

Each function is compiled with Lambda-specific optimizations:

```bash
# Compilation from function directory
cd tmp/status/health_check
go build -o bin/bootstrap \
    -tags lambda.norpc \
    -ldflags="-s -w" \
    main.go

# Compilation environment variables
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0
```

**Applied optimizations:**

- **`-tags lambda.norpc`**: Disables Lambda RPC (reduces size)
- **`-ldflags="-s -w"`**: Removes debug symbols (reduces size)
- **`CGO_ENABLED=0`**: Static compilation (no C dependencies)
- **`GOOS=linux`**: Linux binary (Lambda environment)
- **`GOARCH=amd64`**: x86_64 architecture

#### 3. Tree Shaking and Dead Code Elimination

Go performs **automatic tree shaking** during compilation:

```mermaid
graph LR
    subgraph SourceCode["Source Code"]
        Used[✅ Used Code<br/>handlers.GetHealthCheck<br/>aws.InitializeInfrastructure]
        Unused[❌ Unused Code<br/>handlers.CreateUser<br/>handlers.UpdateUser<br/>handlers.DeleteUser]
    end

    subgraph Compilation["Go Compilation"]
        Analyzer[Go Compiler<br/>Static Analysis]
    end

    subgraph Binary["Final Binary"]
        Included[✅ Only Used Code<br/>~5-10 MB]
    end

    SourceCode --> Analyzer
    Analyzer -->|Tree Shaking| Binary

    style Used fill:#c8e6c9
    style Unused fill:#ffcdd2
    style Included fill:#c8e6c9
```

**Real example:**

For the `health-check` function, the final binary only includes:

- ✅ `handlers.GetHealthCheck` (specific handler)
- ✅ `aws.InitializeInfrastructure` (initialization)
- ✅ `aws.HandleLambdaEvent` (Lambda adapter)
- ✅ Necessary transitive dependencies (GORM, Redis, JWT, etc.)
- ❌ **Does NOT include**: Other handlers (`CreateUser`, `UpdateUser`, etc.)
- ❌ **Does NOT include**: Code from unused modules

**Result:** Binaries of 5-15 MB instead of 50+ MB if everything were included.

### Infrastructure Initialization

#### Initialization Process

Each Lambda function initializes its infrastructure in `init()`:

```go
// main.go de una función Lambda
func init() {
    if !initialized {
        log.Println("Initializing AWS Infrastructure")
        aws.InitializeInfrastructure()
        initialized = true
        log.Println("AWS Infrastructure initialized successfully")
    }
}
```

#### Initialization Flow

```mermaid
sequenceDiagram
    participant Lambda as Lambda Function
    participant Init as InitializeInfrastructure
    participant Config as Config Loader
    participant Secrets as Secrets Manager
    participant Settings as App Settings
    participant DB as Database
    participant Providers as Providers

    Lambda->>Init: aws.InitializeInfrastructure()
    Init->>Config: NewSecretsManagerConfigLoader()
    Config->>Secrets: Load secrets from AWS
    Secrets-->>Config: Config values
    Config->>Settings: Initialize settings
    Settings->>DB: Setup database connection
    Settings->>Providers: Setup JWT, Email, Cache
    Providers-->>Init: ✅ Initialized
    Init-->>Lambda: ✅ Ready
```

#### Initialized Components

1. **Configuration (Settings)**
   ```go
   // Loads from AWS Secrets Manager or environment variables
   settings.AppSettingsInstance.Initialize(config.ToMap())
   ```

2. **Database (GORM)**
   ```go
   database.GoProjectSkeletondb.SetUp(
       host, port, user, password, dbname, ssl, logger
   )
   ```

3. **Providers**
   - **JWT Provider**: Token generation and validation
   - **Email Provider**: Email sending (SMTP)
   - **Cache Provider**: Redis for cache
   - **Logger Provider**: Logging system

4. **Services**
   - **Email Services**: Email services (registration, reset, OTP)

#### Configuration Loading from AWS Secrets Manager

The system can load configuration from **AWS Secrets Manager**:

```go
// If an environment variable is a Secrets Manager ARN
DB_PASSWORD=arn:aws:secretsmanager:us-east-1:123456789:secret:db-password

// The system automatically:
// 1. Detects it's an ARN
// 2. Gets the secret from Secrets Manager
// 3. Uses the value in configuration
```

**Advantages:**

- ✅ **Security**: Secrets not in code or environment variables
- ✅ **Rotation**: Secrets Manager can rotate secrets automatically
- ✅ **Audit**: All secret accesses are audited
- ✅ **Fallback**: If it fails, uses default values

### Deployment Process

#### Complete Deployment Flow

```mermaid
graph TB
    subgraph Local["🖥️ Local Development"]
        Dev[Developer]
        FunctionsJSON[functions.json]
    end

    subgraph Generation["📦 Generation"]
        Generator[Generator Tool]
        Templates[Templates]
        FunctionCode[Function Code]
    end

    subgraph Build["🔨 Build"]
        GoBuild[go build]
        Binary[bootstrap Binary]
        TemplatesCopy[Copied Templates]
        Zip[ZIP Package]
    end

    subgraph AWS["☁️ AWS"]
        Lambda[AWS Lambda]
        APIGateway[API Gateway]
    end

    Dev-->FunctionsJSON
    FunctionsJSON-->|GenerateFunctions| Generator
    Generator-->|Loads templates| Templates
    Templates-->|Generates code| FunctionCode
    FunctionCode-->|go build| GoBuild
    GoBuild-->|Optimized binary| Binary
    Binary-->|Copy templates| TemplatesCopy
    TemplatesCopy-->|Create ZIP| Zip
    Zip-->|Deploy function| Lambda
    Lambda-->|Connect endpoint| APIGateway

    style Dev fill:#e3f2fd
    style Lambda fill:#ff9800
    style APIGateway fill:#ff9800
```

#### Deployment Commands

**1. Generate functions:**
```bash
cd src/infrastructure/clouds/aws/functions
go run main.go generate
```

**2. Build and deploy:**
```bash
# Deploy all functions
go run main.go deploy

# Deploy a specific function
go run main.go deploy --function health-check
```

**3. Internal deployment process:**

```bash
# For each function:
# 1. go mod tidy (resolves dependencies)
cd tmp/status/health_check
go mod tidy

# 2. Compile binary
go build -o bin/bootstrap -tags lambda.norpc main.go

# 3. Copy necessary templates
cp -r ../../../../../../application/shared/templates bin/src/application/shared/templates

# 4. Create ZIP
cd bin
zip -r health-check.zip bootstrap src/

# 5. Deploy to Lambda
aws lambda update-function-code \
    --function-name goprojectskeleton-dev-healthcheck \
    --zip-file fileb://health-check.zip
```

### Size Optimizations

#### Size Comparison

| Function | Size with Everything | Optimized Size | Reduction |
|---------|----------------|-------------------|----------|
| health-check | ~50 MB | ~8 MB | 84% |
| auth-login | ~55 MB | ~12 MB | 78% |
| user-get | ~60 MB | ~15 MB | 75% |

#### Optimization Strategies

1. **Go Tree Shaking**
   - Automatically removes unused code
   - Only includes referenced functions and types

2. **Independent Modules**
   - Each function has its own `go.mod`
   - Dependencies resolved per function

3. **Static Compilation**
   - `CGO_ENABLED=0`: No C dependencies
   - Self-contained binary

4. **Debug Removal**
   - `-ldflags="-s -w"`: Removes symbols
   - Significantly reduces size

5. **Selective Templates**
   - Only necessary templates are copied
   - Not all templates are included

### Lambda Package Structure

```
health-check.zip
├── bootstrap                    # Compiled Go binary (~8 MB)
└── src/
    └── application/
        └── shared/
            └── templates/
                └── emails/      # Only necessary templates
                    ├── otp_en.gohtml
                    └── otp_es.gohtml
```

**Total size:** ~8-10 MB (vs ~50 MB without optimization)

### Architecture Advantages

#### 1. **Fast Deployments**
- Small binaries = fast uploads
- Faster cold start
- Lower storage cost

#### 2. **Improved Security**
- Each function is independent
- Smaller attack surface
- Secrets in Secrets Manager

#### 3. **Scalability**
- Each function scales independently
- Per-function configuration
- Individual optimization

#### 4. **Maintainability**
- Shared code in modules
- Localized changes
- Independent testing

#### 5. **Optimized Cost**
- Smaller size = lower storage cost
- Faster cold start = lower latency
- Lower memory usage

### Complete Architecture Diagram

```mermaid
graph TB
    subgraph Client["🌐 Client"]
        HTTP[HTTP Request]
    end

    subgraph AWSCloud["☁️ AWS Cloud"]
        subgraph APIGW["API Gateway"]
            Route[Route: /api/health-check]
        end

        subgraph Lambda["AWS Lambda"]
            Function[health-check Function<br/>~8 MB]
            Init[InitializeInfrastructure]
            Handler[GetHealthCheck Handler]
        end

        subgraph Services["AWS Services"]
            Secrets[Secrets Manager<br/>Configuration]
            RDS[(RDS PostgreSQL<br/>Database)]
            ElastiCache[(ElastiCache Redis<br/>Cache)]
        end
    end

    HTTP-->Route
    Route-->Function
    Function-->Init
    Init-->Secrets
    Init-->RDS
    Init-->ElastiCache
    Init-->Handler
    Handler-->RDS
    Handler-->ElastiCache

    style Function fill:#ff9800
    style Secrets fill:#4caf50
    style RDS fill:#2196f3
    style ElastiCache fill:#f44336
```

---

## Complete Request Flow

### Request Flow Diagram

```mermaid
sequenceDiagram
    participant Client as HTTP Client
    participant Gin as Gin Router
    participant MW as Middleware<br/>(CORS/Recovery)
    participant AuthMW as Auth Middleware<br/>(if protected)
    participant Router as Router/Wrapper
    participant Handler as HTTP Handler
    participant UseCase as Use Case
    participant Repo as Repository
    participant DB as PostgreSQL
    participant Resolver as Request Resolver

    Client->>Gin: HTTP Request
    Gin->>MW: Applies middlewares
    MW->>AuthMW: If protected route
    AuthMW->>AuthMW: Validates JWT Token
    AuthMW->>AuthMW: Gets User
    AuthMW->>Router: Request with context
    Router->>Router: Extracts params, query, locale
    Router->>Router: Creates HandlerContext
    Router->>Handler: HandlerContext
    Handler->>Handler: Decodes JSON → DTO
    Handler->>UseCase: Executes Use Case
    UseCase->>UseCase: Validates input (DTO.Validate())
    UseCase->>UseCase: Executes Guards (if applicable)
    UseCase->>Repo: Calls Repository
    Repo->>Repo: Converts DTO → DB Model
    Repo->>DB: Executes Query (GORM)
    DB-->>Repo: Returns data
    Repo->>Repo: Converts DB Model → Domain
    Repo->>Repo: Maps errors
    Repo-->>UseCase: Domain Model / Error
    UseCase->>UseCase: Handles errors
    UseCase-->>Handler: UseCaseResult
    Handler->>Resolver: Resolves response
    Resolver->>Resolver: Maps status code
    Resolver->>Resolver: Formats JSON
    Resolver-->>Handler: HTTP Response
    Handler-->>Client: JSON Response
```

### Simplified Flow Diagram

```mermaid
flowchart TD
    Start([HTTP Client]) --> Gin[Gin Router]
    Gin --> MW[Middlewares<br/>CORS/Recovery]
    MW --> Auth{Protected<br/>Route?}
    Auth -->|Yes| AuthMW[Auth Middleware<br/>Validates JWT]
    Auth -->|No| Router
    AuthMW --> Router[Router/Wrapper<br/>Extracts params/query/locale]
    Router --> Handler[HTTP Handler<br/>Decodes JSON → DTO]
    Handler --> UseCase[Use Case<br/>Validates and executes logic]
    UseCase --> Repo[Repository<br/>Converts DTO ↔ Model]
    Repo --> DB[(PostgreSQL<br/>GORM)]
    DB --> Repo
    Repo --> UseCase
    UseCase --> Resolver[Request Resolver<br/>Formats response]
    Resolver --> End([JSON Response])

    style Start fill:#e3f2fd
    style End fill:#c8e6c9
    style DB fill:#fff9c4
    style UseCase fill:#f3e5f5
```

### Detailed Flow with Code

#### 1. HTTP Entry (main.go)

```go
// main.go
func main() {
    infrastructure.Initialize()  // Initializes DB, providers, etc.
    app := buildGinApp()         // Creates Gin application
    loadGinApp(app)               // Configures middlewares and routes
    app.Run("0.0.0.0:8080")      // Starts server
}
```

#### 2. Routing (routes/router.go)

```go
// routes/router.go
func Router(r *gin.RouterGroup) {
    r.POST("/user", wrapHandler(handlers.CreateUser))
    // wrapHandler adapts Gin to HandlerContext
}
```

#### 3. Wrapper (routes/wrap.go)

```go
// routes/wrap.go
func wrapHandler(h func(handlers.HandlerContext)) gin.HandlerFunc {
    return func(c *gin.Context) {
        locale := c.GetHeader("Accept-Language")
        params := extractParams(c)
        query := extractQuery(c)

        hContext := handlers.NewHandlerContext(
            c.Request.Context(),
            &locale,
            params,
            &c.Request.Body,
            query,
            c.Writer,
        )
        h(hContext)
    }
}
```

#### 4. Handler (handlers/user.go)

```go
// handlers/user.go
func CreateUser(ctx HandlerContext) {
    // 1. Decode JSON
    var userCreate dtos.UserCreate
    json.NewDecoder(*ctx.Body).Decode(&userCreate)

    // 2. Create repository
    repo := repositories.NewUserRepository(database.DB, providers.Logger)

    // 3. Create and execute use case
    ucResult := usecases_user.NewCreateUserUseCase(
        providers.Logger,
        repo,
    ).Execute(ctx.c, ctx.Locale, userCreate)

    // 4. Resolve response
    headers := map[HTTPHeaderTypeEnum]string{
        CONTENT_TYPE: string(APPLICATION_JSON),
    }
    NewRequestResolver[models.User]().ResolveDTO(
        ctx.ResponseWriter,
        ucResult,
        headers,
    )
}
```

#### 5. Use Case (application/modules/user/use_cases/create_user.go)

```go
// application/modules/user/use_cases/create_user.go
func (uc *CreateUserUseCase) Execute(
    ctx context.Context,
    locale locales.LocaleTypeEnum,
    input dtos.UserCreate,
) *usecase.UseCaseResult[models.User] {
    result := usecase.NewUseCaseResult[models.User]()

    // 1. Validate input
    uc.validate(input, result)
    if result.HasError() {
        return result
    }

    // 2. Call repository
    res, err := uc.repo.Create(input)
    if err != nil {
        result.SetError(err.Code, err.Context)
        return result
    }

    // 3. Return success
    result.SetData(status.Created, *res, "User created")
    return result
}
```

#### 6. Repository (infrastructure/repositories/user.go)

```go
// infrastructure/repositories/user.go
func (ur *UserRepository) Create(input dtos.UserCreate) (*models.User, *application_errors.ApplicationError) {
    // 1. Convert DTO to GORM model
    userCreate := ur.modelConverter.ToGormCreate(input)

    // 2. Execute query
    if err := ur.DB.Create(userCreate).Error; err != nil {
        return nil, MapOrmError(err)  // Maps DB errors
    }

    // 3. Convert GORM model to domain model
    userModel := ur.modelConverter.ToDomain(userCreate)
    return userModel, nil
}
```

#### 7. Response Resolution (handlers/request_resolver.go)

```go
// handlers/request_resolver.go
func (rr *RequestResolver[D]) ResolveDTO(
    w http.ResponseWriter,
    result *usecase.UseCaseResult[D],
    headersToAdd map[HTTPHeaderTypeEnum]string,
) {
    // 1. Add headers
    rr.setHeaders(w, headersToAdd)

    // 2. Handle errors
    if result.HasError() {
        w.WriteHeader(rr.statusMapping[result.StatusCode])
        json.NewEncoder(w).Encode(map[string]any{
            "details": result.Error,
        })
        return
    }

    // 3. Success response
    w.WriteHeader(rr.statusMapping[result.StatusCode])
    json.NewEncoder(w).Encode(map[string]any{
        "data":    result.Data,
        "details": result.Details,
    })
}
```

### Flow with Pipes (DAG)

For more complex cases that require multiple steps:

```mermaid
sequenceDiagram
    participant Handler as Handler
    participant DAG as DAG Pipe
    participant UC1 as CreateUserAndPassword<br/>Use Case
    participant UC2 as CreateUserSendEmail<br/>Use Case
    participant Repo as User Repository
    participant EmailSvc as Email Service
    participant DB as PostgreSQL
    participant SMTP as SMTP Server

    Handler->>DAG: Execute(userCreate)
    DAG->>UC1: Execute(userCreate)
    UC1->>Repo: CreateWithPassword()
    Repo->>DB: INSERT User + Password
    DB-->>Repo: User created
    Repo-->>UC1: User
    UC1-->>DAG: UseCaseResult[User]

    alt If no error
        DAG->>UC2: Execute(User)
        UC2->>EmailSvc: SendWelcomeEmail()
        EmailSvc->>SMTP: Send email
        SMTP-->>EmailSvc: Email sent
        EmailSvc-->>UC2: Success
        UC2-->>DAG: UseCaseResult[User]
        DAG-->>Handler: UseCaseResult[User]
    else If error
        DAG-->>Handler: UseCaseResult[Error]
    end
```

#### DAG Diagram (Directed Acyclic Graph)

```mermaid
graph LR
    Start([Input:<br/>UserCreate]) --> UC1[Use Case 1:<br/>CreateUserAndPassword]
    UC1 -->|Output: User| UC2[Use Case 2:<br/>CreateUserSendEmail]
    UC2 -->|Output: User| End([Result:<br/>User])

    UC1 -.->|Error| Error[Error Handler]
    UC2 -.->|Error| Error
    Error --> End

    style Start fill:#e3f2fd
    style End fill:#c8e6c9
    style Error fill:#ffcdd2
```

#### Parallel Execution with DAG

```mermaid
graph TB
    Start([Input]) --> Parallel[Parallel DAG]

    Parallel --> UC1[Use Case 1]
    Parallel --> UC2[Use Case 2]
    Parallel --> UC3[Use Case 3]

    UC1 --> Wait[WaitGroup<br/>Waits for all]
    UC2 --> Wait
    UC3 --> Wait

    Wait --> Merge[Merge Results]
    Merge --> End([Results Array])

    style Parallel fill:#fff9c4
    style Wait fill:#e1f5ff
    style Merge fill:#e8f5e9
```

**Example code:**

```go
// Example: Create user and send email
func CreateUserAndPassword(ctx HandlerContext) {
    // 1. Create use cases
    uc_create_user_password := usecases_user.NewCreateUserAndPasswordUseCase(...)
    uc_create_user_email := usecases_user.NewCreateUserSendEmailUseCase(...)

    // 2. Create pipe (DAG)
    pipe := user_pipes.NewCreateUserPipe(
        ctx.c,
        ctx.Locale,
        uc_create_user_password,
        uc_create_user_email,
    )

    // 3. Execute pipe (executes sequentially)
    ucResult := pipe.Execute(userCreate)

    // 4. Resolve response
    NewRequestResolver[models.User]().ResolveDTO(...)
}
```

The DAG executes:
1. `CreateUserAndPasswordUseCase` → returns `User`
2. `CreateUserSendEmailUseCase` → receives `User`, sends email, returns `User`

---

## Virtues and Benefits

### 1. Solid and Scalable Architecture

#### ✅ Clean Architecture
- **Clear separation of responsibilities**: Each layer has a specific purpose
- **Framework independence**: You can change Gin for another framework without affecting business logic
- **Testability**: Easy to test each layer independently

#### ✅ Hexagonal Architecture
- **Total decoupling**: Business logic doesn't know implementation details
- **Ports & Adapters**: Clear interfaces between layers
- **Flexibility**: Change database, providers, etc. without rewriting code

### 2. Developer Productivity

#### ✅ Clear Structure
- **Logical organization**: Easy to find code
- **Consistent conventions**: Same pattern throughout the project
- **Fast onboarding**: New developers understand quickly

#### ✅ Reusability
- **Shared components**: Reusable DTOs, errors, validations
- **Repository base**: `RepositoryBase` reduces code duplication
- **Interchangeable providers**: Easily change implementations

### 3. Maintainability

#### ✅ Clean Code
- **SOLID principles**: Applied consistently
- **DRY (Don't Repeat Yourself)**: Minimal duplication
- **Descriptive names**: Self-documenting code

#### ✅ Extensibility
- **Add functionality**: Without modifying existing code
- **Independent modules**: Each module is self-contained
- **Well-defined interfaces**: Clear contracts

### 4. Testability

#### ✅ Layered Testing
- **Unit tests**: Use cases testable with mocks
- **Integration tests**: Tests with real database
- **E2E tests**: Complete flow tests

#### ✅ Complete Mocks
- **Repository mocks**: Easy to simulate data
- **Provider mocks**: Simulate external services
- **Isolated testing**: No external dependencies

### 5. Security

#### ✅ Robust Authentication
- **Complete JWT**: Access and refresh tokens
- **OTP (2FA)**: Two-factor authentication
- **Secure hash**: Bcrypt for passwords

#### ✅ Validation
- **Multi-layer validation**: DTOs, use cases, repositories
- **Sanitization**: Injection prevention
- **Guards**: Role-based access control

### 6. Performance

#### ✅ Optimizations
- **Redis cache**: Reduces database queries
- **Connection pooling**: Connection reuse
- **Goroutines**: Native Go concurrency
- **Query optimization**: Efficient filters and pagination

### 7. Internationalization

#### ✅ Multi-language
- **i18n support**: Messages in multiple languages
- **Locale per request**: Each request can have its language
- **Centralized messages**: Easy to add new languages

### 8. Documentation

#### ✅ Automatic Swagger
- **Documented API**: Endpoints automatically documented
- **Examples**: Request and response examples
- **Clear types**: Well-defined schemas

### 9. DevOps and Deployment

#### ✅ Complete Docker
- **Multi-service**: Application, DB, Redis, etc.
- **Separate environments**: Dev, test, E2E
- **Hot reload**: Efficient development

#### ✅ Flexible Configuration
- **Environment variables**: Configuration per environment
- **Centralized settings**: Single place for configuration
- **Validation**: Configuration validation at startup

### 10. Scalability

#### ✅ Horizontal
- **Stateless**: Ready for load balancers
- **Distributed cache**: Shared Redis
- **No sessions**: Stateless JWT

#### ✅ Vertical
- **Goroutines**: Efficient concurrency
- **Optimizations**: Cache, pooling, etc.
- **Serverless ready**: Easy migration to serverless

### Benefits for Starting a Project

1. **Time Savings**
   - Ready-to-use structure
   - Implemented patterns
   - Don't start from scratch

2. **Best Practices**
   - Proven architecture
   - Applied design patterns
   - Quality code

3. **Guaranteed Scalability**
   - Designed to grow
   - Easy to add functionality
   - Production ready

4. **Long-term Maintainability**
   - Organized code
   - Easy to understand
   - Easy to modify

5. **Productive Team**
   - Fast onboarding
   - Clear conventions
   - Fewer bugs

---

## 📊 Project Statistics

| Metric | Value |
|---------|-------|
| **Go Files** | ~180+ files |
| **Lines of Code** | ~15,000+ lines |
| **Use Cases** | 20+ use cases |
| **Business Modules** | 4 modules (auth, user, password, status) |
| **Providers** | 7 providers (JWT, Hash, Email, Cache, Logger, Renderer, Status) |
| **Repositories** | 6 repositories |
| **HTTP Handlers** | 15+ endpoints |
| **Tests** | 20+ test files |
| **Templates** | 6+ HTML templates |
| **Supported Languages** | 2 (Spanish, English) |

## Project Structure - Layer by Layer

### Structure Overview

```
GoProjectSkeleton/
├── src/
│   ├── domain/              # 🎯 Domain Layer (Core)
│   │   ├── models/          # Pure business entities
│   │   └── utils/           # Domain utilities
│   ├── application/         # 💼 Application Layer (Business Logic)
│   │   ├── contracts/      # Interfaces (Ports)
│   │   ├── modules/         # Business modules
│   │   └── shared/         # Shared components
│   └── infrastructure/     # 🔧 Infrastructure Layer (Technical Details)
│       ├── server/          # HTTP Server (Gin)
│       ├── database/        # Database (GORM)
│       ├── providers/       # Provider implementations
│       ├── repositories/    # Repository implementations
│       ├── handlers/        # HTTP Handlers
│       ├── config/          # Configuration
│       └── clouds/          # Cloud Adapters (AWS, Azure)
├── docker/                  # 🐳 Docker Configuration
│   ├── docker-compose.dev.yml
│   ├── docker-compose.test.yml
│   ├── docker-compose.e2e.yml
│   └── db/                  # Database configuration
├── tests/                   # 🧪 Project Tests
│   ├── integration/         # Integration tests
│   └── e2e/                 # End-to-end tests (Bruno)
├── src/infrastructure/docs/ # 📚 Independent Swagger Service
│   ├── main.go              # Independent HTTP server for Swagger
│   ├── config/              # Swagger server configuration
│   ├── swagger/             # Generated Swagger files
│   │   ├── swagger.json
│   │   ├── swagger.yaml
│   │   └── docs.go
│   └── go.mod               # Independent module for Swagger
└── IDE/                     # ⚙️ IDE Configuration
    ├── launch.json
    └── tasks.json
```

### Layer 1: Domain

**Responsibility**: Contains pure business entities, without external dependencies.

```
domain/
├── models/                  # Domain models
│   ├── user.go             # User entity
│   ├── role.go             # Role entity
│   ├── password.go         # Password entity
│   ├── one_time_password.go
│   ├── one_time_token.go
│   └── status.go
└── utils/                   # Domain utilities
    └── query_payload.go    # Query utilities
```

#### Characteristics

- **No external dependencies**: Doesn't import frameworks, DB, etc.
- **Pure entities**: Only business logic
- **Domain validations**: Business rules in models
- **Immutability preferred**: Immutable models when possible

#### Example: User Model

```go
// domain/models/user.go
type UserBase struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
    Status   string `json:"status"`
    RoleID   uint   `json:"role_id"`
    OTPLogin bool   `json:"otpLogin"`
}

func (u UserBase) Validate() []string {
    var errs []string
    if u.Name == "" {
        errs = append(errs, "name is required")
    }
    // Más validaciones...
    return errs
}
```

### Layer 2: Application

**Responsibility**: Contains business logic and use cases.

```
application/
├── contracts/               # Contracts (Interfaces)
│   ├── providers/          # Provider interfaces
│   │   ├── cache_provider.go
│   │   ├── email_provider.go
│   │   ├── hash_provider.go
│   │   ├── jwt_provider.go
│   │   └── logger_provider.go
│   └── repositories/       # Repository interfaces
│       ├── user.go
│       ├── password.go
│       └── role.go
├── modules/                # Business modules
│   ├── auth/              # Authentication module
│   ├── user/              # User module
│   ├── password/          # Password module
│   └── status/            # Status module
└── shared/                # Shared components
    ├── DTOs/              # Data Transfer Objects
    ├── errors/            # Error handling
    ├── services/          # Shared services
    ├── templates/         # Templates (emails)
    ├── use_case/          # Use case base
    ├── context/           # Application context
    ├── locales/           # Internationalization
    ├── settings/          # Application configuration
    └── guards/            # Authorization guards
```

#### Characteristics

- **Use cases**: Each functionality is a use case
- **Interfaces**: Defines contracts that infrastructure implements
- **DTOs**: Data transfer objects
- **Services**: Logic shared between modules

#### Example: Use Case

```go
// application/modules/user/use_cases/create_user.go
type CreateUserUseCase struct {
    log  ILoggerProvider
    repo IUserRepository
}

func (uc *CreateUserUseCase) Execute(
    ctx context.Context,
    locale LocaleTypeEnum,
    input UserCreate,
) *UseCaseResult[User] {
    // 1. Validate
    // 2. Execute business logic
    // 3. Call repository
    // 4. Return result
}
```

### Layer 3: Infrastructure

**Responsibility**: Technical implementations (HTTP, DB, etc.).

```
infrastructure/
├── api/                    # HTTP API Layer
│   ├── cmd/
│   │   └── main.go        # Entry point
│   ├── routes/            # Route definitions
│   ├── middlewares/       # HTTP Middlewares
│   └── types.go
├── config/                # Configuration
│   ├── config.go         # Configuration loading
│   └── env.go            # Environment variables
├── database/              # Database
│   └── goprojectskeleton/   # GORM implementation
│       ├── models/       # DB models
│       └── init_db/      # Initialization
├── handlers/              # HTTP Handlers
│   ├── user.go
│   ├── auth.go
│   ├── password.go
│   └── status.go
├── providers/            # Provider implementations
│   ├── jwt_provider.go
│   ├── hash_provider.go
│   ├── email_provider.go
│   └── cache_provider.go
└── repositories/          # Repository implementations
    ├── user.go
    ├── password.go
    └── role.go
```

#### Characteristics

- **Concrete implementations**: GORM, Gin, Redis, etc.
- **Adapters**: Adapt frameworks to interfaces
- **Configuration**: Configuration loading from environment variables

#### Example: Repository

```go
// infrastructure/repositories/user.go
type UserRepository struct {
    RepositoryBase[UserCreate, UserUpdate, User, dbModels.User]
    DB *gorm.DB
}

func (r *UserRepository) Create(input UserCreate) (*User, error) {
    // Implementation with GORM
}
```

---

## Exhaustive Review by Folders

### `/src/domain/` - Domain Layer

#### `/src/domain/models/`

Contains all pure domain entities.

**Main files:**

- **`user.go`**: User entity with domain validations
  - `UserBase`: Base user structure
  - `User`: User with DB metadata
  - `UserWithRole`: User with role information
  - `Validate()`: Business validations

- **`role.go`**: Role entity
  - Defines system roles
  - Priorities and permissions

- **`password.go`**: Password entity
  - Password hashing
  - Expiration
  - Strength validation

- **`one_time_password.go`**: OTP for authentication
  - Code generation
  - TTL and expiration

- **`one_time_token.go`**: Single-use tokens
  - For password reset
  - For account activation

- **`status.go`**: System states
  - User states
  - Application states

- **`db_models_base.go`**: Base for models with DB metadata
  - ID, CreatedAt, UpdatedAt, DeletedAt

- **`models_utils.go`**: Model utilities
  - Email validation
  - Password validation

#### `/src/domain/utils/`

Domain utilities.

- **`query_payload.go`**: Query construction
  - Filters
  - Sorting
  - Pagination

### `/src/application/` - Application Layer

#### `/src/application/contracts/`

Defines all interfaces (contracts) that infrastructure must implement.

##### `/src/application/contracts/providers/`

External provider interfaces:

- **`cache_provider.go`**: Cache interface (Redis)
  - `Get()`, `Set()`, `Delete()`, `Exists()`

- **`email_provider.go`**: Email sending interface
  - `SendEmail()`

- **`hash_provider.go`**: Hashing interface
  - `Hash()`, `Compare()`

- **`jwt_provider.go`**: JWT interface
  - `GenerateAccessToken()`, `GenerateRefreshToken()`, `ParseTokenAndValidate()`

- **`logger_provider.go`**: Logging interface
  - `Info()`, `Error()`, `Debug()`, `Panic()`

- **`renderer_provider.go`**: Template rendering interface
  - `Render()`

- **`status_provider.go`**: System status interface

##### `/src/application/contracts/repositories/`

Repository interfaces:

- **`base.go`**: Base interface for repositories
  - `Create()`, `GetByID()`, `Update()`, `Delete()`, `GetAll()`

- **`user.go`**: User-specific interface
  - `CreateWithPassword()`, `GetUserWithRole()`, `GetByEmailOrPhone()`

- **`password.go`**: Password interface
  - `GetActivePassword()`, `Create()`

- **`role.go`**: Role interface
  - `GetByKey()`, `GetAll()`

- **`one_time_password.go`**: OTP interface
  - `Create()`, `GetByCode()`, `Invalidate()`

- **`one_time_token.go`**: Token interface
  - `Create()`, `GetByToken()`, `Invalidate()`

#### `/src/application/modules/`

Business modules organized by domain.

##### `/src/application/modules/auth/`

Authentication module:

- **`jwt_auth.go`**: JWT authentication
  - Login with email/password
  - Token generation

- **`jwt_auth_refresh.go`**: Token refresh
  - Access token renewal

- **`jwt_auth_otp.go`**: OTP authentication
  - Login with OTP code

- **`jwt_auth_user.go`**: User authentication from token
  - Token validation
  - User retrieval

- **`get_reset_password_token.go`**: Reset token generation
  - Token creation
  - Email sending

- **`get_reset_password_token_email.go`**: Reset email sending

- **`pipe/get_reset_password.go`**: Password reset pipe
  - Orchestrates token generation and email sending

##### `/src/application/modules/user/`

User module:

- **`use_cases/create_user.go`**: Create user
- **`use_cases/create_user_password.go`**: Create user with password
- **`use_cases/create_user_email.go`**: Send welcome email
- **`use_cases/get_user.go`**: Get user
- **`use_cases/get_all_user.go`**: List users (with cache)
- **`use_cases/update_user.go`**: Update user
- **`use_cases/delete_user.go`**: Delete user
- **`use_cases/activate_user.go`**: Activate user
- **`pipes/create_user.go`**: Pipe to create user and send email

##### `/src/application/modules/password/`

Password module:

- **`use_cases/create_password.go`**: Create password
- **`use_cases/create_password_token.go`**: Create reset token
- **`pipes/create_password_token.go`**: Reset pipe

##### `/src/application/modules/status/`

Status module:

- **`use_cases/status.go`**: System health check

#### `/src/application/shared/`

Components shared between modules.

##### `/src/application/shared/DTOs/`

Data Transfer Objects:

- **`base.go`**: Base DTOs
  - `SingleResponse[T]`, `MultipleResponse[T]`

- **`user.go`**: User DTOs
  - `UserCreate`, `UserUpdate`, `UserAndPasswordCreate`

- **`password.go`**: Password DTOs
  - `PasswordCreate`, `PasswordReset`

- **`token.go`**: Token DTOs
  - `Token` (access + refresh)

- **`one_time_password.go`**: OTP DTOs
- **`one_time_token.go`**: Token DTOs

##### `/src/application/shared/errors/`

Error handling:

- **`application_error.go`**: Application error
  - `Code`: Status code
  - `Context`: Error context
  - `ErrMsg`: Error message

##### `/src/application/shared/use_case/`

Use case base:

- **`base.go`**: Base interface and validation
  - `BaseUseCase[Input, Output]`
  - `BaseUseCaseValidation`

- **`use_case_result.go`**: Use case result
  - `UseCaseResult[T]`
  - `SetData()`, `SetError()`, `HasError()`

- **`dag.go`**: Directed Acyclic Graph
  - `DAG`: Sequential execution
  - `UseCaseParallelDag`: Parallel execution

- **`uc_guards.go`**: Authorization guards
  - Permission validation

##### `/src/application/shared/services/`

Shared services:

- **`create_one_time_password.go`**: Create OTP
- **`create_one_time_token.go`**: Create token
- **`create_password.go`**: Create password
- **`emails/`**: Email services
  - `register_user_email.go`
  - `reset_password_email.go`
  - `otp_email.go`

##### `/src/application/shared/templates/`

Email templates:

- **`emails/`**: HTML templates
  - `register_user.gohtml`
  - `reset_password.gohtml`
  - `otp.gohtml`

##### `/src/application/shared/locales/`

Internationalization:

- **`app_messages.go`**: Application messages
- **`messages/`**: Messages by language
  - `en_us.go`, `es_es.go`

##### `/src/application/shared/context/`

Application context:

- **`context.go`**: Context utilities
- **`keys.go`**: Context keys (UserKey, etc.)

##### `/src/application/shared/settings/`

Configuration:

- **`app_settings.go`**: Application configuration
  - Loads from environment variables
  - Type validation

##### `/src/application/shared/guards/`

Authorization guards:

- **`user.go`**: User guards
  - Permission validation

##### `/src/application/shared/defaults/`

Default values:

- **`user.go`**: User defaults
- **`roles.go`**: Default roles
- **`password.go`**: Password configuration

##### `/src/application/shared/mocks/`

Mocks for testing:

- **`dtos/`**: DTO mocks
- **`mock_*.go`**: Provider and repository mocks

### `/src/infrastructure/` - Infrastructure Layer

#### `/src/infrastructure/server/`

HTTP server layer with Gin.

##### `/src/infrastructure/server/cmd/`

- **`main.go`**: Application entry point
  - Infrastructure initialization (`infrastructure.Initialize()`)
  - Gin configuration with graceful shutdown
  - Middleware loading (CORS, Recovery)
  - Route loading (`routes.Router()`)
  - Server startup on configurable port
  - **Note**: Swagger documentation runs as an independent service (see `/src/infrastructure/docs/`)

**Initialization flow:**
```go
1. infrastructure.Initialize()
   ├── Loads configuration (Settings)
   ├── Initializes Logger
   ├── Connects to PostgreSQL (GORM)
   ├── Configures JWT Provider
   ├── Configures Email Provider
   ├── Configures Cache Provider (Redis)
   └── Configures Email Services

2. buildGinApp()
   └── Creates Gin application with graceful shutdown

3. loadGinApp()
   ├── Configures CORS
   ├── Configures Recovery middleware
   └── Loads routes

4. app.Run()
   └── Starts HTTP server

**Note**: Swagger documentation runs as an independent service in `/src/infrastructure/docs/`
```

##### `/src/infrastructure/api/routes/`

- **`router.go`**: Definition of all routes
  - Public routes
  - Private routes (with authentication)
  - Grouping by module

- **`wrap.go`**: Handler wrapper
  - Adapts Gin to `HandlerContext`
  - Extracts parameters, query, locale

##### `/src/infrastructure/api/middlewares/`

- **`auth.go`**: Authentication middleware
  - Extracts JWT token
  - Validates token
  - Injects user into context

- **`query.go`**: Query params middleware
  - Parses filters, sorting, pagination

#### `/src/infrastructure/config/`

- **`config.go`**: Configuration loading
  - Reads environment variables
  - Validates configuration

- **`env.go`**: Environment variable utilities

#### `/src/infrastructure/database/`

##### `/src/infrastructure/database/goprojectskeleton/`

- **`goprojectskeleton.go`**: GORM configuration
  - PostgreSQL connection
  - SSL configuration

- **`models/`**: Database models (GORM)
  - `user.go`, `role.go`, `password.go`, etc.

- **`init_db/`**: Database initialization
  - Automatic migrations
  - Seeds (if applicable)

#### `/src/infrastructure/handlers/`

HTTP handlers (use case to HTTP adapters).

- **`user.go`**: User handlers
  - `CreateUser()`, `GetUser()`, `UpdateUser()`, etc.

- **`auth.go`**: Authentication handlers
  - `Login()`, `RefreshAccessToken()`, `LoginOTP()`, etc.

- **`password.go`**: Password handlers
  - `CreatePassword()`, `CreatePasswordToken()`

- **`status.go`**: Status handler
  - `GetHealthCheck()`

- **`request_resolver.go`**: HTTP response resolution
  - Maps `UseCaseResult` to HTTP response
  - Handles status codes

- **`types.go`**: Handler types
  - `HandlerContext`
  - `Query`

#### `/src/infrastructure/providers/`

Provider implementations.

- **`jwt_provider.go`**: JWT implementation
  - Token generation and validation

- **`hash_provider.go`**: Hashing implementation
  - Bcrypt for passwords

- **`email_provider.go`**: Email implementation
  - SMTP

- **`cache_provider.go`**: Cache implementation
  - Redis

- **`logger_provider.go`**: Logging implementation
  - Structured logging

- **`renderer_provider.go`**: Rendering implementation
  - HTML template rendering

- **`status_provider.go`**: Status implementation

#### `/src/infrastructure/repositories/`

Repository implementations.

- **`base.go`**: Base repository
  - `RepositoryBase`: Generic CRUD implementation
  - `ModelConverter`: Conversion between DTOs and models

- **`user.go`**: User repository
  - Implements `IUserRepository`
  - Specific methods: `CreateWithPassword()`, `GetUserWithRole()`

- **`password.go`**: Password repository
- **`role.go`**: Role repository
- **`one_time_password.go`**: OTP repository
- **`one_time_token.go`**: Token repository

- **`model_converter.go`**: Model converters
- **`orm_error_map.go`**: ORM error mapping

#### `/src/infrastructure/container.go`

Centralized infrastructure initialization:

- **Configuration**: Loads from environment variables
- **Database**: PostgreSQL connection with GORM
- **Providers**: JWT, Email, Cache, Logger
- **Services**: Email services (registration, reset, OTP)

**Initialization order:**
1. Settings (application configuration)
2. Logger (logging system)
3. Database (PostgreSQL connection)
4. JWT Provider (authentication)
5. Email Provider (SMTP)
6. Cache Provider (Redis)
7. Email Services (email services)

#### `/src/infrastructure/clouds/`

Adaptadores para plataformas cloud y serverless.

##### `/src/infrastructure/clouds/aws/`

Implementation for **AWS Lambda**:

- **`init.go`**: AWS infrastructure initialization
  - Loads configuration from AWS Secrets Manager
  - Initializes database, providers and services
  - Optimized for Lambda cold starts

- **`lambda_adapter.go`**: Lambda event adapter
  - Converts API Gateway events to `HandlerContext`
  - Handles HTTP responses
  - Manages errors and status codes

- **`secrets_manager.go`**: AWS Secrets Manager integration
  - Automatic secret loading from ARNs
  - Fallback to environment variables
  - Secret caching for optimization

- **`middleware.go`**: Lambda-specific middlewares
  - JWT authentication
  - Structured logging
  - Error handling

- **`render_provider.go`**: Rendering provider for AWS
  - HTML template rendering
  - Optimized for serverless environment

- **`functions/`**: Lambda function generator
  - **`main.go`**: CLI tool for generation and deployment
  - **`functions.json`**: Lambda function definitions
  - **`utils/generate.go`**: Code generation from templates
  - **`utils/deploy.go`**: Deployment to AWS Lambda

- **`terraform/`**: Infrastructure as code
  - Lambda function definitions
  - API Gateway
  - IAM roles and policies
  - Variables and outputs

##### `/src/infrastructure/clouds/azure/`

Implementation for **Azure Functions**:

- **`init.go`**: Azure infrastructure initialization
  - Loads configuration from Azure Key Vault
  - Initializes database, providers and services
  - Optimized for Azure Functions

- **`http_adapter.go`**: HTTP trigger adapter
  - Converts HTTP requests to `HandlerContext`
  - Handles HTTP responses
  - Manages errors and status codes

- **`vault.go`**: Azure Key Vault integration
  - Automatic secret loading
  - Managed Identity authentication
  - Secret caching

- **`middleware.go`**: Azure-specific middlewares
  - JWT authentication
  - Structured logging
  - Error handling

- **`functions/`**: Azure function generator
  - **`functions.json`**: Function definitions
  - **`generate.go`**: Code generation
  - Templates for Azure Functions

- **`terraform/`**: Infrastructure as code
  - Function App definitions
  - App Service Plans
  - Key Vault integration
  - Variables and outputs

### `/docker/` - Docker Configuration

- **`docker-compose.dev.yml`**: Development services (includes independent Swagger service)
- **`docker-compose.test.yml`**: Testing services
- **`docker-compose.e2e.yml`**: E2E services (includes Bruno for automated E2E testing)
- **`dockerfile.dev`**: Development Dockerfile
- **`dockerfile.swagger`**: Production Dockerfile for independent Swagger service
- **`dockerfile.swagger.debug`**: Development Dockerfile for Swagger with hot reload
- **`dockerfile.e2e`**: E2E Dockerfile
- **`dockerfile.integration`**: Integration Dockerfile
- **`db/`**: Database configuration
  - `Dockerfile`, `create.sql`

### `/tests/` - Tests

- **`integration/`**: Integration tests
  - `main_test.go`: Test setup
  - Repository tests
  - Provider tests

- **`e2e/`**: End-to-end tests
  - `bruno/`: Bruno collection with tests
  - `collections/`: Bruno request collections
  - `environments/`: Bruno environment configurations

### `/src/infrastructure/docs/` - Independent Swagger Service

Independent HTTP service for Swagger documentation.

- **`main.go`**: Independent HTTP server for Swagger UI
  - Configuration from environment variables
  - Server on configurable port (default: 8081)
  - Automatic redirect to `/docs/`
- **`config/`**: Swagger server configuration
  - `config.go`: Environment variable loading
- **`swagger/`**: Generated Swagger files
  - `swagger.json`: Swagger specification (JSON)
  - `swagger.yaml`: Swagger specification (YAML)
  - `docs.go`: Generated Swagger code
- **`go.mod`**: Independent module for Swagger service

---

## Technologies and Dependencies

### Technology Stack

#### 🚀 Language and Runtime
- **Go 1.25**: Programming language
  - Native concurrency with goroutines
  - Static compilation
  - Automatic tree shaking
  - Excellent performance

#### 🌐 Web Framework
- **Gin v1.10.0**: Minimalist and fast HTTP framework
  - High-performance router
  - Middleware chain
  - Automatic JSON binding
  - Request validation

- **gin-contrib/cors**: CORS middleware
  - Flexible origin configuration
  - Credential support
  - Customizable headers

- **gin-contrib/graceful**: Graceful shutdown
  - Orderly server shutdown
  - Completion of in-flight requests
  - Configurable timeout

#### 💾 Database and Persistence
- **GORM v1.25.12**: ORM for Go
  - Automatic migrations
  - Hooks and callbacks
  - Relationships and associations
  - Fluent query builder

- **PostgreSQL (pgx/v5)**: PostgreSQL driver
  - Native connection pooling
  - Transactions
  - Prepared statements
  - Advanced type support

- **Redis (go-redis/v9)**: Redis client for cache
  - Atomic operations
  - Pub/Sub
  - Pipeline support
  - Cluster support

#### 🔐 Authentication and Security
- **golang-jwt/jwt/v5**: JWT implementation
  - Token generation and validation
  - Multiple algorithms (HS256, RS256, etc.)
  - Custom claims
  - Expiration and time validation

- **golang.org/x/crypto**: Cryptographic utilities
  - Bcrypt for password hashing
  - Automatic salt
  - Configurable cost

#### 📚 Documentation and Testing
- **swaggo/swag**: Swagger documentation generation
  - Annotations in Go code
  - Automatic OpenAPI generation
  - Schema validation

- **swaggo/http-swagger**: Independent HTTP server for Swagger UI
  - Completely independent service
  - Interactive UI
  - Browser testing
  - Authentication in Swagger UI
  - Independent deployment in Docker registry

- **stretchr/testify**: Testing framework
  - Enhanced assertions
  - Mocks and suites
  - Test helpers

#### 🛠️ Utilities
- **joho/godotenv**: Environment variable loading
  - Support for `.env` files
  - Variable override
  - Required variable validation

### Main Dependencies

```go
require (
    // Web Framework
    github.com/gin-gonic/gin v1.10.0
    github.com/gin-contrib/cors v1.7.0
    github.com/gin-contrib/graceful v1.0.0

    // Authentication
    github.com/golang-jwt/jwt/v5 v5.3.0

    // Database
    github.com/jackc/pgx/v5 v5.7.2
    gorm.io/gorm v1.25.12
    gorm.io/driver/postgres v1.5.11

    // Cache
    github.com/redis/go-redis/v9 v9.13.0

    // Security
    golang.org/x/crypto v0.41.0

    // Documentation (Independent Swagger Service)
    github.com/swaggo/swag v1.16.6
    github.com/swaggo/http-swagger v1.6.0

    // Testing
    github.com/stretchr/testify v1.10.0

    // Utilities
    github.com/joho/godotenv v1.5.1
)
```

### Dependency Architecture

```mermaid
graph TB
    subgraph Core["Core"]
        Go[Go 1.25]
    end

    subgraph Web["Web Layer"]
        Gin[Gin Framework]
        CORS[CORS Middleware]
        Graceful[Graceful Shutdown]
    end

    subgraph Data["Data Layer"]
        GORM[GORM ORM]
        PostgreSQL[PostgreSQL Driver]
        Redis[Redis Client]
    end

    subgraph Auth["Authentication"]
        JWT[JWT Library]
        Crypto[Crypto Utils]
    end

    subgraph Docs["Documentation"]
        Swag[Swag Generator]
        SwaggerUI[Swagger UI]
    end

    Go --> Gin
    Gin --> CORS
    Gin --> Graceful
    Gin --> GORM
    GORM --> PostgreSQL
    Gin --> Redis
    Gin --> JWT
    JWT --> Crypto
    Swag --> SwaggerUI
    SwaggerUI -.->|Independent Service| SwaggerServer[Swagger HTTP Server<br/>Port: 8081]

    style Go fill:#00ADD8
    style Gin fill:#00ADD8
    style GORM fill:#00ADD8
```

### Dependencias Principales

```go
require (
    github.com/gin-gonic/gin v1.10.0
    github.com/golang-jwt/jwt/v5 v5.3.0
    github.com/jackc/pgx/v5 v5.7.2
    github.com/redis/go-redis/v9 v9.13.0
    github.com/swaggo/gin-swagger v1.6.0
    gorm.io/gorm v1.25.12
    golang.org/x/crypto v0.41.0
)
```

---

## Configuration and Setup

### Environment Variables

The project uses environment variables for all configuration:

```bash
# Application
APP_NAME=goprojectskeleton
APP_ENV=development
APP_PORT=8080
APP_VERSION=0.0.1
APP_DESCRIPTION=Go Project Skeleton
ENABLE_LOG=true
DEBUG_LOG=true

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=goprojectskeleton
DB_SSL=false

# Redis
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_TTL=300

# JWT
JWT_SECRET_KEY=your-secret-key
JWT_ISSUER=your-issuer
JWT_AUDIENCE=your-audience
JWT_ACCESS_TTL=3600
JWT_REFRESH_TTL=86400
JWT_CLOCK_SKEW=60

# Email
MAIL_HOST=localhost
MAIL_PORT=1025
MAIL_FROM=noreply@example.com
MAIL_PASSWORD=password

# Tokens and OTP
ONE_TIME_TOKEN_TTL=15
ONE_TIME_TOKEN_EMAIL_VERIFY_TTL=60
ONE_TIME_PASSWORD_LENGTH=6
ONE_TIME_PASSWORD_TTL=10
FRONTEND_RESET_PASSWORD_URL=http://localhost:3000/reset-password
FRONTEND_ACTIVATE_ACCOUNT_URL=http://localhost:3000/activate-account
```

### Installation

1. **Clone repository**
```bash
git clone <repository-url>
cd GoProjectSkeleton
```

2. **Configure environment variables**
```bash
cp dev.env.example dev.env
# Edit dev.env with your configurations
```

3. **Install dependencies**
```bash
go mod download
```

4. **Run with Docker**
```bash
# Create network
docker network create goprojectskeleton

# Create volume
docker volume create goprojectskeleton-db-data

# Run services
docker-compose -f docker/docker-compose.dev.yml up -d
```

5. **Run application**
```bash
go run src/infrastructure/api/cmd/main.go
```

---

## Business Modules

### 🔐 Authentication Module (`auth`)

**Responsibility**: Complete authentication and authorization management.

#### Features

- ✅ **Login with Email/Password** - Traditional authentication
- ✅ **Login with OTP** - Two-factor authentication
- ✅ **Token Refresh** - Access token renewal
- ✅ **Password Reset** - Recovery via tokens
- ✅ **User Validation** - Verification from JWT token

#### Detailed Use Cases

**`JwtAuthUseCase`** - Main authentication
```go
// Flow:
// 1. Validates credentials (email/phone + password)
// 2. Verifies password with hash
// 3. If OTP enabled → generates and sends OTP
// 4. If OTP disabled → generates JWT tokens
// 5. Returns tokens or indicates OTP was sent
```

**`JwtAuthRefreshUseCase`** - Token renewal
```go
// Flow:
// 1. Validates refresh token
// 2. Verifies expiration and signature
// 3. Generates new access token
// 4. Returns new token
```

**`JwtAuthOtpUseCase`** - OTP authentication
```go
// Flow:
// 1. Validates OTP code
// 2. Verifies expiration and usage
// 3. Invalidates used OTP
// 4. Generates JWT tokens
// 5. Returns tokens
```

**`GetResetPasswordTokenUseCase`** - Reset token generation
```go
// Flow:
// 1. Finds user by email/phone
// 2. Generates unique token
// 3. Creates DB record with expiration
// 4. Sends email with reset link
```

**`JwtAuthUserUseCase`** - User validation from token
```go
// Flow:
// 1. Extracts token from context
// 2. Validates and parses token
// 3. Finds user in DB
// 4. Returns user with role
```

#### Pipes

**`GetResetPasswordPipe`** - Password reset pipe
- Orchestrates token generation and email sending
- Sequential execution with error handling

### 👥 User Module (`user`)

**Responsibility**: Complete user lifecycle management.

#### Features

- ✅ **Complete CRUD** - Create, read, update, delete
- ✅ **Account Activation** - Activation via tokens
- ✅ **Role Management** - Role assignment and validation
- ✅ **Pagination and Filtering** - Efficient queries
- ✅ **Smart Cache** - List caching with Redis
- ✅ **Transactional Emails** - Welcome and reactivation

#### Detailed Use Cases

**`CreateUserUseCase`** - Create basic user
```go
// Flow:
// 1. Validates input data
// 2. Verifies email/phone don't exist
// 3. Creates user with "pending" status
// 4. Returns created user
```

**`CreateUserAndPasswordUseCase`** - Create user with password
```go
// Flow:
// 1. Validates user and password data
// 2. Password hash with Bcrypt
// 3. Creates user and password in transaction
// 4. Returns created user
```

**`CreateUserSendEmailUseCase`** - Send welcome email
```go
// Flow:
// 1. Renders email template
// 2. Sends email with user data
// 3. Handles sending errors
```

**`GetUserUseCase`** - Get user by ID
```go
// Flow:
// 1. Validates ID
// 2. Finds user in DB
// 3. Includes role information
// 4. Returns user with role
```

**`GetAllUserUseCase`** - List users with filters
```go
// Flow:
// 1. Checks cache (Redis)
// 2. If cache hit → returns from cache
// 3. If cache miss → queries DB with filters
// 4. Applies pagination and sorting
// 5. Saves to cache with TTL
// 6. Returns paginated list
```

**`UpdateUserUseCase`** - Update user
```go
// Flow:
// 1. Validates update data
// 2. Finds existing user
// 3. Applies partial changes (PATCH)
// 4. Validates business rules
// 5. Updates in DB
// 6. Returns updated user
```

**`DeleteUserUseCase`** - Delete user (soft delete)
```go
// Flow:
// 1. Finds user
// 2. Verifies permissions (don't delete admin)
// 3. Soft delete (marks as deleted)
// 4. Invalidates related cache
```

**`ActivateUserUseCase`** - Activate user account
```go
// Flow:
// 1. Validates activation token
// 2. Verifies expiration
// 3. Changes status to "active"
// 4. Invalidates used token
```

**`ResendWelcomeEmailUseCase`** - Resend welcome email
```go
// Flow:
// 1. Finds user
// 2. Renders template
// 3. Sends email
```

#### Pipes

**`CreateUserPipe`** - Pipe to create complete user
- Executes sequentially:
  1. `CreateUserAndPasswordUseCase` → Creates user with password
  2. `CreateUserSendEmailUseCase` → Sends welcome email
- Error handling: If any step fails, execution stops

### 🔑 Password Module (`password`)

**Responsibility**: Secure password management.

#### Features

- ✅ **Password Creation** - Secure hash with Bcrypt
- ✅ **Reset Token Generation** - Unique tokens with expiration
- ✅ **Strength Validation** - Secure password rules
- ✅ **Password Expiration** - Temporary passwords

#### Detailed Use Cases

**`CreatePasswordUseCase`** - Create new password
```go
// Flow:
// 1. Validates password strength
// 2. Hash with Bcrypt
// 3. Deactivates previous passwords
// 4. Creates new active password
// 5. Configures expiration if applicable
```

**`CreatePasswordTokenUseCase`** - Create reset token
```go
// Flow:
// 1. Finds user by email/phone
// 2. Generates unique token
// 3. Creates record with expiration
// 4. Sends email with reset link
```

#### Pipes

**`CreatePasswordTokenPipe`** - Password reset pipe
- Orchestrates token creation and email sending

### 📊 Status Module (`status`)

**Responsibility**: System monitoring and health.

#### Features

- ✅ **Health Check** - Overall system status
- ✅ **Service Verification** - DB, Redis, etc.
- ✅ **Version Information** - Application version

#### Detailed Use Cases

**`GetStatusUseCase`** - Get system status
```go
// Flow:
// 1. Verifies PostgreSQL connection
// 2. Verifies Redis connection
// 3. Verifies providers (JWT, Email)
// 4. Returns consolidated status
// 5. Includes version information
```

### 📈 Statistics by Module

| Module | Use Cases | Pipes | Tests | Endpoints |
|--------|-----------|-------|-------|-----------|
| **auth** | 5 | 1 | 5+ | 4 |
| **user** | 9 | 1 | 9+ | 7 |
| **password** | 2 | 1 | 2+ | 2 |
| **status** | 1 | 0 | 1+ | 1 |
| **Total** | **17** | **3** | **17+** | **14** |

---

## API and Endpoints

### Authentication

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|---------------|
| POST | `/api/auth/login` | Login with credentials | No |
| POST | `/api/auth/refresh` | Renew access token | No |
| GET | `/api/auth/login-otp/{otp}` | Login with OTP | No |
| GET | `/api/auth/password-reset/{identifier}` | Request password reset | No |

### Users

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|---------------|
| POST | `/api/user` | Create user | No |
| GET | `/api/user/{id}` | Get user | Yes |
| PATCH | `/api/user/{id}` | Update user | Yes |
| DELETE | `/api/user/{id}` | Delete user | Yes |
| GET | `/api/user` | List users (with filters) | Yes |
| POST | `/api/user-password` | Create user with password | No |
| POST | `/api/user/activate` | Activate user | No |

### Passwords

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|---------------|
| POST | `/api/password` | Create password | Yes |
| POST | `/api/password/reset-token` | Create reset token | No |

### System

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|---------------|
| GET | `/api/health-check` | Health check | No |

**Note**: Swagger documentation is available as an independent service on port 8081 at `http://localhost:8081/docs/`

### Usage Examples

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Create User
```bash
curl -X POST http://localhost:8080/api/user \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "role_id": 2,
    "status": "pending",
    "otpLogin": false
  }'
```

---

## Database and Persistence

### Entity-Relationship Diagram (ERD)

```mermaid
erDiagram
    USER ||--o{ PASSWORD : "tiene"
    USER }o--|| ROLE : "pertenece"
    USER ||--o{ ONE_TIME_PASSWORD : "genera"
    USER ||--o{ ONE_TIME_TOKEN : "genera"

    USER {
        uint id PK
        string name
        string email UK
        string phone UK
        string status
        uint role_id FK
        bool otpLogin
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    ROLE {
        uint id PK
        string key UK
        bool isActive
        int priority
        datetime created_at
        datetime updated_at
    }

    PASSWORD {
        uint id PK
        uint user_id FK
        string hash
        bool isActive
        datetime expiresAt
        datetime created_at
        datetime updated_at
    }

    ONE_TIME_PASSWORD {
        uint id PK
        uint user_id FK
        string code
        bool isUsed
        datetime expiresAt
        datetime created_at
        datetime updated_at
    }

    ONE_TIME_TOKEN {
        uint id PK
        uint user_id FK
        string token
        string type
        bool isUsed
        datetime expiresAt
        datetime created_at
        datetime updated_at
    }
```

### Main Models

#### User
```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique;not null"`
    Phone     string    `gorm:"not null"`
    Status    string    `gorm:"not null"`
    RoleID    uint      `gorm:"not null"`
    OTPLogin  bool      `gorm:"default:false"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

#### Password
```go
type Password struct {
    ID        uint       `gorm:"primaryKey"`
    UserID    uint       `gorm:"not null"`
    Hash      string     `gorm:"not null"`
    IsActive  bool       `gorm:"default:true"`
    ExpiresAt *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Role
```go
type Role struct {
    ID          uint      `gorm:"primaryKey"`
    Key         string    `gorm:"unique;not null"`
    Name        string    `gorm:"not null"`
    Description string
    IsActive    bool      `gorm:"default:true"`
    Priority    int       `gorm:"default:0"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Relationship Diagram

```mermaid
graph LR
    subgraph UserDomain["👤 User"]
        User[User<br/>ID, Name, Email, Phone<br/>Status, RoleID, OTPLogin]
    end

    subgraph AuthDomain["🔐 Authentication"]
        Password[Password<br/>UserID, Hash<br/>IsActive, ExpiresAt]
        OTP[OneTimePassword<br/>UserID, Code<br/>IsUsed, ExpiresAt]
        Token[OneTimeToken<br/>UserID, Token, Type<br/>IsUsed, ExpiresAt]
    end

    subgraph RoleDomain["👥 Roles"]
        Role[Role<br/>ID, Key, Name<br/>IsActive, Priority]
    end

    User -->|1:N| Password
    User -->|1:N| OTP
    User -->|1:N| Token
    User -->|N:1| Role

    style User fill:#e3f2fd
    style Password fill:#fff9c4
    style OTP fill:#fff9c4
    style Token fill:#fff9c4
    style Role fill:#e8f5e9
```

### Migrations

Migrations run automatically when the application starts using GORM AutoMigrate.

### Cache

The system uses Redis to cache frequent queries, especially in user listings.

```mermaid
graph TB
    subgraph App["Application"]
        Handler[Handler]
        UseCase[Use Case]
        Cache[Cache Provider]
    end

    subgraph Storage["Storage"]
        Redis[(Redis<br/>Cache)]
        DB[(PostgreSQL<br/>Database)]
    end

    Handler --> UseCase
    UseCase --> Cache
    Cache -->|Cache Hit| Redis
    Cache -->|Cache Miss| DB
    Cache -->|Write| Redis
    UseCase -->|Direct Query| DB

    style Redis fill:#ffcdd2
    style DB fill:#c8e6c9
```

---

## Authentication and Security

### Authentication Flow Diagram

```mermaid
sequenceDiagram
    participant Client as Client
    participant API as API
    participant AuthUC as Auth Use Case
    participant UserRepo as User Repository
    participant PassRepo as Password Repository
    participant Hash as Hash Provider
    participant JWT as JWT Provider
    participant DB as PostgreSQL

    Client->>API: POST /api/auth/login<br/>{email, password}
    API->>AuthUC: Execute(credentials)
    AuthUC->>UserRepo: GetByEmailOrPhone()
    UserRepo->>DB: SELECT user
    DB-->>UserRepo: User
    UserRepo-->>AuthUC: User

    AuthUC->>PassRepo: GetActivePassword()
    PassRepo->>DB: SELECT password
    DB-->>PassRepo: Password
    PassRepo-->>AuthUC: Password

    AuthUC->>Hash: Compare(password, hash)
    Hash-->>AuthUC: Valid

    AuthUC->>JWT: GenerateAccessToken()
    JWT-->>AuthUC: Access Token
    AuthUC->>JWT: GenerateRefreshToken()
    JWT-->>AuthUC: Refresh Token

    AuthUC-->>API: Tokens
    API-->>Client: {accessToken, refreshToken}
```

### OTP Authentication Flow

```mermaid
sequenceDiagram
    participant Client as Client
    participant API as API
    participant AuthUC as Auth Use Case
    participant OTPUC as OTP Use Case
    participant OTPRepo as OTP Repository
    participant EmailSvc as Email Service
    participant DB as PostgreSQL
    participant SMTP as SMTP

    Client->>API: POST /api/auth/login<br/>{email, password}
    API->>AuthUC: Execute(credentials)
    AuthUC->>AuthUC: Validates credentials
    AuthUC->>AuthUC: Is OTP Login enabled?

    alt OTP Login enabled
        AuthUC->>OTPUC: GenerateOTP()
        OTPUC->>OTPRepo: Create()
        OTPRepo->>DB: INSERT OTP
        OTPUC->>EmailSvc: SendOTPEmail()
        EmailSvc->>SMTP: Send email
        AuthUC-->>API: 204 No Content
        API-->>Client: OTP sent by email
    else OTP Login disabled
        AuthUC->>JWT: GenerateTokens()
        AuthUC-->>API: Tokens
        API-->>Client: {accessToken, refreshToken}
    end

    Note over Client,SMTP: User enters OTP
    Client->>API: GET /api/auth/login-otp/{otp}
    API->>OTPUC: ValidateOTP(otp)
    OTPUC->>OTPRepo: GetByCode()
    OTPRepo->>DB: SELECT OTP
    OTPUC->>OTPUC: Validates expiration
    OTPUC->>JWT: GenerateTokens()
    OTPUC-->>API: Tokens
    API-->>Client: {accessToken, refreshToken}
```

### JWT (JSON Web Tokens)

The system uses JWT for authentication with two types of tokens:

```mermaid
graph TB
    subgraph TokenFlow["Token Flow"]
        Login[Login] --> AccessToken[Access Token<br/>TTL: 1 hour]
        Login --> RefreshToken[Refresh Token<br/>TTL: 24 hours]

        AccessToken -->|Expires| Refresh[Refresh Endpoint]
        Refresh --> NewAccess[New Access Token]

        AccessToken -->|Valid| Protected[Protected Resources]
    end

    subgraph TokenStructure["JWT Structure"]
        Header[Header<br/>alg: HS256<br/>typ: JWT]
        Payload[Payload<br/>iss, aud, sub<br/>iat, exp, typ]
        Signature[Signature<br/>HMAC SHA256]

        Header --> JWT[JWT Token]
        Payload --> JWT
        Signature --> JWT
    end

    style AccessToken fill:#c8e6c9
    style RefreshToken fill:#fff9c4
    style Protected fill:#e3f2fd
```

1. **Access Token**: Access token (TTL: 1 hour by default)
2. **Refresh Token**: Renewal token (TTL: 24 hours by default)

### OTP (One-Time Password)

Two-factor authentication system:

```mermaid
stateDiagram-v2
    [*] --> Login: Credentials
    Login --> CheckOTP: Validate
    CheckOTP --> GenerateOTP: OTP enabled
    CheckOTP --> GenerateJWT: OTP disabled

    GenerateOTP --> SendEmail: Code generated
    SendEmail --> WaitOTP: Email sent
    WaitOTP --> ValidateOTP: User enters code
    ValidateOTP --> GenerateJWT: Valid code
    ValidateOTP --> Expired: Expired code
    ValidateOTP --> Invalid: Invalid code

    GenerateJWT --> [*]: Tokens generated
    Expired --> [*]
    Invalid --> [*]
```

- Automatic OTP code generation
- Email delivery
- Configurable TTL (default: 10 minutes)

### Security

```mermaid
graph TB
    subgraph SecurityLayers["Security Layers"]
        Input[Input Validation<br/>DTOs]
        Auth[Authentication<br/>JWT Middleware]
        Authz[Authorization<br/>Guards/Roles]
        Hash[Password Hashing<br/>Bcrypt]
        CORS[CORS<br/>Cross-Origin]
        Sanitize[Sanitization<br/>SQL Injection]
    end

    Request[HTTP Request] --> Input
    Input --> Auth
    Auth --> Authz
    Authz --> BusinessLogic[Business Logic]

    Password[Password] --> Hash
    Hash --> Storage[(Database)]

    Request --> CORS
    Request --> Sanitize

    style Input fill:#e3f2fd
    style Auth fill:#fff9c4
    style Authz fill:#f3e5f5
    style Hash fill:#c8e6c9
```

- **Password hashing**: Bcrypt with automatic salt
- **Input validation**: On all endpoints
- **Authentication middleware**: For protected routes
- **CORS**: Configured for web security
- **Sanitization**: Injection prevention

---

## Testing

### Test Structure

```
tests/
├── e2e/                    # End-to-end tests
│   ├── bruno/              # Bruno collection
│   │   ├── collections/    # Request collections
│   │   ├── environments/  # Environment configurations
│   │   └── bruno.json      # Bruno configuration
│   └── postman/            # Legacy Postman collection (optional)
└── integration/            # Integration tests
    ├── main_test.go        # Test setup
    ├── user_repository_test.go
    ├── password_repository_test.go
    └── email_provider_test.go
```

### Testing Types

1. **Unit Tests**: Individual unit tests (in each module)
2. **Integration Tests**: Integration tests with database
3. **E2E Tests**: End-to-end tests with Bruno

### Run Tests

```bash
# Unit tests
go test ./src/...

# Integration tests
go test ./tests/integration/...

# E2E tests (requires running services)
# Bruno tests run automatically in docker-compose.e2e.yml
# Or run manually: bruno run --env=dev
```

---

## Docker and Deployment

### Docker Architecture

```mermaid
graph TB
    subgraph DockerNetwork["Docker Network: goprojectskeleton"]
        subgraph AppContainer["Go Application"]
            App[Go Application<br/>Port: 8080<br/>Hot Reload]
        end

        subgraph DBServices["Data Services"]
            PostgreSQL[(PostgreSQL<br/>Port: 5432<br/>Volume: db-data)]
            Redis[(Redis<br/>Port: 6379<br/>Cache)]
        end

        subgraph SwaggerService["Independent Swagger Service"]
            Swagger[Swagger Server<br/>Port: 8081<br/>Independent]
        end

        subgraph DevTools["Development Tools"]
            Mailpit[Mailpit<br/>Port: 8025<br/>Email Testing]
            RedisCommander[Redis Commander<br/>Port: 18081<br/>Redis UI]
        end
    end

    App -->|GORM| PostgreSQL
    App -->|go-redis| Redis
    App -->|SMTP| Mailpit

    Swagger -.->|Documentation| App
    RedisCommander -->|UI| Redis

    style App fill:#e3f2fd
    style PostgreSQL fill:#c8e6c9
    style Redis fill:#ffcdd2
    style Mailpit fill:#fff9c4
    style RedisCommander fill:#f3e5f5
```

### Deployment Diagram

```mermaid
graph TB
    subgraph Dev["🛠️ Development"]
        DevApp[Go App<br/>Hot Reload]
        DevDB[(PostgreSQL<br/>Dev)]
        DevRedis[(Redis<br/>Dev)]
        DevMail[Mailpit]
    end

    subgraph Test["🧪 Testing"]
        TestApp[Go App<br/>Test Mode]
        TestDB[(PostgreSQL<br/>Test)]
        TestRedis[(Redis<br/>Test)]
    end

    subgraph Prod["🚀 Production"]
        LB[Load Balancer]
        App1[App Instance 1]
        App2[App Instance 2]
        App3[App Instance N]
        ProdDB[(PostgreSQL<br/>Primary)]
        ProdDBReplica[(PostgreSQL<br/>Replica)]
        ProdRedis[(Redis<br/>Cluster)]
        ProdSMTP[SMTP Server]
    end

    LB --> App1
    LB --> App2
    LB --> App3

    App1 --> ProdDB
    App2 --> ProdDB
    App3 --> ProdDB

    App1 --> ProdDBReplica
    App2 --> ProdDBReplica
    App3 --> ProdDBReplica

    App1 --> ProdRedis
    App2 --> ProdRedis
    App3 --> ProdRedis

    App1 --> ProdSMTP
    App2 --> ProdSMTP
    App3 --> ProdSMTP

    style Dev fill:#e1f5ff
    style Test fill:#fff4e1
    style Prod fill:#e8f5e9
```

### Docker Services

The project includes Docker configuration for development:

- **Application**: Go server with hot reload (port 8080)
- **Swagger**: Independent documentation service (port 8081)
- **PostgreSQL**: Main database
- **Redis**: Cache and sessions
- **Mailpit**: Email server for development
- **Redis Commander**: Web interface for Redis (port 18081)

**E2E Testing Services** (docker-compose.e2e.yml):
- **Application**: Go server for E2E testing
- **PostgreSQL**: Test database
- **Redis**: Test cache
- **Mailpit**: Email testing
- **Bruno**: Automated E2E test runner (executes tests automatically)

### Docker Commands

```bash
# Development
docker-compose -f docker/docker-compose.dev.yml up -d

# Testing
docker-compose -f docker/docker-compose.test.yml up -d

# E2E Testing (Bruno runs automatically)
docker-compose -f docker/docker-compose.e2e.yml up -d
# Bruno tests will execute automatically once the API is healthy
```

### Independent Swagger Deployment

The Swagger service can be deployed completely independently from the main application:

#### Image Build

```bash
# Build Swagger image
docker build -f docker/dockerfile.swagger -t your-registry/swagger:latest .

# Tag for specific version
docker tag your-registry/swagger:latest your-registry/swagger:v1.0.0
```

#### Push to Docker Registry

```bash
# Push to Docker Hub
docker push your-registry/swagger:latest

# Push to other registries (e.g., AWS ECR, Google GCR, Azure ACR)
# AWS ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker tag your-registry/swagger:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/swagger:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/swagger:latest
```

#### Container Execution

```bash
# Run with environment variables
docker run -d \
  -p 8081:8081 \
  -e SWAGGER_PORT=8081 \
  -e API_HOST=api.example.com:8080 \
  -e API_TITLE="My API Documentation" \
  -e API_VERSION="1.0" \
  -e API_DESCRIPTION="API documentation" \
  -e API_BASE_PATH="/api" \
  --name swagger-docs \
  your-registry/swagger:latest

# Or using a .env file
docker run -d \
  -p 8081:8081 \
  --env-file swagger.env \
  --name swagger-docs \
  your-registry/swagger:latest
```

#### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: swagger-docs
spec:
  replicas: 1
  selector:
    matchLabels:
      app: swagger-docs
  template:
    metadata:
      labels:
        app: swagger-docs
    spec:
      containers:
      - name: swagger
        image: your-registry/swagger:latest
        ports:
        - containerPort: 8081
        env:
        - name: SWAGGER_PORT
          value: "8081"
        - name: API_HOST
          value: "api.example.com:8080"
        - name: API_TITLE
          value: "My API Documentation"
---
apiVersion: v1
kind: Service
metadata:
  name: swagger-docs-service
spec:
  selector:
    app: swagger-docs
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8081
  type: LoadBalancer
```

#### Advantages of Independent Deployment

- ✅ **Independent scalability**: Scale Swagger without affecting the main application
- ✅ **Updates without downtime**: Update documentation without restarting the API
- ✅ **Separation of concerns**: Documentation separated from business logic
- ✅ **Different environments**: Different documentation versions for dev/staging/prod
- ✅ **CDN and caching**: Serve documentation from CDN for better performance

---

## GitHub Actions Deployment

**GoProjectSkeleton** includes a comprehensive GitHub Actions workflow for automated deployment to AWS and Azure cloud platforms. The workflow supports infrastructure provisioning with Terraform and automated function deployment.

### Overview

The deployment workflow (`deploy.yml`) provides:

- ✅ **Multi-cloud support**: Deploy to AWS or Azure
- ✅ **Environment management**: Separate deployments for development, staging, and production
- ✅ **Terraform integration**: Infrastructure as Code with plan, apply, and destroy actions
- ✅ **Automated function deployment**: Deploy Lambda/Functions after infrastructure (AWS only)
- ✅ **Plan artifacts**: Upload Terraform plans for review
- ✅ **Safety checks**: Prevent accidental production destruction

### Deployment Flow

```mermaid
graph TB
    subgraph Trigger["🚀 Workflow Trigger"]
        Manual[Manual Dispatch<br/>GitHub Actions UI]
    end

    subgraph Inputs["📋 Workflow Inputs"]
        Cloud[Cloud Provider<br/>AWS or Azure]
        Env[Environment<br/>dev/staging/prod]
        Action[Terraform Action<br/>plan/apply/destroy]
        DeployFunc[Deploy Functions<br/>AWS only]
    end

    subgraph Setup["⚙️ Setup Phase"]
        Checkout[Checkout Code]
        GoSetup[Setup Go 1.25.5]
        TfSetup[Install Terraform 1.14.1]
        Creds[Configure Cloud Credentials]
    end

    subgraph Build["🔨 Build Phase"]
        Deps[Download Dependencies]
        GenAWS[Generate AWS Functions<br/>if AWS]
        GenAzure[Generate Azure Functions<br/>if Azure]
        Tfvars[Create terraform.tfvars]
    end

    subgraph Terraform["🏗️ Terraform Phase"]
        Init[Terraform Init]
        Validate[Terraform Validate]
        Plan[Terraform Plan<br/>if not destroy]
        UploadPlan[Upload Plan Artifact<br/>if plan action]
        Apply[Terraform Apply<br/>if apply action]
        Destroy[Terraform Destroy<br/>if destroy action]
    end

    subgraph Deploy["📦 Deploy Phase"]
        DeployLambda[Deploy Lambda Functions<br/>AWS only]
        Output[Terraform Output]
    end

    Manual --> Cloud
    Manual --> Env
    Manual --> Action
    Manual --> DeployFunc

    Cloud --> Checkout
    Env --> Checkout
    Action --> Checkout

    Checkout --> GoSetup
    GoSetup --> TfSetup
    TfSetup --> Creds
    Creds --> Deps

    Deps --> GenAWS
    Deps --> GenAzure
    GenAWS --> Tfvars
    GenAzure --> Tfvars

    Tfvars --> Init
    Init --> Validate
    Validate --> Plan
    Plan --> UploadPlan
    Plan --> Apply
    Apply --> DeployLambda
    DeployLambda --> Output

    Validate --> Destroy
    Destroy --> Output

    style Manual fill:#e3f2fd
    style Cloud fill:#fff9c4
    style Env fill:#fff9c4
    style Action fill:#fff9c4
    style Apply fill:#c8e6c9
    style Destroy fill:#ffcdd2
    style DeployLambda fill:#ff9800
```

### Workflow Inputs

When triggering the workflow manually, you'll be prompted for:

| Input | Description | Options | Default |
|-------|-------------|---------|---------|
| `cloud` | Cloud provider | `aws`, `azure` | Required |
| `environment` | Target environment | `development`, `staging`, `production` | `development` |
| `terraform_action` | Terraform operation | `plan`, `apply`, `destroy` | `apply` |
| `deploy_functions` | Deploy functions after Terraform (AWS only) | `true`, `false` | `true` |

### GitHub Secrets Configuration

Before using the deployment workflow, you need to configure GitHub Secrets. Secrets are environment-specific and can be set at the repository or environment level.

#### Setting Up GitHub Secrets

1. **Navigate to Repository Settings**
   - Go to your GitHub repository
   - Click on **Settings** → **Secrets and variables** → **Actions**

2. **Create Environment Secrets (Recommended)**
   - Click on **Environments** in the left sidebar
   - Create environments: `development`, `staging`, `production`
   - Add secrets to each environment as needed

3. **Create Repository Secrets (Alternative)**
   - Add secrets at the repository level (available to all environments)

#### Required Secrets

##### AWS Secrets

| Secret Name | Description | Example | Required For |
|-------------|-------------|---------|--------------|
| `AWS_ACCESS_KEY_ID` | AWS access key ID | `AKIAIOSFODNN7EXAMPLE` | AWS deployments |
| `AWS_SECRET_ACCESS_KEY` | AWS secret access key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` | AWS deployments |
| `AWS_REGION` | AWS region | `us-east-1` | AWS deployments (optional, defaults to `us-east-1`) |
| `PROJECT_NAME` | Project name for resource naming | `go-project-skeleton` | AWS deployments (optional, defaults to `go-project-skeleton`) |
| `TFVARS` | Terraform variables file content | See below | All AWS deployments |

##### Azure Secrets

| Secret Name | Description | Example | Required For |
|-------------|-------------|---------|--------------|
| `AZURE_CREDENTIALS` | Azure service principal JSON | See below | Azure deployments |
| `TFVARS` | Terraform variables file content | See below | All Azure deployments |

#### Creating Secrets

##### AWS Access Keys

1. **Create IAM User** (if not exists):
   ```bash
   aws iam create-user --user-name github-actions-deploy
   ```

2. **Attach Policies**:
   ```bash
   aws iam attach-user-policy \
     --user-name github-actions-deploy \
     --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
   ```
   > **Note**: For production, use least-privilege policies. Create custom policies with only required permissions.

3. **Create Access Key**:
   ```bash
   aws iam create-access-key --user-name github-actions-deploy
   ```

4. **Add to GitHub Secrets**:
   - Copy `AccessKeyId` → `AWS_ACCESS_KEY_ID`
   - Copy `SecretAccessKey` → `AWS_SECRET_ACCESS_KEY`

##### Azure Service Principal

1. **Create Service Principal**:
   ```bash
   az ad sp create-for-rbac --name github-actions-deploy \
     --role contributor \
     --scopes /subscriptions/{subscription-id} \
     --sdk-auth
   ```

2. **Copy the JSON output** and add it to GitHub Secret `AZURE_CREDENTIALS`:
   ```json
   {
     "clientId": "xxx",
     "clientSecret": "xxx",
     "subscriptionId": "xxx",
     "tenantId": "xxx",
     "activeDirectoryEndpointUrl": "https://login.microsoftonline.com",
     "resourceManagerEndpointUrl": "https://management.azure.com/",
     "activeDirectoryGraphResourceId": "https://graph.windows.net/",
     "sqlManagementEndpointUrl": "https://management.core.windows.net:8443/",
     "galleryEndpointUrl": "https://gallery.azure.com/",
     "managementEndpointUrl": "https://management.core.windows.net/"
   }
   ```

##### Terraform Variables (TFVARS)

The `TFVARS` secret contains the content of your `terraform.tfvars` file. This should include all required Terraform variables for your infrastructure.

**Example for AWS:**
```hcl
# terraform.tfvars content
project_name = "go-project-skeleton"
environment = "development"
region = "us-east-1"
db_instance_class = "db.t3.micro"
lambda_memory_size = 512
# ... other variables
```

**Example for Azure:**
```hcl
# terraform.tfvars content
project_name = "go-project-skeleton"
environment = "development"
location = "eastus"
app_service_plan_sku = "B1"
# ... other variables
```

**To create the secret:**
1. Create your `terraform.tfvars` file locally
2. Copy the entire content
3. Add to GitHub Secret `TFVARS` (paste the entire content)

### Workflow Steps

#### 1. Setup Phase
- **Checkout code**: Clones the repository
- **Setup Go**: Installs Go 1.25.5
- **Install Terraform**: Installs Terraform 1.14.1
- **Configure credentials**: Sets up AWS or Azure credentials based on selected cloud

#### 2. Build Phase
- **Download dependencies**: Runs `make deps`
- **Generate functions**:
  - AWS: Runs `make build-aws-functions` (if `deploy_functions` is true)
  - Azure: Runs `make build-azure-functions`
- **Create terraform.tfvars**: Creates the file from `TFVARS` secret

#### 3. Terraform Phase
- **Terraform Init**: Initializes Terraform backend
- **Terraform Validate**: Validates Terraform configuration
- **Terraform Plan**: Creates execution plan (if action is not `destroy`)
  - Uploads plan artifact for review (if action is `plan`)
- **Terraform Apply**: Applies infrastructure changes (if action is `apply`)
- **Terraform Destroy**: Destroys infrastructure (if action is `destroy`)
  - **Safety**: Destroy fails in production unless explicitly allowed

#### 4. Deploy Phase (AWS only)
- **Deploy Lambda Functions**: Runs `make deploy-aws` to deploy all Lambda functions
- **Terraform Output**: Displays infrastructure outputs

### Usage Examples

#### Plan Infrastructure Changes

1. Go to **Actions** tab in GitHub
2. Select **Deploy to Cloud** workflow
3. Click **Run workflow**
4. Fill inputs:
   - Cloud: `aws`
   - Environment: `staging`
   - Terraform action: `plan`
   - Deploy functions: `false`
5. Click **Run workflow**

**Result**: Creates a Terraform plan and uploads it as an artifact. Review the plan before applying.

#### Deploy to Staging

1. Go to **Actions** tab
2. Select **Deploy to Cloud** workflow
3. Click **Run workflow**
4. Fill inputs:
   - Cloud: `aws`
   - Environment: `staging`
   - Terraform action: `apply`
   - Deploy functions: `true`
5. Click **Run workflow**

**Result**:
- Provisions infrastructure with Terraform
- Deploys all Lambda functions
- Displays infrastructure outputs

#### Destroy Development Environment

1. Go to **Actions** tab
2. Select **Deploy to Cloud** workflow
3. Click **Run workflow**
4. Fill inputs:
   - Cloud: `aws`
   - Environment: `development`
   - Terraform action: `destroy`
5. Click **Run workflow`

**Result**: Destroys all infrastructure in the development environment.

> **⚠️ Warning**: Destroy operations in production will fail by default. The workflow uses `continue-on-error: ${{ inputs.environment != 'production' }}` to prevent accidental production destruction.

### Environment Protection

GitHub Environments can be configured with protection rules:

1. **Required Reviewers**: Require approval before deployment
2. **Wait Timer**: Add a delay before deployment
3. **Deployment Branches**: Restrict which branches can deploy

**To configure:**
1. Go to **Settings** → **Environments**
2. Click on an environment (e.g., `production`)
3. Add protection rules as needed

### Best Practices

1. **Use Environment Secrets**: Store secrets per environment for better security
2. **Review Plans**: Always run `plan` before `apply` in production
3. **Use Protection Rules**: Enable required reviewers for production
4. **Monitor Deployments**: Check workflow runs regularly
5. **Rotate Credentials**: Regularly rotate access keys and secrets
6. **Least Privilege**: Use IAM roles/policies with minimum required permissions
7. **Version Control**: Keep Terraform code in version control
8. **Backup State**: Ensure Terraform state is backed up (S3, Azure Storage)

### Troubleshooting

#### Common Issues

**Issue**: "AWS credentials not found"
- **Solution**: Ensure `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are set in GitHub Secrets

**Issue**: "Terraform plan fails"
- **Solution**: Check `TFVARS` secret content matches expected format

**Issue**: "Lambda deployment fails"
- **Solution**: Ensure `PROJECT_NAME` and `AWS_REGION` are set correctly

**Issue**: "Destroy fails in production"
- **Solution**: This is by design. Modify workflow if production destruction is needed.

### Workflow File Location

The deployment workflow is located at:
```
.github/workflows/deploy.yml
```

---

## Development Guide

### Adding New Functionality

#### Flow Diagram for Adding Functionality

```mermaid
flowchart TD
    Start([New Functionality]) --> Domain[1. Domain Layer<br/>Create Model]
    Domain --> Contracts[2. Application Contracts<br/>Define Interfaces]
    Contracts --> UseCase[3. Application Use Case<br/>Implement Logic]
    UseCase --> Repo[4. Infrastructure Repository<br/>Implement with GORM]
    Repo --> Handler[5. Infrastructure Handler<br/>Adapt HTTP]
    Handler --> Route[6. Routes<br/>Define Endpoint]
    Route --> Tests[7. Tests<br/>Write Tests]
    Tests --> End([✅ Completed])

    style Domain fill:#e8f5e9
    style Contracts fill:#fff9c4
    style UseCase fill:#e3f2fd
    style Repo fill:#f3e5f5
    style Handler fill:#ffebee
    style Route fill:#e1f5ff
    style Tests fill:#fff4e1
```

#### Detailed Steps

1. **Create model in `domain/models/`**
```go
type NewEntity struct {
    // Fields
}
```

2. **Define interface in `application/contracts/repositories/`**
```go
type INewEntityRepository interface {
    Create(input NewEntityCreate) (*NewEntity, error)
    // More methods
}
```

3. **Implement use case in `application/modules/newmodule/use_cases/`**
```go
type CreateNewEntityUseCase struct {
    repo INewEntityRepository
}

func (uc *CreateNewEntityUseCase) Execute(...) *UseCaseResult[NewEntity] {
    // Business logic
}
```

4. **Create repository in `infrastructure/repositories/`**
```go
type NewEntityRepository struct {
    RepositoryBase[...]
}

func (r *NewEntityRepository) Create(...) (*NewEntity, error) {
    // Implementation with GORM
}
```

5. **Add handler in `infrastructure/handlers/`**
```go
func CreateNewEntity(ctx HandlerContext) {
    // Adapt HTTP to use case
}
```

6. **Define route in `infrastructure/api/routes/router.go`**
```go
r.POST("/new-entity", wrapHandler(handlers.CreateNewEntity))
```

7. **Write tests**

### Code Conventions

- **Names**: camelCase for variables, PascalCase for types
- **Comments**: Document public functions
- **Errors**: Use centralized error system
- **Logging**: Use centralized logger
- **Validation**: Validate input on all endpoints

---

## 🎯 Best Practices and Conventions

### Code Conventions

#### Naming
- **Variables**: `camelCase` (e.g., `userName`, `isActive`)
- **Types and Structs**: `PascalCase` (e.g., `User`, `UserRepository`)
- **Interfaces**: `I` + `PascalCase` (e.g., `IUserRepository`, `IHashProvider`)
- **Constants**: `UPPER_SNAKE_CASE` (e.g., `USER_STATUS_ACTIVE`)
- **Packages**: `lowercase` (e.g., `user`, `auth`)

#### File Structure
- **Use Cases**: `{action}_{entity}.go` (e.g., `create_user.go`, `get_user.go`)
- **Repositories**: `{entity}.go` (e.g., `user.go`, `password.go`)
- **Handlers**: `{entity}.go` or `{module}.go` (e.g., `user.go`, `auth.go`)
- **Tests**: `{file}_test.go` (e.g., `create_user_test.go`)

#### Code Organization
```go
// 1. Imports (std, third-party, local)
import (
    "context"
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/simon3640/goprojectskeleton/src/domain/models"
)

// 2. Types and Structs
type UserRepository struct {
    // ...
}

// 3. Constructors
func NewUserRepository(...) *UserRepository {
    // ...
}

// 4. Public methods
func (r *UserRepository) Create(...) {
    // ...
}

// 5. Private methods
func (r *UserRepository) validate(...) {
    // ...
}
```

### Applied Design Principles

#### 1. Single Responsibility Principle (SRP)
Each component has a single responsibility:
- **Use Cases**: Specific business logic
- **Repositories**: Data access
- **Handlers**: HTTP adaptation
- **Providers**: External services

#### 2. Dependency Inversion Principle (DIP)
Inner layers define interfaces that outer layers implement:
```go
// Application defines the interface
type IUserRepository interface {
    Create(input UserCreate) (*User, error)
}

// Infrastructure implements the interface
type UserRepository struct {
    DB *gorm.DB
}
```

#### 3. Open/Closed Principle (OCP)
Open for extension, closed for modification:
- New providers without modifying existing code
- New use cases without affecting others
- New repositories following interfaces

#### 4. Interface Segregation Principle (ISP)
Specific and small interfaces:
```go
// Instead of a large interface
type IProvider interface {
    // 50 methods...
}

// Specific interfaces
type IHashProvider interface {
    Hash(password string) (string, error)
    Compare(hashed, plain string) bool
}
```

#### 5. Don't Repeat Yourself (DRY)
Reuse through:
- `RepositoryBase` for common CRUD operations
- `BaseUseCaseValidation` for common validation
- Shared services in `application/shared/services`

### Error Handling

#### Error Structure
```go
type ApplicationError struct {
    Code    status.ApplicationStatusCode
    Context string
    ErrMsg  string
}
```

#### Error Strategy
1. **Domain Layer**: Business errors (validations)
2. **Application Layer**: Application errors (`ApplicationError`)
3. **Infrastructure Layer**: Technical errors (mapped to `ApplicationError`)

### Logging

#### Logging Levels
- **Info**: General flow information
- **Error**: Errors requiring attention
- **Debug**: Detailed information for debugging
- **Panic**: Critical errors that stop the application

#### Usage Example
```go
providers.Logger.Info("User created successfully", map[string]interface{}{
    "user_id": user.ID,
    "email": user.Email,
})

providers.Logger.Error("Failed to create user", err)
```

### Testing

#### Test Structure
```go
func TestCreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    useCase := NewCreateUserUseCase(logger, mockRepo)

    // Act
    result := useCase.Execute(ctx, locale, input)

    // Assert
    assert.NoError(t, result.Error)
    assert.NotNil(t, result.Data)
}
```

#### Test Types
1. **Unit Tests**: Use cases with mocks
2. **Integration Tests**: Repositories with real DB
3. **E2E Tests**: Complete flows with Bruno

## Conclusion

**GoProjectSkeleton** provides a solid, scalable, and maintainable foundation for developing enterprise applications in Go. The Clean Architecture, along with advanced design patterns, provides an excellent foundation for system growth.

### 🎯 Key Strengths

- ✅ **Solid architecture** and well-structured
- ✅ **Scalability** horizontal and vertical
- ✅ **Serverless ready** - Easy migration to serverless (AWS Lambda, Azure Functions)
- ✅ **Complete testing** in multiple layers
- ✅ **Exhaustive documentation** with Swagger
- ✅ **Complete Docker** for development and production
- ✅ **Multi-cloud** - Support for AWS and Azure
- ✅ **Security** - JWT, OTP, secure password hashing
- ✅ **Internationalization** - Multi-language support
- ✅ **Optimization** - Cache, tree shaking, connection pooling

### 🚀 Ideal Use Cases

- **Enterprise RESTful APIs**
- **Scalable microservices**
- **Serverless applications** (AWS Lambda, Azure Functions)
- **Robust authentication systems**
- **High-concurrency APIs**
- **Projects requiring long-term maintainability**

### 📚 Next Steps

1. **Explore the Documentation**
   - Review Swagger at `http://localhost:8081/docs/` (independent service)
   - Read code examples in each module

2. **Run Tests**
   ```bash
   go test ./src/...
   go test ./tests/integration/...
   ```

3. **Adapt to Your Needs**
   - Customize domain models
   - Add new business modules
   - Configure providers according to your services

4. **Deploy**
   - Development: Docker Compose
   - Production: Traditional monolith or Serverless
   - Cloud: AWS Lambda or Azure Functions

### 🤝 Contributions

This project follows Go best practices and clean architecture. When contributing:

1. Maintain layer separation
2. Follow code conventions
3. Write tests for new features
4. Document important changes
5. Update Swagger for new endpoints

### 📞 Support

For more information:
- **Swagger Documentation**: `/docs/`
- **Example Tests**: `tests/`
- **Source Code**: Explore `src/` to see implementations

---

<div align="center">
  <p>Made with ❤️ using Go and Clean Architecture</p>
  <p>⭐ If this project is useful to you, consider giving it a star</p>
</div>
