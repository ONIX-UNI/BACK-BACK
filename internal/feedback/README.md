# feedback/

## Propósito

El módulo `feedback` es responsable de la **recolección, gestión y consulta de opiniones, calificaciones y comentarios** de los usuarios sobre los servicios prestados por el consultorio jurídico.

Su objetivo es proporcionar **insumos medibles y trazables** para evaluar la calidad del servicio y apoyar procesos de mejora continua.

---

## Responsabilidades

Este módulo se encarga de:

- Recepción de feedback de usuarios
- Registro de calificaciones y comentarios
- Asociación del feedback a servicios, turnos o atenciones
- Consulta y listado de feedback
- Emisión de eventos para análisis o notificaciones

---

## Qué NO hace este módulo

- No gestiona usuarios ni autenticación
- No responde PQRS
- No genera reportes agregados complejos
- No envía notificaciones directamente
- No toma decisiones de negocio basadas en el feedback

---

## Estructura interna

```text
feedback/
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

## Tipos de feedback

El módulo soporta distintos tipos de retroalimentación, por ejemplo:

- Calificación numérica (ej. 1 a 5)
- Comentario textual
- Feedback anónimo (si aplica)
- Feedback asociado a:
  - un turno
  - un documento
  - una atención específica

---

## Flujo principal (ejemplo: envío de feedback)

1. El usuario envía su feedback
2. El handler valida el input
3. El servicio:
   - valida reglas (ej. una sola vez por atención)
   - registra el feedback

4. Se persiste la información
5. Se emite un evento de feedback recibido

---

## Modelos principales

- Feedback
- TipoFeedback
- Calificación
- ContextoFeedback (turno, servicio, etc.)

---

## Eventos

### Eventos emitidos

- `FeedbackSubmitted`
- `FeedbackUpdated` (si se permite edición)

Estos eventos pueden ser consumidos por:

- reportes
- notificaciones
- análisis de calidad

---

## Reglas importantes

- El feedback no modifica estados de otros dominios
- Puede ser anónimo según configuración
- La edición, si existe, es controlada por reglas claras
- El servicio es la única capa que valida duplicidad

---

## Dependencias

El módulo depende de:

- Repositorio de persistencia (interface)
- Emisor de eventos (interface)
- Configuración de reglas de feedback

Las dependencias concretas se inyectan desde capas superiores.

---

## Escalabilidad

Este diseño permite:

- Manejar grandes volúmenes de feedback
- Incorporar análisis automático en el futuro
- Integración con sistemas de métricas
- Uso del feedback como insumo para reportes

---

## Relación con otros módulos

- **appointments**: feedback sobre turnos
- **notifications**: confirmación de recepción
- **reports**: agregación y análisis
- **auth**: validación de permisos (externo)

---

## Resumen

- `feedback` captura la voz del usuario
- Es un módulo pasivo, no decisorio
- Está diseñado para análisis y mejora continua
- Mantiene bajo acoplamiento con el resto del sistema

Este módulo aporta visibilidad y control de calidad.
