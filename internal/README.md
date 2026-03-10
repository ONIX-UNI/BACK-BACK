# internal/

## Propósito

La carpeta `internal/` contiene **todo el núcleo del negocio del sistema**.  
Aquí vive la lógica que define **qué hace** el sistema y **cómo se comporta**, independientemente de la infraestructura, frameworks o mecanismos de despliegue.

Siguiendo las convenciones de Go, todo lo que esté dentro de `internal/` **no puede ser importado desde fuera del proyecto**, lo que protege el dominio y evita dependencias indebidas.

---

## Principios de diseño

La estructura de `internal/` está guiada por los siguientes principios:

- **Arquitectura modular por dominio**
- **Alto desacoplamiento entre módulos**
- **Escalabilidad horizontal**
- **Facilidad de lectura y onboarding**
- **Consistencia estructural entre módulos**

Cada módulo representa un **dominio funcional claro** del sistema.

---

## Estructura general

```text
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

Cada carpeta corresponde a un **módulo de negocio independiente**.

---

## ¿Qué es un módulo?

Un módulo dentro de `internal/`:

- Encapsula un **dominio funcional específico**
- Contiene su propia:
  - lógica de negocio
  - contratos de entrada/salida
  - acceso a datos
  - eventos

- No conoce detalles de otros módulos
- Puede evolucionar a microservicio sin reestructuración mayor

---

## Estructura estándar de un módulo

Todos los módulos dentro de `internal/` siguen la **misma estructura base** para garantizar consistencia:

```text
internal/<modulo>/
├── handlers/        # Entradas HTTP (Fiber)
├── service/         # Lógica de negocio
├── repository/      # Persistencia y consultas
├── models/          # Entidades del dominio
├── dto/             # Request / Response / Payloads
├── events/          # Definición y emisión de eventos
├── routes.go        # Registro de rutas del módulo
└── module.go        # Inicialización y dependencias
```

### Responsabilidad por carpeta

| Carpeta       | Responsabilidad                                  |
| ------------- | ------------------------------------------------ |
| `handlers/`   | Adaptadores HTTP. No contienen lógica de negocio |
| `service/`    | Reglas del negocio y orquestación                |
| `repository/` | Acceso a PostgreSQL u otros stores               |
| `models/`     | Entidades centrales del dominio                  |
| `dto/`        | Contratos externos (input/output)                |
| `events/`     | Eventos Kafka/Rabbit emitidos o consumidos       |
| `routes.go`   | Declaración de endpoints                         |
| `module.go`   | Wiring interno del módulo                        |

---

## Reglas estrictas dentro de `internal/`

### 1. No importar otros módulos directamente

❌ Incorrecto:

```go
import "internal/documents/service"
```

✔ Correcto:

- Comunicación mediante eventos
- Interfaces expuestas
- Orquestación en capas superiores

---

### 2. Los handlers no contienen lógica de negocio

Los `handlers` solo deben:

- Validar input
- Convertir DTOs
- Invocar servicios
- Retornar respuestas HTTP

Toda la lógica vive en `service/`.

---

### 3. El dominio no depende de infraestructura

Los módulos:

- **No conocen Kafka, RabbitMQ, MinIO ni SMTP**
- Solo dependen de **interfaces**
- La infraestructura se inyecta desde `pkg/`

---

## Relación con otras carpetas

```text
cmd/        → inicia el sistema
internal/   → define el negocio
pkg/        → provee infraestructura
docs/       → explica el sistema
```

`internal/` **nunca depende de `cmd/`**, y solo depende de `pkg/` mediante abstracciones.

---

## Escalabilidad futura

Esta estructura permite:

- Separar módulos como microservicios
- Escalar workers de forma independiente
- Reemplazar tecnologías sin reescribir dominio
- Incorporar nuevos módulos sin impacto lateral

---

## ¿Dónde documentar cada módulo?

Cada módulo dentro de `internal/` debe tener su propio:

```text
internal/<modulo>/README.md
```

Ese README debe explicar:

- Propósito del módulo
- Responsabilidades
- Flujo principal
- Eventos emitidos/consumidos
- Dependencias

---

## Resumen

- `internal/` es el **corazón del sistema**
- Contiene **únicamente lógica de negocio**
- Está organizada por **dominios funcionales**
- Prioriza claridad, escalabilidad y mantenibilidad

```

```
