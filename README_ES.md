<div align="center">
  <img src="logo.png" alt="GoProjectSkeleton logo" height="400">

  *Una base para proyectos de arquitectura limpia en Go*

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
# 1. Clonar el repositorio
git clone <repository-url>
cd GoProjectSkeleton

# 2. Configurar variables de entorno
cp dev.env.example dev.env
# Editar dev.env con tus configuraciones

# 3. Iniciar servicios con Docker
docker network create goprojectskeleton
docker volume create goprojectskeleton-db-data
docker-compose -f docker/docker-compose.dev.yml up -d

# 4. Ejecutar la aplicación
go run src/infrastructure/server/cmd/main.go

# 5. Acceder a la documentación Swagger (servicio independiente)
# Swagger se ejecuta como un servicio separado en el puerto 8081
# http://localhost:8081/docs/
```

## 📋 Tabla de Contenidos

1. [Introducción](#introducción)
2. [Ambiente de Desarrollo Completo](#-ambiente-de-desarrollo-completo)
3. [Arquitectura del Proyecto](#arquitectura-del-proyecto)
4. [Escalabilidad y Serverless](#escalabilidad-y-serverless)
5. [Flujo Completo de Request](#flujo-completo-de-request)
6. [Ejecución de Tareas en Background](#-ejecución-de-tareas-en-background)
7. [Observabilidad](#-observabilidad)
8. [Virtudes y Beneficios](#virtudes-y-beneficios)
9. [Estructura del Proyecto - Capa por Capa](#estructura-del-proyecto---capa-por-capa)
10. [Revisión Exhaustiva por Carpetas](#revisión-exhaustiva-por-carpetas)
11. [Tecnologías y Dependencias](#tecnologías-y-dependencias)
12. [Configuración y Setup](#configuración-y-setup)
13. [Módulos de Negocio](#módulos-de-negocio)
14. [API y Endpoints](#api-y-endpoints)
15. [Base de Datos y Persistencia](#base-de-datos-y-persistencia)
16. [Autenticación y Seguridad](#autenticación-y-seguridad)
17. [Testing](#testing)
18. [Docker y Despliegue](#docker-y-despliegue)
19. [Despliegue con GitHub Actions](#despliegue-con-github-actions)
20. [Guía de Desarrollo](#guía-de-desarrollo)

---

## Introducción

**GoProjectSkeleton** es un esqueleto de proyecto Go profesional que implementa **Clean Architecture** y **Arquitectura Hexagonal** (Ports & Adapters). Este proyecto proporciona una base sólida, escalable y mantenible para construir aplicaciones empresariales en Go.

### Filosofía del Proyecto

La filosofía central de **GoProjectSkeleton** es que el **dominio** y la **lógica de aplicación** (business logic) deben ser completamente independientes de la **infraestructura** (frameworks, librerías, bases de datos, etc.). Esto permite:

- **Intercambiabilidad**: Cambiar cualquier componente de infraestructura sin afectar el negocio
- **Testabilidad**: Probar la lógica de negocio sin dependencias externas
- **Escalabilidad**: Migrar a arquitecturas serverless o microservicios sin reescribir código
- **Mantenibilidad**: Código limpio y bien estructurado

### Características Principales

#### 🏗️ Arquitectura y Diseño
- ✅ **Clean Architecture** - Separación clara de responsabilidades en capas (Domain, Application, Infrastructure)
- ✅ **Arquitectura Hexagonal** - Desacoplamiento total de dependencias externas mediante Ports & Adapters
- ✅ **Inversión de Dependencias** - Las capas internas no dependen de las externas
- ✅ **SOLID Principles** - Principios de diseño aplicados consistentemente
- ✅ **Repository Pattern** - Abstracción del acceso a datos
- ✅ **Use Case Pattern** - Encapsulación de lógica de negocio
- ✅ **Factory Pattern** - Creación de instancias con inyección de dependencias

#### 🔐 Autenticación y Seguridad
- ✅ **Autenticación JWT Completa** - Access tokens y refresh tokens con configuración flexible
- ✅ **OTP (One-Time Password)** - Autenticación de dos factores con códigos temporales
- ✅ **Sistema de Contraseñas Seguro** - Hash con Bcrypt, reset de contraseñas con tokens
- ✅ **Guards y Autorización** - Control de acceso basado en roles y permisos
- ✅ **Validación Multi-capa** - Validación en DTOs, casos de uso y repositorios
- ✅ **CORS Configurado** - Seguridad para aplicaciones web

#### 👥 Gestión de Usuarios
- ✅ **CRUD Completo** - Crear, leer, actualizar y eliminar usuarios
- ✅ **Gestión de Roles** - Sistema de roles con prioridades
- ✅ **Estados de Usuario** - Pending, Active, Inactive, Suspended, Deleted
- ✅ **Activación de Cuentas** - Sistema de activación mediante tokens
- ✅ **Paginación y Filtrado** - Consultas eficientes con Query Payload
- ✅ **Cache Inteligente** - Cache de listados con Redis

#### 🌐 Internacionalización y Comunicación
- ✅ **Internacionalización (i18n)** - Soporte multiidioma (Español, Inglés, extensible)
- ✅ **Sistema de Emails** - Templates HTML con renderizado dinámico
- ✅ **Emails Transaccionales** - Registro, reset de contraseña, OTP
- ✅ **Locale por Request** - Cada request puede tener su idioma

#### 📚 Documentación y Testing
- ✅ **Servicio Swagger Independiente** - Swagger se ejecuta como un servicio separado y desplegable (puerto 8081)
- ✅ **Listo para Docker Registry** - Swagger puede construirse y desplegarse independientemente
- ✅ **Testing Completo** - Unitarios, integración y E2E
- ✅ **Mocks Completos** - Mocks de repositorios y providers para testing
- ✅ **Bruno Collection** - Colección lista para pruebas E2E

#### 🛠️ Ambiente de Desarrollo
- ✅ **Configuración IDE Completa** - Configuración preestablecida de VS Code/IDE para debugging
- ✅ **Hot Reload con Air** - Recarga automática de código al cambiar archivos
- ✅ **Debugging Remoto** - Debugger Delve configurado para Docker y desarrollo local
- ✅ **Setup Docker de Desarrollo** - Ambiente de desarrollo completo con todas las herramientas preinstaladas
- ✅ **Servicio Swagger Independiente** - Swagger se ejecuta en contenedor separado con hot reload
- ✅ **Herramientas de Desarrollo Incluidas** - Air, Delve, Swag preconfiguradas y listas para usar
- ✅ **Testing de Emails** - Mailpit integrado para pruebas de email en desarrollo
- ✅ **UI de Gestión Redis** - Redis Commander para inspección de cache

#### 🐳 DevOps y Despliegue
- ✅ **Docker Completo** - Multi-servicio para desarrollo, test y E2E
- ✅ **Serverless Ready** - Soporte para AWS Lambda y Azure Functions
- ✅ **Terraform** - Infraestructura como código para AWS y Azure
- ✅ **Secrets Management** - Integración con AWS Secrets Manager y Azure Key Vault
- ✅ **Hot Reload** - Desarrollo eficiente con recarga automática

#### 📊 Observabilidad
- ✅ **OpenTelemetry** - Trazado distribuido e instrumentación de métricas
- ✅ **Prometheus** - Recolección y almacenamiento de métricas
- ✅ **Jaeger** - Visualización de trazas distribuidas
- ✅ **Grafana** - Dashboards y monitoreo en tiempo real
- ✅ **Logging Estructurado** - Logs contextuales con correlación de trazas

#### ⚡ Rendimiento y Escalabilidad
- ✅ **Cache con Redis** - Optimización de rendimiento con TTL configurable
- ✅ **Connection Pooling** - Reutilización de conexiones a base de datos
- ✅ **Sistema de Pipes (DAG)** - Orquestación de casos de uso secuenciales
- ✅ **Paralelización** - Ejecución concurrente de casos de uso con goroutines
- ✅ **Stateless Design** - Listo para escalabilidad horizontal
- ✅ **Tree Shaking** - Optimización automática de binarios en serverless

---

## 🛠️ Ambiente de Desarrollo Completo

**GoProjectSkeleton** viene con un **ambiente de desarrollo completamente configurado** que permite a los desarrolladores comenzar a codificar inmediatamente sin perder tiempo en configuración. Todo está preconfigurado y listo para usar.

### Características del Ambiente de Desarrollo

#### ✅ Configuración IDE Preestablecida

El proyecto incluye archivos de configuración completos del IDE en el directorio `IDE/`:

- **`launch.json`**: Configuraciones de debugging para VS Code
  - **Attach to Docker Go (Delve)**: Configuración de debugging remoto para contenedores Docker
  - **Launch Go Program**: Debugging local con generación automática de Swagger
  - Variables de entorno preconfiguradas para desarrollo
  - Mapeo automático de rutas para debugging en Docker

- **`tasks.json`**: Tareas de compilación
  - Generación automática de documentación Swagger
  - Comandos de compilación preconfigurados

**Uso:**
```bash
# Copiar configuración del IDE a tu workspace
cp -r IDE/.vscode .  # Para VS Code
# O usar la carpeta IDE/ directamente
```

#### ✅ Hot Reload con Air

El proyecto usa **Air** para recarga automática de código durante el desarrollo:

```toml
# Configuración .air.toml
[build]
cmd = "swag init && go build -gcflags 'all=-N -l' -o /app/tmp/main"
bin = "/app/tmp/main"
full_bin = "dlv exec /app/tmp/main --headless --listen=:40000"
include_ext = ["go", "tpl", "tmpl", "html"]
```

**Características:**
- ✅ Recompilación automática al cambiar archivos
- ✅ Integrado con debugger Delve
- ✅ Observa archivos Go, templates y HTML
- ✅ Logs con colores y timestamps
- ✅ Modo polling para mejor compatibilidad con sistemas de archivos

**Cómo funciona:**
1. Air observa cambios en archivos
2. Regenera automáticamente la documentación Swagger
3. Recompila la aplicación con símbolos de debug (flags `-N -l`)
4. Reinicia la aplicación con debugger Delve adjunto
5. Listo para debugging remoto en el puerto 40000

#### ✅ Debugging Remoto con Delve

**Delve (dlv)** está preconfigurado para debugging local y remoto:

**Debugging Remoto en Docker:**
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

**Características:**
- ✅ Depurar aplicación ejecutándose en Docker
- ✅ Establecer breakpoints en VS Code
- ✅ Inspeccionar variables y call stack
- ✅ Ejecutar paso a paso
- ✅ Mapeo de rutas para debugging sin problemas

**Flujo de Debugging:**
1. Iniciar ambiente de desarrollo Docker: `docker-compose -f docker/docker-compose.dev.yml up -d`
2. La aplicación inicia con Delve en el puerto 40000
3. Adjuntar debugger desde VS Code usando "Attach to Docker Go (Delve)"
4. Establecer breakpoints y depurar como si fuera local

#### ✅ Setup Docker de Desarrollo Completo

El `docker-compose.dev.yml` incluye todo lo necesario para desarrollo:

**Servicios:**
- **Aplicación**: Aplicación Go con hot reload y debugging
- **PostgreSQL**: Base de datos en puerto 5436
- **Redis**: Servidor de cache
- **Mailpit**: Testing de emails (puerto 8025 para UI, 1025 para SMTP)
- **Redis Commander**: UI de gestión de Redis (puerto 8081)

**Características de Desarrollo:**
```yaml
services:
  goprojectskeleton:
    ports:
      - "8080:8080"    # Aplicación
      - "40000:40000"  # Debugger Delve
    volumes:
      - ../src:/app/src              # Sincronización de código en vivo
      - ../dev.env:/app/.env:ro       # Variables de entorno
      - ../.air.toml:/app/.air.toml   # Configuración de Air
    command: ["air", "-c", ".air.toml"]  # Hot reload
```

**Herramientas Preinstaladas en Docker:**
- ✅ **Air** (v1.61.7) - Hot reload
- ✅ **Delve** (dlv) - Debugger
- ✅ **Swag** - Generador de Swagger
- ✅ Todas las dependencias de Go

#### ✅ Servicio de Documentación Swagger Independiente

La documentación Swagger es ahora un **servicio completamente independiente** que se ejecuta por separado de la aplicación principal:

**Arquitectura:**
- **Servidor Independiente**: Swagger se ejecuta en su propio servidor HTTP (puerto 8081)
- **Despliegue Separado**: Puede desplegarse independientemente usando `dockerfile.swagger`
- **Listo para Docker Registry**: Puede subirse a cualquier Docker registry y desplegarse por separado
- **Sin Dependencias**: Completamente desacoplado de la aplicación principal

**Desarrollo:**
- El servicio Swagger se ejecuta en un contenedor Docker separado
- Regenera automáticamente la documentación al cambiar el código
- Disponible en `http://localhost:8081/docs/`
- Configurado mediante variables de entorno (ver `dev.env.example`)

**Despliegue en Producción:**
- Construir: `docker build -f docker/dockerfile.swagger -t your-registry/swagger:latest .`
- Subir: `docker push your-registry/swagger:latest`
- Ejecutar: `docker run -p 8081:8081 -e API_HOST=your-api-host your-registry/swagger:latest`

**Configuración:**
- `SWAGGER_PORT`: Puerto para el servidor Swagger (por defecto: 8081)
- `API_HOST`: Host de la API principal (ej: `api.example.com:8080`)
- `API_TITLE`: Título de la API en Swagger UI
- `API_VERSION`: Versión de la API
- `API_DESCRIPTION`: Descripción de la API
- `API_BASE_PATH`: Ruta base para los endpoints de la API

#### ✅ Integración de Herramientas de Desarrollo

**Testing de Emails con Mailpit:**
- Web UI: `http://localhost:8025`
- SMTP: `localhost:1025`
- Probar todos los emails transaccionales sin servidor SMTP real

**Gestión de Redis:**
- Redis Commander: `http://localhost:8081`
- Inspección visual de cache
- Navegación y edición de claves

**Acceso a Base de Datos:**
- PostgreSQL: `localhost:5436`
- Acceso directo a base de datos para debugging
- Datos persistentes con volúmenes Docker

### Inicio Rápido de Desarrollo

**1. Iniciar Ambiente de Desarrollo:**
```bash
# Crear red y volumen
docker network create goprojectskeleton
docker volume create goprojectskeleton-db-data

# Iniciar todos los servicios
docker-compose -f docker/docker-compose.dev.yml up -d
```

**2. Configurar IDE:**
```bash
# Copiar configuración del IDE (si usas VS Code)
cp -r IDE/.vscode .
```

**3. Iniciar Debugging:**
- Abrir VS Code
- Ir a Run and Debug (F5)
- Seleccionar "Attach to Docker Go (Delve)"
- ¡Establecer breakpoints y comenzar a depurar!

**4. Flujo de Desarrollo:**
- Editar código en directorio `src/`
- Air detecta automáticamente los cambios
- La aplicación se recompila y reinicia
- El debugger se readjunta automáticamente
- Ver cambios inmediatamente

### Arquitectura del Ambiente de Desarrollo

```mermaid
graph TB
    subgraph Developer["👨‍💻 Desarrollador"]
        IDE[VS Code/IDE<br/>Breakpoints & Debugging]
        Code[Código Fuente<br/>src/]
    end

    subgraph Docker["🐳 Contenedor Docker"]
        Air[Air<br/>Hot Reload Watcher]
        Delve[Delve Debugger<br/>Puerto 40000]
        App[Aplicación Go<br/>Puerto 8080]
    end

    subgraph Services["🔧 Servicios de Desarrollo"]
        DB[(PostgreSQL<br/>Puerto 5436)]
        Redis[(Redis)]
        Mailpit[Mailpit<br/>Puerto 8025/1025]
        RedisUI[Redis Commander<br/>Puerto 8081]
    end

    Code -->|Cambios de Archivos| Air
    Air -->|Recompilar & Reiniciar| Delve
    Delve -->|Sesión de Debug| App
    IDE -->|Adjuntar Debugger| Delve
    App --> DB
    App --> Redis
    App -->|Enviar Emails| Mailpit

    style Air fill:#e3f2fd
    style Delve fill:#fff9c4
    style App fill:#c8e6c9
```

### Archivos de Configuración de Desarrollo

| Archivo | Propósito |
|---------|----------|
| `IDE/launch.json` | Configuraciones de debugging de VS Code |
| `IDE/tasks.json` | Tareas de compilación (generación de Swagger) |
| `.air.toml` | Configuración de hot reload de Air |
| `docker/dockerfile.dev` | Imagen Docker de desarrollo |
| `docker/docker-compose.dev.yml` | Stack completo de desarrollo |
| `dev.env.example` | Plantilla de variables de entorno de desarrollo |

### Beneficios para Desarrolladores

1. **Tiempo de Configuración Cero**
   - Todo preconfigurado
   - Comenzar a codificar inmediatamente
   - No se necesita instalación manual de herramientas

2. **Debugging Productivo**
   - Debugging remoto en Docker
   - Breakpoints funcionan sin problemas
   - Inspección de variables y call stack

3. **Ciclo de Desarrollo Rápido**
   - Hot reload al cambiar archivos
   - No se necesitan reinicios manuales
   - Ver cambios instantáneamente

4. **Herramientas Completas**
   - Testing de emails sin SMTP
   - UI de inspección de Redis
   - Acceso a base de datos listo

5. **Consistencia del Equipo**
   - Mismo ambiente para todos los desarrolladores
   - No hay problemas de "funciona en mi máquina"
   - Docker asegura consistencia

---

## ☁️ Capacidades Cloud y Serverless

**GoProjectSkeleton** está diseñado para funcionar tanto en entornos tradicionales como en arquitecturas serverless modernas.

### 🚀 AWS Lambda

El proyecto incluye soporte completo para **AWS Lambda** con:

- ✅ **Generación Automática de Funciones** - Sistema de generación desde `functions.json`
- ✅ **Módulos Independientes** - Cada función Lambda tiene su propio módulo Go
- ✅ **Tree Shaking Optimizado** - Binarios de 5-15 MB vs 50+ MB sin optimización
- ✅ **AWS Secrets Manager** - Carga automática de secretos desde Secrets Manager
- ✅ **Lambda Adapter** - Adaptador para eventos de API Gateway
- ✅ **Terraform** - Infraestructura como código lista para desplegar
- ✅ **Compilación Optimizada** - Flags específicos para Lambda (`lambda.norpc`, `-ldflags="-s -w"`)

**Estructura de funciones Lambda:**
```
src/infrastructure/clouds/aws/
├── functions.json          # Definición de funciones
├── init.go                 # Inicialización AWS
├── lambda_adapter.go       # Adaptador Lambda
├── secrets_manager.go      # Gestión de secretos
├── terraform/              # Infraestructura Terraform
└── functions/              # Generador de funciones
```

### 🔷 Azure Functions

Soporte completo para **Azure Functions** con:

- ✅ **HTTP Adapter** - Adaptador para Azure Functions HTTP triggers
- ✅ **Azure Key Vault** - Integración con Azure Key Vault para secretos
- ✅ **Terraform** - Infraestructura como código para Azure
- ✅ **Módulos Independientes** - Cada función tiene su propio módulo
- ✅ **Generación Automática** - Sistema de generación desde `functions.json`

**Estructura de funciones Azure:**
```
src/infrastructure/clouds/azure/
├── functions.json          # Definición de funciones
├── init.go                 # Inicialización Azure
├── http_adapter.go         # Adaptador HTTP
├── vault.go                # Integración Key Vault
├── terraform/              # Infraestructura Terraform
└── functions/              # Generador de funciones
```

### 📊 Comparación de Arquitecturas

| Característica | Monolito Tradicional | AWS Lambda | Azure Functions |
|---------------|---------------------|------------|-----------------|
| **Inicialización** | Una vez al inicio | Por función | Por función |
| **Escalabilidad** | Manual/Horizontal | Automática | Automática |
| **Costo** | Fijo | Por uso | Por uso |
| **Cold Start** | N/A | ~100-500ms | ~200-800ms |
| **Tamaño Binario** | ~50 MB | ~5-15 MB | ~5-15 MB |
| **Gestión Secretos** | Variables de entorno | Secrets Manager | Key Vault |
| **Despliegue** | Docker/VM | ZIP a Lambda | ZIP a Functions |

---

## Arquitectura del Proyecto

### Visión General de la Arquitectura

El proyecto implementa **Clean Architecture** con tres capas principales:

```mermaid
graph TB
    subgraph Infrastructure["🔧 Capa de Infraestructura"]
        API["API Layer<br/>(Gin HTTP)"]
        DB["Database<br/>(GORM/Redis)"]
        Providers["Providers<br/>(JWT/Email)"]
        Handlers["Handlers<br/>(HTTP)"]
        Repos["Repositories<br/>(GORM)"]
        Middlewares["Middlewares<br/>(Auth/CORS)"]
    end

    subgraph Application["💼 Capa de Aplicación"]
        UseCases["Use Cases<br/>(Business Logic)"]
        Services["Services<br/>(Email/OTP)"]
        Contracts["Contracts<br/>(Interfaces)"]
        Pipes["Pipes<br/>(DAG)"]
        DTOs["DTOs<br/>(Transfer)"]
        Errors["Errors<br/>(Handling)"]
    end

    subgraph Domain["🎯 Capa de Dominio"]
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

### Diagrama de Arquitectura Hexagonal

```mermaid
graph LR
    subgraph External["🌐 Mundo Externo"]
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

### Principios de Diseño

#### 1. Inversión de Dependencias (Dependency Inversion)

