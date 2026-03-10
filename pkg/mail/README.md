# pkg/mail

Cliente SMTP reutilizable para envio de correos.

## Compatibilidad con Google Workspace (cuentas institucionales)

Este cliente soporta tres modos:

- `plain`: usuario + app password.
- `xoauth2`: usuario + access token OAuth2.
- `none`: sin autenticacion (relay SMTP institucional).

Para Google, usa:

- `smtp.gmail.com:587` con `STARTTLS`, o
- `smtp.gmail.com:465` con TLS implicito.

En ambos casos, TLS debe estar habilitado.

## Variables de entorno

- `SMTP_HOST` (default: `smtp.gmail.com`)
- `SMTP_PORT` (default: `587`)
- `SMTP_AUTH_METHOD` (`plain`, `xoauth2`, `none`)
- `SMTP_USERNAME`
- `SMTP_PASSWORD` (solo `plain`)
- `SMTP_ACCESS_TOKEN` (solo `xoauth2`)
- `SMTP_FROM_NAME`
- `SMTP_FROM_ADDRESS` (si no se define, usa `SMTP_USERNAME`)
- `SMTP_USE_STARTTLS` (default: `true`)
- `SMTP_USE_IMPLICIT_TLS` (default: `false`)
- `SMTP_REQUIRE_TLS` (default: `true`)
- `SMTP_TIMEOUT_SECONDS` (default: `10`)
- `SMTP_INSECURE_SKIP_VERIFY` (default: `false`, solo desarrollo)

## Worker

El entrypoint `cmd/workers/mailer/main.go` carga esta configuracion y valida
el cliente SMTP al iniciar.
