# pkg/

## Propósito

La carpeta `pkg/` contiene **componentes de infraestructura compartidos** por todo el sistema.

Aquí viven las implementaciones técnicas necesarias para operar el backend, **sin lógica de negocio**.  
Su objetivo es **soportar al dominio**, no definirlo.

Todo lo que esté en `pkg/` debe ser:

- Reutilizable
- Desacoplado
- Intercambiable
- Independiente del negocio

---

## Principios de diseño

- **Infraestructura ≠ dominio**
- **Código reutilizable**
- **Dependencias explícitas**
- **Bajo acoplamiento**
- **Alta cohesión por responsabilidad**

El dominio (`internal/`) depende de `pkg/` **solo mediante interfaces**.

---

## Estructura general

```text
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

Cada carpeta cumple una **responsabilidad técnica específica**.

---

## Descripción por carpeta

### `config/`

Responsable de la **carga y validación de configuración** del sistema.

Contiene:

- Lectura de variables de entorno
- Structs de configuración tipados
- Validaciones de valores requeridos

No contiene:

- Lógica de negocio
- Defaults implícitos peligrosos

---

### `database/`

Provee acceso a **bases de datos relacionales**, principalmente PostgreSQL.

Contiene:

- Inicialización de conexiones
- Pool de conexiones
- Manejo de transacciones
- Health checks

No contiene:

- Queries de negocio
- Modelos del dominio

---

### `messaging/`

Abstrae los **sistemas de mensajería** usados por el sistema.

Soporta:

- Kafka
- RabbitMQ

Contiene:

- Publishers
- Consumers
- Configuración de tópicos/colas
- Manejo de reconexión y retries

No contiene:

- Lógica de eventos de negocio
- Definición de payloads del dominio

---

### `storage/`

Provee acceso a **almacenamiento de objetos**.

Pensado para:

- MinIO
- S3 compatible

Contiene:

- Clientes de storage
- Operaciones básicas (put, get, delete)
- Manejo de buckets

No contiene:

- Reglas de documentos
- Validaciones de negocio

---

### `mail/`

Encapsula el **envío de correos electrónicos**.

Contiene:

- Cliente SMTP
- Configuración de servidor
- Envío básico de mensajes

No contiene:

- Plantillas de correo
- Decisiones de cuándo enviar

---

### `logger/`

Provee un **sistema de logging centralizado**.

Contiene:

- Inicialización del logger
- Niveles de log
- Formato estructurado

Debe ser:

- Usado por todos los módulos
- Configurable por entorno

---

### `errors/`

Define **errores tipados y reutilizables**.

Contiene:

- Errores base
- Helpers para wrapping
- Clasificación de errores

Objetivo:

- Manejo consistente de errores
- Traducción a respuestas HTTP

---

### `utils/`

Contiene **helpers genéricos** y funciones utilitarias.

Ejemplos:

- Manejo de fechas
- Generación de IDs
- Conversión de tipos

Regla:

- Si empieza a crecer o especializarse → crear paquete dedicado

---

## Reglas estrictas en `pkg/`

### 1. Prohibido importar `internal/`

❌ Incorrecto:

```go
import "internal/auth/models"
```

✔ Correcto:

- `pkg/` es completamente agnóstico al dominio

---

### 2. Nada de lógica de negocio

Si un paquete empieza a:

- tomar decisiones de negocio
- validar reglas funcionales

➡️ pertenece a `internal/`, no a `pkg/`.

---

### 3. Interfaces primero

Siempre que sea posible:

- Definir interfaces
- Inyectar implementaciones
- Permitir reemplazos tecnológicos

---

## Relación con otras carpetas

```text
cmd/        → inicia y ensambla
internal/   → define el negocio
pkg/        → ejecuta la infraestructura
```

`pkg/` es utilizado por:

- API
- Workers
- Tests
- Herramientas internas

---

## Escalabilidad y mantenibilidad

Esta estructura permite:

- Cambiar Kafka por otro broker
- Reemplazar MinIO por S3 real
- Modificar logging sin tocar dominio
- Reusar infraestructura en otros proyectos

---

## Resumen

- `pkg/` es la **base técnica del sistema**
- No contiene reglas de negocio
- Está organizada por responsabilidad técnica
- Facilita reemplazos, pruebas y escalabilidad

Cualquier código aquí debe poder sobrevivir **fuera de este proyecto**.

```

```
