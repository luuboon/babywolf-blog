# Práctica de Laboratorio 1
## Selección de Framework para Desarrollo Web
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Repositorio:** https://github.com/luuboon/babywolf-blog
**Autor:** Abraham

---

## 1. Descripción del Proyecto y Necesidades del Cliente

### 1.1 Perfil del proyecto
**BabyWolf Blog** es una plataforma de blog con estética *Dark Gaming / Cyberpunk Minimal*, orientada a audiencias interesadas en videojuegos, hardware y cultura pop. Permite a usuarios registrados publicar artículos, dejar comentarios y administrar contenido desde un panel dedicado.

### 1.2 Requerimientos funcionales identificados
| ID | Requerimiento |
|----|--------------|
| RF-01 | Listado y lectura pública de posts por slug y categoría |
| RF-02 | Búsqueda de posts en tiempo real |
| RF-03 | Sistema de comentarios por post |
| RF-04 | Autenticación de usuarios (registro, login, 2FA TOTP) |
| RF-05 | Panel de administración: gestión de posts, usuarios y comentarios |
| RF-06 | Editor de posts con soporte de imágenes (upload a Storage) |
| RF-07 | Notificaciones por correo en cada inicio de sesión (Resend API) |
| RF-08 | Perfil de usuario editable con avatar |

### 1.3 Requerimientos no funcionales identificados
| ID | Requerimiento |
|----|--------------|
| RNF-01 | SEO: el contenido debe ser indexable por motores de búsqueda |
| RNF-02 | Responsivo: móvil, tablet y desktop |
| RNF-03 | Seguridad: JWT validado en backend, RLS en base de datos |
| RNF-04 | Despliegue continuo automático desde GitHub |
| RNF-05 | Control de versiones obligatorio |
| RNF-06 | Rendimiento: bajo tiempo de respuesta bajo carga concurrente |

---

## 2. Análisis y Selección de Framework

### 2.1 Frontend — Opciones evaluadas

| Criterio | **Angular 21+ (Elegido)** | React + Vite | Vue + Nuxt |
|----------|--------------------------|-------------|------------|
| SSR / SSG nativo | ✅ Angular Universal | ⚠️ Requiere configuración | ✅ Nuxt nativo |
| SEO | ✅ SSR incluido | ⚠️ Manual | ✅ |
| Tipado estático | ✅ TypeScript nativo | ⚠️ Opcional | ⚠️ Opcional |
| Arquitectura escalable | ✅ Módulos + Standalones | ⚠️ Depende del desarrollador | ✅ |
| Manejo de estado reactivo | ✅ Signals + RxJS | ⚠️ Múltiples opciones externas | ✅ Pinia |
| Lazy loading integrado | ✅ `loadComponent()` nativo | ⚠️ Manual | ✅ |
| Inyección de dependencias | ✅ DI nativo | ❌ No tiene DI nativo | ❌ |

**Decisión Frontend: ✅ Angular 21+ (Standalone Components + SSR)**

### 2.2 Backend — Opciones evaluadas

| Criterio | **Go 1.25 + Gin (Elegido)** | Node.js + Express | Python + FastAPI |
|----------|-----------------------------|-------------------|-----------------|
| Rendimiento / concurrencia | ⭐⭐⭐⭐⭐ Goroutines nativas | ⭐⭐⭐ Event loop | ⭐⭐⭐⭐ |
| Tipado estático | ✅ | ❌ (TS opcional) | ⚠️ |
| Binario de despliegue único | ✅ | ❌ | ❌ |
| ORM maduro | ✅ GORM | ✅ Prisma | ✅ SQLAlchemy |

**Decisión Backend: ✅ Go 1.25.7 + Gin + GORM**

### 2.3 Infraestructura real del proyecto

| Servicio | Rol |
|---------|-----|
| **Supabase PostgreSQL 15** | Base de datos, RLS, triggers |
| **Supabase Auth** | JWT, 2FA TOTP, sesiones |
| **Supabase Storage** | Imágenes (S3-compatible) |
| **Resend API** | Alertas de inicio de sesión |
| **Vercel** | Deploy frontend con CD automático |

---

## 3. Arquitectura real del sistema

```
Vercel Edge (Angular SSR)
        │ HTTP/REST
Go 1.25 + Gin API
        │
Supabase (PostgreSQL + Auth + Storage)
```

---

## 4. Arquitectura limpia (Clean Architecture)

**Backend:**
```
internal/domain/entities/       ← Post, User, Comment
internal/application/usecases/  ← ManagePostsUseCase
internal/infrastructure/        ← GORM, Supabase, middleware JWT
internal/interfaces/controllers/ ← PostController, UserController
```

**Frontend:**
```
features/posts/domain/          ← post.model.ts, post.repository.ts
features/posts/application/     ← usecases
features/posts/infrastructure/  ← supabase-post.repository.ts
features/posts/presentation/    ← components, pages, state
```

---

## 5. Estrategia de Versionamiento

**Git Flow adaptado:**
- `main` → producción (CD automático a Vercel)
- `feature/*` → nuevas funcionalidades
- `hotfix/*` → correcciones urgentes

**Convención de commits (Conventional Commits):**
```
feat(posts): agregar filtro por categoría
fix(auth): corregir redirección tras login
refactor(middleware): extraer validación de rol
```

**SemVer:** `MAJOR.MINOR.PATCH` — ejemplo: `v1.2.0`

---

## 6. Stack tecnológico final

| Capa | Tecnología | Versión |
|------|-----------|---------|
| Frontend | Angular + SSR | 21.x |
| Lenguaje FE | TypeScript | 5.9 |
| Backend | Go + Gin | 1.25.7 |
| ORM | GORM | 1.31 |
| BaaS | Supabase | v2 |
| Email | Resend API | — |
| Deploy | Vercel | — |
| Versión | Git + GitHub | — |

---
*Práctica de Laboratorio 1 — Desarrollo y Gestión de Software*