Las capas internas (Domain y Application) **nunca** dependen de las capas externas (Infrastructure). En su lugar, definen **interfaces (contratos)** que la infraestructura implementa.

**Ejemplo:**
```go
// Domain/Application define la interfaz
type IUserRepository interface {
    Create(input UserCreate) (*User, error)
    GetByID(id uint) (*User, error)
}

// Infrastructure implementa la interfaz
type UserRepository struct {
    DB *gorm.DB
}

func (r *UserRepository) Create(input UserCreate) (*User, error) {
    // Implementación con GORM
}
```

#### 2. Separación de Responsabilidades (Separation of Concerns)

Cada capa tiene una responsabilidad única y bien definida:

- **Domain**: Entidades de negocio puras, sin dependencias
- **Application**: Lógica de negocio y casos de uso
- **Infrastructure**: Implementaciones técnicas (HTTP, BD, etc.)

#### 3. Testabilidad

Gracias a las interfaces, podemos testear la lógica de negocio sin necesidad de bases de datos o servicios reales:

```go
// En tests, usamos mocks
mockRepo := &MockUserRepository{}
useCase := NewCreateUserUseCase(logger, mockRepo)
```

#### 4. Extensibilidad

Agregar nuevas funcionalidades es simple y no afecta código existente:

1. Crear modelo en `domain/models/`
2. Definir interfaz en `application/contracts/`
3. Implementar caso de uso en `application/modules/`
4. Crear repositorio en `infrastructure/repositories/`
5. Agregar handler en `infrastructure/handlers/`

### Patrones de Diseño Implementados

#### 1. Repository Pattern

Abstrae el acceso a datos, permitiendo cambiar de base de datos sin afectar la lógica de negocio.

```go
type IUserRepository interface {
    IRepositoryBase[UserCreate, UserUpdate, User, UserInDB]
    CreateWithPassword(input UserAndPasswordCreate) (*User, error)
    GetUserWithRole(id uint) (*UserWithRole, error)
    GetByEmailOrPhone(emailOrPhone string) (*User, error)
}
```

#### 2. Use Case Pattern

Encapsula la lógica de negocio en casos de uso reutilizables y testables.

```go
type CreateUserUseCase struct {
    log  ILoggerProvider
    repo IUserRepository
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, locale LocaleTypeEnum, input UserCreate) *UseCaseResult[User] {
    // Lógica de negocio aquí
}
```

#### 3. DAG (Directed Acyclic Graph) Pattern

Orquesta múltiples casos de uso en secuencia o paralelo.

```go
// Ejecución secuencial
dag := NewDag(NewStep(createUserUC), locale, ctx)
dag = Then(dag, NewStep(sendEmailUC))
result := dag.Execute(input)

// Ejecución paralela
parallelDag := NewUseCaseParallelDag[Input, Output]()
parallelDag.Usecases = []BaseUseCase{uc1, uc2, uc3}
result := parallelDag.Execute(ctx, locale, input)
```

#### 4. Factory Pattern

Crea instancias de casos de uso con dependencias inyectadas.

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

Permite intercambiar implementaciones mediante interfaces (providers).

```go
// Interfaz
type IHashProvider interface {
    Hash(password string) (string, error)
    Compare(hashed, plain string) bool
}

// Implementaciones intercambiables
type BcryptHashProvider struct {}
type Argon2HashProvider struct {}
```

---

## Escalabilidad y Serverless

### Capacidades de Escalabilidad

**GoProjectSkeleton** está diseñado para escalar tanto **horizontalmente** como **verticalmente**, y puede migrarse fácilmente a arquitecturas **serverless** o **monolito serverless**.

### Escalabilidad Horizontal

#### Características que Facilitan la Escalabilidad Horizontal

1. **Stateless Design**
   - La aplicación no mantiene estado en memoria
   - Cada request es independiente
   - Perfecto para load balancers

2. **Cache Distribuido (Redis)**
   - Cache compartido entre instancias
   - No hay dependencias de sesión local

3. **Base de Datos Externa**
   - PostgreSQL independiente de la aplicación
   - Múltiples instancias pueden conectarse

4. **Sin Estado de Sesión**
   - Autenticación basada en JWT (stateless)
   - No requiere sticky sessions

#### Implementación para Escalabilidad Horizontal

```go
// Cada handler es stateless
func CreateUser(ctx HandlerContext) {
    // No hay estado compartido
    // Cada request es independiente
    // Puede ejecutarse en cualquier instancia
}
```

### Escalabilidad Vertical

#### Optimizaciones Implementadas

1. **Goroutines Nativas**
   - Concurrencia eficiente de Go
   - Múltiples requests procesados simultáneamente

2. **Connection Pooling**
   - Reutilización de conexiones a BD
   - Configuración optimizada de GORM

3. **Cache Inteligente**
   - Reduce consultas a base de datos
   - TTL configurable

4. **Query Optimization**
   - Filtros y ordenamiento eficientes
   - Paginación para control de memoria

### Migración a Monolito Serverless

**GoProjectSkeleton** puede migrarse fácilmente a una arquitectura **monolito serverless** (como AWS Lambda, Google Cloud Functions, Azure Functions) gracias a:

#### 1. Arquitectura Desacoplada

```go
// La lógica de negocio no depende de Gin
type CreateUserUseCase struct {
    // No hay referencias a HTTP
    // Solo interfaces
}
```

#### 2. Inicialización Modular

```go
// container.go - Inicialización separada
func Initialize() {
    // Configuración
    // Base de datos
    // Providers
    // Servicios
}
```

#### 3. Handlers Independientes

Los handlers pueden adaptarse fácilmente a diferentes frameworks o entornos serverless:

```go
// Handler actual (Gin)
func CreateUser(ctx HandlerContext) { ... }

// Adaptación para Lambda
func CreateUserLambda(ctx context.Context, event APIGatewayEvent) (Response, error) {
    // Misma lógica, diferente wrapper
}
```

#### 4. Sin Estado Global

- No hay variables globales de estado
- Todo se pasa por contexto o inyección de dependencias

#### Pasos para Migración Serverless

1. **Extraer Lógica de Negocio**
   ```go
   // Ya está hecho - los Use Cases son independientes
   ```

2. **Crear Adapter para Serverless**
   ```go
   // lambda/handlers/user.go
   func CreateUserHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
       // Inicializar infraestructura
       infrastructure.Initialize()

       // Adaptar evento a DTO
       var userCreate dtos.UserCreate
       json.Unmarshal([]byte(event.Body), &userCreate)

       // Ejecutar caso de uso (mismo código)
       ucResult := usecases_user.NewCreateUserUseCase(...).Execute(...)

       // Adaptar respuesta
       return adaptResponse(ucResult), nil
   }
   ```

3. **Configurar Variables de Entorno**
   - Usar el mismo sistema de configuración
   - Adaptar a variables de entorno del proveedor cloud

4. **Mantener Misma Lógica**
   - Los Use Cases no cambian
   - Los Repositories no cambian
   - Solo cambia la capa de infraestructura HTTP

### Escalabilidad en Producción

#### Diagrama de Escalabilidad Horizontal

```mermaid
graph TB
    subgraph Clients["Clientes"]
        Web[Web App]
        Mobile[Mobile App]
        API_Client[API Clients]
    end

    subgraph LB["Load Balancer"]
        Nginx[Nginx/HAProxy<br/>Health Checks]
    end

    subgraph AppInstances["Instancias de Aplicación"]
        App1[App Instance 1<br/>Stateless]
        App2[App Instance 2<br/>Stateless]
        App3[App Instance N<br/>Stateless]
    end

    subgraph DataLayer["Capa de Datos"]
        RedisCluster[(Redis Cluster<br/>Cache Distribuido)]
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

#### Configuración Recomendada

1. **Load Balancer**
   - Distribuir tráfico entre múltiples instancias
   - Health checks en `/api/health-check`

2. **Base de Datos**
   - Read replicas para consultas
   - Connection pooling optimizado
   - Índices en columnas frecuentes

3. **Cache**
   - Redis cluster para alta disponibilidad
   - Estrategias de cache (cache-aside, write-through)

4. **Monitoring**
   - Métricas de rendimiento
   - Logging estructurado
   - Alertas de errores

---

## Despliegue e Inicialización AWS Serverless

### Arquitectura de Módulos Go para Serverless

**GoProjectSkeleton** utiliza una arquitectura de **múltiples módulos Go** para optimizar el despliegue en AWS Lambda. Cada función serverless tiene su propio módulo independiente que solo incluye las dependencias necesarias para ejecutarse.

#### Estructura de Módulos

```mermaid
graph TB
    subgraph RootModule["Módulo Raíz: goprojectskeleton"]
        Domain[Domain Models]
        Application[Application Layer]
        Infrastructure[Infrastructure Core]
    end

    subgraph AWSModule["Módulo AWS: goprojectskeleton/aws"]
        AWSInit[AWS Init]
        LambdaAdapter[Lambda Adapter]
        SecretsManager[Secrets Manager]
    end

    subgraph FunctionModules["Módulos de Funciones"]
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

#### Resolución de Dependencias

Cada función Lambda tiene su propio `go.mod` que utiliza **replace directives** para apuntar a los módulos locales:

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

**Ventajas de esta arquitectura:**

1. **Optimización de Binarios**: Go solo compila el código que realmente se usa
2. **Dependencias Mínimas**: Cada función solo incluye lo necesario
3. **Tree Shaking Automático**: Go elimina código no utilizado en tiempo de compilación
4. **Módulos Independientes**: Cada función puede evolucionar independientemente

### Proceso de Generación y Compilación

#### 1. Generación de Funciones

Las funciones Lambda se generan automáticamente desde `functions.json`:

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

**Proceso de generación:**

```mermaid
sequenceDiagram
    participant Dev as Desarrollador
    participant Gen as Generator
    participant Template as Templates
    participant Func as Función Lambda

    Dev->>Gen: GenerateFunctions(functions.json)
    Gen->>Gen: Lee functions.json
    loop Para cada función
        Gen->>Template: Carga main.go.tmpl
        Gen->>Template: Carga go.mod.tmpl
        Gen->>Func: Genera main.go
        Gen->>Func: Genera go.mod
        Gen->>Func: Crea estructura de directorios
    end
    Gen-->>Dev: ✅ Funciones generadas
```

**Estructura generada:**

```
tmp/
└── status/
    └── health_check/
        ├── main.go          # Handler Lambda generado
        ├── go.mod           # Módulo independiente
        └── bin/             # Directorio de compilación
            ├── bootstrap    # Binario compilado
            └── src/         # Templates copiados
```

#### 2. Compilación Optimizada

Cada función se compila con optimizaciones específicas para Lambda:

```bash
# Compilación desde el directorio de la función
cd tmp/status/health_check
go build -o bin/bootstrap \
    -tags lambda.norpc \
    -ldflags="-s -w" \
    main.go

# Variables de entorno de compilación
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0
```

**Optimizaciones aplicadas:**

- **`-tags lambda.norpc`**: Desactiva RPC de Lambda (reduce tamaño)
- **`-ldflags="-s -w"`**: Elimina símbolos de debug (reduce tamaño)
- **`CGO_ENABLED=0`**: Compilación estática (sin dependencias C)
- **`GOOS=linux`**: Binario para Linux (entorno Lambda)
- **`GOARCH=amd64`**: Arquitectura x86_64

#### 3. Tree Shaking y Eliminación de Código Muerto

Go realiza **tree shaking automático** durante la compilación:

```mermaid
graph LR
    subgraph SourceCode["Código Fuente"]
        Used[✅ Código Usado<br/>handlers.GetHealthCheck<br/>aws.InitializeInfrastructure]
        Unused[❌ Código No Usado<br/>handlers.CreateUser<br/>handlers.UpdateUser<br/>handlers.DeleteUser]
    end

    subgraph Compilation["Compilación Go"]
        Analyzer[Go Compiler<br/>Análisis Estático]
    end

    subgraph Binary["Binario Final"]
        Included[✅ Solo Código Usado<br/>~5-10 MB]
    end

    SourceCode --> Analyzer
    Analyzer -->|Tree Shaking| Binary

    style Used fill:#c8e6c9
    style Unused fill:#ffcdd2
    style Included fill:#c8e6c9
```

**Ejemplo real:**

Para la función `health-check`, el binario final solo incluye:

- ✅ `handlers.GetHealthCheck` (handler específico)
- ✅ `aws.InitializeInfrastructure` (inicialización)
- ✅ `aws.HandleLambdaEvent` (adaptador Lambda)
- ✅ Dependencias transitivas necesarias (GORM, Redis, JWT, etc.)
- ❌ **NO incluye**: Otros handlers (`CreateUser`, `UpdateUser`, etc.)
- ❌ **NO incluye**: Código de módulos no utilizados

**Resultado:** Binarios de 5-15 MB en lugar de 50+ MB si se incluyera todo.

### Inicialización de Infraestructura

#### Proceso de Inicialización

Cada función Lambda inicializa su infraestructura en el `init()`:

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

#### Flujo de Inicialización

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

#### Componentes Inicializados

1. **Configuración (Settings)**
   ```go
   // Carga desde AWS Secrets Manager o variables de entorno
   settings.AppSettingsInstance.Initialize(config.ToMap())
   ```

2. **Base de Datos (GORM)**
   ```go
   database.GoProjectSkeletondb.SetUp(
       host, port, user, password, dbname, ssl, logger
   )
   ```

3. **Providers**
   - **JWT Provider**: Generación y validación de tokens
   - **Email Provider**: Envío de emails (SMTP)
   - **Cache Provider**: Redis para cache
   - **Logger Provider**: Sistema de logging

4. **Servicios**
   - **Email Services**: Servicios de email (registro, reset, OTP)

#### Carga de Configuración desde AWS Secrets Manager

El sistema puede cargar configuración desde **AWS Secrets Manager**:

```go
// Si una variable de entorno es un ARN de Secrets Manager
DB_PASSWORD=arn:aws:secretsmanager:us-east-1:123456789:secret:db-password

// El sistema automáticamente:
// 1. Detecta que es un ARN
// 2. Obtiene el secreto de Secrets Manager
// 3. Usa el valor en la configuración
```

**Ventajas:**

- ✅ **Seguridad**: Secretos no en código o variables de entorno
- ✅ **Rotación**: Secrets Manager puede rotar secretos automáticamente
- ✅ **Auditoría**: Todas las accesos a secretos son auditados
- ✅ **Fallback**: Si falla, usa valores por defecto

### Proceso de Despliegue

#### Flujo Completo de Despliegue

```mermaid
graph TB
    subgraph Local["🖥️ Desarrollo Local"]
        Dev[Desarrollador]
        FunctionsJSON[functions.json]
    end

    subgraph Generation["📦 Generación"]
        Generator[Generator Tool]
        Templates[Templates]
        FunctionCode[Código de Funciones]
    end

    subgraph Build["🔨 Compilación"]
        GoBuild[go build]
        Binary[Binario bootstrap]
        TemplatesCopy[Templates copiados]
        Zip[ZIP Package]
    end

    subgraph AWS["☁️ AWS"]
        Lambda[AWS Lambda]
        APIGateway[API Gateway]
    end

    Dev-->FunctionsJSON
    FunctionsJSON-->|GenerateFunctions| Generator
    Generator-->|Carga templates| Templates
    Templates-->|Genera código| FunctionCode
    FunctionCode-->|go build| GoBuild
    GoBuild-->|Binario optimizado| Binary
    Binary-->|Copia templates| TemplatesCopy
    TemplatesCopy-->|Crea ZIP| Zip
    Zip-->|Deploy function| Lambda
    Lambda-->|Conecta endpoint| APIGateway

    style Dev fill:#e3f2fd
    style Lambda fill:#ff9800
    style APIGateway fill:#ff9800
```

#### Comandos de Despliegue

**1. Generar funciones:**
```bash
cd src/infrastructure/clouds/aws/functions
go run main.go generate
```

**2. Compilar y desplegar:**
```bash
# Desplegar todas las funciones
go run main.go deploy

# Desplegar una función específica
go run main.go deploy --function health-check
```

**3. Proceso interno de despliegue:**

```bash
# Para cada función:
# 1. go mod tidy (resuelve dependencias)
cd tmp/status/health_check
go mod tidy

# 2. Compilar binario
go build -o bin/bootstrap -tags lambda.norpc main.go

# 3. Copiar templates necesarios
cp -r ../../../../../../application/shared/templates bin/src/application/shared/templates

# 4. Crear ZIP
cd bin
zip -r health-check.zip bootstrap src/

# 5. Desplegar a Lambda
aws lambda update-function-code \
    --function-name goprojectskeleton-dev-healthcheck \
    --zip-file fileb://health-check.zip
```

### Optimizaciones de Tamaño

#### Comparación de Tamaños

| Función | Tamaño con Todo | Tamaño Optimizado | Reducción |
|---------|----------------|-------------------|----------|
| health-check | ~50 MB | ~8 MB | 84% |
| auth-login | ~55 MB | ~12 MB | 78% |
| user-get | ~60 MB | ~15 MB | 75% |

#### Estrategias de Optimización

1. **Tree Shaking de Go**
   - Elimina código no utilizado automáticamente
   - Solo incluye funciones y tipos referenciados

2. **Módulos Independientes**
   - Cada función tiene su propio `go.mod`
   - Dependencias resueltas por función

3. **Compilación Estática**
   - `CGO_ENABLED=0`: Sin dependencias C
   - Binario autocontenido

4. **Eliminación de Debug**
   - `-ldflags="-s -w"`: Elimina símbolos
   - Reduce tamaño significativamente

5. **Templates Selectivos**
   - Solo se copian templates necesarios
   - No se incluyen todos los templates

### Estructura del Paquete Lambda

```
health-check.zip
├── bootstrap                    # Binario Go compilado (~8 MB)
└── src/
    └── application/
        └── shared/
            └── templates/
                └── emails/      # Solo templates necesarios
                    ├── otp_en.gohtml
                    └── otp_es.gohtml
```

**Tamaño total:** ~8-10 MB (vs ~50 MB sin optimización)

### Ventajas de la Arquitectura

#### 1. **Despliegues Rápidos**
- Binarios pequeños = uploads rápidos
- Cold start más rápido
- Menor costo de almacenamiento

#### 2. **Seguridad Mejorada**
- Cada función es independiente
- Menor superficie de ataque
- Secretos en Secrets Manager

#### 3. **Escalabilidad**
- Cada función escala independientemente
- Configuración por función
- Optimización individual

#### 4. **Mantenibilidad**
- Código compartido en módulos
- Cambios localizados
- Testing independiente

#### 5. **Costo Optimizado**
- Menor tamaño = menor costo de almacenamiento
- Cold start más rápido = menor latencia
- Menor uso de memoria

### Diagrama de Arquitectura Completa

```mermaid
graph TB
    subgraph Client["🌐 Cliente"]
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
            Secrets[Secrets Manager<br/>Configuración]
            RDS[(RDS PostgreSQL<br/>Base de Datos)]
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

## Flujo Completo de Request

### Diagrama de Flujo de Request

```mermaid
sequenceDiagram
    participant Client as Cliente HTTP
    participant Gin as Gin Router
    participant MW as Middleware<br/>(CORS/Recovery)
    participant AuthMW as Auth Middleware<br/>(si protegido)
    participant Router as Router/Wrapper
    participant Handler as HTTP Handler
    participant UseCase as Use Case
    participant Repo as Repository
    participant DB as PostgreSQL
    participant Resolver as Request Resolver

    Client->>Gin: HTTP Request
    Gin->>MW: Aplica middlewares
    MW->>AuthMW: Si ruta protegida
    AuthMW->>AuthMW: Valida JWT Token
    AuthMW->>AuthMW: Obtiene Usuario
    AuthMW->>Router: Request con contexto
    Router->>Router: Extrae params, query, locale
    Router->>Router: Crea HandlerContext
    Router->>Handler: HandlerContext
    Handler->>Handler: Decodifica JSON → DTO
    Handler->>UseCase: Ejecuta Use Case
    UseCase->>UseCase: Valida input (DTO.Validate())
    UseCase->>UseCase: Ejecuta Guards (si aplica)
    UseCase->>Repo: Llama Repository
    Repo->>Repo: Convierte DTO → Modelo BD
    Repo->>DB: Ejecuta Query (GORM)
    DB-->>Repo: Retorna datos
    Repo->>Repo: Convierte Modelo BD → Dominio
    Repo->>Repo: Mapea errores
    Repo-->>UseCase: Modelo de Dominio / Error
    UseCase->>UseCase: Maneja errores
    UseCase-->>Handler: UseCaseResult
    Handler->>Resolver: Resuelve respuesta
    Resolver->>Resolver: Mapea status code
    Resolver->>Resolver: Formatea JSON
    Resolver-->>Handler: HTTP Response
    Handler-->>Client: JSON Response
```

### Diagrama de Flujo Simplificado

```mermaid
flowchart TD
    Start([Cliente HTTP]) --> Gin[Gin Router]
    Gin --> MW[Middlewares<br/>CORS/Recovery]
    MW --> Auth{¿Ruta<br/>Protegida?}
    Auth -->|Sí| AuthMW[Auth Middleware<br/>Valida JWT]
    Auth -->|No| Router
    AuthMW --> Router[Router/Wrapper<br/>Extrae params/query/locale]
    Router --> Handler[HTTP Handler<br/>Decodifica JSON → DTO]
    Handler --> UseCase[Use Case<br/>Valida y ejecuta lógica]
    UseCase --> Repo[Repository<br/>Convierte DTO ↔ Modelo]
    Repo --> DB[(PostgreSQL<br/>GORM)]
    DB --> Repo
    Repo --> UseCase
    UseCase --> Resolver[Request Resolver<br/>Formatea respuesta]
    Resolver --> End([Respuesta JSON])

    style Start fill:#e3f2fd
    style End fill:#c8e6c9
    style DB fill:#fff9c4
    style UseCase fill:#f3e5f5
