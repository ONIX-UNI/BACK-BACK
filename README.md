Estamos desarrollando el **backend de un sistema para un consultorio jurídico**, implementado en **Go** utilizando el framework **Go Fiber**.
En esta primera etapa nos concentraremos exclusivamente en la arquitectura y los módulos del backend.

El sistema estará compuesto por los siguientes **módulos funcionales**:

### Módulos principales

- **Autenticación**
  Gestión de usuarios, roles y control de acceso.

- **Gestión de turnos**
  Sistema de turnos similar al de una entidad bancaria, donde el usuario toma una ficha y el sistema administra la cola, los estados y las notificaciones asociadas.

- **Gestión de documentos**
  Cargue, almacenamiento y consulta de documentos, con soporte para almacenamiento distribuido.

- **Módulo de feedback**
  Recepción y gestión de opiniones, calificaciones o comentarios de los usuarios sobre el servicio.

- **Módulo de plantillas**
  Permite rellenar dinámicamente diferentes archivos **PDF** a partir de formularios definidos, generando documentos personalizados con la información ingresada por el usuario.

- **Módulo de reportes**
  Generación de reportes del sistema, con opción de descarga en formato **CSV**.

### Infraestructura y servicios externos

El sistema se integrará con los siguientes servicios:

- **PostgreSQL** como base de datos relacional.
- **Kafka** para mensajería y procesamiento de eventos.
- **MinIO (open source)** como almacenamiento de objetos compatible con S3.
- **SMTP** para el envío de correos electrónicos.
- **RabbitMQ** para colas de mensajes y comunicación asíncrona.

### Workers (procesos en segundo plano)

Se contará con workers independientes para manejar tareas asincrónicas, organizados en los siguientes módulos:

- **Envío de correos electrónicos**
- **Gestión de PQRS** (Peticiones, Quejas, Reclamos y Sugerencias)
- **Notificaciones**, incluyendo:
  - Actualización de turnos
  - Confirmación de documentos cargados
  - Notificación de documentos generados a partir de plantillas completadas

---

## Principios que guían la estructura

Antes del árbol, dejo claros los criterios (esto explica por qué está organizada así):

- **Arquitectura modular por dominio** (no por capa técnica).
- **Cada módulo es autocontenido**: handlers, lógica, repositorios, contratos.
- **Infraestructura desacoplada** del dominio.
- **Workers separados del API HTTP** (pueden escalar distinto).
- **Preparada para microservicios**, aunque inicialmente sea un monolito modular.
- **Go idiomático**, sin sobre–ingeniería temprana.

---

## Estructura general del proyecto

```bash
SICOU/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── workers/
│       ├── mailer/
│       │   └── main.go
│       ├── notifications/
│       │   └── main.go
│       └── pqrs/
│           └── main.go
│
├── internal/
│   ├── auth/
│   ├── appointments/
│   ├── documents/
│   ├── feedback/
│   ├── templates/
│   ├── reports/
│   ├── notifications/
│   └── pqrs/
│
├── pkg/
│   ├── config/
│   ├── database/
│   ├── messaging/
│   ├── storage/
│   ├── mail/
│   ├── logger/
│   ├── errors/
│   └── utils/
│
├── migrations/
│
├── scripts/
│
├── deployments/
│
├── docs/
│
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

---

## 1️⃣ `cmd/` – puntos de entrada (entrypoints)

```bash
cmd/
├── api/
│   └── main.go
└── workers/
    ├── mailer/
    ├── notifications/
    └── pqrs/
```

### ¿Por qué así?

- Cada carpeta dentro de `cmd/` representa **un binario independiente**.
- Permite:
  - Escalar API y workers por separado
  - Deploys independientes
  - Menor acoplamiento

Ejemplos:

- `cmd/api` → servidor Fiber (HTTP)
- `cmd/workers/mailer` → worker SMTP
- `cmd/workers/notifications` → eventos Kafka/Rabbit
- `cmd/workers/pqrs` → procesamiento PQRS

---

## 2️⃣ `internal/` – dominio del negocio (núcleo del sistema)

Aquí vive **toda la lógica del sistema**.
Nada fuera de este repo debería importar esto (por eso `internal`).

```bash
internal/
├── auth/
├── appointments/
├── documents/
├── feedback/
├── templates/
├── reports/
├── notifications/
└── pqrs/
```

Cada módulo sigue **la misma estructura interna** (consistencia > creatividad).

### Estructura estándar por módulo

```bash
internal/auth/
├── handlers/        # HTTP handlers (Fiber)
├── service/         # lógica de negocio
├── repository/      # acceso a datos
├── models/          # entidades del dominio
├── dto/             # request/response structs
├── events/          # eventos Kafka/Rabbit
└── module.go        # registro de rutas
```

> **Regla clave:**
> Un módulo **no importa otro módulo directamente**.
> La comunicación es vía:
>
> - eventos
> - interfaces
> - capa de aplicación

---

## 3️⃣ Módulos principales (visión rápida)

### 🔐 `auth`

- Usuarios
- Roles
- JWT / sesiones
- Control de acceso

### 🧾 `appointments`

- Turnos tipo banco
- Estados de ficha
- Colas
- Notificaciones de avance

### 📁 `documents`

- Cargue
- Validación
- Metadatos
- Integración con MinIO

### 🗣️ `feedback`

- Opiniones
- Calificaciones
- Comentarios

### 🧩 `templates`

- Formularios dinámicos
- Relleno de PDFs
- Generación de documentos finales

### 📊 `reports`

- Consultas agregadas
- Exportación CSV

### 🔔 `notifications`

- Orquestación de notificaciones
- Consumo de eventos
- Producción de mensajes

### 📮 `pqrs`

- Recepción
- Flujo de atención
- Estados
- Integración con notificaciones y correo

---

## 4️⃣ `pkg/` – infraestructura compartida (reutilizable)

```bash
pkg/
├── config/
├── database/
├── messaging/
├── storage/
├── mail/
├── logger/
├── errors/
└── utils/
```

### Responsabilidad clara

| Carpeta     | Contiene                               |
| ----------- | -------------------------------------- |
| `config`    | Carga de env, structs de configuración |
| `database`  | PostgreSQL (pool, tx, health)          |
| `messaging` | Kafka + RabbitMQ                       |
| `storage`   | MinIO / S3                             |
| `mail`      | SMTP                                   |
| `logger`    | Logging central                        |
| `errors`    | Errores tipados                        |
| `utils`     | Helpers genéricos                      |

> **Nada aquí conoce el negocio.**
> Solo infraestructura.

---

## 5️⃣ `migrations/`

```bash
migrations/
├── 0001_init.sql
├── 0002_users.sql
└── 0003_appointments.sql
```

- SQL puro
- Versionadas
- Independientes del código Go

---
