# cmd/workers/

## Propósito

La carpeta `cmd/workers/` contiene los **procesos en segundo plano (workers)** del sistema.

Cada worker es un **binario independiente**, diseñado para ejecutar **tareas asíncronas**, consumir eventos o procesar colas sin bloquear el API principal.

Los workers permiten:

- Desacoplar operaciones pesadas
- Escalar componentes de forma independiente
- Mejorar la resiliencia del sistema

---

## Principios de diseño

- **Un worker = una responsabilidad**
- **Procesos independientes**
- **Comunicación basada en eventos**
- **Escalabilidad horizontal**
- **Tolerancia a fallos**

Los workers **no exponen endpoints HTTP**.

---

## Estructura general

```text
cmd/workers/
├── mailer/
│   └── main.go
├── notifications/
│   └── main.go
├── pqrs/
│   └── main.go
└── README.md
```

Cada carpeta representa un **worker autónomo**.

---

## ¿Qué hace un worker?

Un worker:

- Escucha eventos (Kafka / RabbitMQ)
- Procesa mensajes de una cola
- Ejecuta tareas asíncronas
- Interactúa con servicios externos
- Reporta resultados mediante eventos

---

## Flujo general de un worker

1. Inicio del proceso
2. Carga de configuración
3. Inicialización de infraestructura necesaria
4. Suscripción a colas o tópicos
5. Procesamiento continuo de mensajes
6. Manejo de errores y reintentos

---

## Workers disponibles

### `mailer/`

Responsable del **envío de correos electrónicos**.

Procesa tareas como:

- Confirmación de acciones
- Notificación de eventos
- Respuestas a PQRS

Consume:

- eventos de notificación
- mensajes de cola de correo

---

### `notifications/`

Encargado de **ejecutar notificaciones** definidas por el módulo `internal/notifications`.

Procesa:

- decisiones de notificación
- selección de canal
- despacho al worker correspondiente

Puede actuar como **router** de notificaciones.

---

### `pqrs/`

Procesa tareas asíncronas relacionadas con PQRS, como:

- escalamiento automático
- recordatorios
- cierre automático por SLA

Consume eventos del módulo `pqrs`.

---

## Qué NO deben hacer los workers

- No contienen lógica de negocio central
- No modifican estados directamente en el dominio
- No exponen APIs HTTP
- No acceden a otros workers directamente
- No toman decisiones funcionales

---

## Dependencias

Cada worker puede depender de:

- `pkg/messaging` (Kafka / RabbitMQ)
- `pkg/mail` (SMTP)
- `pkg/logger`
- `pkg/config`
- Interfaces del dominio (`internal/`)

Las dependencias se inyectan al inicio.

---

## Manejo de errores

- Los errores deben ser **controlados y registrados**
- Se soportan reintentos
- Mensajes fallidos pueden ir a DLQ
- El fallo de un worker **no afecta al API**

---

## Escalabilidad

Esta arquitectura permite:

- Escalar workers de forma independiente
- Asignar más recursos a tareas pesadas
- Reintentos sin impacto al usuario
- Procesamiento paralelo

---

## Relación con el resto del sistema

```text
cmd/workers/ → ejecución asíncrona
internal/    → reglas de negocio
pkg/         → infraestructura
cmd/api/     → entrada HTTP
```

Los workers **consumen eventos**, no llaman directamente al API.

---

## Resumen

- `cmd/workers/` contiene procesos asíncronos
- Cada worker es independiente
- El sistema es resiliente y escalable
- El API permanece liviano y estable

Esta carpeta es clave para la robustez del backend.