```

### Flujo Detallado con Código

#### 1. Entrada HTTP (main.go)

```go
// main.go
func main() {
    infrastructure.Initialize()  // Inicializa BD, providers, etc.
    app := buildGinApp()         // Crea aplicación Gin
    loadGinApp(app)               // Configura middlewares y rutas
    app.Run("0.0.0.0:8080")      // Inicia servidor
}
```

#### 2. Ruteo (routes/router.go)

```go
// routes/router.go
func Router(r *gin.RouterGroup) {
    r.POST("/user", wrapHandler(handlers.CreateUser))
    // wrapHandler adapta Gin a HandlerContext
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
    // 1. Decodificar JSON
    var userCreate dtos.UserCreate
    json.NewDecoder(*ctx.Body).Decode(&userCreate)

    // 2. Crear repositorio
    repo := repositories.NewUserRepository(database.DB, providers.Logger)

    // 3. Crear y ejecutar caso de uso
    ucResult := usecases_user.NewCreateUserUseCase(
        providers.Logger,
        repo,
    ).Execute(ctx.c, ctx.Locale, userCreate)

    // 4. Resolver respuesta
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

#### 5. Caso de Uso (application/modules/user/use_cases/create_user.go)

```go
// application/modules/user/use_cases/create_user.go
func (uc *CreateUserUseCase) Execute(
    ctx context.Context,
    locale locales.LocaleTypeEnum,
    input dtos.UserCreate,
) *usecase.UseCaseResult[models.User] {
    result := usecase.NewUseCaseResult[models.User]()

    // 1. Validar input
    uc.validate(input, result)
    if result.HasError() {
        return result
    }

    // 2. Llamar a repositorio
    res, err := uc.repo.Create(input)
    if err != nil {
        result.SetError(err.Code, err.Context)
        return result
    }

    // 3. Retornar éxito
    result.SetData(status.Created, *res, "User created")
    return result
}
```

#### 6. Repositorio (infrastructure/repositories/user.go)

```go
// infrastructure/repositories/user.go
func (ur *UserRepository) Create(input dtos.UserCreate) (*models.User, *application_errors.ApplicationError) {
    // 1. Convertir DTO a modelo GORM
    userCreate := ur.modelConverter.ToGormCreate(input)

    // 2. Ejecutar query
    if err := ur.DB.Create(userCreate).Error; err != nil {
        return nil, MapOrmError(err)  // Mapea errores de BD
    }

    // 3. Convertir modelo GORM a modelo de dominio
    userModel := ur.modelConverter.ToDomain(userCreate)
    return userModel, nil
}
```

#### 7. Resolución de Respuesta (handlers/request_resolver.go)

```go
// handlers/request_resolver.go
func (rr *RequestResolver[D]) ResolveDTO(
    w http.ResponseWriter,
    result *usecase.UseCaseResult[D],
    headersToAdd map[HTTPHeaderTypeEnum]string,
) {
    // 1. Agregar headers
    rr.setHeaders(w, headersToAdd)

    // 2. Manejar errores
    if result.HasError() {
        w.WriteHeader(rr.statusMapping[result.StatusCode])
        json.NewEncoder(w).Encode(map[string]any{
            "details": result.Error,
        })
        return
    }

    // 3. Respuesta exitosa
    w.WriteHeader(rr.statusMapping[result.StatusCode])
    json.NewEncoder(w).Encode(map[string]any{
        "data":    result.Data,
        "details": result.Details,
    })
}
```

### Flujo con Pipes (DAG)

Para casos más complejos que requieren múltiples pasos:

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
    DB-->>Repo: User creado
    Repo-->>UC1: User
    UC1-->>DAG: UseCaseResult[User]

    alt Si no hay error
        DAG->>UC2: Execute(User)
        UC2->>EmailSvc: SendWelcomeEmail()
        EmailSvc->>SMTP: Enviar email
        SMTP-->>EmailSvc: Email enviado
        EmailSvc-->>UC2: Success
        UC2-->>DAG: UseCaseResult[User]
        DAG-->>Handler: UseCaseResult[User]
    else Si hay error
        DAG-->>Handler: UseCaseResult[Error]
    end
```

#### Diagrama de DAG (Directed Acyclic Graph)

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

#### Ejecución Paralela con DAG

```mermaid
graph TB
    Start([Input]) --> Parallel[Parallel DAG]

    Parallel --> UC1[Use Case 1]
    Parallel --> UC2[Use Case 2]
    Parallel --> UC3[Use Case 3]

    UC1 --> Wait[WaitGroup<br/>Espera todos]
    UC2 --> Wait
    UC3 --> Wait

    Wait --> Merge[Merge Results]
    Merge --> End([Results Array])

    style Parallel fill:#fff9c4
    style Wait fill:#e1f5ff
    style Merge fill:#e8f5e9
```

**Código de ejemplo:**

```go
// Ejemplo: Crear usuario y enviar email
func CreateUserAndPassword(ctx HandlerContext) {
    // 1. Crear casos de uso
    uc_create_user_password := usecases_user.NewCreateUserAndPasswordUseCase(...)
    uc_create_user_email := usecases_user.NewCreateUserSendEmailUseCase(...)

    // 2. Crear pipe (DAG)
    pipe := user_pipes.NewCreateUserPipe(
        ctx.c,
        ctx.Locale,
        uc_create_user_password,
        uc_create_user_email,
    )

    // 3. Ejecutar pipe (ejecuta secuencialmente)
    ucResult := pipe.Execute(userCreate)

    // 4. Resolver respuesta
    NewRequestResolver[models.User]().ResolveDTO(...)
}
```

El DAG ejecuta:
1. `CreateUserAndPasswordUseCase` → retorna `User`
2. `CreateUserSendEmailUseCase` → recibe `User`, envía email, retorna `User`

---

## ⚡ Ejecución de Tareas en Background

**GoProjectSkeleton** proporciona tres mecanismos principales para ejecutar tareas en background, cada uno diseñado para diferentes casos de uso. Esta sección explica cada opción, cuándo usarla y cómo implementarla.

### Visión General

El proyecto ofrece tres niveles de abstracción para ejecución en background:

1. **BackgroundExecutor (Workers)** - Pool de workers de bajo nivel para tareas genéricas
2. **BackgroundService** - Abstracción de alto nivel para servicios de negocio
3. **DAG con ThenBackground** - Ejecución de casos de uso en background después de operaciones exitosas

```mermaid
graph TB
    subgraph Niveles["Niveles de Abstracción"]
        Low[BackgroundExecutor<br/>Bajo Nivel<br/>Pool de Workers]
        Mid[BackgroundService<br/>Nivel Medio<br/>Servicios de Negocio]
        High[DAG ThenBackground<br/>Alto Nivel<br/>Casos de Uso]
    end

    Low -->|Usado por| Mid
    Mid -->|Usado por| High
    Low -->|Directo| DirectUse[Uso Directo<br/>Tareas Genéricas]

    style Low fill:#ffcdd2
    style Mid fill:#fff9c4
    style High fill:#c8e6c9
```

### 1. BackgroundExecutor (Workers)

**BackgroundExecutor** es un pool de workers configurable que ejecuta tareas genéricas en background. Es la capa más baja de abstracción y proporciona control total sobre la ejecución.

#### Características

- ✅ **Pool de Workers Configurable**: Número de workers y tamaño de cola personalizables
- ✅ **Gestión de Contexto**: Soporte para cancelación mediante context
- ✅ **Recuperación de Panics**: Los panics en tareas no crashean la aplicación
- ✅ **Thread-Safe**: Seguro para uso concurrente
- ✅ **Singleton Pattern**: Instancia única global disponible

#### Inicialización

El executor se inicializa durante el arranque de la aplicación:

```go
// En infrastructure/container.go
ctx := context.Background()
workers.InitializeBackgroundExecutor(
    ctx,
    settings.AppSettingsInstance.BackgroundWorkers,  // Número de workers (default: 4)
    settings.AppSettingsInstance.BackgroundQueueSize, // Tamaño de cola (default: 100)
)

// Obtener la instancia singleton
executor := workers.GetBackgroundExecutor()
```

#### Uso Básico

```go
import (
    "context"
    "github.com/simon3640/goprojectskeleton/src/application/shared/workers"
)

// Obtener el executor
executor := workers.GetBackgroundExecutor()

// Enviar una tarea
err := executor.Submit(func(ctx context.Context) {
    // Tu lógica aquí
    // Esta función se ejecutará en un worker del pool
    doSomething()
})

if err != nil {
    // Manejar error (ej: cola llena)
    log.Printf("Error al enviar tarea: %v", err)
}
```

#### Ejemplo Completo: Procesamiento de Imágenes

```go
type ImageProcessor struct {
    executor *workers.BackgroundExecutor
}

func NewImageProcessor() *ImageProcessor {
    return &ImageProcessor{
        executor: workers.GetBackgroundExecutor(),
    }
}

func (p *ImageProcessor) ProcessImageAsync(imagePath string) error {
    return p.executor.Submit(func(ctx context.Context) {
        // Verificar cancelación
        select {
        case <-ctx.Done():
            log.Printf("Procesamiento cancelado: %s", imagePath)
            return
        default:
        }

        // Procesar imagen
        if err := processImage(imagePath); err != nil {
            log.Printf("Error procesando imagen %s: %v", imagePath, err)
            return
        }

        // Notificar completado
        log.Printf("Imagen procesada: %s", imagePath)
    })
}

func processImage(path string) error {
    // Lógica de procesamiento
    return nil
}
```

#### Control de Ciclo de Vida

```go
executor := workers.GetBackgroundExecutor()

// Iniciar workers (se hace automáticamente al hacer Submit, pero puedes hacerlo manualmente)
executor.Start()

// Esperar a que todas las tareas en cola se completen
executor.Wait()

// Detener el executor (cancela contexto y cierra workers)
executor.Stop()
```

#### Configuración Recomendada

| Escenario | Workers | Queue Size | Razón |
|-----------|---------|------------|-------|
| **Desarrollo** | 2-4 | 50-100 | Recursos limitados |
| **Producción Ligera** | 4-8 | 100-200 | Carga moderada |
| **Producción Media** | 8-16 | 200-500 | Carga media |
| **Producción Alta** | 16-32 | 500-1000 | Alta concurrencia |

### 2. BackgroundService

**BackgroundService** es una abstracción de alto nivel diseñada para ejecutar servicios de negocio en background. Proporciona una interfaz tipada y estructurada para servicios que necesitan contexto de aplicación.

#### Características

- ✅ **Interfaz Tipada**: Servicios con tipos de entrada específicos
- ✅ **Contexto de Aplicación**: Acceso a AppContext y Locale
- ✅ **Factory Pattern**: Gestión centralizada de ejecución
- ✅ **Fire-and-Forget**: Ejecución asíncrona sin bloquear
- ✅ **Fallback Automático**: Si no hay executor, usa goroutine simple

#### Definir un BackgroundService

```go
import (
    app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
    "github.com/simon3640/goprojectskeleton/src/application/shared/locales"
    "github.com/simon3640/goprojectskeleton/src/application/shared/services"
)

// Definir el servicio
type SendWelcomeEmailService struct {
    emailProvider contracts.IEmailProvider
}

func NewSendWelcomeEmailService(emailProvider contracts.IEmailProvider) *SendWelcomeEmailService {
    return &SendWelcomeEmailService{
        emailProvider: emailProvider,
    }
}

// Implementar la interfaz BackgroundService
func (s *SendWelcomeEmailService) Execute(
    ctx *app_context.AppContext,
    locale locales.LocaleTypeEnum,
    input UserEmailData,
) error {
    // Lógica del servicio
    email := renderWelcomeEmail(input, locale)
    return s.emailProvider.SendEmail(email)
}

func (s *SendWelcomeEmailService) Name() string {
    return "SendWelcomeEmailService"
}

// Tipo de entrada del servicio
type UserEmailData struct {
    UserID  uint
    Email   string
    Name    string
}
```

#### Ejecutar un BackgroundService

```go
import (
    "github.com/simon3640/goprojectskeleton/src/application/shared/services"
)

// Obtener el factory singleton (inicializado en container.go)
factory := services.GetBackgroundServiceFactory()

// Crear el servicio
emailService := NewSendWelcomeEmailService(emailProvider)

// Ejecutar en background
input := UserEmailData{
    UserID: user.ID,
    Email: user.Email,
    Name:  user.Name,
}

err := services.ExecuteService(
    factory,
    emailService,
    appCtx,
    locale,
    input,
)

if err != nil {
    // Error al encolar (ej: cola llena)
    log.Printf("Error al encolar servicio: %v", err)
}
```

#### Ejemplo Completo: Servicio de Notificaciones

```go
// Servicio de notificaciones
type NotificationService struct {
    emailProvider contracts.IEmailProvider
    cacheProvider contracts.ICacheProvider
}

func NewNotificationService(
    emailProvider contracts.IEmailProvider,
    cacheProvider contracts.ICacheProvider,
) *NotificationService {
    return &NotificationService{
        emailProvider: emailProvider,
        cacheProvider: cacheProvider,
    }
}

type NotificationInput struct {
    UserID    uint
    Type      string // "welcome", "password_reset", etc.
    Data      map[string]interface{}
}

func (s *NotificationService) Execute(
    ctx *app_context.AppContext,
    locale locales.LocaleTypeEnum,
    input NotificationInput,
) error {
    // 1. Verificar si el usuario tiene notificaciones deshabilitadas
    key := fmt.Sprintf("user:%d:notifications:disabled", input.UserID)
    disabled, _ := s.cacheProvider.Exists(key)
    if disabled {
        return nil // Usuario deshabilitó notificaciones
    }

    // 2. Renderizar template según tipo
    var template string
    switch input.Type {
    case "welcome":
        template = renderWelcomeEmail(input.Data, locale)
    case "password_reset":
        template = renderPasswordResetEmail(input.Data, locale)
    default:
        return fmt.Errorf("tipo de notificación desconocido: %s", input.Type)
    }

    // 3. Enviar email
    return s.emailProvider.SendEmail(template)
}

func (s *NotificationService) Name() string {
    return "NotificationService"
}

// Uso en un handler
func CreateUserHandler(ctx HandlerContext) {
    // ... crear usuario ...

    // Enviar notificación en background
    notificationService := NewNotificationService(emailProvider, cacheProvider)
    input := NotificationInput{
        UserID: user.ID,
        Type:   "welcome",
        Data: map[string]interface{}{
            "name": user.Name,
            "email": user.Email,
        },
    }

    services.ExecuteService(
        services.GetBackgroundServiceFactory(),
        notificationService,
        ctx.AppContext,
        ctx.Locale,
        input,
    )
}
```

### 3. DAG con ThenBackground

**ThenBackground** permite ejecutar casos de uso en background después de que un DAG se ejecute exitosamente. Es ideal para tareas que deben ejecutarse después de operaciones principales pero no deben bloquear la respuesta.

#### Características

- ✅ **Integración con DAG**: Se ejecuta automáticamente después de éxito
- ✅ **Fire-and-Forget**: No bloquea la respuesta principal
- ✅ **Múltiples Background Steps**: Puedes agregar múltiples tareas en background
- ✅ **Respeto de Contexto**: Las tareas pueden ser canceladas
- ✅ **Manejo de Errores**: Los errores se registran pero no afectan el resultado principal

#### Uso Básico

```go
import (
    "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

// Crear DAG principal
dag := use_case.NewDag(
    appCtx,
    use_case.NewStep(createUserUseCase),
    locale,
    executor, // BackgroundExecutor opcional
)

// Agregar paso en background
dag = use_case.ThenBackground(
    dag,
    use_case.NewStep(sendWelcomeEmailUseCase),
    "send-welcome-email", // Nombre para logging
)

// Ejecutar DAG
result := dag.Execute(userCreate)
// El email se enviará en background si la creación fue exitosa
```

#### Ejemplo Completo: Crear Usuario con Tareas en Background

```go
func CreateUserWithBackgroundTasks(ctx HandlerContext) {
    // 1. Casos de uso principales (síncronos)
    createUserUC := usecases_user.NewCreateUserAndPasswordUseCase(...)
    activateUserUC := usecases_user.NewActivateUserUseCase(...)

    // 2. Casos de uso en background
    sendWelcomeEmailUC := usecases_user.NewCreateUserSendEmailUseCase(...)
    sendNotificationUC := usecases_user.NewSendUserNotificationUseCase(...)
    updateAnalyticsUC := usecases_user.NewUpdateUserAnalyticsUseCase(...)

    // 3. Obtener executor (opcional, si es nil usa goroutines)
    executor := workers.GetBackgroundExecutor()

    // 4. Construir DAG
    dag := use_case.NewDag(
        ctx.AppContext,
        use_case.NewStep(createUserUC),
        ctx.Locale,
        executor,
    )

    // 5. Agregar paso síncrono
    dag = use_case.Then(dag, use_case.NewStep(activateUserUC))

    // 6. Agregar pasos en background (se ejecutan solo si todo fue exitoso)
    dag = use_case.ThenBackground(
        dag,
        use_case.NewStep(sendWelcomeEmailUC),
        "send-welcome-email",
    )
    dag = use_case.ThenBackground(
        dag,
        use_case.NewStep(sendNotificationUC),
        "send-notification",
    )
    dag = use_case.ThenBackground(
        dag,
        use_case.NewStep(updateAnalyticsUC),
        "update-analytics",
    )

    // 7. Ejecutar DAG
    result := dag.Execute(userCreate)

    // 8. Resolver respuesta (las tareas en background se ejecutan asíncronamente)
    NewRequestResolver[models.User]().ResolveDTO(
        ctx.ResponseWriter,
        result,
        headers,
    )
}
```

#### Flujo de Ejecución

```mermaid
sequenceDiagram
    participant Handler as Handler
    participant DAG as DAG
    participant UC1 as CreateUser<br/>(Síncrono)
    participant UC2 as ActivateUser<br/>(Síncrono)
    participant BG1 as SendEmail<br/>(Background)
    participant BG2 as Notification<br/>(Background)

    Handler->>DAG: Execute(input)
    DAG->>UC1: Execute(input)
    UC1-->>DAG: Success[User]
    DAG->>UC2: Execute(User)
    UC2-->>DAG: Success[User]
    DAG-->>Handler: Return Result[User]

    Note over Handler: Respuesta enviada al cliente

    par Tareas en Background
        DAG->>BG1: Execute(User) [Async]
        DAG->>BG2: Execute(User) [Async]
    end

    BG1-->>DAG: (Fire-and-forget)
    BG2-->>DAG: (Fire-and-forget)
```

#### Ejecutar y Esperar Background Tasks

Si necesitas esperar a que las tareas en background se completen (útil en tests o operaciones críticas):

```go
// Ejecutar y esperar con timeout
result := dag.ExecuteWithBackground(input, 30*time.Second)

// O esperar indefinidamente
result := dag.ExecuteWithBackground(input, 0)
```

### Comparación y Guía de Selección

#### Tabla Comparativa

| Característica | BackgroundExecutor | BackgroundService | DAG ThenBackground |
|----------------|-------------------|-------------------|---------------------|
| **Nivel de Abstracción** | Bajo | Medio | Alto |
| **Tipado** | Genérico (`func(ctx)`) | Tipado (`BackgroundService[Input]`) | Tipado (Use Cases) |
| **Contexto** | `context.Context` | `AppContext` + `Locale` | `AppContext` + `Locale` |
| **Integración con DAG** | No | No | Sí |
| **Casos de Uso** | Tareas genéricas | Servicios de negocio | Casos de uso después de DAG |
| **Control de Flujo** | Manual | Manual | Automático (después de éxito) |
| **Manejo de Errores** | Manual | Logging automático | Logging automático |

#### Cuándo Usar Cada Opción

##### Usa BackgroundExecutor cuando:

- ✅ Necesitas ejecutar tareas genéricas sin estructura de negocio
- ✅ Requieres control total sobre la ejecución
- ✅ Las tareas no están relacionadas con casos de uso
- ✅ Ejemplos: procesamiento de archivos, limpieza de cache, sincronización de datos

```go
// Ejemplo: Limpieza periódica de cache
executor.Submit(func(ctx context.Context) {
    cleanExpiredCacheEntries()
})
```

##### Usa BackgroundService cuando:

- ✅ Tienes un servicio de negocio bien definido
- ✅ Necesitas contexto de aplicación (AppContext, Locale)
- ✅ El servicio tiene entrada tipada
- ✅ Quieres reutilizar el servicio en múltiples lugares
- ✅ Ejemplos: envío de emails, notificaciones, reportes

```go
// Ejemplo: Servicio de reportes
reportService := NewGenerateReportService(...)
services.ExecuteService(factory, reportService, appCtx, locale, reportInput)
```

##### Usa DAG ThenBackground cuando:

- ✅ Necesitas ejecutar casos de uso después de una operación exitosa
- ✅ Las tareas están relacionadas con el flujo principal
- ✅ Quieres que se ejecuten automáticamente solo si el DAG fue exitoso
- ✅ Necesitas múltiples tareas en background relacionadas
- ✅ Ejemplos: enviar emails después de crear usuario, actualizar analytics después de transacción

```go
// Ejemplo: Flujo completo con background
dag = Then(dag, NewStep(mainUseCase))
dag = ThenBackground(dag, NewStep(emailUseCase), "email")
dag = ThenBackground(dag, NewStep(analyticsUseCase), "analytics")
```

### Ejemplos Prácticos Completos

#### Ejemplo 1: Sistema de Notificaciones Completo

```go
// 1. Definir servicio de notificaciones
type UserNotificationService struct {
    emailProvider contracts.IEmailProvider
    logger        contracts.ILoggerProvider
}

type NotificationData struct {
    UserID  uint
    Type    string
    Subject string
    Body    string
}

func (s *UserNotificationService) Execute(
    ctx *app_context.AppContext,
    locale locales.LocaleTypeEnum,
    input NotificationData,
) error {
    s.logger.Info("Enviando notificación", map[string]interface{}{
        "user_id": input.UserID,
        "type":    input.Type,
    })

    email := &Email{
        To:      getUserEmail(input.UserID),
        Subject: input.Subject,
        Body:    input.Body,
    }

    return s.emailProvider.SendEmail(email)
}

func (s *UserNotificationService) Name() string {
    return "UserNotificationService"
}

// 2. Usar en un handler
func UpdateUserHandler(ctx HandlerContext) {
    // ... actualizar usuario ...

    // Enviar notificación en background
    notificationService := NewUserNotificationService(emailProvider, logger)
    notificationData := NotificationData{
        UserID:  user.ID,
        Type:    "profile_updated",
        Subject: "Tu perfil ha sido actualizado",
        Body:    renderNotificationBody(user, ctx.Locale),
    }

    services.ExecuteService(
        services.GetBackgroundServiceFactory(),
        notificationService,
        ctx.AppContext,
        ctx.Locale,
        notificationData,
    )
}
```

#### Ejemplo 2: Procesamiento Asíncrono de Archivos

```go
// Usando BackgroundExecutor directamente
type FileProcessor struct {
    executor *workers.BackgroundExecutor
}

func (p *FileProcessor) ProcessFileAsync(filePath string) error {
    executor := workers.GetBackgroundExecutor()

    return executor.Submit(func(ctx context.Context) {
        // Verificar cancelación
        select {
        case <-ctx.Done():
            log.Printf("Procesamiento cancelado: %s", filePath)
            return
        default:
        }

        // Procesar archivo
        if err := processFile(filePath); err != nil {
            log.Printf("Error procesando archivo %s: %v", filePath, err)
            return
        }

        // Actualizar estado
        updateFileStatus(filePath, "processed")
    })
}
```

#### Ejemplo 3: DAG Completo con Múltiples Background Tasks

```go
func CompleteUserRegistration(ctx HandlerContext) {
    // Casos de uso principales
    createUserUC := usecases_user.NewCreateUserAndPasswordUseCase(...)
    createProfileUC := usecases_user.NewCreateUserProfileUseCase(...)

    // Casos de uso en background
    sendWelcomeEmailUC := usecases_user.NewSendWelcomeEmailUseCase(...)
    sendVerificationEmailUC := usecases_user.NewSendVerificationEmailUseCase(...)
    createUserAnalyticsUC := usecases_user.NewCreateUserAnalyticsUseCase(...)
    notifyAdminsUC := usecases_user.NewNotifyAdminsNewUserUseCase(...)

    executor := workers.GetBackgroundExecutor()

    // Construir DAG
    dag := use_case.NewDag(
        ctx.AppContext,
        use_case.NewStep(createUserUC),
        ctx.Locale,
        executor,
    )

    // Paso síncrono
    dag = use_case.Then(dag, use_case.NewStep(createProfileUC))

    // Múltiples tareas en background
    dag = use_case.ThenBackground(dag, use_case.NewStep(sendWelcomeEmailUC), "welcome-email")
    dag = use_case.ThenBackground(dag, use_case.NewStep(sendVerificationEmailUC), "verification-email")
    dag = use_case.ThenBackground(dag, use_case.NewStep(createUserAnalyticsUC), "analytics")
    dag = use_case.ThenBackground(dag, use_case.NewStep(notifyAdminsUC), "admin-notification")

    // Ejecutar
    result := dag.Execute(userRegistrationInput)

    // Responder inmediatamente (las tareas en background se ejecutan asíncronamente)
    NewRequestResolver[models.User]().ResolveDTO(ctx.ResponseWriter, result, headers)
}
```

### Mejores Prácticas

#### 1. Selección de Mecanismo

- **Usa el nivel más alto posible**: DAG ThenBackground > BackgroundService > BackgroundExecutor
- **Mantén la consistencia**: Si ya usas DAG, usa ThenBackground para tareas relacionadas
- **Separa responsabilidades**: BackgroundExecutor para infraestructura, BackgroundService para negocio

#### 2. Manejo de Errores

```go
// ✅ Correcto: Logging en background service
func (s *MyService) Execute(ctx *app_context.AppContext, locale locales.LocaleTypeEnum, input Input) error {
    if err := doSomething(); err != nil {
        // Log pero no propagues el error (fire-and-forget)
        s.logger.Error("Error en background service", err)
        return err // Se loguea pero no afecta al caller
    }
    return nil
}

// ❌ Incorrecto: Panic en background
func (s *MyService) Execute(...) error {
    if err := doSomething(); err != nil {
        panic(err) // No hacer panic, usar logging
    }
}
```

#### 3. Gestión de Contexto

```go
// ✅ Correcto: Respetar cancelación
func (s *MyService) Execute(ctx *app_context.AppContext, locale locales.LocaleTypeEnum, input Input) error {
    // Verificar cancelación antes de operaciones largas
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Operación que puede tardar
    return longRunningOperation()
}
```

#### 4. Configuración de Workers

```go
// En settings o configuración
BackgroundWorkers: 8,      // Ajustar según carga
BackgroundQueueSize: 200,  // Ajustar según picos de tráfico
```

#### 5. Testing

```go
// En tests, puedes usar un executor pequeño
func TestMyService(t *testing.T) {
    ctx := context.Background()
    executor := workers.NewBackgroundExecutor(ctx, 2, 10)
    executor.Start()
    defer executor.Stop()

    factory := services.NewBackgroundServiceFactory(
        services.NewBackgroundExecutorAdapter(executor),
    )

    // Test tu servicio
    service := NewMyService(...)
    err := services.ExecuteService(factory, service, appCtx, locale, input)
    assert.NoError(t, err)

    // Esperar a que se complete
    executor.Wait()
}
```

### Resumen

| Necesidad | Solución Recomendada | Razón |
|-----------|---------------------|-------|
| Tarea genérica simple | `BackgroundExecutor` | Control directo, sin abstracciones |
| Servicio de negocio reutilizable | `BackgroundService` | Tipado, contexto de aplicación |
| Tarea después de DAG exitoso | `DAG ThenBackground` | Integración automática, flujo claro |
| Múltiples tareas relacionadas | `DAG ThenBackground` | Agregar múltiples pasos fácilmente |
| Procesamiento de archivos | `BackgroundExecutor` | Tareas de infraestructura |
| Envío de emails/notificaciones | `BackgroundService` | Servicios de negocio bien definidos |

---

## 📊 Observabilidad

**GoProjectSkeleton** incluye un stack completo de observabilidad con **OpenTelemetry**, **Prometheus**, **Jaeger** y **Grafana**. El sistema proporciona trazado distribuido, recolección de métricas y logging estructurado en todas las capas de la aplicación incluyendo Casos de Uso, DAGs y Servicios en Background.

### Visión General

El sistema de observabilidad sigue estos principios:

1. **Instrumentación Siempre Activa**: Todos los componentes están instrumentados por defecto
2. **Fallback No-Op**: Cuando la observabilidad está deshabilitada, se usan implementaciones no-op
3. **Cumplimiento de Arquitectura Limpia**: Los contratos de observabilidad viven en la capa de aplicación
4. **Propagación Automática de Trazas**: Las trazas se propagan a través de límites de contexto

### Arquitectura de Observabilidad

```mermaid
graph TB
    subgraph Application["📱 Aplicación"]
        UC[Casos de Uso<br/>Spans Automáticos]
        DAG[DAG Steps<br/>Ejecución Paralela/Secuencial]
        BG[Servicios Background<br/>Trazas Asíncronas]
        HTTP[Handlers HTTP<br/>Métricas de Request]
    end

    subgraph Instrumentation["🔧 Capa de Instrumentación"]
        Tracer[ITracer<br/>Gestión de Spans]
        Metrics[IMetricsCollector<br/>Contadores/Histogramas]
        Logger[ILoggerProvider<br/>Logs Estructurados]
    end

    subgraph Infrastructure["🏗️ Infraestructura"]
        OTEL[OpenTelemetry SDK<br/>Exportadores]
        PROM[Prometheus<br/>Almacenamiento de Métricas]
        JAEGER[Jaeger<br/>Backend de Trazas]
    end

    subgraph Visualization["📊 Visualización"]
        GRAFANA[Grafana<br/>Dashboards]
        JAEGER_UI[Jaeger UI<br/>Explorador de Trazas]
        PROM_UI[Prometheus UI<br/>Explorador de Métricas]
    end

    UC --> Tracer
    DAG --> Tracer
    BG --> Tracer
    HTTP --> Metrics

    Tracer --> OTEL
    Metrics --> OTEL
    Logger --> OTEL

    OTEL --> PROM
    OTEL --> JAEGER

    PROM --> GRAFANA
    JAEGER --> JAEGER_UI
    PROM --> PROM_UI

    style UC fill:#e3f2fd
    style DAG fill:#e3f2fd
    style BG fill:#e3f2fd
    style GRAFANA fill:#c8e6c9
    style JAEGER fill:#fff9c4
    style PROM fill:#ffcdd2
```

### Configuración

La observabilidad se configura a través de variables de entorno:

```bash
# Configuración OpenTelemetry
OTEL_ENABLED=true                           # Habilitar/deshabilitar observabilidad
OTEL_SERVICE_NAME=goprojectskeleton         # Nombre del servicio para trazas
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317  # Endpoint del colector OTLP
OTEL_EXPORTER_OTLP_INSECURE=true            # Usar conexión insegura (desarrollo)
OTEL_TRACES_SAMPLER=always_on               # Estrategia de muestreo de trazas
OTEL_METRICS_EXPORTER=prometheus            # Exportador de métricas
OTEL_LOGS_EXPORTER=otlp                     # Exportador de logs

# Configuración de Prometheus
PROMETHEUS_PORT=9090                        # Puerto del servidor Prometheus
METRICS_PATH=/metrics                       # Ruta de exposición de métricas

# Configuración del Colector OpenTelemetry
OTEL_COLLECTOR_HOST=otel-collector          # Host del colector
OTEL_COLLECTOR_GRPC_PORT=4317               # Puerto gRPC del colector
OTEL_COLLECTOR_HTTP_PORT=4318               # Puerto HTTP del colector
```

### Componentes de Observabilidad

#### Interfaz ITracer

El tracer proporciona gestión de spans para trazado distribuido:

```go
// application/contracts/observability/tracer.go
type ITracer interface {
    // StartSpan crea un nuevo span con el nombre dado
    StartSpan(ctx context.Context, name string, opts ...SpanOption) (context.Context, ISpan)

    // StartSpanWithParent crea un span hijo de un span padre
    StartSpanWithParent(ctx context.Context, parent ISpan, name string, opts ...SpanOption) (context.Context, ISpan)

    // ExtractSpanContext extrae el contexto del span de los carriers (headers HTTP)
    ExtractSpanContext(ctx context.Context, carrier map[string]string) context.Context

    // InjectSpanContext inyecta el contexto del span en carriers para propagación
    InjectSpanContext(ctx context.Context, carrier map[string]string)
}
```

#### Interfaz ISpan

Los spans representan operaciones individuales dentro de una traza:

```go
// application/contracts/observability/span.go
type ISpan interface {
    // End completa el span
    End()

    // SetStatus establece el estado del span (OK, Error)
    SetStatus(code SpanStatusCode, description string)

    // SetAttributes agrega atributos clave-valor al span
    SetAttributes(attrs ...SpanAttribute)

    // RecordError registra un error en el span
    RecordError(err error)

    // AddEvent agrega un evento con timestamp al span
    AddEvent(name string, attrs ...SpanAttribute)

    // SpanContext retorna el contexto del span para propagación
    SpanContext() SpanContext
}
```

#### Interfaz IMetricsCollector

El recolector de métricas proporciona contadores, gauges e histogramas:

```go
// application/contracts/observability/metrics_collector.go
type IMetricsCollector interface {
    // Counter incrementa un contador
    Counter(name string, value float64, labels ...MetricLabel)

    // Gauge establece un valor de gauge
    Gauge(name string, value float64, labels ...MetricLabel)

    // Histogram registra un valor en un histograma
    Histogram(name string, value float64, labels ...MetricLabel)

    // Timer registra una duración
    Timer(name string, duration time.Duration, labels ...MetricLabel)
}
```

#### Interfaz ILoggerProvider (Mejorado)

El logger soporta logs estructurados con correlación de trazas:

```go
// application/contracts/providers/logger_provider.go
type ILoggerProvider interface {
    // Métodos básicos de logging
    Info(msg string, fields ...map[string]interface{})
    Error(msg string, err error, fields ...map[string]interface{})
    Debug(msg string, fields ...map[string]interface{})
    Warn(msg string, fields ...map[string]interface{})

    // WithContext crea un logger con contexto de traza
    WithContext(ctx context.Context) ILoggerProvider

    // WithFields crea un logger con campos predeterminados
    WithFields(fields map[string]interface{}) ILoggerProvider
}
```

### Observabilidad en Casos de Uso

Los Casos de Uso se instrumentan automáticamente con trazado y métricas:

```go
// application/modules/user/use_cases/create_user.go
type CreateUserUseCase struct {
    log         contracts.ILoggerProvider
    repo        contracts.IUserRepository
    tracer      observability.ITracer
    metrics     observability.IMetricsCollector
}

func (uc *CreateUserUseCase) Execute(
    ctx context.Context,
    locale locales.LocaleTypeEnum,
    input dtos.UserCreate,
) *usecase.UseCaseResult[models.User] {
    // Iniciar span para este caso de uso
    ctx, span := uc.tracer.StartSpan(ctx, "CreateUserUseCase.Execute",
        observability.WithSpanKind(observability.SpanKindInternal),
        observability.WithAttributes(
            observability.String("user.email", input.Email),
            observability.String("locale", string(locale)),
        ),
    )
    defer span.End()

    result := usecase.NewUseCaseResult[models.User]()
    startTime := time.Now()

    // Validar entrada
    uc.validate(ctx, input, result)
    if result.HasError() {
        span.SetStatus(observability.SpanStatusError, "validation failed")
        span.RecordError(fmt.Errorf("validation error: %v", result.Error))
        uc.metrics.Counter("usecase.create_user.validation_errors", 1,
            observability.Label("error_code", string(result.StatusCode)),
        )
        return result
    }

    // Crear usuario
    user, err := uc.repo.Create(input)
    if err != nil {
        span.SetStatus(observability.SpanStatusError, err.ErrMsg)
        span.RecordError(fmt.Errorf(err.ErrMsg))
        uc.metrics.Counter("usecase.create_user.errors", 1)
        result.SetError(err.Code, err.Context)
        return result
    }

    // Registrar éxito
    span.SetStatus(observability.SpanStatusOK, "user created")
    span.SetAttributes(observability.Int64("user.id", int64(user.ID)))

    // Registrar métricas
    uc.metrics.Counter("usecase.create_user.success", 1)
    uc.metrics.Timer("usecase.create_user.duration", time.Since(startTime))

    result.SetData(status.Created, *user, "User created successfully")
    return result
}
```

### Observabilidad en DAG

El sistema DAG proporciona instrumentación automática para pasos secuenciales, paralelos y en background:

#### Steps Secuenciales

```go
// Los steps secuenciales crean spans hijos automáticamente
dag := use_case.NewDag(
    appCtx,
    use_case.NewStep(createUserUC),
    locale,
    executor,
)
dag = use_case.Then(dag, use_case.NewStep(sendEmailUC))

// Jerarquía de spans resultante:
// DAG.Execute
// ├── Step[0]: CreateUserUseCase
// └── Step[1]: SendEmailUseCase
```

#### Steps Paralelos

```go
// Los steps paralelos crean spans hermanos con el mismo padre
dag := use_case.NewDag(appCtx, use_case.NewStep(mainUC), locale, executor)
dag = use_case.ThenParallel(dag,
    use_case.NewStep(notifyUC),
    use_case.NewStep(analyticsUC),
    use_case.NewStep(auditUC),
)

// Jerarquía de spans resultante:
// DAG.Execute
// ├── Step[0]: MainUseCase
// └── Parallel
//     ├── ParallelStep[0]: NotifyUseCase
//     ├── ParallelStep[1]: AnalyticsUseCase
//     └── ParallelStep[2]: AuditUseCase
```

#### Steps en Background

```go
// Los steps en background usan span links para mantener la correlación de trazas
dag = use_case.ThenBackground(
    dag,
    use_case.NewStep(sendWelcomeEmailUC),
    "send-welcome-email",
)

// Jerarquía de spans resultante:
// DAG.Execute (completa inmediatamente)
// ├── Step[0]: CreateUserUseCase
// └── BackgroundStep: SendWelcomeEmailUseCase
//     └── [Link al span padre del DAG]
```

### Observabilidad en Servicios Background

Los servicios en background se instrumentan automáticamente con span links:

```go
// application/shared/services/background_service.go
type ObservableBackgroundService[I any] struct {
    service BackgroundService[I]
    tracer  observability.ITracer
    metrics observability.IMetricsCollector
}

func (s *ObservableBackgroundService[I]) Execute(
    ctx *app_context.AppContext,
    locale locales.LocaleTypeEnum,
    input I,
) error {
    // Crear span con link al span padre (si existe)
    spanCtx, span := s.tracer.StartSpan(ctx.Context(), s.service.Name(),
        observability.WithSpanKind(observability.SpanKindInternal),
        observability.WithFollowsFrom(ctx.Context()), // Span link, no hijo
    )
    defer span.End()

    startTime := time.Now()

    // Ejecutar el servicio
    err := s.service.Execute(
        app_context.WithContext(ctx, spanCtx),
        locale,
        input,
    )

    // Registrar métricas
    duration := time.Since(startTime)
    s.metrics.Timer("background_service.duration", duration,
        observability.Label("service", s.service.Name()),
    )

    if err != nil {
        span.SetStatus(observability.SpanStatusError, err.Error())
        span.RecordError(err)
        s.metrics.Counter("background_service.errors", 1,
            observability.Label("service", s.service.Name()),
        )
        return err
    }

    span.SetStatus(observability.SpanStatusOK, "completed")
    s.metrics.Counter("background_service.success", 1,
        observability.Label("service", s.service.Name()),
    )

    return nil
}
```

### Instrumentación HTTP

Los handlers HTTP se instrumentan automáticamente para métricas de request:

```go
// infrastructure/server/middlewares/observability.go
func ObservabilityMiddleware(metrics observability.IMetricsCollector, tracer observability.ITracer) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extraer contexto de traza de headers entrantes
        ctx := tracer.ExtractSpanContext(c.Request.Context(), extractHeaders(c))

        // Iniciar span para el request HTTP
        ctx, span := tracer.StartSpan(ctx, fmt.Sprintf("HTTP %s %s", c.Request.Method, c.FullPath()),
            observability.WithSpanKind(observability.SpanKindServer),
            observability.WithAttributes(
                observability.String("http.method", c.Request.Method),
                observability.String("http.url", c.Request.URL.String()),
                observability.String("http.user_agent", c.Request.UserAgent()),
            ),
        )
        defer span.End()

        // Actualizar contexto del request
        c.Request = c.Request.WithContext(ctx)

        startTime := time.Now()

        // Procesar request
        c.Next()

        // Registrar métricas
        duration := time.Since(startTime)
        statusCode := c.Writer.Status()

        span.SetAttributes(
            observability.Int("http.status_code", statusCode),
            observability.Int64("http.response_size", int64(c.Writer.Size())),
        )

        if statusCode >= 400 {
            span.SetStatus(observability.SpanStatusError, fmt.Sprintf("HTTP %d", statusCode))
        } else {
            span.SetStatus(observability.SpanStatusOK, "")
        }

        metrics.Histogram("http_request_duration_seconds", duration.Seconds(),
            observability.Label("method", c.Request.Method),
            observability.Label("path", c.FullPath()),
            observability.Label("status", fmt.Sprintf("%d", statusCode)),
        )

        metrics.Counter("http_requests_total", 1,
            observability.Label("method", c.Request.Method),
            observability.Label("path", c.FullPath()),
            observability.Label("status", fmt.Sprintf("%d", statusCode)),
        )
    }
}
```

### Stack de Grafana

El proyecto incluye una configuración preconfigurada de Grafana con:

#### Fuentes de Datos

```yaml
# docker/grafana/provisioning/datasources/datasources.yaml
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true

  - name: Jaeger
    type: jaeger
    access: proxy
    url: http://jaeger:16686
