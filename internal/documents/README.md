## 📁 `internal/documents/README.md`

# documents/

## Propósito

El módulo `documents` gestiona el **cargue, almacenamiento, consulta y estado de documentos** dentro del sistema.

Funciona como el **registro y orquestador de documentos**, independientemente de dónde estén almacenados físicamente.

---

## Responsabilidades

Este módulo se encarga de:

- Registro de documentos
- Validación de metadata
- Versionado lógico
- Asociación a usuarios, turnos o PQRS
- Consulta y listado de documentos

---

## Qué NO hace este módulo

- No almacena archivos directamente
- No conoce MinIO ni S3 de forma concreta
- No genera PDFs
- No envía correos de confirmación

---

## Estructura interna

```text
documents/
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

## Flujo principal (ejemplo: cargue de documento)

1. El cliente solicita cargar un documento
2. El sistema valida:
   - tipo
   - tamaño
   - relación con dominio

3. Se registra el documento
4. Se delega el almacenamiento al servicio externo
5. Se emite evento de documento cargado

---

## Modelos principales

- Documento
- TipoDocumento
- EstadoDocumento
- Relación (usuario / turno / PQRS)

---

## Estados típicos de documento

- `UPLOADED`
- `PROCESSING`
- `AVAILABLE`
- `REJECTED`
- `ARCHIVED`

---

## Eventos

Eventos relevantes:

- `DocumentUploaded`
- `DocumentValidated`
- `DocumentAvailable`

Consumidos por:

- notificaciones
- plantillas
- auditoría

---

## Consideraciones de diseño

- El documento es **metadata + referencia**
- El storage es reemplazable
- El dominio nunca depende de S3/MinIO
- Eventos garantizan desacoplamiento

---

## Escalabilidad

Este módulo permite:

- Múltiples backends de storage
- Procesamiento asíncrono
- Integración con OCR o validaciones futuras

```

```
