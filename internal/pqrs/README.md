# pqrs/

## Propósito

El módulo `pqrs` gestiona el ciclo completo de **Peticiones, Quejas, Reclamos y Sugerencias (PQRS)** del sistema.

Su objetivo es ofrecer un **canal formal, trazable y auditable** para que los usuarios se comuniquen con el consultorio jurídico y para que el sistema administre su atención de manera estructurada.

---

## Responsabilidades

Este módulo se encarga de:

- Recepción de PQRS desde distintos canales
- Clasificación por tipo (P, Q, R, S)
- Gestión de estados y flujos de atención
- Asignación de responsables (si aplica)
- Registro de respuestas
- Emisión de eventos asociados al ciclo de vida

---

## Qué NO hace este módulo

- No envía correos ni notificaciones directamente
- No gestiona usuarios (solo referencia IDs)
- No define turnos ni documentos
- No implementa colas ni workers
- No conoce SMTP, Kafka ni RabbitMQ

---

## Estructura interna

```text
pqrs/
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

## Tipos de PQRS

El módulo maneja explícitamente los siguientes tipos:

- **P** – Petición
- **Q** – Queja
- **R** – Reclamo
- **S** – Sugerencia

Cada tipo puede compartir flujo, pero con reglas específicas.

---

## Estados típicos de una PQRS

Ejemplo de estados controlados por el dominio:

- `CREATED`
- `RECEIVED`
- `IN_REVIEW`
- `RESPONDED`
- `CLOSED`
- `CANCELLED`

Las transiciones de estado están **estrictamente controladas** por el servicio.

---

## Flujo principal (ejemplo: creación de PQRS)

1. El usuario envía una PQRS
2. El handler valida el input
3. El servicio:
   - clasifica el tipo
   - asigna estado inicial
   - registra la PQRS

4. Se emite un evento de PQRS creada
5. Otros módulos reaccionan (notificaciones, seguimiento)

---

## Respuesta a una PQRS

1. Un actor autorizado responde la PQRS
2. Se valida el estado actual
3. Se registra la respuesta
4. Se cambia el estado correspondiente
5. Se emite evento de actualización

---

## Modelos principales

- PQRS
- TipoPQRS
- EstadoPQRS
- RespuestaPQRS
- HistorialPQRS

---

## Eventos

### Eventos emitidos

- `PQRSCreated`
- `PQRSUpdated`
- `PQRSResponded`
- `PQRSClosed`

### Eventos consumidos (opcional)

- Eventos de usuarios
- Eventos administrativos

Los eventos permiten:

- Notificar al usuario
- Auditoría
- Métricas y seguimiento

---

## Reglas importantes

- Toda PQRS debe ser trazable
- Ninguna PQRS se elimina físicamente
- Los cambios de estado son auditables
- El servicio es la única capa que modifica el estado

---

## Dependencias

El módulo depende de:

- Repositorio de persistencia (interface)
- Emisor de eventos (interface)
- Configuración de flujos y SLA (opcional)

Las implementaciones concretas se inyectan externamente.

---

## Escalabilidad y control

Este diseño permite:

- Manejar altos volúmenes de PQRS
- Automatizar respuestas en el futuro
- Integración con flujos legales formales
- Medición de tiempos de atención (SLA)

---

## Relación con otros módulos

- **notifications**: notificación al usuario
- **documents**: adjuntos asociados
- **auth**: validación de roles y permisos
- **workers**: ejecución asíncrona de acciones

---

## Resumen

- `pqrs` centraliza la comunicación formal del usuario
- Controla flujos, estados y respuestas
- Es auditable y extensible
- Está diseñado para crecer sin romper el dominio

Este módulo es crítico para la confianza y trazabilidad del sistema.