```

#### Dashboards Incluidos

1. **Métricas de API** (`api-metrics.json`)
   - Tasa de requests por endpoint
   - Latencia de respuesta (p50, p90, p99)
   - Tasa de errores por código de estado
   - Volumen de requests por método

2. **Métricas de Casos de Uso** (`usecase-metrics.json`)
   - Tasa de éxito/error por caso de uso
   - Duración de ejecución
   - Errores de validación
   - Distribución de throughput

3. **Métricas de Background** (`background-metrics.json`)
   - Tamaño de cola de servicios background
   - Duración de procesamiento
   - Tasa de errores
   - Utilización de workers

4. **Métricas de DAG** (`dag-metrics.json`)
   - Duración de ejecución del DAG
   - Tiempo de ejecución de pasos paralelos
   - Rendimiento de tareas background
   - Propagación de errores

### Configuración Docker

El stack de observabilidad está incluido en `docker-compose.dev.yml`:

```yaml
services:
  # Colector OpenTelemetry
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel/otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4317:4317"   # gRPC OTLP
      - "4318:4318"   # HTTP OTLP
      - "8888:8888"   # Métricas del colector

  # Prometheus
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  # Jaeger
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686" # Jaeger UI
      - "14268:14268" # Collector HTTP
      - "14250:14250" # Collector gRPC

  # Grafana
  grafana:
    image: grafana/grafana:latest
    volumes:
      - ./grafana/provisioning:/etc/grafana/provisioning
      - ./grafana/dashboards:/var/lib/grafana/dashboards
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
```

### Crear un Caso de Uso con Observabilidad

Aquí tienes una guía completa para crear un caso de uso instrumentado:

```go
package use_cases

