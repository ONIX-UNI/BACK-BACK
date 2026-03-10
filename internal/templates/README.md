# templates/

## Propósito

El módulo `templates` es responsable de la **generación dinámica de documentos**, principalmente **PDF**, a partir de **formularios definidos** y datos ingresados por el usuario.

Su función principal es **transformar información estructurada en documentos legales formales**, siguiendo plantillas predefinidas.

---

## Responsabilidades

Este módulo se encarga de:

- Definir plantillas de documentos (estructura lógica)
- Gestionar formularios asociados a cada plantilla
- Validar datos ingresados por el usuario
- Rellenar plantillas con información dinámica
- Generar documentos finales (ej. PDF)
- Emitir eventos cuando un documento es generado

---

## Qué NO hace este módulo

- No almacena archivos físicamente
- No gestiona usuarios ni permisos
- No envía notificaciones
- No administra turnos ni PQRS
- No conoce MinIO, S3 ni sistemas de archivos concretos

---

## Estructura interna

```text
templates/
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

## Conceptos clave del dominio

### Plantilla

Representa un **modelo base de documento**.

Ejemplos:

- Poder legal
- Derecho de petición
- Tutela
- Contrato simple

Contiene:

- Identificador
- Tipo de documento
- Versión
- Estructura lógica del contenido
- Campos dinámicos

---

### Formulario

Define **qué datos debe ingresar el usuario** para completar una plantilla.

Incluye:

- Nombre del campo
- Tipo (texto, número, fecha, selección, etc.)
- Reglas de validación
- Obligatorio / opcional

---

### Documento generado

Resultado final de:

- Plantilla + Datos del formulario

El documento generado:

- Tiene identidad propia
- Puede ser versionado
- Es trazable mediante eventos

---

## Flujo principal (ejemplo: generación de documento)

1. El usuario selecciona una plantilla
2. El sistema devuelve el formulario asociado
3. El usuario completa los datos
4. El handler valida el input
5. El servicio:
   - valida reglas de negocio
   - rellena la plantilla
   - genera el documento

6. Se registra el documento generado
7. Se emite un evento de documento generado

---

## Eventos

Eventos comunes emitidos por este módulo:

- `TemplateSelected`
- `TemplateCompleted`
- `DocumentGeneratedFromTemplate`

Estos eventos son consumidos por:

- módulo de documentos
- notificaciones
- workers de correo

---

## Dependencias

El módulo depende únicamente de **interfaces**, tales como:

- Motor de generación de documentos (PDF)
- Repositorio de plantillas
- Servicio de persistencia de metadata

Las implementaciones concretas se inyectan desde capas superiores.

---

## Reglas importantes

- Una plantilla puede tener múltiples versiones
- Un documento generado siempre referencia la versión exacta de la plantilla
- Los datos ingresados deben validarse antes de la generación
- El servicio es la única capa que conoce el flujo completo

---

## Escalabilidad y extensibilidad

Este diseño permite:

- Soportar múltiples formatos (PDF, DOCX, etc.)
- Añadir nuevos tipos de plantillas sin impacto lateral
- Incorporar firmas digitales en el futuro
- Integración con flujos legales automatizados

---

## Relación con otros módulos

- **documents**: recibe y gestiona el documento generado
- **notifications**: informa al usuario del resultado
- **auth**: valida permisos de acceso (externo al módulo)

---

## Resumen

- `templates` convierte datos en documentos legales
- Es un módulo puramente de **orquestación y reglas**
- No depende de infraestructura concreta
- Está diseñado para crecer sin romper contratos

Este módulo es clave para la automatización documental del sistema.
