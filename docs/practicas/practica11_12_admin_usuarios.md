# Práctica de Laboratorio 11-12
## Módulo administrativo de gestión de usuarios
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Objetivo

Implementar un módulo administrativo para gestionar usuarios aplicando principios de seguridad informática y administración de privilegios: alta, edición, cambio de contraseña, asignación de roles, activación/desactivación de cuentas, auditoría y controles de seguridad.

## 2. Arquitectura de la solución

El proyecto ya tenía dos piezas construidas por separado que esta práctica conecta:

- **Frontend Angular** (`frontend/babywolf-frontend`) hablando directo con **Supabase** (Postgres + Auth) para todo lo que un usuario hace sobre **sus propios datos** (perfil, contraseña propia, posts).
- **Backend Go** (`backend/`) con la `SUPABASE_SERVICE_ROLE_KEY`, que es la única pieza autorizada a actuar **sobre otros usuarios** (crear cuentas, resetear contraseñas ajenas, cambiar roles, activar/desactivar). Esto es necesario porque la tabla `public.users` tiene Row Level Security y sus políticas (`supabase_setup_and_seed.sql`) sólo permiten a cada usuario leer/editar su propia fila — un admin no puede, por RLS, modificar la fila de otro usuario directamente desde el navegador sin la service role key.

```
Angular (anon key, respeta RLS)  ─────►  Supabase (Postgres + Auth)
        │
        └── operaciones sobre OTROS usuarios
                  │
                  ▼
        Backend Go (service_role key, RLS bypass)  ─────►  Supabase Admin API / Postgres directo
```

## 3. Parte 1 — Administración de usuarios

| Operación | Dónde vive | Archivo |
|---|---|---|
| Alta de usuario | Backend (Supabase Admin API) | `usecases/manage_users.go` → `CreateUserByAdmin`, `controllers/user_controller.go` → `CreateUser` |
| Edición por el propio usuario (username, avatar) | Frontend, sobre el propio perfil | `features/user/profile/profile.page.ts` |
| Edición por el admin (username de otro usuario) | Backend | `UpdateUsername` (usecase) / botón "✏️ Editar" en `user-management.component` |
| Eliminación lógica | Backend, columna `active` | `SetActive` (usecase) / `toggleActive()` (frontend) |
| Activar/desactivar | Backend | Mismo endpoint que eliminación lógica: `PATCH /api/admin/users/:id/active` |

La "eliminación" nunca es un `DELETE` físico: se marca `active = false` en `public.users`, conservando historial y evitando romper llaves foráneas (`posts.author_id`, `comments.author_id`).

## 4. Parte 2 — Gestión de roles

- `PATCH /api/admin/users/:id/role` (`UpdateRole` en `manage_users.go`) valida que quien ejecuta la acción sea `admin` **consultando la base de datos**, no el JWT, evitando que un token desactualizado otorgue privilegios indebidos.
- El rol se muestra en la UI (`user-management.component.html`) con badges y se puede alternar `user` ⇄ `admin` con un botón.
- Los permisos asociados a cada rol están codificados en el propio flujo de la app: `admin` accede a `/admin/*` (guard `AuthGuard` + verificación de rol en `AuthService.isAdmin$`); `user` sólo accede a `/profile` y `/editor/*` sobre su propio contenido.

## 5. Parte 3 — Auditoría

Tabla `public.audit_log` (`backend/sql/004_security_audit.sql`), con columnas: `user_id`, `email`, `action`, `ip_address`, `details`, `created_at`.

Eventos registrados:

| Acción | Disparador |
|---|---|
| `login_success` / `login_failed` | `AuthService.signIn()` vía RPC `record_successful_login` / `record_failed_login` |
| `logout` | `AuthService.signOut()` vía RPC `record_logout` |
| `password_change` | `AuthService.changePassword()` vía RPC `record_password_change` |
| `user_created` | Backend Go, `CreateUserByAdmin` |
| `user_activated` / `user_deactivated` | Backend Go, `SetActive` |
| `role_changed` | Backend Go, `UpdateRole` |

