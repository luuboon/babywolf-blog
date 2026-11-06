# Scripts de la base (Supabase)

Se ejecutan **una vez, en orden**, en Supabase → SQL Editor. Todos son
repetibles. Estado en producción y reglas para cambios futuros:
[`SEGURIDAD_CUENTAS.md`](https://github.com/luuboon/babywolf-ecosystem/blob/main/docs/SEGURIDAD_CUENTAS.md).

| Script | Para qué |
|---|---|
| `001`–`004` | Esquema, categorías, almacenamiento, auditoría y bloqueo por intentos |
| `005_soporte.sql` | Mesa de ayuda |
| `006_incidencias_rls.sql` | Reportes con dueño y RLS |
| `007_app_lectores.sql` | Reportes con prioridad y versión, admin los atiende, comentarios con reglas (aplicado 06/11/2026) |
| `008_permisos_usuarios.sql` | **Imprescindible.** Sin él, cualquier usuario con sesión puede ponerse `role = 'admin'` (aplicado 06/11/2026) |

Pruebas contra Postgres con los roles de Supabase: `../pruebas/README.md`.
Una tabla nueva sin RLS, o un UPDATE sin permisos por columna, deja la base
abierta a la llave pública: revisa antes de aplicar.
