# notifications/

## Propósito

El módulo `notifications` es responsable de la **orquestación de notificaciones del sistema**.

Actúa como un **centro de coordinación**, reaccionando a eventos del dominio y decidiendo **qué notificar, a quién y por qué canal**, sin ejecutar directamente el envío.

---

## Responsabilidades

Este módulo se encarga de:

- Escuchar eventos del sistema (turnos, documentos, plantillas, PQRS)
- Determinar reglas de notificación
- Construir mensajes de notificación
- Publicar tareas de notificación en colas
- Centralizar la lógica de comunicación con usuarios

---

## Qué NO hace este módulo

- No envía correos directamente
- No envía mensajes SMS, push o WhatsApp
- No implementa clientes SMTP ni servicios externos
- No gestiona estados de turnos o documentos
- No accede directamente a infraestructura de mensajería

---

## Estructura interna

```text
notifications/
├── handlers/
├── service/
├── repository/
├── models/
├── dto/
├── events/
├── routes.go
└── module.go
```

---

## Rol dentro del sistema

`notifications` funciona como un **traductor de eventos de negocio** a **intenciones de notificación**.

Ejemplo:

- Evento: `AppointmentStatusChanged`
- Acción: generar notificación de actualización de turno
- Canal: correo electrónico (decisión delegada al worker)

---

## Tipos de notificaciones

El módulo define notificaciones para eventos como:

- Actualización de turnos
- Confirmación de documentos cargados
- Documentos generados desde plantillas
- Respuestas a PQRS
- Eventos administrativos del sistema

---

## Flujo principal (ejemplo: notificación de turno)

1. Se recibe evento `AppointmentStatusChanged`
2. El servicio evalúa reglas:
   - tipo de usuario
   - estado del turno

3. Se construye el payload de notificación
4. Se publica una tarea en la cola correspondiente
5. Se registra el intento de notificación (opcional)

---

## Eventos

### Eventos consumidos

- `AppointmentCreated`
- `AppointmentStatusChanged`
- `DocumentUploaded`
- `DocumentGeneratedFromTemplate`
- `PQRSUpdated`

### Eventos emitidos

- `NotificationQueued`
- `NotificationFailed` (opcional)
- `NotificationSent` (confirmación asíncrona)

---

## Modelos principales

- Notificación
- TipoNotificación
- Canal (EMAIL, PUSH, etc.)
- EstadoNotificación

---

## Reglas importantes

- Las notificaciones son **eventualmente consistentes**
- El envío es siempre asíncrono
- La lógica de decisión vive en el servicio
- El módulo no asume éxito inmediato

---

## Dependencias

Este módulo depende de:

- Bus de eventos (interfaces)
- Repositorio de notificaciones (opcional)
- Configuración de reglas

Las dependencias concretas (Kafka, RabbitMQ) se inyectan externamente.

---

## Escalabilidad

Este diseño permite:

- Añadir nuevos canales sin tocar el dominio
- Escalar workers de envío de forma independiente
- Reintentos y DLQ sin impacto en el API
- Observabilidad de notificaciones

---

## Relación con otros módulos

- **appointments**: eventos de turnos
- **documents**: carga y disponibilidad de documentos
- **templates**: documentos generados
- **pqrs**: respuestas y cambios de estado
- **workers**: ejecución real del envío

---

## Resumen

- `notifications` decide **qué notificar**
- Los workers deciden **cómo notificar**
- El módulo es puramente orquestador
- Garantiza desacoplamiento y escalabilidad

Es un punto clave para la experiencia del usuario.
