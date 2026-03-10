# cmd/api/

## Propósito

La carpeta `cmd/api/` contiene el **punto de entrada del servidor HTTP** del sistema.

Aquí se **inicializa, configura y ensambla** la aplicación backend expuesta vía API, utilizando **Go Fiber** como framework web.

Este módulo **no contiene lógica de negocio**; su responsabilidad es **orquestar dependencias** y levantar el servidor.

---

## Responsabilidades

Este módulo se encarga de:

- Cargar configuración del sistema
- Inicializar infraestructura compartida
- Construir e inyectar dependencias
- Registrar middlewares globales
- Registrar rutas de todos los módulos
- Iniciar el servidor HTTP

---

## Qué NO hace este módulo

- No implementa reglas de negocio
- No accede directamente a la base de datos para operaciones funcionales
- No define flujos de dominio
- No contiene lógica específica de módulos
- No procesa tareas asíncronas

---

## Estructura típica

```text
cmd/api/
├── main.go
├── server.go        # Inicialización de Fiber
├── routes.go        # Registro global de rutas
└── middleware.go    # Middlewares compartidos
```

> La estructura puede variar, pero las responsabilidades deben mantenerse claras.

---

## Flujo de arranque

1. Ejecución de `main.go`
2. Carga de variables de entorno
3. Inicialización de configuración
4. Inicialización de infraestructura:
   - base de datos
   - mensajería
   - storage
   - logger

5. Construcción de módulos de dominio
6. Registro de rutas
7. Inicio del servidor HTTP

---

## Registro de módulos

Cada módulo de `internal/` expone una función de registro, por ejemplo:

- `auth.RegisterRoutes(app)`
- `appointments.RegisterRoutes(app)`
- `documents.RegisterRoutes(app)`

El API **no conoce la lógica interna** del módulo, solo su contrato público.

---

## Middlewares globales

Ejemplos de middlewares definidos aquí:

- Logging de requests
- Manejo centralizado de errores
- Autenticación / autorización
- CORS
- Rate limiting (si aplica)

Los middlewares:

- Son transversales
- No contienen lógica de negocio

---

## Manejo de errores

- Los errores del dominio se traducen a HTTP aquí
- Se centraliza el formato de respuestas
- Se evita duplicar lógica de error en handlers

---

## Configuración

Toda configuración proviene de:

- Variables de entorno
- Archivos `.env`
- Structs definidos en `pkg/config`

El API **no define configuración propia**, solo consume.

---

## Escalabilidad

Este diseño permite:

- Escalar el API de forma horizontal
- Separar workers sin impacto
- Migrar a microservicios en el futuro
- Añadir gateways o balanceadores sin cambios internos

---

## Relación con otras carpetas

```text
cmd/api/   → arranque HTTP
internal/  → lógica de negocio
pkg/       → infraestructura
```

`cmd/api/` **depende de todas**, pero ninguna depende de él.

---

## Resumen

- `cmd/api/` es el **orquestador del backend**
- Solo ensambla y arranca
- No contiene reglas de negocio
- Es el punto de integración del sistema

Cualquier cambio en infraestructura o despliegue comienza aquí.
