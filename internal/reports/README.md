# reports/

## Propósito

El módulo `reports` es responsable de la **generación, consulta y exportación de reportes del sistema**.

Su función es **agregar y transformar información de distintos dominios** en vistas útiles para análisis, control operativo y toma de decisiones, sin alterar el estado del sistema.

---

## Responsabilidades

Este módulo se encarga de:

- Definir reportes disponibles en el sistema
- Ejecutar consultas agregadas
- Aplicar filtros y rangos de tiempo
- Generar salidas en formatos descargables (CSV)
- Exponer endpoints de consulta de reportes

---

## Qué NO hace este módulo

- No modifica datos de negocio
- No gestiona estados de otros módulos
- No ejecuta lógica operativa
- No envía notificaciones
- No reemplaza sistemas de BI externos

---

## Estructura interna

```text
reports/
├── handlers/
├── service/
├── repository/
├── models/
├── dto/
├── routes.go
└── module.go
```

> Nota: este módulo **no emite eventos** por defecto, ya que es de lectura y análisis.

---

## Tipos de reportes

Ejemplos de reportes soportados:

- Turnos atendidos por período
- PQRS por tipo y estado
- Documentos cargados / generados
- Feedback promedio por servicio
- Actividad general del sistema

Cada reporte tiene:

- definición clara
- filtros permitidos
- formato de salida

---

## Flujo principal (ejemplo: generación de reporte)

1. El usuario solicita un reporte
2. El handler valida parámetros (fechas, filtros)
3. El servicio:
   - valida reglas de acceso
   - ejecuta consultas agregadas

4. Se transforma el resultado
5. Se devuelve:
   - JSON (consulta)
   - CSV (descarga)

---

## Modelos principales

- Reporte
- FiltroReporte
- ResultadoReporte
- MetadataReporte

---

## Reglas importantes

- Los reportes son **read-only**
- Las consultas deben ser eficientes
- No se exponen datos sensibles sin control
- El servicio centraliza validaciones

---

## Dependencias

El módulo depende de:

- Repositorios de lectura
- Conexiones a base de datos
- Utilidades de exportación (CSV)

No depende directamente de otros módulos de dominio.

---

## Escalabilidad y rendimiento

Este diseño permite:

- Optimizar consultas de forma independiente
- Usar vistas materializadas en el futuro
- Mover reportes pesados a procesos batch
- Integrar con herramientas externas de análisis

---

## Relación con otros módulos

- **appointments**: datos de turnos
- **pqrs**: métricas de atención
- **documents**: actividad documental
- **feedback**: indicadores de calidad
- **auth**: control de acceso

La relación es **siempre de lectura**.

---

## Resumen

- `reports` transforma datos en información útil
- Es un módulo pasivo y seguro
- No altera el dominio
- Está diseñado para crecer sin impactar la operación

Con este módulo se completa la visión analítica del sistema.
