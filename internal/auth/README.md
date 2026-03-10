## 📁 `internal/auth/README.md`

# auth/

## Propósito

El módulo `auth` es responsable de la **autenticación, autorización y gestión de identidades** del sistema.

Define **quién es el usuario**, **qué puede hacer** y **bajo qué contexto** accede a los recursos del sistema.

---

## Responsabilidades

Este módulo se encarga de:

- Registro y gestión de usuarios
- Autenticación (login / logout)
- Emisión y validación de tokens (JWT u otro mecanismo)
- Gestión de roles y permisos
- Control de acceso a endpoints protegidos

---

## Qué NO hace este módulo

- No gestiona lógica de negocio de otros dominios
- No envía correos directamente
- No conoce detalles de infraestructura (SMTP, Kafka, DB concreta)
- No define reglas de turnos, documentos o PQRS

---

## Estructura interna

```text
auth/
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

## Flujo principal (ejemplo: login)

1. El cliente envía credenciales al endpoint de autenticación
2. El handler valida el input
3. El servicio:
   - verifica credenciales
   - valida estado del usuario
   - genera token

4. Se retorna el token y metadata del usuario
5. Se emite un evento de autenticación (opcional)

---

## Modelos principales

- Usuario
- Rol
- Permiso
- Sesión / Token

---

## Eventos

Eventos típicos emitidos por este módulo:

- `UserAuthenticated`
- `UserCreated`
- `RoleAssigned`

Estos eventos permiten desacoplar:

- notificaciones
- auditoría
- métricas

---

## Dependencias

- Repositorio de persistencia (interface)
- Servicio de hashing (interface)
- Generador/validador de tokens (interface)

Todas las dependencias se inyectan desde capas superiores.

---

## Consideraciones de seguridad

- No exponer datos sensibles en DTOs
- Tokens con expiración definida
- Roles y permisos evaluados en middleware
- Auditoría mediante eventos

```

```
