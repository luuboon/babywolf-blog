# Práctica de Laboratorio 2
## Display de Página Web — Especificaciones de Distribución
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Design System: "Dark Gaming / Cyberpunk Minimal"

| Token | Valor | Uso |
|-------|-------|-----|
| Space Void (body) | `#1a1a2e` | Fondo base |
| Elevation (cards/nav) | `#16213e` | Superficies elevadas |
| Action Red (acento) | `#e94560` | CTAs, hover, glows |
| Text White Alpha | `rgba(255,255,255,0.9)` | Texto principal |

**Tipografía:** `Inter, system-ui, -apple-system, ...` (fuentes del SO)

**Efectos:**
- Glassmorphism en Header: `backdrop-filter: blur(10px)`
- Neón en botones: `box-shadow: 0 4px 15px rgba(233,69,96,0.3)`
- Skeleton loaders en `/profile` y Admin (no spinners)

---

## 2. Breakpoints y dispositivos

| Dispositivo | Resolución | Breakpoint | Ajustes |
|-------------|-----------|------------|---------|
| Wearable | < 400px | `< 400px` | Solo título + CTA; sin animaciones |
| Móvil | 375–428px | `400–767px` | Menú hamburguesa, 1 columna |
| Tablet | 768–1023px | `768–1023px` | 2 columnas, sidebar colapsable |
| Desktop | 1024–1439px | `1024–1439px` | Layout completo, 3 columnas |
| Smart TV | ≥ 1440px | `≥ 1440px` | Fuente ≥ 24px, focus-visible para control remoto |

---

## 3. Rutas del frontend (app.routes.ts)

| Ruta | Componente | Guard |
|------|-----------|-------|
| `/` | HomeComponent | — |
| `/posts` | PostsPage | — |
| `/posts/:slug` | PostDetailPage | — |
| `/category/:name` | CategoryPostsPage | — |
| `/search` | SearchResultsPage | — |
| `/login` | LoginComponent | — |
| `/register` | RegisterComponent | — |
| `/profile` | ProfilePage | AuthGuard |
| `/editor/new` | UserEditorComponent | AuthGuard |
| `/editor/:id` | UserEditorComponent | AuthGuard |
| `/admin` | AdminDashboard | AuthGuard |
| `/admin/posts` | PostManagement | AuthGuard |
| `/admin/users` | UserManagement | AuthGuard |
| `/admin/comments` | CommentManagement | AuthGuard |
| `/admin/tareas` | TareasComponent | AuthGuard |
| `/admin/dom-demo` | DomDemoComponent | AuthGuard |

---

## 4. Contenido dinámico

| Elemento | Fuente | Comportamiento |
|----------|--------|---------------|
| Grid de posts | `GET /api/posts` | SSR + filtro client-side |
| Post individual | `GET /api/posts/:slug` | SSR completo para SEO |
| Búsqueda | `SearchStateService` + Supabase | Reactiva con Angular Signals |
| Auth state | `AuthService` | Signal reactivo global |
| Skeleton loaders | — | Reemplazan spinners en /profile y Admin |

---

## 5. Popups implementados

| Tipo | Disparador | Contenido |
|------|-----------|-----------|
| Login | `/login` (ruta propia) | Email + contraseña |
| Registro | `/register` | Email + contraseña + username |
| 2FA TOTP | Post-login si está activado | 6 dígitos (Google Authenticator) |
| Confirmación | Admin → borrar post/usuario | Confirm/cancel |

---
*Práctica de Laboratorio 2 — Desarrollo y Gestión de Software*