Consulta: `/admin/audit-log` (`AuditLogComponent`), tabla ordenada por fecha descendente, sólo visible para admins (política RLS `"Admins can read audit log"`).

## 6. Parte 4 — Controles de seguridad implementados

| # | Control | Implementación |
|---|---|---|
| 1 | Contraseñas cifradas | Delegado a Supabase Auth (bcrypt), nunca se almacena texto plano en `public.users` |
| 2 | Validación de formularios | `Validators.pattern(STRONG_PASSWORD_PATTERN)` en registro/cambio de contraseña; `binding:"required,email"` / `min=8` en los DTOs del backend Go |
| 3 | Protección CSRF | La API usa JWT Bearer (`Authorization` header), no cookies de sesión — el vector clásico de CSRF (envío automático de cookies por el navegador) no aplica |
| 4 | Expiración de sesión | JWT de Supabase con `exp` + `SessionTimeoutService` (cierra sesión tras 20 min de inactividad del usuario) |
| 5 | Bloqueo temporal tras intentos fallidos | Funciones RPC `check_account_lock` / `record_failed_login` en Postgres: bloquean la cuenta 15 minutos tras 5 intentos fallidos |
| 6 | Política de contraseña segura | `STRONG_PASSWORD_PATTERN`: mínimo 8 caracteres, mayúscula, minúscula, número y símbolo — aplicada en registro, cambio propio, alta por admin y reset por admin |

## 7. Manual técnico (resumen de despliegue)

1. Ejecutar en Supabase SQL Editor, en orden: `backend/sql/001_initial_schema.sql` (o `supabase_setup_and_seed.sql` si es la primera vez), `002_add_category.sql`, `003_storage_policies.sql`, `004_security_audit.sql`, y `supabase_auth_trigger.sql`.
2. Backend Go: copiar `backend/.env.example` a `.env`, completar `SUPABASE_URL`, `SUPABASE_SERVICE_ROLE_KEY`, `JWT_SECRET`, `DATABASE_URL`. Ejecutar `go run cmd/api/main.go` (por defecto en `:8080`).
3. Frontend: en `src/environments/environment.ts`, `backendUrl` debe apuntar al backend Go (`http://localhost:8080/api` en desarrollo). Ejecutar `npm start` / `ng serve`.
4. El primer usuario admin se promueve manualmente en Supabase: `update public.users set role = 'admin' where email = '...';`

## 8. Manual de usuario

**Administrador** (`/admin/users`):
- **Nuevo Usuario**: botón "+ Nuevo Usuario", completar email/usuario/contraseña/rol, "Crear usuario".
- **Editar usuario**: botón "✏️ Editar" en la fila, modificar el nombre de usuario y "Guardar".
- **Cambiar rol**: botón "Hacer Admin" / "Quitar Admin" en la fila del usuario.
- **Resetear contraseña**: botón "🔑 Reset Password", escribir la nueva contraseña y "Guardar".
- **Activar/Desactivar**: botón "🚫 Desactivar" / "✅ Activar" (eliminación lógica).
- **Bitácora**: botón "📜 Bitácora" en la parte superior, o desde el menú lateral.

**Usuario normal** (`/profile`):
- **Editar perfil**: cambiar nombre de usuario o foto y "Actualizar Perfil".
- **Cambiar contraseña**: sección "🔑 Cambiar Contraseña", escribir y confirmar la nueva contraseña.
- **Iniciar sesión**: `/login`; tras 5 intentos fallidos la cuenta se bloquea 15 minutos.

## 9. Evidencias

> Anexar aquí:
> 1. Capturas del flujo de alta de usuario desde `/admin/users`.
> 2. Captura de un cambio de rol y su reflejo en `/admin/audit-log`.
> 3. Captura de una cuenta desactivada intentando iniciar sesión (mensaje de error).
> 4. Captura de 5 intentos fallidos de login mostrando el bloqueo temporal.
> 5. Captura de la tabla `audit_log` en el Supabase Table Editor.
> 6. Video demostrativo recorriendo las 4 partes (alta, roles, auditoría, seguridad).

---
*Práctica de Laboratorio 11-12 — Desarrollo y Gestión de Software*