import (
    "context"
    "fmt"
    "time"

    "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
    "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
    "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
    "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs/dtos"
    "github.com/simon3640/goprojectskeleton/src/application/shared/locales"
    "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
    "github.com/simon3640/goprojectskeleton/src/domain/models"
)

type CreateOrderUseCase struct {
    log         providers.ILoggerProvider
    repo        repositories.IOrderRepository
    tracer      observability.ITracer
    metrics     observability.IMetricsCollector
}

func NewCreateOrderUseCase(
    log providers.ILoggerProvider,
    repo repositories.IOrderRepository,
    tracer observability.ITracer,
    metrics observability.IMetricsCollector,
) *CreateOrderUseCase {
    return &CreateOrderUseCase{
        log:     log,
        repo:    repo,
        tracer:  tracer,
        metrics: metrics,
    }
}

func (uc *CreateOrderUseCase) Execute(
    ctx context.Context,
    locale locales.LocaleTypeEnum,
    input dtos.OrderCreate,
) *use_case.UseCaseResult[models.Order] {
    // 1. Iniciar span con atributos
    ctx, span := uc.tracer.StartSpan(ctx, "CreateOrderUseCase.Execute",
        observability.WithSpanKind(observability.SpanKindInternal),
        observability.WithAttributes(
            observability.String("order.customer_id", input.CustomerID),
            observability.Float64("order.total", input.Total),
            observability.Int("order.items_count", len(input.Items)),
        ),
    )
    defer span.End()

    result := use_case.NewUseCaseResult[models.Order]()
    startTime := time.Now()

    // 2. Usar logger con contexto para correlación de trazas
    log := uc.log.WithContext(ctx)
    log.Info("Creating new order", map[string]interface{}{
        "customer_id": input.CustomerID,
        "items_count": len(input.Items),
    })

    // 3. Agregar evento para operaciones significativas
    span.AddEvent("validating_input")

    // 4. Validar entrada
    if err := input.Validate(); err != nil {
        span.SetStatus(observability.SpanStatusError, "validation failed")
        span.RecordError(err)
        uc.metrics.Counter("usecase.create_order.validation_errors", 1,
            observability.Label("reason", err.Error()),
        )
        result.SetValidationError(err.Error())
        return result
    }

    // 5. Agregar evento para llamada al repositorio
    span.AddEvent("calling_repository")

    // 6. Crear orden (el repositorio también puede crear su propio span hijo)
    order, repoErr := uc.repo.Create(input)
    if repoErr != nil {
        span.SetStatus(observability.SpanStatusError, repoErr.ErrMsg)
        span.RecordError(fmt.Errorf(repoErr.ErrMsg))
        uc.metrics.Counter("usecase.create_order.errors", 1,
            observability.Label("error_type", "repository"),
        )
        log.Error("Failed to create order", fmt.Errorf(repoErr.ErrMsg), nil)
        result.SetError(repoErr.Code, repoErr.Context)
        return result
    }

    // 7. Registrar métricas de éxito
    duration := time.Since(startTime)
    uc.metrics.Counter("usecase.create_order.success", 1)
    uc.metrics.Timer("usecase.create_order.duration", duration)
    uc.metrics.Histogram("order.total_amount", input.Total,
        observability.Label("currency", input.Currency),
    )

    // 8. Establecer estado de éxito y agregar atributos del resultado
    span.SetStatus(observability.SpanStatusOK, "order created")
    span.SetAttributes(
        observability.String("order.id", order.ID),
        observability.String("order.status", order.Status),
    )

    log.Info("Order created successfully", map[string]interface{}{
        "order_id": order.ID,
        "duration_ms": duration.Milliseconds(),
    })

    result.SetData(status.Created, *order, "Order created successfully")
    return result
}
```

### Implementación No-Op

Cuando la observabilidad está deshabilitada, se usan implementaciones no-op:

```go
// application/shared/observability/noop/noop_tracer.go
type NoopTracer struct{}

func (t *NoopTracer) StartSpan(ctx context.Context, name string, opts ...observability.SpanOption) (context.Context, observability.ISpan) {
    return ctx, &NoopSpan{}
}

type NoopSpan struct{}

func (s *NoopSpan) End()                                                    {}
func (s *NoopSpan) SetStatus(code observability.SpanStatusCode, desc string) {}
func (s *NoopSpan) SetAttributes(attrs ...observability.SpanAttribute)       {}
func (s *NoopSpan) RecordError(err error)                                   {}
func (s *NoopSpan) AddEvent(name string, attrs ...observability.SpanAttribute) {}
```

### Acceder a los Dashboards

Una vez que los servicios estén corriendo:

| Servicio | URL | Descripción |
|----------|-----|-------------|
| **Grafana** | `http://localhost:3001` | Dashboards y alertas (admin/admin) |
| **Jaeger UI** | `http://localhost:16686` | Explorador de trazas |
| **Prometheus** | `http://localhost:9090` | Consultas de métricas |

### Mejores Prácticas

1. **Siempre propagar contexto**: Pasar `context.Context` a través de todas las capas
2. **Usar nombres de span significativos**: `{Componente}.{Operación}` (ej: `UserRepository.Create`)
3. **Agregar atributos relevantes**: Incluir IDs de negocio, estados, conteos
4. **Registrar errores apropiadamente**: Usar `span.RecordError()` para excepciones
5. **Usar eventos para hitos**: Agregar eventos para operaciones significativas
6. **Mantener cardinalidad baja**: Evitar etiquetas de alta cardinalidad en métricas
7. **Usar muestreo**: Configurar muestreo apropiado para producción

---

## Virtudes y Beneficios

### 1. Arquitectura Sólida y Escalable

#### ✅ Clean Architecture
- **Separación clara de responsabilidades**: Cada capa tiene un propósito específico
- **Independencia de frameworks**: Puedes cambiar Gin por otro framework sin afectar el negocio
- **Testabilidad**: Fácil de testear cada capa independientemente

#### ✅ Arquitectura Hexagonal
- **Desacoplamiento total**: La lógica de negocio no conoce detalles de implementación
- **Ports & Adapters**: Interfaces claras entre capas
- **Flexibilidad**: Cambiar base de datos, proveedores, etc. sin reescribir código

### 2. Productividad del Desarrollador

#### ✅ Estructura Clara
- **Organización lógica**: Fácil encontrar código
- **Convenciones consistentes**: Mismo patrón en todo el proyecto
- **Onboarding rápido**: Nuevos desarrolladores entienden rápido

#### ✅ Reutilización
- **Componentes compartidos**: DTOs, errores, validaciones reutilizables
- **Base de repositorios**: `RepositoryBase` reduce código duplicado
- **Providers intercambiables**: Cambiar implementaciones fácilmente

### 3. Mantenibilidad

#### ✅ Código Limpio
- **Principios SOLID**: Aplicados consistentemente
- **DRY (Don't Repeat Yourself)**: Mínima duplicación
- **Nombres descriptivos**: Código auto-documentado

#### ✅ Extensibilidad
- **Agregar funcionalidades**: Sin modificar código existente
- **Módulos independientes**: Cada módulo es autocontenido
- **Interfaces bien definidas**: Contratos claros

### 4. Testabilidad

#### ✅ Testing en Capas
- **Unit tests**: Casos de uso testables con mocks
- **Integration tests**: Tests con base de datos real
- **E2E tests**: Tests completos del flujo

#### ✅ Mocks Completos
- **Mocks de repositorios**: Fácil simular datos
- **Mocks de providers**: Simular servicios externos
- **Testing aislado**: Sin dependencias externas

### 5. Seguridad

#### ✅ Autenticación Robusta
- **JWT completo**: Access y refresh tokens
- **OTP (2FA)**: Autenticación de dos factores
- **Hash seguro**: Bcrypt para contraseñas

#### ✅ Validación
- **Validación en múltiples capas**: DTOs, casos de uso, repositorios
- **Sanitización**: Prevención de inyecciones
- **Guards**: Control de acceso basado en roles

### 6. Rendimiento

#### ✅ Optimizaciones
- **Cache con Redis**: Reduce consultas a BD
- **Connection pooling**: Reutilización de conexiones
- **Goroutines**: Concurrencia nativa de Go
- **Query optimization**: Filtros y paginación eficientes

### 7. Internacionalización

#### ✅ Multiidioma
- **Soporte i18n**: Mensajes en múltiples idiomas
- **Locale por request**: Cada request puede tener su idioma
- **Mensajes centralizados**: Fácil agregar nuevos idiomas

### 8. Documentación

#### ✅ Swagger Automático
- **API documentada**: Endpoints documentados automáticamente
- **Ejemplos**: Ejemplos de requests y responses
- **Tipos claros**: Esquemas bien definidos

### 9. DevOps y Despliegue

#### ✅ Docker Completo
- **Multi-servicio**: Aplicación, BD, Redis, etc.
- **Entornos separados**: Dev, test, E2E
- **Hot reload**: Desarrollo eficiente

#### ✅ Configuración Flexible
- **Variables de entorno**: Configuración por ambiente
- **Settings centralizados**: Un solo lugar para configuración
- **Validación**: Validación de configuración al inicio

### 10. Escalabilidad

#### ✅ Horizontal
- **Stateless**: Listo para load balancers
- **Cache distribuido**: Redis compartido
- **Sin sesiones**: JWT stateless

#### ✅ Vertical
- **Goroutines**: Concurrencia eficiente
- **Optimizaciones**: Cache, pooling, etc.
- **Serverless ready**: Fácil migración a serverless

### 11. Observabilidad

#### ✅ Trazado Distribuido
- **OpenTelemetry**: Instrumentación estándar de la industria
- **Jaeger**: Visualización y análisis de trazas
- **Propagación de contexto**: Trazas a través de límites de servicios

#### ✅ Métricas
- **Prometheus**: Recolección y almacenamiento de métricas
- **Grafana**: Dashboards y alertas
- **Métricas personalizadas**: Métricas de casos de uso y negocio

#### ✅ Logging
- **Logging estructurado**: Logs en formato JSON
- **Correlación de trazas**: Logs conectados a trazas
- **Niveles configurables**: Info, Debug, Error, Warn

### Beneficios para Iniciar un Proyecto

1. **Ahorro de Tiempo**
   - Estructura lista para usar
   - Patrones implementados
   - No empezar desde cero

2. **Mejores Prácticas**
   - Arquitectura probada
   - Patrones de diseño aplicados
   - Código de calidad

3. **Escalabilidad Garantizada**
   - Diseñado para crecer
   - Fácil agregar funcionalidades
   - Listo para producción

4. **Mantenibilidad a Largo Plazo**
   - Código organizado
   - Fácil de entender
   - Fácil de modificar

5. **Equipo Productivo**
   - Onboarding rápido
   - Convenciones claras
   - Menos bugs

---

## 📊 Estadísticas del Proyecto

| Métrica | Valor |
|---------|-------|
| **Archivos Go** | ~180+ archivos |
| **Líneas de Código** | ~15,000+ líneas |
| **Casos de Uso** | 20+ casos de uso |
| **Módulos de Negocio** | 4 módulos (auth, user, password, status) |
| **Providers** | 7 providers (JWT, Hash, Email, Cache, Logger, Renderer, Status) |
| **Repositorios** | 6 repositorios |
| **Handlers HTTP** | 15+ endpoints |
| **Tests** | 20+ archivos de test |
| **Templates** | 6+ templates HTML |
| **Idiomas Soportados** | 2 (Español, Inglés) |
| **Componentes Observabilidad** | 4 (OpenTelemetry, Prometheus, Jaeger, Grafana) |
| **Dashboards Grafana** | 4 dashboards preconfigurados |

## Estructura del Proyecto - Capa por Capa

### Visión General de la Estructura

```
GoProjectSkeleton/
├── src/
│   ├── domain/              # 🎯 Capa de Dominio (Núcleo)
│   │   ├── models/          # Entidades de negocio puras
│   │   └── utils/           # Utilidades de dominio
│   ├── application/         # 💼 Capa de Aplicación (Lógica de Negocio)
│   │   ├── contracts/      # Interfaces (Ports)
│   │   ├── modules/         # Módulos de negocio
│   │   └── shared/         # Componentes compartidos
│   └── infrastructure/     # 🔧 Capa de Infraestructura (Detalles Técnicos)
│       ├── server/          # Servidor HTTP (Gin)
│       ├── database/        # Base de datos (GORM)
│       ├── providers/       # Implementaciones de providers
│       ├── repositories/    # Implementaciones de repositorios
│       ├── handlers/        # Handlers HTTP
│       ├── config/          # Configuración
│       └── clouds/          # Adaptadores Cloud (AWS, Azure)
├── docker/                  # 🐳 Configuración Docker
│   ├── docker-compose.dev.yml
│   ├── docker-compose.test.yml
│   ├── docker-compose.e2e.yml
│   ├── db/                  # Configuración de base de datos
│   ├── grafana/             # Dashboards y datasources de Grafana
│   ├── otel/                # Configuración del colector OpenTelemetry
│   └── prometheus/          # Configuración de Prometheus
├── tests/                   # 🧪 Tests del proyecto
│   ├── integration/         # Tests de integración
│   └── e2e/                 # Tests end-to-end (Bruno)
├── src/infrastructure/docs/ # 📚 Servicio Swagger Independiente
│   ├── main.go              # Servidor HTTP independiente para Swagger
│   ├── config/              # Configuración del servidor Swagger
│   ├── swagger/             # Archivos generados de Swagger
│   │   ├── swagger.json
│   │   ├── swagger.yaml
│   │   └── docs.go
│   └── go.mod               # Módulo independiente para Swagger
└── IDE/                     # ⚙️ Configuración del IDE
    ├── launch.json
    └── tasks.json
```

### Capa 1: Domain (Dominio)

**Responsabilidad**: Contiene las entidades de negocio puras, sin dependencias externas.

```
domain/
├── models/                  # Modelos de dominio
│   ├── user.go             # Entidad User
│   ├── role.go             # Entidad Role
│   ├── password.go         # Entidad Password
│   ├── one_time_password.go
│   ├── one_time_token.go
│   └── status.go
└── utils/                   # Utilidades de dominio
    └── query_payload.go    # Utilidades para queries
```

#### Características

- **Sin dependencias externas**: No importa frameworks, BD, etc.
- **Entidades puras**: Solo lógica de negocio
- **Validaciones de dominio**: Reglas de negocio en los modelos
- **Inmutabilidad preferida**: Modelos inmutables cuando es posible

#### Ejemplo: Modelo User

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

### Capa 2: Application (Aplicación)

**Responsabilidad**: Contiene la lógica de negocio y casos de uso.

```
application/
├── contracts/               # Contratos (Interfaces)
│   ├── providers/          # Interfaces de proveedores
│   │   ├── cache_provider.go
│   │   ├── email_provider.go
│   │   ├── hash_provider.go
│   │   ├── jwt_provider.go
│   │   └── logger_provider.go
│   └── repositories/       # Interfaces de repositorios
│       ├── user.go
│       ├── password.go
│       └── role.go
├── modules/                # Módulos de negocio
│   ├── auth/              # Módulo de autenticación
│   ├── user/              # Módulo de usuarios
│   ├── password/          # Módulo de contraseñas
│   └── status/            # Módulo de estado
└── shared/                # Componentes compartidos
    ├── DTOs/              # Data Transfer Objects
    ├── errors/            # Manejo de errores
    ├── services/          # Servicios compartidos
    ├── templates/         # Plantillas (emails)
    ├── use_case/          # Base de casos de uso
    ├── context/           # Contexto de aplicación
    ├── locales/           # Internacionalización
    ├── settings/          # Configuración de aplicación
    └── guards/            # Guards de autorización
```

#### Características

- **Casos de uso**: Cada funcionalidad es un caso de uso
- **Interfaces**: Define contratos que la infraestructura implementa
- **DTOs**: Objetos de transferencia de datos
- **Servicios**: Lógica compartida entre módulos

#### Ejemplo: Caso de Uso

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
    // 1. Validar
    // 2. Ejecutar lógica de negocio
    // 3. Llamar a repositorio
    // 4. Retornar resultado
}
```

### Capa 3: Infrastructure (Infraestructura)

**Responsabilidad**: Implementaciones técnicas (HTTP, BD, etc.).

```
infrastructure/
├── api/                    # Capa de API HTTP
│   ├── cmd/
│   │   └── main.go        # Punto de entrada
│   ├── routes/            # Definición de rutas
│   ├── middlewares/       # Middlewares HTTP
│   └── types.go
├── config/                # Configuración
│   ├── config.go         # Carga de configuración
│   └── env.go            # Variables de entorno
├── database/              # Base de datos
│   └── goprojectskeleton/   # Implementación GORM
│       ├── models/       # Modelos de BD
│       └── init_db/      # Inicialización
├── handlers/              # Handlers HTTP
│   ├── user.go
│   ├── auth.go
│   ├── password.go
│   └── status.go
├── providers/            # Implementaciones de proveedores
│   ├── jwt_provider.go
│   ├── hash_provider.go
│   ├── email_provider.go
│   └── cache_provider.go
└── repositories/          # Implementaciones de repositorios
    ├── user.go
    ├── password.go
    └── role.go
```

#### Características

- **Implementaciones concretas**: GORM, Gin, Redis, etc.
- **Adaptadores**: Adaptan frameworks a interfaces
- **Configuración**: Carga de configuración desde variables de entorno

#### Ejemplo: Repositorio

```go
// infrastructure/repositories/user.go
type UserRepository struct {
    RepositoryBase[UserCreate, UserUpdate, User, dbModels.User]
    DB *gorm.DB
}

