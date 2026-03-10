## 📁 `internal/appointments/README.md`

# appointments/

## Propósito

El módulo `appointments` gestiona el **sistema de turnos**, similar al de una entidad bancaria.

Controla la **creación, avance y finalización de fichas**, así como los estados asociados y las notificaciones derivadas.

---

## Responsabilidades

Este módulo se encarga de:

- Creación de turnos (fichas)
- Administración de colas
- Transiciones de estado del turno
- Priorización (si aplica)
- Exposición del estado actual de la cola

---

## Qué NO hace este módulo

- No envía notificaciones directamente
- No gestiona usuarios (solo referencia IDs)
- No maneja persistencia fuera de su dominio
- No conoce SMTP, Kafka ni WebSockets

---

## Estructura interna

```text
appointments/
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

## Estados típicos de un turno

Ejemplo de estados manejados:

- `CREATED`
- `WAITING`
- `IN_PROGRESS`
- `COMPLETED`
- `CANCELLED`

Las transiciones están controladas **exclusivamente** por el servicio.

---

## Flujo principal (ejemplo: tomar turno)

1. Usuario solicita un turno
2. El sistema genera una ficha
3. Se asigna posición en la cola
4. Se persiste el estado
5. Se emite un evento de turno creado

---

## Eventos

Eventos comunes emitidos:

- `AppointmentCreated`
- `AppointmentStatusChanged`
- `AppointmentCompleted`

Estos eventos son consumidos por:

- módulo de notificaciones
- workers
- métricas

---

## Reglas importantes

- Las transiciones de estado son atómicas
- No se puede saltar estados inválidos
- El orden de la cola es consistente
- El servicio es la única fuente de verdad

---

## Escalabilidad

El diseño permite:

- Escalar colas por tipo de servicio
- Manejar múltiples puntos de atención
- Migrar a un sistema distribuido sin reescritura del dominio