func (r *UserRepository) Create(input UserCreate) (*User, error) {
    // Implementación con GORM
}
```

---

## Revisión Exhaustiva por Carpetas

### `/src/domain/` - Capa de Dominio

#### `/src/domain/models/`

Contiene todas las entidades de dominio puras.

**Archivos principales:**

- **`user.go`**: Entidad User con validaciones de dominio
  - `UserBase`: Estructura base del usuario
  - `User`: Usuario con metadatos de BD
  - `UserWithRole`: Usuario con información de rol
  - `Validate()`: Validaciones de negocio

- **`role.go`**: Entidad Role
  - Define roles del sistema
  - Prioridades y permisos

- **`password.go`**: Entidad Password
  - Hash de contraseñas
  - Expiración
  - Validación de fortaleza

- **`one_time_password.go`**: OTP para autenticación
  - Generación de códigos
  - TTL y expiración

- **`one_time_token.go`**: Tokens de un solo uso
  - Para reset de contraseña
  - Para activación de cuenta

- **`status.go`**: Estados del sistema
  - Estados de usuarios
  - Estados de aplicación

- **`db_models_base.go`**: Base para modelos con metadatos de BD
  - ID, CreatedAt, UpdatedAt, DeletedAt

- **`models_utils.go`**: Utilidades para modelos
  - Validación de email
  - Validación de contraseña

#### `/src/domain/utils/`

Utilidades de dominio.

- **`query_payload.go`**: Construcción de queries
  - Filtros
  - Ordenamiento
  - Paginación

### `/src/application/` - Capa de Aplicación

#### `/src/application/contracts/`

Define todas las interfaces (contratos) que la infraestructura debe implementar.

##### `/src/application/contracts/providers/`

Interfaces de proveedores externos:

- **`cache_provider.go`**: Interfaz para cache (Redis)
  - `Get()`, `Set()`, `Delete()`, `Exists()`

- **`email_provider.go`**: Interfaz para envío de emails
  - `SendEmail()`

- **`hash_provider.go`**: Interfaz para hashing
  - `Hash()`, `Compare()`

- **`jwt_provider.go`**: Interfaz para JWT
  - `GenerateAccessToken()`, `GenerateRefreshToken()`, `ParseTokenAndValidate()`

- **`logger_provider.go`**: Interfaz para logging
  - `Info()`, `Error()`, `Debug()`, `Panic()`

- **`renderer_provider.go`**: Interfaz para renderizado de templates
  - `Render()`

- **`status_provider.go`**: Interfaz para estado del sistema

##### `/src/application/contracts/observability/`

Interfaces de observabilidad:

- **`tracer.go`**: Interfaz para trazado distribuido
  - `StartSpan()`, `StartSpanWithParent()`, `ExtractSpanContext()`, `InjectSpanContext()`

- **`span.go`**: Interfaz para spans individuales
  - `End()`, `SetStatus()`, `SetAttributes()`, `RecordError()`, `AddEvent()`

- **`metrics_collector.go`**: Interfaz para métricas
  - `Counter()`, `Gauge()`, `Histogram()`, `Timer()`

- **`logger.go`**: Interfaz extendida para logging con trazas
  - `WithContext()`, `WithFields()`

- **`clock.go`**: Interfaz para abstracción de tiempo
  - `Now()`, `Since()`

##### `/src/application/contracts/repositories/`

Interfaces de repositorios:

- **`base.go`**: Interfaz base para repositorios
  - `Create()`, `GetByID()`, `Update()`, `Delete()`, `GetAll()`

- **`user.go`**: Interfaz específica de usuarios
  - `CreateWithPassword()`, `GetUserWithRole()`, `GetByEmailOrPhone()`

- **`password.go`**: Interfaz de contraseñas
  - `GetActivePassword()`, `Create()`

- **`role.go`**: Interfaz de roles
  - `GetByKey()`, `GetAll()`

- **`one_time_password.go`**: Interfaz de OTP
  - `Create()`, `GetByCode()`, `Invalidate()`

- **`one_time_token.go`**: Interfaz de tokens
  - `Create()`, `GetByToken()`, `Invalidate()`

#### `/src/application/modules/`

Módulos de negocio organizados por dominio.

##### `/src/application/modules/auth/`

Módulo de autenticación:

- **`jwt_auth.go`**: Autenticación con JWT
  - Login con email/contraseña
  - Generación de tokens

- **`jwt_auth_refresh.go`**: Refresh de tokens
  - Renovación de access token

- **`jwt_auth_otp.go`**: Autenticación con OTP
  - Login con código OTP

- **`jwt_auth_user.go`**: Autenticación de usuario desde token
  - Validación de token
  - Obtención de usuario

- **`get_reset_password_token.go`**: Generación de token de reset
  - Creación de token
  - Envío de email

- **`get_reset_password_token_email.go`**: Envío de email de reset

- **`pipe/get_reset_password.go`**: Pipe para reset de contraseña
  - Orquesta generación de token y envío de email

##### `/src/application/modules/user/`

Módulo de usuarios:

- **`use_cases/create_user.go`**: Crear usuario
- **`use_cases/create_user_password.go`**: Crear usuario con contraseña
- **`use_cases/create_user_email.go`**: Enviar email de bienvenida
- **`use_cases/get_user.go`**: Obtener usuario
- **`use_cases/get_all_user.go`**: Listar usuarios (con cache)
- **`use_cases/update_user.go`**: Actualizar usuario
- **`use_cases/delete_user.go`**: Eliminar usuario
- **`use_cases/activate_user.go`**: Activar usuario
- **`pipes/create_user.go`**: Pipe para crear usuario y enviar email

##### `/src/application/modules/password/`

Módulo de contraseñas:

- **`use_cases/create_password.go`**: Crear contraseña
- **`use_cases/create_password_token.go`**: Crear token de reset
- **`pipes/create_password_token.go`**: Pipe para reset

##### `/src/application/modules/status/`

Módulo de estado:

- **`use_cases/status.go`**: Health check del sistema

#### `/src/application/shared/`

Componentes compartidos entre módulos.

##### `/src/application/shared/DTOs/`

Data Transfer Objects:

- **`base.go`**: DTOs base
  - `SingleResponse[T]`, `MultipleResponse[T]`

- **`user.go`**: DTOs de usuarios
  - `UserCreate`, `UserUpdate`, `UserAndPasswordCreate`

- **`password.go`**: DTOs de contraseñas
  - `PasswordCreate`, `PasswordReset`

- **`token.go`**: DTOs de tokens
  - `Token` (access + refresh)

- **`one_time_password.go`**: DTOs de OTP
- **`one_time_token.go`**: DTOs de tokens

##### `/src/application/shared/errors/`

Manejo de errores:

- **`application_error.go`**: Error de aplicación
  - `Code`: Código de estado
  - `Context`: Contexto del error
  - `ErrMsg`: Mensaje de error

##### `/src/application/shared/use_case/`

Base para casos de uso:

- **`base.go`**: Interfaz base y validación
  - `BaseUseCase[Input, Output]`
  - `BaseUseCaseValidation`

- **`use_case_result.go`**: Resultado de caso de uso
  - `UseCaseResult[T]`
  - `SetData()`, `SetError()`, `HasError()`

- **`dag.go`**: Directed Acyclic Graph
  - `DAG`: Ejecución secuencial
  - `UseCaseParallelDag`: Ejecución paralela

- **`uc_guards.go`**: Guards de autorización
  - Validación de permisos

##### `/src/application/shared/services/`

Servicios compartidos:

- **`create_one_time_password.go`**: Crear OTP
- **`create_one_time_token.go`**: Crear token
- **`create_password.go`**: Crear contraseña
- **`emails/`**: Servicios de email
  - `register_user_email.go`
  - `reset_password_email.go`
  - `otp_email.go`

##### `/src/application/shared/templates/`

Plantillas de email:

- **`emails/`**: Templates HTML
  - `register_user.gohtml`
  - `reset_password.gohtml`
  - `otp.gohtml`

##### `/src/application/shared/locales/`

Internacionalización:

- **`app_messages.go`**: Mensajes de la aplicación
- **`messages/`**: Mensajes por idioma
  - `en_us.go`, `es_es.go`

##### `/src/application/shared/context/`

Contexto de aplicación:

- **`context.go`**: Utilidades de contexto
- **`keys.go`**: Claves de contexto (UserKey, etc.)

##### `/src/application/shared/settings/`

Configuración:

- **`app_settings.go`**: Configuración de la aplicación
  - Carga desde variables de entorno
  - Validación de tipos

##### `/src/application/shared/guards/`

Guards de autorización:

- **`user.go`**: Guards de usuario
  - Validación de permisos

##### `/src/application/shared/defaults/`

Valores por defecto:

- **`user.go`**: Valores por defecto de usuarios
- **`roles.go`**: Roles por defecto
- **`password.go`**: Configuración de contraseñas

##### `/src/application/shared/mocks/`

Mocks para testing:

- **`dtos/`**: Mocks de DTOs
- **`mock_*.go`**: Mocks de providers y repositorios

### `/src/infrastructure/` - Capa de Infraestructura

#### `/src/infrastructure/server/`

Capa de servidor HTTP con Gin.

##### `/src/infrastructure/server/cmd/`

- **`main.go`**: Punto de entrada de la aplicación
  - Inicialización de infraestructura (`infrastructure.Initialize()`)
  - Configuración de Gin con graceful shutdown
  - Carga de middlewares (CORS, Recovery)
  - Carga de rutas (`routes.Router()`)
  - Inicio del servidor en puerto configurable
  - **Nota**: La documentación Swagger se ejecuta como un servicio independiente (ver `/src/infrastructure/docs/`)

**Flujo de inicialización:**
```go
1. infrastructure.Initialize()
   ├── Carga configuración (Settings)
   ├── Inicializa Logger
   ├── Conecta a PostgreSQL (GORM)
   ├── Configura JWT Provider
   ├── Configura Email Provider
   ├── Configura Cache Provider (Redis)
   └── Configura Email Services

2. buildGinApp()
   └── Crea aplicación Gin con graceful shutdown

3. loadGinApp()
   ├── Configura CORS
   ├── Configura Recovery middleware
   └── Carga rutas

4. app.Run()
   └── Inicia servidor HTTP

**Nota**: La documentación Swagger se ejecuta como un servicio independiente en `/src/infrastructure/docs/`
```

##### `/src/infrastructure/api/routes/`

- **`router.go`**: Definición de todas las rutas
  - Rutas públicas
  - Rutas privadas (con autenticación)
  - Agrupación por módulo

- **`wrap.go`**: Wrapper de handlers
  - Adapta Gin a `HandlerContext`
  - Extrae parámetros, query, locale

##### `/src/infrastructure/api/middlewares/`

- **`auth.go`**: Middleware de autenticación
  - Extrae token JWT
  - Valida token
  - Inyecta usuario en contexto

- **`query.go`**: Middleware de query params
  - Parsea filtros, ordenamiento, paginación

#### `/src/infrastructure/config/`

- **`config.go`**: Carga de configuración
  - Lee variables de entorno
  - Valida configuración

- **`env.go`**: Utilidades de variables de entorno

#### `/src/infrastructure/database/`

##### `/src/infrastructure/database/goprojectskeleton/`

- **`goprojectskeleton.go`**: Configuración de GORM
  - Conexión a PostgreSQL
  - Configuración de SSL

- **`models/`**: Modelos de base de datos (GORM)
  - `user.go`, `role.go`, `password.go`, etc.

- **`init_db/`**: Inicialización de BD
  - Migraciones automáticas
  - Seeds (si aplica)

#### `/src/infrastructure/handlers/`

Handlers HTTP (adaptadores de casos de uso a HTTP).

- **`user.go`**: Handlers de usuarios
  - `CreateUser()`, `GetUser()`, `UpdateUser()`, etc.

- **`auth.go`**: Handlers de autenticación
  - `Login()`, `RefreshAccessToken()`, `LoginOTP()`, etc.

- **`password.go`**: Handlers de contraseñas
  - `CreatePassword()`, `CreatePasswordToken()`

- **`status.go`**: Handler de estado
  - `GetHealthCheck()`

- **`request_resolver.go`**: Resolución de respuestas HTTP
  - Mapea `UseCaseResult` a respuesta HTTP
  - Maneja códigos de estado

- **`types.go`**: Tipos de handlers
  - `HandlerContext`
  - `Query`

#### `/src/infrastructure/providers/`

Implementaciones de proveedores.

- **`jwt_provider.go`**: Implementación de JWT
  - Generación y validación de tokens

- **`hash_provider.go`**: Implementación de hashing
  - Bcrypt para contraseñas

- **`email_provider.go`**: Implementación de email
  - SMTP

- **`cache_provider.go`**: Implementación de cache
  - Redis

- **`logger_provider.go`**: Implementación de logging
  - Logging estructurado

- **`renderer_provider.go`**: Implementación de renderizado
  - Renderizado de templates HTML

- **`status_provider.go`**: Implementación de estado

#### `/src/infrastructure/otel/`

Implementación de OpenTelemetry para observabilidad.

- **`otel_init.go`**: Inicialización del SDK de OpenTelemetry
  - Configuración de exportadores (OTLP, Prometheus)
  - Configuración de muestreo de trazas
  - Configuración de proveedores de métricas

- **`otel_tracer.go`**: Implementación del tracer con OpenTelemetry
  - Implementa `ITracer`
  - Gestión de spans con OpenTelemetry SDK
  - Propagación de contexto W3C

- **`otel_span.go`**: Implementación de spans con OpenTelemetry
  - Implementa `ISpan`
  - Wrapper sobre spans de OpenTelemetry
  - Conversión de atributos y estados

- **`otel_metrics.go`**: Implementación de métricas con OpenTelemetry
  - Implementa `IMetricsCollector`
  - Contadores, histogramas y gauges
  - Integración con Prometheus

- **`otel_logger.go`**: Logger con correlación de trazas
  - Implementa `ILoggerProvider` extendido
  - Inyección automática de trace_id y span_id
  - Logs estructurados en formato JSON

#### `/src/infrastructure/repositories/`

Implementaciones de repositorios.

- **`base.go`**: Repositorio base
  - `RepositoryBase`: Implementación genérica de CRUD
  - `ModelConverter`: Conversión entre DTOs y modelos

- **`user.go`**: Repositorio de usuarios
  - Implementa `IUserRepository`
  - Métodos específicos: `CreateWithPassword()`, `GetUserWithRole()`

- **`password.go`**: Repositorio de contraseñas
- **`role.go`**: Repositorio de roles
- **`one_time_password.go`**: Repositorio de OTP
- **`one_time_token.go`**: Repositorio de tokens

- **`model_converter.go`**: Convertidores de modelos
- **`orm_error_map.go`**: Mapeo de errores de ORM

#### `/src/infrastructure/container.go`

Inicialización centralizada de infraestructura:

- **Configuración**: Carga desde variables de entorno
- **Base de datos**: Conexión a PostgreSQL con GORM
- **Providers**: JWT, Email, Cache, Logger
- **Servicios**: Email services (registro, reset, OTP)

**Orden de inicialización:**
1. Settings (configuración de aplicación)
2. Logger (sistema de logging)
3. Database (conexión PostgreSQL)
4. JWT Provider (autenticación)
5. Email Provider (SMTP)
6. Cache Provider (Redis)
7. Email Services (servicios de email)

#### `/src/infrastructure/clouds/`

Adaptadores para plataformas cloud y serverless.

##### `/src/infrastructure/clouds/aws/`

Implementación para **AWS Lambda**:

- **`init.go`**: Inicialización de infraestructura AWS
  - Carga configuración desde AWS Secrets Manager
  - Inicializa base de datos, providers y servicios
  - Optimizado para cold starts de Lambda

- **`lambda_adapter.go`**: Adaptador para eventos Lambda
  - Convierte eventos API Gateway a `HandlerContext`
  - Maneja respuestas HTTP
  - Gestiona errores y códigos de estado

- **`secrets_manager.go`**: Integración con AWS Secrets Manager
  - Carga automática de secretos desde ARNs
  - Fallback a variables de entorno
  - Cache de secretos para optimización

- **`middleware.go`**: Middlewares específicos para Lambda
  - Autenticación JWT
  - Logging estructurado
  - Manejo de errores

- **`render_provider.go`**: Provider de renderizado para AWS
  - Renderizado de templates HTML
  - Optimizado para entorno serverless

- **`functions/`**: Generador de funciones Lambda
  - **`main.go`**: Herramienta CLI para generar y desplegar
  - **`functions.json`**: Definición de funciones Lambda
  - **`utils/generate.go`**: Generación de código desde templates
  - **`utils/deploy.go`**: Despliegue a AWS Lambda

- **`terraform/`**: Infraestructura como código
  - Definición de funciones Lambda
  - API Gateway
  - IAM roles y políticas
  - Variables y outputs

##### `/src/infrastructure/clouds/azure/`

Implementación para **Azure Functions**:

- **`init.go`**: Inicialización de infraestructura Azure
  - Carga configuración desde Azure Key Vault
  - Inicializa base de datos, providers y servicios
  - Optimizado para Azure Functions

- **`http_adapter.go`**: Adaptador para HTTP triggers
  - Convierte requests HTTP a `HandlerContext`
  - Maneja respuestas HTTP
  - Gestiona errores y códigos de estado

- **`vault.go`**: Integración con Azure Key Vault
  - Carga automática de secretos
  - Autenticación con Managed Identity
  - Cache de secretos

- **`middleware.go`**: Middlewares específicos para Azure
  - Autenticación JWT
  - Logging estructurado
  - Manejo de errores

- **`functions/`**: Generador de funciones Azure
  - **`functions.json`**: Definición de funciones
  - **`generate.go`**: Generación de código
  - Templates para Azure Functions

- **`terraform/`**: Infraestructura como código
  - Definición de Function Apps
  - App Service Plans
  - Key Vault integration
  - Variables y outputs

### `/docker/` - Configuración Docker

- **`docker-compose.dev.yml`**: Servicios de desarrollo (incluye servicio Swagger independiente)
- **`docker-compose.test.yml`**: Servicios de testing
- **`docker-compose.e2e.yml`**: Servicios de E2E (incluye Bruno para tests E2E automatizados)
- **`dockerfile.dev`**: Dockerfile de desarrollo
- **`dockerfile.swagger`**: Dockerfile de producción para servicio Swagger independiente
- **`dockerfile.swagger.debug`**: Dockerfile de desarrollo para servicio Swagger con hot reload
- **`dockerfile.e2e`**: Dockerfile de E2E
- **`dockerfile.integration`**: Dockerfile de integración
- **`db/`**: Configuración de base de datos
  - `Dockerfile`, `create.sql`
- **`grafana/`**: Configuración de Grafana
  - `provisioning/datasources/`: Configuración de fuentes de datos (Prometheus, Jaeger)
  - `dashboards/`: Dashboards JSON preconfigurados
- **`otel/`**: Configuración de OpenTelemetry
  - `otel-collector-config.yaml`: Configuración del colector OTLP
- **`prometheus/`**: Configuración de Prometheus
  - `prometheus.yml`: Configuración de scraping de métricas

### `/tests/` - Tests

- **`integration/`**: Tests de integración
  - `main_test.go`: Setup de tests
  - Tests de repositorios
  - Tests de providers

- **`e2e/`**: Tests end-to-end
  - `bruno/`: Colección Bruno con tests
  - `collections/`: Colecciones de requests de Bruno
  - `environments/`: Configuraciones de entornos de Bruno

### `/src/infrastructure/docs/` - Servicio Swagger Independiente

Servicio HTTP independiente para documentación Swagger.

- **`main.go`**: Servidor HTTP independiente para Swagger UI
  - Configuración desde variables de entorno
  - Servidor en puerto configurable (por defecto: 8081)
  - Redirección automática a `/docs/`
- **`config/`**: Configuración del servidor Swagger
  - `config.go`: Carga de variables de entorno
- **`swagger/`**: Archivos generados de Swagger
  - `swagger.json`: Especificación Swagger (JSON)
  - `swagger.yaml`: Especificación Swagger (YAML)
  - `docs.go`: Código generado de Swagger
- **`go.mod`**: Módulo independiente para el servicio Swagger

---

## Tecnologías y Dependencias

### Stack Tecnológico

#### 🚀 Lenguaje y Runtime
- **Go 1.25**: Lenguaje de programación
  - Concurrencia nativa con goroutines
  - Compilación estática
  - Tree shaking automático
  - Excelente rendimiento

#### 🌐 Framework Web
- **Gin v1.10.0**: Framework HTTP minimalista y rápido
  - Router de alto rendimiento
  - Middleware chain
  - JSON binding automático
  - Validación de requests

- **gin-contrib/cors**: Middleware CORS
  - Configuración flexible de orígenes
  - Soporte para credenciales
  - Headers personalizables

- **gin-contrib/graceful**: Shutdown graceful
  - Cierre ordenado del servidor
  - Finalización de requests en curso
  - Timeout configurable

#### 💾 Base de Datos y Persistencia
- **GORM v1.25.12**: ORM para Go
  - Migraciones automáticas
  - Hooks y callbacks
  - Relaciones y asociaciones
  - Query builder fluido

- **PostgreSQL (pgx/v5)**: Driver de PostgreSQL
  - Connection pooling nativo
  - Transacciones
  - Prepared statements
  - Soporte para tipos avanzados

- **Redis (go-redis/v9)**: Cliente Redis para cache
  - Operaciones atómicas
  - Pub/Sub
  - Pipeline support
  - Cluster support

#### 🔐 Autenticación y Seguridad
- **golang-jwt/jwt/v5**: Implementación de JWT
  - Generación y validación de tokens
  - Múltiples algoritmos (HS256, RS256, etc.)
  - Claims personalizados
  - Expiración y validación de tiempo

- **golang.org/x/crypto**: Utilidades criptográficas
  - Bcrypt para hash de contraseñas
  - Salt automático
  - Cost configurable

#### 📚 Documentación y Testing
- **swaggo/swag**: Generación de documentación Swagger
  - Anotaciones en código Go
  - Generación automática de OpenAPI
  - Validación de esquemas

- **swaggo/http-swagger**: Servidor HTTP independiente para Swagger UI
  - Servicio completamente independiente
  - UI interactiva
  - Pruebas desde navegador
  - Autenticación en Swagger UI
  - Despliegue independiente en Docker registry

- **stretchr/testify**: Framework de testing
  - Assertions mejoradas
  - Mocks y suites
  - Test helpers

#### 🛠️ Utilidades
- **joho/godotenv**: Carga de variables de entorno
  - Soporte para `.env` files
  - Override de variables
  - Validación de variables requeridas

### Dependencias Principales

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

### Arquitectura de Dependencias

```mermaid
graph TB
    subgraph Core["Núcleo"]
        Go[Go 1.25]
    end

    subgraph Web["Capa Web"]
        Gin[Gin Framework]
        CORS[CORS Middleware]
        Graceful[Graceful Shutdown]
    end

    subgraph Data["Capa de Datos"]
        GORM[GORM ORM]
        PostgreSQL[PostgreSQL Driver]
        Redis[Redis Client]
    end

    subgraph Auth["Autenticación"]
        JWT[JWT Library]
        Crypto[Crypto Utils]
    end

    subgraph Docs["Documentación"]
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
    Gin --> Swag
    Swag --> SwaggerUI

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
    github.com/swaggo/http-swagger v1.6.0
    gorm.io/gorm v1.25.12
    golang.org/x/crypto v0.41.0
)
```

---

## Configuración y Setup

### Variables de Entorno

El proyecto utiliza variables de entorno para toda la configuración:

```bash
# Aplicación
APP_NAME=goprojectskeleton
APP_ENV=development
APP_PORT=8080
APP_VERSION=0.0.1
APP_DESCRIPTION=Go Project Skeleton
ENABLE_LOG=true
DEBUG_LOG=true

# Base de Datos
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

# Tokens y OTP
ONE_TIME_TOKEN_TTL=15
ONE_TIME_TOKEN_EMAIL_VERIFY_TTL=60
ONE_TIME_PASSWORD_LENGTH=6
ONE_TIME_PASSWORD_TTL=10
FRONTEND_RESET_PASSWORD_URL=http://localhost:3000/reset-password
FRONTEND_ACTIVATE_ACCOUNT_URL=http://localhost:3000/activate-account
```

### Instalación

1. **Clonar repositorio**
```bash
git clone <repository-url>
cd GoProjectSkeleton
```

2. **Configurar variables de entorno**
```bash
cp dev.env.example dev.env
# Editar dev.env con tus configuraciones
```

3. **Instalar dependencias**
```bash
go mod download
```

4. **Ejecutar con Docker**
```bash
# Crear red
docker network create goprojectskeleton

# Crear volumen
docker volume create goprojectskeleton-db-data

# Ejecutar servicios
docker-compose -f docker/docker-compose.dev.yml up -d
```

5. **Ejecutar aplicación**
```bash
go run src/infrastructure/api/cmd/main.go
```

---

## Módulos de Negocio

### 🔐 Módulo de Autenticación (`auth`)

**Responsabilidad**: Gestión completa de autenticación y autorización.

#### Funcionalidades

- ✅ **Login con Email/Contraseña** - Autenticación tradicional
- ✅ **Login con OTP** - Autenticación de dos factores
- ✅ **Refresh de Tokens** - Renovación de access tokens
- ✅ **Reset de Contraseña** - Recuperación mediante tokens
- ✅ **Validación de Usuario** - Verificación desde JWT token

#### Casos de Uso Detallados

**`JwtAuthUseCase`** - Autenticación principal
```go
// Flujo:
// 1. Valida credenciales (email/phone + password)
// 2. Verifica contraseña con hash
// 3. Si OTP activado → genera y envía OTP
// 4. Si OTP desactivado → genera tokens JWT
// 5. Retorna tokens o indica que se envió OTP
```

**`JwtAuthRefreshUseCase`** - Renovación de tokens
```go
// Flujo:
// 1. Valida refresh token
// 2. Verifica expiración y firma
// 3. Genera nuevo access token
// 4. Retorna nuevo token
```

**`JwtAuthOtpUseCase`** - Autenticación con OTP
```go
// Flujo:
// 1. Valida código OTP
// 2. Verifica expiración y uso
// 3. Invalida OTP usado
// 4. Genera tokens JWT
// 5. Retorna tokens
```

**`GetResetPasswordTokenUseCase`** - Generación de token de reset
```go
// Flujo:
// 1. Busca usuario por email/phone
// 2. Genera token único
// 3. Crea registro en BD con expiración
// 4. Envía email con link de reset
```

**`JwtAuthUserUseCase`** - Validación de usuario desde token
```go
// Flujo:
// 1. Extrae token del contexto
// 2. Valida y parsea token
// 3. Busca usuario en BD
// 4. Retorna usuario con rol
```

#### Pipes

**`GetResetPasswordPipe`** - Pipe para reset de contraseña
- Orquesta generación de token y envío de email
- Ejecución secuencial con manejo de errores

### 👥 Módulo de Usuarios (`user`)

**Responsabilidad**: Gestión completa del ciclo de vida de usuarios.

#### Funcionalidades

- ✅ **CRUD Completo** - Crear, leer, actualizar, eliminar
- ✅ **Activación de Cuentas** - Activación mediante tokens
- ✅ **Gestión de Roles** - Asignación y validación de roles
- ✅ **Paginación y Filtrado** - Consultas eficientes
- ✅ **Cache Inteligente** - Cache de listados con Redis
- ✅ **Emails Transaccionales** - Bienvenida y reactivación

#### Casos de Uso Detallados

**`CreateUserUseCase`** - Crear usuario básico
```go
// Flujo:
// 1. Valida datos de entrada
// 2. Verifica que email/phone no existan
// 3. Crea usuario con estado "pending"
// 4. Retorna usuario creado
```

**`CreateUserAndPasswordUseCase`** - Crear usuario con contraseña
```go
// Flujo:
// 1. Valida datos de usuario y contraseña
// 2. Hash de contraseña con Bcrypt
// 3. Crea usuario y contraseña en transacción
// 4. Retorna usuario creado
```

**`CreateUserSendEmailUseCase`** - Enviar email de bienvenida
```go
// Flujo:
// 1. Renderiza template de email
// 2. Envía email con datos del usuario
// 3. Maneja errores de envío
```

**`GetUserUseCase`** - Obtener usuario por ID
```go
// Flujo:
// 1. Valida ID
// 2. Busca usuario en BD
// 3. Incluye información de rol
// 4. Retorna usuario con rol
```

**`GetAllUserUseCase`** - Listar usuarios con filtros
```go
// Flujo:
// 1. Verifica cache (Redis)
// 2. Si cache hit → retorna desde cache
// 3. Si cache miss → consulta BD con filtros
// 4. Aplica paginación y ordenamiento
// 5. Guarda en cache con TTL
// 6. Retorna lista paginada
```

**`UpdateUserUseCase`** - Actualizar usuario
```go
// Flujo:
// 1. Valida datos de actualización
// 2. Busca usuario existente
// 3. Aplica cambios parciales (PATCH)
// 4. Valida reglas de negocio
// 5. Actualiza en BD
// 6. Retorna usuario actualizado
```

**`DeleteUserUseCase`** - Eliminar usuario (soft delete)
```go
// Flujo:
// 1. Busca usuario
// 2. Verifica permisos (no eliminar admin)
// 3. Soft delete (marca como deleted)
// 4. Invalida cache relacionado
```

**`ActivateUserUseCase`** - Activar cuenta de usuario
```go
// Flujo:
// 1. Valida token de activación
// 2. Verifica expiración
// 3. Cambia estado a "active"
// 4. Invalida token usado
```

**`ResendWelcomeEmailUseCase`** - Reenviar email de bienvenida
```go
// Flujo:
// 1. Busca usuario
// 2. Renderiza template
// 3. Envía email
```

#### Pipes

**`CreateUserPipe`** - Pipe para crear usuario completo
- Ejecuta secuencialmente:
  1. `CreateUserAndPasswordUseCase` → Crea usuario con contraseña
  2. `CreateUserSendEmailUseCase` → Envía email de bienvenida
- Manejo de errores: Si falla cualquier paso, se detiene la ejecución

### 🔑 Módulo de Contraseñas (`password`)

**Responsabilidad**: Gestión segura de contraseñas.

#### Funcionalidades

- ✅ **Creación de Contraseñas** - Hash seguro con Bcrypt
- ✅ **Generación de Tokens de Reset** - Tokens únicos con expiración
- ✅ **Validación de Fortaleza** - Reglas de contraseña segura
- ✅ **Expiración de Contraseñas** - Contraseñas temporales

#### Casos de Uso Detallados

**`CreatePasswordUseCase`** - Crear nueva contraseña
```go
// Flujo:
// 1. Valida fortaleza de contraseña
// 2. Hash con Bcrypt
// 3. Desactiva contraseñas anteriores
// 4. Crea nueva contraseña activa
// 5. Configura expiración si aplica
```

**`CreatePasswordTokenUseCase`** - Crear token de reset
```go
// Flujo:
// 1. Busca usuario por email/phone
// 2. Genera token único
// 3. Crea registro con expiración
// 4. Envía email con link de reset
```

#### Pipes

**`CreatePasswordTokenPipe`** - Pipe para reset de contraseña
- Orquesta creación de token y envío de email

### 📊 Módulo de Estado (`status`)

**Responsabilidad**: Monitoreo y salud del sistema.

#### Funcionalidades

- ✅ **Health Check** - Estado general del sistema
- ✅ **Verificación de Servicios** - BD, Redis, etc.
- ✅ **Información de Versión** - Versión de la aplicación

#### Casos de Uso Detallados

**`GetStatusUseCase`** - Obtener estado del sistema
```go
// Flujo:
// 1. Verifica conexión a PostgreSQL
// 2. Verifica conexión a Redis
// 3. Verifica providers (JWT, Email)
// 4. Retorna estado consolidado
// 5. Incluye información de versión
```

### 📈 Estadísticas por Módulo

| Módulo | Casos de Uso | Pipes | Tests | Endpoints |
|--------|--------------|-------|-------|-----------|
| **auth** | 5 | 1 | 5+ | 4 |
| **user** | 9 | 1 | 9+ | 7 |
| **password** | 2 | 1 | 2+ | 2 |
| **status** | 1 | 0 | 1+ | 1 |
| **Total** | **17** | **3** | **17+** | **14** |

---

## API y Endpoints

### Autenticación

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| POST | `/api/auth/login` | Login con credenciales | No |
| POST | `/api/auth/refresh` | Renovar token de acceso | No |
| GET | `/api/auth/login-otp/{otp}` | Login con OTP | No |
| GET | `/api/auth/password-reset/{identifier}` | Solicitar reset de contraseña | No |

### Usuarios

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| POST | `/api/user` | Crear usuario | No |
| GET | `/api/user/{id}` | Obtener usuario | Sí |
| PATCH | `/api/user/{id}` | Actualizar usuario | Sí |
| DELETE | `/api/user/{id}` | Eliminar usuario | Sí |
| GET | `/api/user` | Listar usuarios (con filtros) | Sí |
| POST | `/api/user-password` | Crear usuario con contraseña | No |
| POST | `/api/user/activate` | Activar usuario | No |

### Contraseñas

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| POST | `/api/password` | Crear contraseña | Sí |
| POST | `/api/password/reset-token` | Crear token de reset | No |

### Sistema

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| GET | `/api/health-check` | Health check | No |

**Nota**: La documentación Swagger está disponible como un servicio independiente en el puerto 8081 en `http://localhost:8081/docs/`

### Ejemplos de Uso

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Crear Usuario
```bash
curl -X POST http://localhost:8080/api/user \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "phone": "+1234567890",
    "role_id": 2,
    "status": "pending",
    "otpLogin": false
  }'
```

---

## Base de Datos y Persistencia

### Diagrama Entidad-Relación (ERD)

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

### Modelos Principales

#### User (Usuario)
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

#### Password (Contraseña)
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

#### Role (Rol)
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

### Diagrama de Relaciones

```mermaid
graph LR
    subgraph UserDomain["👤 Usuario"]
        User[User<br/>ID, Name, Email, Phone<br/>Status, RoleID, OTPLogin]
    end

    subgraph AuthDomain["🔐 Autenticación"]
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

### Migraciones

Las migraciones se ejecutan automáticamente al iniciar la aplicación usando GORM AutoMigrate.

### Cache

El sistema utiliza Redis para cachear consultas frecuentes, especialmente en listados de usuarios.

```mermaid
graph TB
    subgraph App["Aplicación"]
        Handler[Handler]
        UseCase[Use Case]
        Cache[Cache Provider]
    end

    subgraph Storage["Almacenamiento"]
        Redis[(Redis<br/>Cache)]
        DB[(PostgreSQL<br/>Base de Datos)]
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

## Autenticación y Seguridad

### Diagrama de Flujo de Autenticación

```mermaid
sequenceDiagram
    participant Client as Cliente
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

### Flujo de Autenticación con OTP

```mermaid
sequenceDiagram
    participant Client as Cliente
    participant API as API
    participant AuthUC as Auth Use Case
    participant OTPUC as OTP Use Case
    participant OTPRepo as OTP Repository
    participant EmailSvc as Email Service
    participant DB as PostgreSQL
    participant SMTP as SMTP

    Client->>API: POST /api/auth/login<br/>{email, password}
    API->>AuthUC: Execute(credentials)
    AuthUC->>AuthUC: Valida credenciales
    AuthUC->>AuthUC: ¿OTP Login activado?

    alt OTP Login activado
        AuthUC->>OTPUC: GenerateOTP()
        OTPUC->>OTPRepo: Create()
        OTPRepo->>DB: INSERT OTP
        OTPUC->>EmailSvc: SendOTPEmail()
        EmailSvc->>SMTP: Enviar email
        AuthUC-->>API: 204 No Content
        API-->>Client: OTP enviado por email
    else OTP Login desactivado
        AuthUC->>JWT: GenerateTokens()
        AuthUC-->>API: Tokens
        API-->>Client: {accessToken, refreshToken}
    end

    Note over Client,SMTP: Usuario ingresa OTP
    Client->>API: GET /api/auth/login-otp/{otp}
    API->>OTPUC: ValidateOTP(otp)
    OTPUC->>OTPRepo: GetByCode()
    OTPRepo->>DB: SELECT OTP
    OTPUC->>OTPUC: Valida expiración
    OTPUC->>JWT: GenerateTokens()
    OTPUC-->>API: Tokens
    API-->>Client: {accessToken, refreshToken}
```

### JWT (JSON Web Tokens)

El sistema utiliza JWT para autenticación con dos tipos de tokens:

```mermaid
graph TB
    subgraph TokenFlow["Flujo de Tokens"]
        Login[Login] --> AccessToken[Access Token<br/>TTL: 1 hora]
        Login --> RefreshToken[Refresh Token<br/>TTL: 24 horas]

        AccessToken -->|Expira| Refresh[Refresh Endpoint]
        Refresh --> NewAccess[New Access Token]

        AccessToken -->|Válido| Protected[Protected Resources]
    end

    subgraph TokenStructure["Estructura JWT"]
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

1. **Access Token**: Token de acceso (TTL: 1 hora por defecto)
2. **Refresh Token**: Token de renovación (TTL: 24 horas por defecto)

### OTP (One-Time Password)

Sistema de autenticación de dos factores:

```mermaid
stateDiagram-v2
    [*] --> Login: Credenciales
    Login --> CheckOTP: Validar
    CheckOTP --> GenerateOTP: OTP activado
    CheckOTP --> GenerateJWT: OTP desactivado

    GenerateOTP --> SendEmail: Código generado
    SendEmail --> WaitOTP: Email enviado
    WaitOTP --> ValidateOTP: Usuario ingresa código
    ValidateOTP --> GenerateJWT: Código válido
    ValidateOTP --> Expired: Código expirado
    ValidateOTP --> Invalid: Código inválido

    GenerateJWT --> [*]: Tokens generados
    Expired --> [*]
    Invalid --> [*]
```

- Generación automática de códigos OTP
- Envío por email
- TTL configurable (por defecto: 10 minutos)

### Seguridad

```mermaid
graph TB
    subgraph SecurityLayers["Capas de Seguridad"]
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

- **Hash de contraseñas**: Bcrypt con salt automático
- **Validación de entrada**: En todos los endpoints
- **Middleware de autenticación**: Para rutas protegidas
- **CORS**: Configurado para seguridad web
- **Sanitización**: Prevención de inyecciones

---

## Testing

### Estructura de Tests

```
tests/
├── e2e/                    # Tests end-to-end
│   ├── bruno/              # Colección Bruno
│   │   ├── collections/    # Colecciones de requests
│   │   ├── environments/  # Configuraciones de entornos
│   │   └── bruno.json      # Configuración de Bruno
│   └── postman/            # Colección Postman legacy (opcional)
└── integration/            # Tests de integración
    ├── main_test.go        # Setup de tests
    ├── user_repository_test.go
    ├── password_repository_test.go
    └── email_provider_test.go
```

### Tipos de Testing

1. **Unit Tests**: Tests de unidades individuales (en cada módulo)
2. **Integration Tests**: Tests de integración con base de datos
3. **E2E Tests**: Tests end-to-end con Bruno

### Ejecutar Tests

```bash
# Tests unitarios
go test ./src/...

# Tests de integración
go test ./tests/integration/...

# Tests E2E (requiere servicios corriendo)
# Los tests de Bruno se ejecutan automáticamente en docker-compose.e2e.yml
# O ejecutar manualmente: bruno run --env=dev
```

---

## Docker y Despliegue

### Arquitectura Docker

```mermaid
graph TB
    subgraph DockerNetwork["Docker Network: goprojectskeleton"]
        subgraph AppContainer["Aplicación Go"]
            App[Go Application<br/>Port: 8080<br/>Hot Reload]
        end

        subgraph DBServices["Servicios de Datos"]
            PostgreSQL[(PostgreSQL<br/>Port: 5432<br/>Volume: db-data)]
            Redis[(Redis<br/>Port: 6379<br/>Cache)]
        end

        subgraph SwaggerService["Servicio Swagger Independiente"]
            Swagger[Swagger Server<br/>Port: 8081<br/>Independent]
        end

        subgraph Observability["📊 Observabilidad"]
            OTELCollector[OTEL Collector<br/>Port: 4317/4318]
            Prometheus[(Prometheus<br/>Port: 9090)]
            Jaeger[Jaeger<br/>Port: 16686]
            Grafana[Grafana<br/>Port: 3001]
        end

        subgraph DevTools["Herramientas de Desarrollo"]
            Mailpit[Mailpit<br/>Port: 8025<br/>Email Testing]
            RedisCommander[Redis Commander<br/>Port: 18081<br/>Redis UI]
        end
    end

    App -->|GORM| PostgreSQL
    App -->|go-redis| Redis
    App -->|SMTP| Mailpit
    App -->|OTLP| OTELCollector

    OTELCollector --> Prometheus
    OTELCollector --> Jaeger
    Prometheus --> Grafana
    Jaeger --> Grafana

    Swagger -.->|Documentation| App
    RedisCommander -->|UI| Redis

    style App fill:#e3f2fd
    style PostgreSQL fill:#c8e6c9
    style Redis fill:#ffcdd2
    style Mailpit fill:#fff9c4
    style RedisCommander fill:#f3e5f5
    style OTELCollector fill:#fff9c4
    style Prometheus fill:#ffcdd2
    style Jaeger fill:#e3f2fd
    style Grafana fill:#c8e6c9
```

### Diagrama de Despliegue

```mermaid
graph TB
    subgraph Dev["🛠️ Desarrollo"]
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

    subgraph Prod["🚀 Producción"]
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

### Servicios Docker

El proyecto incluye configuración Docker para desarrollo:

- **Aplicación**: Servidor Go con hot reload (puerto 8080)
- **Swagger**: Servicio independiente de documentación (puerto 8081)
- **PostgreSQL**: Base de datos principal
- **Redis**: Cache y sesiones
- **Mailpit**: Servidor de email para desarrollo
- **Redis Commander**: Interfaz web para Redis (puerto 18081)

**Servicios de Observabilidad**:
- **OTEL Collector**: Colector OpenTelemetry (puertos 4317/4318)
- **Prometheus**: Almacenamiento de métricas (puerto 9090)
- **Jaeger**: Backend de trazas distribuidas (puerto 16686)
- **Grafana**: Dashboards y visualización (puerto 3001)

**Servicios de Testing E2E** (docker-compose.e2e.yml):
- **Aplicación**: Servidor Go para tests E2E
- **PostgreSQL**: Base de datos de pruebas
- **Redis**: Cache de pruebas
- **Mailpit**: Testing de email
- **Bruno**: Ejecutor automatizado de tests E2E (ejecuta tests automáticamente)

### Comandos Docker

```bash
# Desarrollo
docker-compose -f docker/docker-compose.dev.yml up -d

# Testing
docker-compose -f docker/docker-compose.test.yml up -d

# E2E Testing (Bruno se ejecuta automáticamente)
docker-compose -f docker/docker-compose.e2e.yml up -d
# Los tests de Bruno se ejecutarán automáticamente una vez que la API esté lista
```

### Despliegue Independiente de Swagger

El servicio Swagger puede desplegarse completamente independiente de la aplicación principal:

#### Construcción de la Imagen

```bash
# Construir imagen de Swagger
docker build -f docker/dockerfile.swagger -t your-registry/swagger:latest .

# Etiquetar para versión específica
docker tag your-registry/swagger:latest your-registry/swagger:v1.0.0
```

#### Push a Docker Registry

```bash
# Push a Docker Hub
docker push your-registry/swagger:latest

# Push a otros registries (ej: AWS ECR, Google GCR, Azure ACR)
# AWS ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker tag your-registry/swagger:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/swagger:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/swagger:latest
```

#### Ejecución del Contenedor

```bash
# Ejecutar con variables de entorno
docker run -d \
  -p 8081:8081 \
  -e SWAGGER_PORT=8081 \
  -e API_HOST=api.example.com:8080 \
  -e API_TITLE="Mi Documentación API" \
  -e API_VERSION="1.0" \
  -e API_DESCRIPTION="Documentación de la API" \
  -e API_BASE_PATH="/api" \
  --name swagger-docs \
  your-registry/swagger:latest

# O usando un archivo .env
docker run -d \
  -p 8081:8081 \
  --env-file swagger.env \
  --name swagger-docs \
  your-registry/swagger:latest
```

#### Despliegue en Kubernetes

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
          value: "Mi Documentación API"
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

#### Ventajas del Despliegue Independiente

- ✅ **Escalabilidad independiente**: Escalar Swagger sin afectar la aplicación principal
- ✅ **Actualizaciones sin downtime**: Actualizar documentación sin reiniciar la API
- ✅ **Separación de responsabilidades**: Documentación separada de la lógica de negocio
- ✅ **Diferentes entornos**: Diferentes versiones de documentación para dev/staging/prod
- ✅ **CDN y caching**: Servir documentación desde CDN para mejor rendimiento

---

## Despliegue con GitHub Actions

**GoProjectSkeleton** incluye un flujo de trabajo completo de GitHub Actions para el despliegue automatizado a plataformas cloud de AWS y Azure. El flujo de trabajo soporta el aprovisionamiento de infraestructura con Terraform y el despliegue automatizado de funciones.

### Resumen

El flujo de trabajo de despliegue (`deploy.yml`) proporciona:

- ✅ **Soporte multi-cloud**: Desplegar a AWS o Azure
- ✅ **Gestión de entornos**: Despliegues separados para desarrollo, staging y producción
- ✅ **Integración con Terraform**: Infraestructura como Código con acciones plan, apply y destroy
- ✅ **Despliegue automatizado de funciones**: Desplegar Lambda/Functions después de la infraestructura (solo AWS)
- ✅ **Artefactos de plan**: Subir planes de Terraform para revisión
- ✅ **Verificaciones de seguridad**: Prevenir destrucción accidental de producción

### Flujo de Despliegue

```mermaid
graph TB
    subgraph Trigger["🚀 Activación del Flujo"]
        Manual[Ejecución Manual<br/>GitHub Actions UI]
    end

    subgraph Inputs["📋 Entradas del Flujo"]
        Cloud[Proveedor Cloud<br/>AWS o Azure]
        Env[Entorno<br/>dev/staging/prod]
        Action[Acción Terraform<br/>plan/apply/destroy]
        DeployFunc[Desplegar Funciones<br/>Solo AWS]
    end

    subgraph Setup["⚙️ Fase de Configuración"]
        Checkout[Checkout Código]
        GoSetup[Configurar Go 1.25.5]
        TfSetup[Instalar Terraform 1.14.1]
        Creds[Configurar Credenciales Cloud]
    end

    subgraph Build["🔨 Fase de Construcción"]
        Deps[Descargar Dependencias]
        GenAWS[Generar Funciones AWS<br/>si AWS]
        GenAzure[Generar Funciones Azure<br/>si Azure]
        Tfvars[Crear terraform.tfvars]
    end

    subgraph Terraform["🏗️ Fase Terraform"]
        Init[Terraform Init]
        Validate[Terraform Validate]
        Plan[Terraform Plan<br/>si no es destroy]
        UploadPlan[Subir Artefacto Plan<br/>si acción plan]
        Apply[Terraform Apply<br/>si acción apply]
        Destroy[Terraform Destroy<br/>si acción destroy]
    end

    subgraph Deploy["📦 Fase de Despliegue"]
        DeployLambda[Desplegar Funciones Lambda<br/>Solo AWS]
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

### Entradas del Flujo de Trabajo

Al activar el flujo de trabajo manualmente, se te solicitará:

| Entrada | Descripción | Opciones | Por Defecto |
|---------|-------------|----------|-------------|
| `cloud` | Proveedor cloud | `aws`, `azure` | Requerido |
| `environment` | Entorno objetivo | `development`, `staging`, `production` | `development` |
| `terraform_action` | Operación Terraform | `plan`, `apply`, `destroy` | `apply` |
| `deploy_functions` | Desplegar funciones después de Terraform (solo AWS) | `true`, `false` | `true` |

### Configuración de Secretos de GitHub

Antes de usar el flujo de trabajo de despliegue, necesitas configurar los Secretos de GitHub. Los secretos son específicos del entorno y pueden configurarse a nivel de repositorio o de entorno.

#### Configuración de Secretos de GitHub

1. **Navegar a Configuración del Repositorio**
   - Ve a tu repositorio de GitHub
   - Haz clic en **Settings** → **Secrets and variables** → **Actions**

2. **Crear Secretos de Entorno (Recomendado)**
   - Haz clic en **Environments** en la barra lateral izquierda
   - Crea entornos: `development`, `staging`, `production`
   - Agrega secretos a cada entorno según sea necesario

3. **Crear Secretos de Repositorio (Alternativa)**
   - Agrega secretos a nivel de repositorio (disponibles para todos los entornos)

#### Secretos Requeridos

##### Secretos AWS

| Nombre del Secreto | Descripción | Ejemplo | Requerido Para |
|-------------------|-------------|---------|----------------|
| `AWS_ACCESS_KEY_ID` | ID de clave de acceso AWS | `AKIAIOSFODNN7EXAMPLE` | Despliegues AWS |
| `AWS_SECRET_ACCESS_KEY` | Clave secreta de acceso AWS | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` | Despliegues AWS |
| `AWS_REGION` | Región AWS | `us-east-1` | Despliegues AWS (opcional, por defecto `us-east-1`) |
| `PROJECT_NAME` | Nombre del proyecto para nombrar recursos | `go-project-skeleton` | Despliegues AWS (opcional, por defecto `go-project-skeleton`) |
| `TFVARS` | Contenido del archivo de variables de Terraform | Ver abajo | Todos los despliegues AWS |

##### Secretos Azure

| Nombre del Secreto | Descripción | Ejemplo | Requerido Para |
|-------------------|-------------|---------|----------------|
| `AZURE_CREDENTIALS` | JSON de entidad de servicio de Azure | Ver abajo | Despliegues Azure |
| `TFVARS` | Contenido del archivo de variables de Terraform | Ver abajo | Todos los despliegues Azure |

#### Creación de Secretos

##### Claves de Acceso AWS

1. **Crear Usuario IAM** (si no existe):
   ```bash
   aws iam create-user --user-name github-actions-deploy
   ```

2. **Adjuntar Políticas**:
   ```bash
   aws iam attach-user-policy \
     --user-name github-actions-deploy \
     --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
   ```
   > **Nota**: Para producción, usa políticas de menor privilegio. Crea políticas personalizadas con solo los permisos necesarios.

3. **Crear Clave de Acceso**:
   ```bash
   aws iam create-access-key --user-name github-actions-deploy
   ```

4. **Agregar a Secretos de GitHub**:
   - Copia `AccessKeyId` → `AWS_ACCESS_KEY_ID`
   - Copia `SecretAccessKey` → `AWS_SECRET_ACCESS_KEY`

##### Entidad de Servicio Azure

1. **Crear Entidad de Servicio**:
   ```bash
   az ad sp create-for-rbac --name github-actions-deploy \
     --role contributor \
     --scopes /subscriptions/{subscription-id} \
     --sdk-auth
   ```

2. **Copia la salida JSON** y agrégala al Secreto de GitHub `AZURE_CREDENTIALS`:
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

##### Variables de Terraform (TFVARS)

El secreto `TFVARS` contiene el contenido de tu archivo `terraform.tfvars`. Debe incluir todas las variables de Terraform requeridas para tu infraestructura.

**Ejemplo para AWS:**
```hcl
# Contenido de terraform.tfvars
project_name = "go-project-skeleton"
environment = "development"
region = "us-east-1"
db_instance_class = "db.t3.micro"
lambda_memory_size = 512
# ... otras variables
```

**Ejemplo para Azure:**
```hcl
# Contenido de terraform.tfvars
project_name = "go-project-skeleton"
environment = "development"
location = "eastus"
app_service_plan_sku = "B1"
# ... otras variables
```

**Para crear el secreto:**
1. Crea tu archivo `terraform.tfvars` localmente
2. Copia todo el contenido
3. Agrega al Secreto de GitHub `TFVARS` (pega todo el contenido)

### Pasos del Flujo de Trabajo

#### 1. Fase de Configuración
- **Checkout código**: Clona el repositorio
- **Configurar Go**: Instala Go 1.25.5
- **Instalar Terraform**: Instala Terraform 1.14.1
- **Configurar credenciales**: Configura credenciales de AWS o Azure según el cloud seleccionado

#### 2. Fase de Construcción
- **Descargar dependencias**: Ejecuta `make deps`
- **Generar funciones**:
  - AWS: Ejecuta `make build-aws-functions` (si `deploy_functions` es true)
  - Azure: Ejecuta `make build-azure-functions`
- **Crear terraform.tfvars**: Crea el archivo desde el secreto `TFVARS`

#### 3. Fase Terraform
- **Terraform Init**: Inicializa el backend de Terraform
- **Terraform Validate**: Valida la configuración de Terraform
- **Terraform Plan**: Crea un plan de ejecución (si la acción no es `destroy`)
  - Sube artefacto de plan para revisión (si la acción es `plan`)
- **Terraform Apply**: Aplica cambios de infraestructura (si la acción es `apply`)
- **Terraform Destroy**: Destruye la infraestructura (si la acción es `destroy`)
  - **Seguridad**: Destroy falla en producción a menos que se permita explícitamente

#### 4. Fase de Despliegue (Solo AWS)
- **Desplegar Funciones Lambda**: Ejecuta `make deploy-aws` para desplegar todas las funciones Lambda
- **Terraform Output**: Muestra las salidas de la infraestructura

### Ejemplos de Uso

#### Planificar Cambios de Infraestructura

1. Ve a la pestaña **Actions** en GitHub
2. Selecciona el flujo de trabajo **Deploy to Cloud**
3. Haz clic en **Run workflow**
4. Completa las entradas:
   - Cloud: `aws`
   - Environment: `staging`
   - Terraform action: `plan`
   - Deploy functions: `false`
5. Haz clic en **Run workflow**

**Resultado**: Crea un plan de Terraform y lo sube como artefacto. Revisa el plan antes de aplicar.

#### Desplegar a Staging

1. Ve a la pestaña **Actions**
2. Selecciona el flujo de trabajo **Deploy to Cloud**
3. Haz clic en **Run workflow**
4. Completa las entradas:
   - Cloud: `aws`
   - Environment: `staging`
   - Terraform action: `apply`
   - Deploy functions: `true`
5. Haz clic en **Run workflow**

**Resultado**:
- Aprovisiona infraestructura con Terraform
- Despliega todas las funciones Lambda
- Muestra las salidas de la infraestructura

#### Destruir Entorno de Desarrollo

1. Ve a la pestaña **Actions**
2. Selecciona el flujo de trabajo **Deploy to Cloud**
3. Haz clic en **Run workflow**
4. Completa las entradas:
   - Cloud: `aws`
   - Environment: `development`
   - Terraform action: `destroy`
5. Haz clic en **Run workflow**

**Resultado**: Destruye toda la infraestructura en el entorno de desarrollo.

> **⚠️ Advertencia**: Las operaciones de destrucción en producción fallarán por defecto. El flujo de trabajo usa `continue-on-error: ${{ inputs.environment != 'production' }}` para prevenir la destrucción accidental de producción.

### Protección de Entornos

Los Entornos de GitHub pueden configurarse con reglas de protección:

1. **Revisores Requeridos**: Requerir aprobación antes del despliegue
2. **Temporizador de Espera**: Agregar un retraso antes del despliegue
3. **Ramas de Despliegue**: Restringir qué ramas pueden desplegar

**Para configurar:**
1. Ve a **Settings** → **Environments**
2. Haz clic en un entorno (ej: `production`)
3. Agrega reglas de protección según sea necesario

### Mejores Prácticas

1. **Usar Secretos de Entorno**: Almacenar secretos por entorno para mejor seguridad
2. **Revisar Planes**: Siempre ejecutar `plan` antes de `apply` en producción
3. **Usar Reglas de Protección**: Habilitar revisores requeridos para producción
4. **Monitorear Despliegues**: Revisar ejecuciones del flujo de trabajo regularmente
5. **Rotar Credenciales**: Rotar claves de acceso y secretos regularmente
6. **Menor Privilegio**: Usar roles/políticas IAM con permisos mínimos requeridos
7. **Control de Versiones**: Mantener código de Terraform en control de versiones
8. **Respaldo de Estado**: Asegurar que el estado de Terraform esté respaldado (S3, Azure Storage)

### Solución de Problemas

#### Problemas Comunes

**Problema**: "Credenciales AWS no encontradas"
- **Solución**: Asegúrate de que `AWS_ACCESS_KEY_ID` y `AWS_SECRET_ACCESS_KEY` estén configurados en Secretos de GitHub

**Problema**: "Terraform plan falla"
- **Solución**: Verifica que el contenido del secreto `TFVARS` coincida con el formato esperado

**Problema**: "Despliegue de Lambda falla"
- **Solución**: Asegúrate de que `PROJECT_NAME` y `AWS_REGION` estén configurados correctamente

**Problema**: "Destroy falla en producción"
- **Solución**: Esto es por diseño. Modifica el flujo de trabajo si se necesita destrucción en producción.

### Ubicación del Archivo del Flujo de Trabajo

El flujo de trabajo de despliegue se encuentra en:
```
.github/workflows/deploy.yml
```

---

## Guía de Desarrollo

### Agregar Nueva Funcionalidad

#### Diagrama de Flujo para Agregar Funcionalidad

```mermaid
flowchart TD
    Start([Nueva Funcionalidad]) --> Domain[1. Domain Layer<br/>Crear Modelo]
    Domain --> Contracts[2. Application Contracts<br/>Definir Interfaces]
    Contracts --> UseCase[3. Application Use Case<br/>Implementar Lógica]
    UseCase --> Repo[4. Infrastructure Repository<br/>Implementar con GORM]
    Repo --> Handler[5. Infrastructure Handler<br/>Adaptar HTTP]
    Handler --> Route[6. Routes<br/>Definir Endpoint]
    Route --> Tests[7. Tests<br/>Escribir Tests]
    Tests --> End([✅ Completado])

    style Domain fill:#e8f5e9
    style Contracts fill:#fff9c4
    style UseCase fill:#e3f2fd
    style Repo fill:#f3e5f5
    style Handler fill:#ffebee
    style Route fill:#e1f5ff
    style Tests fill:#fff4e1
```

#### Pasos Detallados

1. **Crear modelo en `domain/models/`**
```go
type NewEntity struct {
    // Campos
}
```

2. **Definir interfaz en `application/contracts/repositories/`**
```go
type INewEntityRepository interface {
    Create(input NewEntityCreate) (*NewEntity, error)
    // Más métodos
}
```

3. **Implementar caso de uso en `application/modules/newmodule/use_cases/`**
```go
type CreateNewEntityUseCase struct {
    repo INewEntityRepository
}

func (uc *CreateNewEntityUseCase) Execute(...) *UseCaseResult[NewEntity] {
    // Lógica de negocio
}
```

4. **Crear repositorio en `infrastructure/repositories/`**
```go
type NewEntityRepository struct {
    RepositoryBase[...]
}

func (r *NewEntityRepository) Create(...) (*NewEntity, error) {
    // Implementación con GORM
}
```

5. **Agregar handler en `infrastructure/handlers/`**
```go
func CreateNewEntity(ctx HandlerContext) {
    // Adaptar HTTP a caso de uso
}
```

6. **Definir ruta en `infrastructure/api/routes/router.go`**
```go
r.POST("/new-entity", wrapHandler(handlers.CreateNewEntity))
```

7. **Escribir tests**

### Convenciones de Código

- **Nombres**: camelCase para variables, PascalCase para tipos
- **Comentarios**: Documentar funciones públicas
- **Errores**: Usar el sistema de errores centralizado
- **Logging**: Usar el logger centralizado
- **Validación**: Validar entrada en todos los endpoints

---

## 🎯 Mejores Prácticas y Convenciones

### Convenciones de Código

#### Nomenclatura
- **Variables**: `camelCase` (ej: `userName`, `isActive`)
- **Tipos y Structs**: `PascalCase` (ej: `User`, `UserRepository`)
- **Interfaces**: `I` + `PascalCase` (ej: `IUserRepository`, `IHashProvider`)
- **Constantes**: `UPPER_SNAKE_CASE` (ej: `USER_STATUS_ACTIVE`)
- **Paquetes**: `lowercase` (ej: `user`, `auth`)

#### Estructura de Archivos
- **Use Cases**: `{action}_{entity}.go` (ej: `create_user.go`, `get_user.go`)
- **Repositories**: `{entity}.go` (ej: `user.go`, `password.go`)
- **Handlers**: `{entity}.go` o `{module}.go` (ej: `user.go`, `auth.go`)
- **Tests**: `{file}_test.go` (ej: `create_user_test.go`)

#### Organización de Código
```go
// 1. Imports (std, third-party, local)
import (
    "context"
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/simon3640/goprojectskeleton/src/domain/models"
)

// 2. Types y Structs
type UserRepository struct {
    // ...
}

// 3. Constructors
func NewUserRepository(...) *UserRepository {
    // ...
}

// 4. Métodos públicos
func (r *UserRepository) Create(...) {
    // ...
}

// 5. Métodos privados
func (r *UserRepository) validate(...) {
    // ...
}
```

### Principios de Diseño Aplicados

#### 1. Single Responsibility Principle (SRP)
Cada componente tiene una única responsabilidad:
- **Use Cases**: Lógica de negocio específica
- **Repositories**: Acceso a datos
- **Handlers**: Adaptación HTTP
- **Providers**: Servicios externos

#### 2. Dependency Inversion Principle (DIP)
Las capas internas definen interfaces que las externas implementan:
```go
// Application define la interfaz
type IUserRepository interface {
    Create(input UserCreate) (*User, error)
}

// Infrastructure implementa la interfaz
type UserRepository struct {
    DB *gorm.DB
}
```

#### 3. Open/Closed Principle (OCP)
Abierto para extensión, cerrado para modificación:
- Nuevos providers sin modificar código existente
- Nuevos casos de uso sin afectar otros
- Nuevos repositorios siguiendo interfaces

#### 4. Interface Segregation Principle (ISP)
Interfaces específicas y pequeñas:
```go
// En lugar de una interfaz grande
type IProvider interface {
    // 50 métodos...
}

// Interfaces específicas
type IHashProvider interface {
    Hash(password string) (string, error)
    Compare(hashed, plain string) bool
}
```

#### 5. Don't Repeat Yourself (DRY)
Reutilización mediante:
- `RepositoryBase` para operaciones CRUD comunes
- `BaseUseCaseValidation` para validación común
- Servicios compartidos en `application/shared/services`

### Manejo de Errores

#### Estructura de Errores
```go
type ApplicationError struct {
    Code    status.ApplicationStatusCode
    Context string
    ErrMsg  string
}
```

#### Estrategia de Errores
1. **Domain Layer**: Errores de negocio (validaciones)
2. **Application Layer**: Errores de aplicación (`ApplicationError`)
3. **Infrastructure Layer**: Errores técnicos (mapeados a `ApplicationError`)

### Logging

#### Niveles de Logging
- **Info**: Información general del flujo
- **Error**: Errores que requieren atención
- **Debug**: Información detallada para debugging
- **Panic**: Errores críticos que detienen la aplicación

#### Ejemplo de Uso
```go
providers.Logger.Info("User created successfully", map[string]interface{}{
    "user_id": user.ID,
    "email": user.Email,
})

providers.Logger.Error("Failed to create user", err)
```

### Testing

#### Estructura de Tests
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

#### Tipos de Tests
1. **Unit Tests**: Casos de uso con mocks
2. **Integration Tests**: Repositorios con BD real
3. **E2E Tests**: Flujos completos con Bruno

## Conclusión

**GoProjectSkeleton** proporciona una base sólida, escalable y mantenible para desarrollar aplicaciones empresariales en Go. La arquitectura Clean Architecture, junto con patrones de diseño avanzados, proporciona una base excelente para el crecimiento del sistema.

### 🎯 Puntos Fuertes

- ✅ **Arquitectura sólida** y bien estructurada
- ✅ **Escalabilidad** horizontal y vertical
- ✅ **Serverless ready** - Fácil migración a serverless (AWS Lambda, Azure Functions)
- ✅ **Testing completo** en múltiples capas
- ✅ **Documentación** exhaustiva con Swagger
- ✅ **Docker completo** para desarrollo y producción
- ✅ **Multi-cloud** - Soporte para AWS y Azure
- ✅ **Seguridad** - JWT, OTP, hash seguro de contraseñas
- ✅ **Internacionalización** - Soporte multiidioma
- ✅ **Optimización** - Cache, tree shaking, connection pooling
- ✅ **Observabilidad** - OpenTelemetry, Prometheus, Jaeger, Grafana

### 🚀 Casos de Uso Ideales

- **APIs RESTful** empresariales
- **Microservicios** escalables
- **Aplicaciones Serverless** (AWS Lambda, Azure Functions)
- **Sistemas de autenticación** robustos
- **APIs con alta concurrencia**
- **Proyectos que requieren mantenibilidad** a largo plazo

### 📚 Próximos Pasos

1. **Explorar la Documentación**
   - Revisar Swagger en `http://localhost:8081/docs/` (servicio independiente)
   - Leer ejemplos de código en cada módulo

2. **Ejecutar Tests**
   ```bash
   go test ./src/...
   go test ./tests/integration/...
   ```

3. **Explorar Observabilidad**
   - Acceder a Grafana en `http://localhost:3001` (admin/admin)
   - Explorar trazas en Jaeger en `http://localhost:16686`
   - Consultar métricas en Prometheus en `http://localhost:9090`
   - Revisar dashboards preconfigurados

4. **Adaptar a tus Necesidades**
   - Personalizar modelos de dominio
   - Agregar nuevos módulos de negocio
   - Configurar providers según tus servicios

5. **Desplegar**
   - Desarrollo: Docker Compose
   - Producción: Monolito tradicional o Serverless
   - Cloud: AWS Lambda o Azure Functions

### 🤝 Contribuciones

Este proyecto sigue las mejores prácticas de Go y arquitectura limpia. Al contribuir:

1. Mantén la separación de capas
2. Sigue las convenciones de código
3. Escribe tests para nuevas funcionalidades
4. Documenta cambios importantes
5. Actualiza Swagger para nuevos endpoints

### 📞 Soporte

Para más información:
- **Documentación Swagger**: `http://localhost:8081/docs/` (servicio independiente)
- **Tests de Ejemplo**: `tests/`
- **Código Fuente**: Explora `src/` para ver implementaciones

---

<div align="center">
  <p>Hecho con ❤️ usando Go y Clean Architecture</p>
  <p>⭐ Si este proyecto te es útil, considera darle una estrella</p>
</div>
