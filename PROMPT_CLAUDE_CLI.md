Estoy en la raíz del repositorio `babywolf-blog`. Necesito que hagas lo siguiente de forma completa y sin preguntar. Lee los archivos existentes antes de modificarlos.

---

## TAREA GENERAL

Integrar 6 prácticas de laboratorio al proyecto. Implica:
1. Crear carpeta `docs/` con documentación en markdown
2. Crear dos componentes Angular en `/admin`: gestor de tareas (Lab 6) y personalizador DOM (Lab 5)
3. Registrar las rutas nuevas en `app.routes.ts`
4. Agregar los links en `admin-layout.ts`

---

## PASO 1 — Crear estructura de documentación

Crea estas carpetas y archivos:

```
docs/practicas/
docs/prototipos/     ← vacía por ahora, solo créala
```

---

### docs/practicas/practica1_seleccion_framework.md

```markdown
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
```

---

### docs/practicas/practica2_especificaciones_display.md

```markdown
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
```

---

### docs/practicas/practica3_librerias_compatibilidad.md

```markdown
# Práctica de Laboratorio 3
## Librerías, Compatibilidad de Navegadores y Desempeño
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Dependencias de producción (package.json)

| Librería | Versión | Rol |
|---------|---------|-----|
| `@angular/core` | 21.1.0 | Framework principal — Signals, DI, change detection |
| `@angular/common` | 21.1.0 | Pipes, `@if`, `@for`, `AsyncPipe` |
| `@angular/router` | 21.1.0 | Routing SPA con lazy loading |
| `@angular/forms` | 21.1.0 | Formularios reactivos |
| `@angular/animations` | 21.1.3 | Animaciones declarativas |
| `@angular/platform-server` | 21.1.0 | SSR |
| `@angular/ssr` | 21.1.2 | Bundle SSR + Express |
| `@angular/material` | 21.1.3 | MatToolbar, MatButton (uso mínimo) |
| `@angular/cdk` | ~21.1.3 | Overlays, accesibilidad |
| `@supabase/supabase-js` | ^2.100.0 | Auth, queries PostgreSQL, Storage, Realtime |
| `rxjs` | ~7.8.0 | Observables, BehaviorSubject |
| `express` | ^5.1.0 | Servidor Node.js para SSR |
| `tslib` | ^2.3.0 | Helpers TypeScript compilado |

## 2. Dependencias de desarrollo

| Librería | Versión | Rol |
|---------|---------|-----|
| `@angular/cli` | ^21.1.2 | Toolchain de construcción |
| `typescript` | ~5.9.2 | Tipado estático |
| `vitest` | ^4.0.8 | Pruebas unitarias (reemplaza Karma) |
| `jsdom` | ^27.1.0 | DOM virtual para pruebas |

---

## 3. Matriz de compatibilidad

| Dispositivo | Angular 21 | Material | Supabase JS | RxJS | Web Push | Veredicto |
|------------|-----------|----------|-------------|------|----------|-----------|
| Chrome 115+ Desktop | ✅ | ✅ | ✅ | ✅ | ✅ | Completo |
| Firefox 115+ | ✅ | ✅ | ✅ | ✅ | ✅ | Completo |
| Safari 16.4+ | ✅ | ✅ | ✅ | ✅ | Solo PWA | Completo* |
| Edge 115+ | ✅ | ✅ | ✅ | ✅ | ✅ | Completo |
| Chrome Android | ✅ | ✅ | ✅ | ✅ | ✅ | Completo |
| Safari iOS 16.4+ | ✅ | ✅ | ✅ | ✅ | Solo PWA | Completo* |
| Smart TV Tizen 6+ | SSR HTML | Básico | Sin Realtime | ✅ | ❌ | Parcial |
| Smart TV WebOS 6+ | SSR HTML | Básico | Sin Realtime | ✅ | ❌ | Parcial |
| Android TV | ✅ | ✅ | ✅ | ✅ | ❌ | Aceptable |
| Wear OS | ❌ Pantalla | ❌ | ❌ | ✅ | ❌ | No viable |
| Apple Watch | ❌ ES2022 | ❌ | ❌ | ❌ | ❌ | No compatible |

*Sin push web en navegador Safari; funciona si se instala como PWA.

---

## 4. Notificaciones — Sistema implementado

El proyecto **no usa Web Push**. Usa **Resend API vía email**, implementado en `auth.service.ts`:

```typescript
// Al detectar SIGNED_IN en onAuthStateChange:
fetch('/api/notify-login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email: session.user.email, timestamp: new Date().toISOString() })
});
```

**Ventaja:** Funciona en 100% de dispositivos sin importar el navegador.

| Tipo | Desktop | iOS Safari | Android | Smart TV | Wearable |
|------|---------|-----------|---------|----------|---------|
| Email (Resend) ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Web Push (no impl.) | ✅ | Solo PWA | ✅ | ❌ | ❌ |

---

## 5. Desempeño móvil — Estrategias implementadas

| Técnica | Implementación | Impacto |
|---------|--------------|--------|
| SSR | `@angular/ssr` + `server.ts` | LCP mejorado |
| Lazy loading | `loadComponent()` en cada ruta | Bundle inicial reducido |
| Compilación AOT | `@angular/build` (Vite) | Bundle 30–40% más pequeño |
| Signals + `computed()` | `PostStateService`, `SearchStateService` | Re-renders granulares |
| Skeleton loaders | Admin, `/profile` | Percepción de velocidad |
| `finalize()` en RxJS | `PostStateService.loadPosts()` | Evita memory leaks |

---
*Práctica de Laboratorio 3 — Desarrollo y Gestión de Software*
```

---

### docs/practicas/practica4_buscador_pruebas.md

```markdown
# Práctica de Laboratorio 4
## Buscador Interno del Sitio Web — Implementación y Pruebas
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Arquitectura del buscador

El buscador es **client-side**: descarga todos los posts publicados desde Supabase y filtra en memoria.

```
Header input (input event)
    │
SearchStateService.setSearchQuery(query)  [Angular Signal]
    │
Navegación a /search?q=query
    │
SearchResultsPage.ngOnInit()
route.queryParams.subscribe(params)
    │
PostRepository.getPosts()  ← Supabase: SELECT * WHERE published=true
    │
posts.filter(p =>
  p.title.toLowerCase().includes(q)    ||
  p.excerpt.toLowerCase().includes(q)  ||
  p.content.toLowerCase().includes(q)  ||
  p.category.toLowerCase().includes(q)
)
    │
Renderizado con @for
```

## 2. Archivos involucrados

| Archivo | Rol |
|---------|-----|
| `core/services/search-state.service.ts` | Signal global del query |
| `core/layout/header/header.ts` | Captura input, llama `setSearchQuery()` |
| `features/search/search-results.ts` | Lee queryParams, filtra, renderiza |
| `features/posts/infrastructure/repositories/supabase-post.repository.ts` | `getPosts()` |

## 3. Casos de prueba

| ID | Query | Esperado | Campo | Resultado |
|----|-------|---------|-------|-----------|
| TC-01 | `rtx` | ≥ 1 resultado | título/contenido | ✅ PASA |
| TC-02 | `gaming` | ≥ 1 resultado | categoría/contenido | ✅ PASA |
| TC-03 | `ANGULAR` | ≥ 1 resultado | case-insensitive | ✅ PASA |
| TC-04 | `2026` | ≥ 1 resultado | título/excerpt | ✅ PASA |
| TC-05 | `hardware` | ≥ 1 resultado | categoría | ✅ PASA |
| TC-06 | ` ` (espacio) | 0 resultados | trim() vacío | ✅ PASA |
| TC-07 | `xkwqz` | 0 resultados | sin coincidencias | ✅ PASA |
| TC-08 | `a` | ≥ 1 resultado | carácter único | ✅ PASA |
| TC-09 | `mexico` | Condicional | sin normalización de acentos | ⚠️ |
| TC-10 | `borrador` | 0 resultados | RLS: published=false | ✅ PASA |

## 4. Limitaciones identificadas

| Limitación | Impacto |
|-----------|--------|
| Sin normalización de acentos ("mexico" ≠ "México") | Bajo-medio |
| Carga todos los posts en memoria | Medio-alto con 1000+ posts |
| Sin debounce en el input del header | Bajo |
| Sin paginación de resultados | Bajo |

## 5. Mejora recomendada (full-text search en Supabase)

```typescript
// En lugar de traer todos los posts y filtrar client-side:
this.sb.client
  .from('posts')
  .select('*')
  .textSearch('title', query, { type: 'websearch' })
  .eq('published', true)
```

---
*Práctica de Laboratorio 4 — Desarrollo y Gestión de Software*
```

---

### docs/practicas/practica5_dom_aplicacion.md

```markdown
# Práctica de Laboratorio 5
## Aplicación del DOM en Sitios Web
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Objetivo

Desarrollar una página web dinámica aplicando el DOM para seleccionar, modificar, crear y actualizar elementos de la interfaz mediante eventos del usuario.

## 2. Contexto

El componente `DomDemoComponent` en `/admin/dom-demo` es un **"Personalizador Dinámico"** del blog que demuestra en vivo cada operación del DOM mediante Angular — mostrando la equivalencia entre las APIs nativas y el framework.

---

## 3. Operaciones del DOM implementadas

### 3.1 Selección de elementos

| Método vanilla | Equivalente Angular |
|----------------|---------------------|
| `document.getElementById()` | `inject(ElementRef)` / `viewChild()` |
| `document.querySelector()` | `viewChild('ref')` |
| `document.querySelectorAll()` | `@for` loop sobre signal |
| `element.closest()` | Event delegation en template |

### 3.2 Modificación de contenido

| Método / Propiedad | Uso en el componente |
|--------------------|----------------------|
| `textContent` | `{{ tituloBlog() }}` — actualiza en tiempo real |
| `setAttribute()` | `[attr.data-id]`, `[attr.aria-label]` |
| `dataset` | `[attr.data-category]`, `[attr.data-priority]` |

### 3.3 Modificación de estilos

| Método | Uso |
|--------|-----|
| `classList.toggle('dark')` | `[class.dark]="temaOscuro()"` |
| `style.setProperty('--accent')` | `[style.--accent]="colorAcento()"` |
| `element.style.fontSize` | `[style.font-size.px]="tamFuente()"` |
| `classList.toggle('destacada')` | `[class.destacada]="card.destacada"` |

### 3.4 Creación y eliminación de estructura

| Método vanilla | Equivalente Angular |
|----------------|---------------------|
| `createElement + appendChild` | `cards.update(arr => [...arr, newCard])` + `@for` |
| `element.remove()` | `cards.update(arr => arr.filter(...))` |
| `cloneNode(true) + insertBefore` | spread `{...card}` + `splice` en signal |
| `replaceWith(input)` | `card.editando = true` → `@if` muestra `<input>` |
| `prepend()` | `bannerVisible.set(true)` → `@if` al inicio del template |

### 3.5 Eventos implementados

| Evento | Elemento | Acción |
|--------|----------|--------|
| `click` | Botones | Crear card, eliminar, duplicar, tema |
| `input` | Input texto | `actualizarTitulo()` en tiempo real |
| `input` | Slider range | `cambiarFuente()` |
| `change` | Color picker | `cambiarAcento()` |
| `dblclick` | Título card | `iniciarEdicion()` |
| `keydown` Enter/Escape | Input edición | `guardarEdicion()` / `cancelarEdicion()` |
| `blur` | Input edición | `guardarEdicion()` automático |

---

## 4. Comparación vanilla ↔ Angular

| Operación vanilla | Angular en el proyecto |
|-------------------|----------------------|
| `createElement + appendChild` | `@for (card of cards(); track card.id)` |
| `textContent = valor` | `{{ tituloBlog() }}` |
| `classList.toggle('dark')` | `[class.dark]="temaOscuro()"` |
| `addEventListener('click', fn)` | `(click)="fn()"` |
| `input.value` + evento `input` | `[(ngModel)]` + `(input)` |
| Actualizar contador manualmente | `computed(() => cards().length)` |
| `style.width = pct + '%'` | `[style.width]="progreso() + '%'"` |

**La diferencia clave:** en vanilla el desarrollador sincroniza manualmente el estado con el DOM; Angular lo hace automáticamente con change detection y Signals.

---
*Práctica de Laboratorio 5 — Desarrollo y Gestión de Software*
```

---

### docs/practicas/practica6_gestor_tareas_dom.md

```markdown
# Práctica de Laboratorio 6
## Funciones del DOM — Administrador de Tareas
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog (contexto editorial)
**Autor:** Abraham

---

## 1. Objetivo

Crear una aplicación web interactiva (administrador de tareas) usando funciones y métodos del DOM para creación, eliminación y modificación de elementos HTML, integrando eventos, validaciones y actualización dinámica.

## 2. Descripción

El componente `TareasComponent` en `/admin/tareas` es un **Gestor de Tareas Editoriales** para el flujo de trabajo del blog: tareas como "escribir reseña RTX 5090", "revisar borrador de Gloria", etc.

---

## 3. Funcionalidades implementadas

| # | Funcionalidad | Operación DOM (vanilla) | Equivalente Angular |
|---|--------------|------------------------|---------------------|
| 1 | Agregar tarea | `createElement + appendChild` | `signal.update(arr => [...arr, nueva])` |
| 2 | Eliminar tarea | `element.remove()` | `signal.update(arr => arr.filter(...))` |
| 3 | Marcar completada | `classList.toggle('done')` | `[class.completada]="t.completada"` |
| 4 | Edición in-place | `replaceWith(input)` | `t.editando = true` → `@if` muestra input |
| 5 | Filtros | `classList` + `querySelectorAll` | `computed()` filtra el signal |
| 6 | Prioridad | `dataset.priority` | `[attr.data-priority]="t.prioridad"` |
| 7 | Contadores | `textContent` manual | `{{ total() }}`, `{{ pendientes() }}` |
| 8 | Barra de progreso | `style.width = pct + '%'` | `[style.width]="progreso() + '%'"` |
| 9 | Limpiar completadas | iteración + `remove()` | `signal.update(arr => arr.filter(!done))` |
| 10 | Estado vacío | `createElement` condicional | `@if (total() === 0)` |

---

## 4. Validaciones implementadas

| Validación | Regla | Respuesta |
|-----------|-------|-----------|
| Campo vacío | `trim() === ''` | Error: "La tarea no puede estar vacía." |
| Longitud mínima | `< 3 caracteres` | Error: "Mínimo 3 caracteres." |
| Longitud máxima | `> 80 caracteres` | Error: "Máximo 80 caracteres." |
| Duplicados | Comparación case-insensitive | Error: "Esa tarea ya existe." |
| Edición vacía | `guardar('')` | Restaura texto original |

```typescript
private validar(texto: string): string {
  const t = texto.trim();
  if (!t)             return 'La tarea no puede estar vacía.';
  if (t.length < 3)   return 'Mínimo 3 caracteres.';
  if (t.length > 80)  return 'Máximo 80 caracteres.';
  const existe = this.tareas().some(
    ta => ta.texto.toLowerCase() === t.toLowerCase()
  );
  if (existe) return 'Esa tarea ya existe.';
  return '';
}
```

---

## 5. Eventos integrados

| Evento | Elemento | Acción |
|--------|----------|--------|
| `click` | Botón "Agregar" | Valida y crea tarea |
| `keydown` Enter | Input | Valida y crea tarea |
| `input` | Input | Limpia mensaje de error |
| `change` | Checkbox | `toggleCompletar()` |
| `dblclick` | Texto de tarea | `iniciarEdicion()` |
| `keydown` Enter | Input edición | `guardarEdicion()` |
| `keydown` Escape | Input edición | `cancelarEdicion()` |
| `blur` | Input edición | `guardarEdicion()` automático |
| `click` | Chips de filtro | `setFiltro()` |
| `click` | "Limpiar completadas" | Elimina en lote |

---

## 6. Actualización dinámica

Todos los nodos dependientes se actualizan automáticamente via `computed()`:

```typescript
total       = computed(() => this.tareas().length);
completadas = computed(() => this.tareas().filter(t => t.completada).length);
pendientes  = computed(() => this.total() - this.completadas());
progreso    = computed(() =>
  this.total() ? Math.round((this.completadas() / this.total()) * 100) : 0
);
tareasFiltradas = computed(() => {
  const f = this.filtro();
  return this.tareas().filter(t =>
    f === 'todas'       ? true :
    f === 'pendientes'  ? !t.completada :
                           t.completada
  );
});
```

---
*Práctica de Laboratorio 6 — Desarrollo y Gestión de Software*
```

---

## PASO 2 — Crear componente TareasComponent (Lab 6)

Crea estos tres archivos en `frontend/babywolf-frontend/src/app/features/admin/tareas/`:

### tareas.component.ts

```typescript
import {
  Component,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

export type Prioridad = 'alta' | 'media' | 'baja';
export type Filtro    = 'todas' | 'pendientes' | 'completadas';

export interface Tarea {
  id: number;
  texto: string;
  prioridad: Prioridad;
  completada: boolean;
  editando: boolean;
  textoEdicion: string;
}

@Component({
  selector: 'app-tareas',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './tareas.component.html',
  styleUrls: ['./tareas.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TareasComponent {

  private nextId = 4;

  tareas = signal<Tarea[]>([
    this.nuevaTarea('Escribir reseña de la RTX 5090',              'alta'),
    this.nuevaTarea('Revisar borrador de Gloria sobre indies',      'media'),
    this.nuevaTarea('Subir imágenes de portada a Supabase Storage', 'baja'),
  ]);

  filtro       = signal<Filtro>('todas');
  errorMsg     = signal<string>('');
  nuevoTexto   = '';
  nuevaPrioridad: Prioridad = 'media';

  total       = computed(() => this.tareas().length);
  completadas = computed(() => this.tareas().filter(t => t.completada).length);
  pendientes  = computed(() => this.total() - this.completadas());
  progreso    = computed(() =>
    this.total() ? Math.round((this.completadas() / this.total()) * 100) : 0
  );

  tareasFiltradas = computed(() => {
    const f = this.filtro();
    return this.tareas().filter(t =>
      f === 'todas'       ? true :
      f === 'pendientes'  ? !t.completada :
                             t.completada
    );
  });

  private validar(texto: string): string {
    const t = texto.trim();
    if (!t)             return 'La tarea no puede estar vacía.';
    if (t.length < 3)   return 'Mínimo 3 caracteres.';
    if (t.length > 80)  return 'Máximo 80 caracteres.';
    const existe = this.tareas().some(
      ta => ta.texto.toLowerCase() === t.toLowerCase()
    );
    if (existe) return 'Esa tarea ya existe.';
    return '';
  }

  agregar(): void {
    const error = this.validar(this.nuevoTexto);
    if (error) { this.errorMsg.set(error); return; }
    this.errorMsg.set('');
    this.tareas.update(arr => [
      ...arr,
      this.nuevaTarea(this.nuevoTexto.trim(), this.nuevaPrioridad),
    ]);
    this.nuevoTexto = '';
  }

  toggleCompletar(id: number): void {
    this.tareas.update(arr =>
      arr.map(t => t.id === id ? { ...t, completada: !t.completada } : t)
    );
  }

  eliminar(id: number): void {
    this.tareas.update(arr => arr.filter(t => t.id !== id));
  }

  limpiarCompletadas(): void {
    this.tareas.update(arr => arr.filter(t => !t.completada));
  }

  iniciarEdicion(tarea: Tarea): void {
    this.tareas.update(arr =>
      arr.map(t => ({
        ...t,
        editando:     t.id === tarea.id,
        textoEdicion: t.id === tarea.id ? t.texto : t.textoEdicion,
      }))
    );
  }

  guardarEdicion(tarea: Tarea): void {
    const nuevo = tarea.textoEdicion.trim();
    this.tareas.update(arr =>
      arr.map(t =>
        t.id === tarea.id
          ? { ...t, texto: nuevo || t.texto, editando: false }
          : t
      )
    );
  }

  cancelarEdicion(id: number): void {
    this.tareas.update(arr =>
      arr.map(t => t.id === id ? { ...t, editando: false } : t)
    );
  }

  onEditKeydown(event: KeyboardEvent, tarea: Tarea): void {
    if (event.key === 'Enter')  this.guardarEdicion(tarea);
    if (event.key === 'Escape') this.cancelarEdicion(tarea.id);
  }

  setFiltro(f: Filtro): void { this.filtro.set(f); }
  limpiarError(): void { if (this.errorMsg()) this.errorMsg.set(''); }
  onEnter(event: KeyboardEvent): void { if (event.key === 'Enter') this.agregar(); }
  trackById(_: number, t: Tarea): number { return t.id; }

  private nuevaTarea(texto: string, prioridad: Prioridad): Tarea {
    return {
      id: this.nextId++,
      texto,
      prioridad,
      completada: false,
      editando: false,
      textoEdicion: '',
    };
  }
}
```

### tareas.component.html

```html
<section class="tareas-page">

  <div class="page-header">
    <h1 class="page-title">📋 Tareas Editoriales</h1>
    <p class="page-sub">Gestiona el flujo de trabajo del blog</p>
  </div>

  <!-- Estadísticas -->
  <div class="stats-row">
    <div class="stat-card">
      <span class="stat-num">{{ total() }}</span>
      <span class="stat-lbl">Total</span>
    </div>
    <div class="stat-card pending">
      <span class="stat-num">{{ pendientes() }}</span>
      <span class="stat-lbl">Pendientes</span>
    </div>
    <div class="stat-card done">
      <span class="stat-num">{{ completadas() }}</span>
      <span class="stat-lbl">Completadas</span>
    </div>
    <div class="stat-card pct">
      <span class="stat-num">{{ progreso() }}%</span>
      <span class="stat-lbl">Progreso</span>
    </div>
  </div>

  <!-- Barra de progreso -->
  <div class="progress-track"
       role="progressbar"
       [attr.aria-valuenow]="progreso()"
       aria-valuemin="0" aria-valuemax="100">
    <div class="progress-fill" [style.width]="progreso() + '%'"></div>
  </div>

  <!-- Agregar tarea -->
  <div class="add-row">
    <div class="input-wrap" [class.error]="errorMsg()">
      <input
        type="text"
        class="task-input"
        placeholder="Nueva tarea editorial..."
        maxlength="80"
        [(ngModel)]="nuevoTexto"
        (keydown)="onEnter($event)"
        (input)="limpiarError()"
        aria-label="Texto de la nueva tarea"
      />
    </div>
    <select class="priority-select" [(ngModel)]="nuevaPrioridad" aria-label="Prioridad">
      <option value="alta">🔴 Alta</option>
      <option value="media">🟡 Media</option>
      <option value="baja">🟢 Baja</option>
    </select>
    <button class="btn-neo btn-add" (click)="agregar()" type="button">+ Agregar</button>
  </div>

  @if (errorMsg()) {
    <div class="error-msg" role="alert" aria-live="polite">⚠ {{ errorMsg() }}</div>
  }

  <!-- Filtros -->
  <div class="filters" role="group" aria-label="Filtrar tareas">
    <button class="chip" [class.active]="filtro() === 'todas'"       (click)="setFiltro('todas')"       type="button">Todas</button>
    <button class="chip" [class.active]="filtro() === 'pendientes'"  (click)="setFiltro('pendientes')"  type="button">Pendientes <span class="chip-count">{{ pendientes() }}</span></button>
    <button class="chip" [class.active]="filtro() === 'completadas'" (click)="setFiltro('completadas')" type="button">Completadas <span class="chip-count">{{ completadas() }}</span></button>
    @if (completadas() > 0) {
      <button class="chip chip-clear" (click)="limpiarCompletadas()" type="button">🗑 Limpiar completadas</button>
    }
  </div>

  <!-- Lista -->
  <ul class="task-list" role="list" aria-label="Lista de tareas">

    @if (tareasFiltradas().length === 0 && total() === 0) {
      <li class="empty-state">
        <span class="empty-icon">🐺</span>
        Sin tareas editoriales. Agrega la primera arriba.
      </li>
    } @else if (tareasFiltradas().length === 0) {
      <li class="empty-state">No hay tareas en "{{ filtro() }}".</li>
    }

    @for (tarea of tareasFiltradas(); track tarea.id) {
      <li class="task-item"
          [class.completada]="tarea.completada"
          [attr.data-priority]="tarea.prioridad">

        <input type="checkbox"
               class="task-check"
               [id]="'chk-' + tarea.id"
               [checked]="tarea.completada"
               (change)="toggleCompletar(tarea.id)"
               [attr.aria-label]="'Completar: ' + tarea.texto"/>

        <span class="priority-dot" [attr.data-p]="tarea.prioridad"></span>

        @if (tarea.editando) {
          <input type="text"
                 class="task-edit-input"
                 [(ngModel)]="tarea.textoEdicion"
                 (keydown)="onEditKeydown($event, tarea)"
                 (blur)="guardarEdicion(tarea)"
                 maxlength="80"
                 autofocus/>
        } @else {
          <span class="task-text"
                [class.tachado]="tarea.completada"
                title="Doble clic para editar"
                (dblclick)="iniciarEdicion(tarea)">
            {{ tarea.texto }}
          </span>
        }

        <div class="task-actions">
          @if (!tarea.editando && !tarea.completada) {
            <button class="btn-icon btn-edit" (click)="iniciarEdicion(tarea)" type="button" aria-label="Editar">✏️</button>
          }
          <button class="btn-icon btn-del" (click)="eliminar(tarea.id)" type="button" [attr.aria-label]="'Eliminar: ' + tarea.texto">✕</button>
        </div>

      </li>
    }
  </ul>

</section>
```

### tareas.component.scss

```scss
.tareas-page {
  max-width: 760px;
  margin: 0 auto;
  padding: 32px 24px;
}

.page-header { margin-bottom: 28px; }
.page-title { font-family: var(--font-heading); font-size: 1.6rem; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 0.88rem; color: #666; }

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.stat-card {
  background: var(--color-surface, #f6f2ec);
  border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b);
  border-radius: 4px;
  box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b);
  padding: 14px 12px;
  text-align: center;

  .stat-num { display: block; font-family: var(--font-heading); font-size: 1.6rem; font-weight: 700; line-height: 1; }
  .stat-lbl { display: block; font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.06em; color: #666; margin-top: 4px; }
  &.pending .stat-num { color: #b45309; }
  &.done    .stat-num { color: #16a34a; }
  &.pct     .stat-num { color: #1d4ed8; }
}

.progress-track {
  height: 8px;
  background: #e0ddd8;
  border: 1.5px solid var(--color-primary, #1b1b1b);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 28px;
}

.progress-fill {
  height: 100%;
  background: var(--color-primary, #1b1b1b);
  border-radius: 4px;
  transition: width 0.35s ease;
}

.add-row { display: flex; gap: 10px; margin-bottom: 8px; flex-wrap: wrap; }

.input-wrap {
  flex: 1;
  min-width: 200px;
  border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b);
  border-radius: 4px;
  box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b);
  background: white;
  &.error { border-color: #dc2626; box-shadow: 4px 4px 0 #dc2626; }
}

.task-input { width: 100%; padding: 10px 14px; border: none; outline: none; font-family: inherit; font-size: 0.9rem; background: transparent; }

.priority-select {
  padding: 10px 12px;
  border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b);
  border-radius: 4px;
  box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b);
  font-size: 0.85rem;
  background: white;
  cursor: pointer;
}

.btn-neo {
  padding: 10px 18px;
  border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b);
  border-radius: 4px;
  font-family: var(--font-heading);
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b);
  transition: transform 0.1s ease, box-shadow 0.1s ease;
  &:hover { transform: translate(-2px, -2px); box-shadow: 6px 6px 0 var(--color-primary, #1b1b1b); }
  &:active { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--color-primary, #1b1b1b); }
}

.btn-add { background: var(--color-primary, #1b1b1b); color: white; white-space: nowrap; }

.error-msg {
  font-size: 0.82rem; color: #dc2626; font-weight: 600; margin-bottom: 14px;
  padding: 6px 10px; background: #fef2f2; border: 1.5px solid #dc2626; border-radius: 4px;
  animation: shake 0.25s ease;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25%       { transform: translateX(-4px); }
  75%       { transform: translateX(4px); }
}

.filters { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 20px; }

.chip {
  padding: 5px 14px; border: 1.5px solid #ccc; border-radius: 20px;
  font-size: 0.78rem; font-weight: 600; background: white; cursor: pointer;
  transition: all 0.15s ease; display: flex; align-items: center; gap: 6px;
  &:hover { border-color: var(--color-primary, #1b1b1b); }
  &.active { background: var(--color-primary, #1b1b1b); border-color: var(--color-primary, #1b1b1b); color: white; }
  &-count { background: rgba(255,255,255,0.25); border-radius: 10px; padding: 0 6px; font-size: 0.72rem; }
  &.active .chip-count { background: rgba(255,255,255,0.2); }
  &-clear { margin-left: auto; border-color: #dc2626; color: #dc2626; &:hover { background: #dc2626; border-color: #dc2626; color: white; } }
}

.task-list { list-style: none; display: flex; flex-direction: column; gap: 10px; }

.task-item {
  display: flex; align-items: center; gap: 12px;
  background: white;
  border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b);
  border-radius: 4px;
  box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b);
  padding: 12px 14px;
  border-left: 6px solid #ccc;
  &[data-priority="alta"]  { border-left-color: #dc2626; }
  &[data-priority="media"] { border-left-color: #ca8a04; }
  &[data-priority="baja"]  { border-left-color: #16a34a; }
  &.completada { opacity: 0.55; background: #f9f9f9; }
}

.task-check { width: 18px; height: 18px; cursor: pointer; accent-color: var(--color-primary, #1b1b1b); flex-shrink: 0; }

.priority-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0;
  &[data-p="alta"]  { background: #dc2626; }
  &[data-p="media"] { background: #ca8a04; }
  &[data-p="baja"]  { background: #16a34a; }
}

.task-text { flex: 1; font-size: 0.9rem; cursor: default; &.tachado { text-decoration: line-through; color: #999; } }
.task-edit-input { flex: 1; padding: 4px 8px; border: 1.5px solid var(--color-primary, #1b1b1b); border-radius: 4px; font-size: 0.9rem; font-family: inherit; outline: none; &:focus { border-color: #1d4ed8; } }
.task-actions { display: flex; gap: 4px; flex-shrink: 0; }
.btn-icon { background: none; border: 1.5px solid transparent; border-radius: 4px; padding: 4px 6px; font-size: 0.8rem; cursor: pointer; transition: all 0.12s ease; &:hover { border-color: #ccc; background: #f5f5f5; } &.btn-del:hover { border-color: #dc2626; color: #dc2626; background: #fef2f2; } }

.empty-state { text-align: center; padding: 40px 20px; color: #888; font-size: 0.9rem; border: 2px dashed #ccc; border-radius: 4px; list-style: none; .empty-icon { display: block; font-size: 2rem; margin-bottom: 8px; } }

@media (max-width: 600px) {
  .stats-row { grid-template-columns: repeat(2, 1fr); }
  .add-row { flex-direction: column; }
  .priority-select, .btn-add { width: 100%; }
}
```

---

## PASO 3 — Crear componente DomDemoComponent (Lab 5)

Crea estos tres archivos en `frontend/babywolf-frontend/src/app/features/admin/dom-demo/`:

### dom-demo.component.ts

```typescript
import {
  Component,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

export interface PostCard {
  id: number;
  titulo: string;
  categoria: string;
  destacada: boolean;
}

export interface LogEntry {
  id: number;
  metodo: string;
  detalle: string;
  ts: string;
}

@Component({
  selector: 'app-dom-demo',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './dom-demo.component.html',
  styleUrls: ['./dom-demo.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DomDemoComponent {

  private nextId = 4;
  private logId  = 1;

  tituloBlog  = signal('BabyWolf Blog');
  temaOscuro  = signal(false);
  colorAcento = signal('#e94560');
  tamFuente   = signal(16);
  bannerVisible = signal(false);
  bannerTexto   = signal('');
  cardEditandoId = signal<number | null>(null);
  cardEditTexto  = '';
  nuevoTitulo    = '';
  nuevaCategoria = 'gaming';

  cards = signal<PostCard[]>([
    { id: 1, titulo: 'Análisis: RTX 5090',     categoria: 'hardware',   destacada: false },
    { id: 2, titulo: 'Top 10 juegos indie',     categoria: 'gaming',     destacada: false },
    { id: 3, titulo: 'Angular Signals vs RxJS', categoria: 'tutoriales', destacada: false },
  ]);

  log = signal<LogEntry[]>([]);

  totalCards = computed(() => this.cards().length);

  actualizarTitulo(valor: string): void {
    this.tituloBlog.set(valor || 'BabyWolf Blog');
    this.registrar('textContent', `tituloBlog = "${this.tituloBlog()}"`);
  }

  toggleTema(): void {
    this.temaOscuro.update(v => !v);
    this.registrar('classList.toggle("dark")', `tema → ${this.temaOscuro() ? 'oscuro' : 'claro'}`);
  }

  cambiarAcento(color: string): void {
    this.colorAcento.set(color);
    this.registrar('style.setProperty("--accent")', `color → ${color}`);
  }

  cambiarFuente(tam: number): void {
    this.tamFuente.set(tam);
    this.registrar('style.fontSize', `${tam}px`);
  }

  agregarCard(): void {
    if (!this.nuevoTitulo.trim()) return;
    const card: PostCard = { id: this.nextId++, titulo: this.nuevoTitulo.trim(), categoria: this.nuevaCategoria, destacada: false };
    this.cards.update(arr => [...arr, card]);
    this.registrar('createElement + appendChild', `card "${card.titulo}" → #postsGrid`);
    this.nuevoTitulo = '';
  }

  onEnter(e: KeyboardEvent): void { if (e.key === 'Enter') this.agregarCard(); }

  eliminarCard(id: number): void {
    const card = this.cards().find(c => c.id === id);
    this.cards.update(arr => arr.filter(c => c.id !== id));
    this.registrar('element.remove()', `card "${card?.titulo}" eliminada`);
  }

  duplicarCard(card: PostCard): void {
    const clon: PostCard = { ...card, id: this.nextId++, titulo: card.titulo + ' (copia)' };
    this.cards.update(arr => {
      const idx = arr.findIndex(c => c.id === card.id);
      const next = [...arr];
      next.splice(idx + 1, 0, clon);
      return next;
    });
    this.registrar('cloneNode(true) + insertBefore', `"${card.titulo}" duplicada`);
  }

  toggleDestacada(id: number): void {
    this.cards.update(arr => arr.map(c => c.id === id ? { ...c, destacada: !c.destacada } : c));
    const c = this.cards().find(c => c.id === id);
    this.registrar('classList.toggle("destacada")', `"${c?.titulo}" → destacada: ${c?.destacada}`);
  }

  iniciarEdicion(card: PostCard): void {
    this.cardEditandoId.set(card.id);
    this.cardEditTexto = card.titulo;
    this.registrar('replaceWith(input)', `edición in-place de "${card.titulo}"`);
  }

  guardarEdicion(id: number): void {
    const nuevo = this.cardEditTexto.trim();
    if (nuevo) {
      this.cards.update(arr => arr.map(c => c.id === id ? { ...c, titulo: nuevo } : c));
      this.registrar('input.replaceWith(span)', `título guardado: "${nuevo}"`);
    }
    this.cardEditandoId.set(null);
  }

  cancelarEdicion(): void { this.cardEditandoId.set(null); }

  onEditKeydown(e: KeyboardEvent, id: number): void {
    if (e.key === 'Enter')  this.guardarEdicion(id);
    if (e.key === 'Escape') this.cancelarEdicion();
  }

  mostrarBanner(): void {
    this.bannerTexto.set('🎉 Banner creado dinámicamente con prepend() / signal');
    this.bannerVisible.set(true);
    this.registrar('prepend()', 'banner insertado al inicio del contenido');
    setTimeout(() => this.cerrarBanner(), 4000);
  }

  cerrarBanner(): void {
    this.bannerVisible.set(false);
    this.registrar('element.remove()', 'banner eliminado');
  }

  limpiarLog(): void { this.log.set([]); }

  trackById = (_: number, item: { id: number }) => item.id;

  private registrar(metodo: string, detalle: string): void {
    const entry: LogEntry = { id: this.logId++, metodo, detalle, ts: new Date().toLocaleTimeString('es-MX') };
    this.log.update(arr => [entry, ...arr].slice(0, 20));
  }
}
```

### dom-demo.component.html

```html
<div class="dom-page"
     [class.dark]="temaOscuro()"
     [style.font-size.px]="tamFuente()"
     [style.--accent]="colorAcento()">

  @if (bannerVisible()) {
    <div class="banner-top" role="alert">
      <span>{{ bannerTexto() }}</span>
      <button class="banner-close" (click)="cerrarBanner()">✕</button>
    </div>
  }

  <header class="demo-header">
    <h1 class="blog-title">{{ tituloBlog() }}</h1>
    <span class="card-count">{{ totalCards() }} cards</span>
  </header>

  <div class="controls-grid">

    <div class="control-group">
      <label class="ctrl-label">
        <span class="ctrl-method">textContent</span>
        Título del blog (evento input)
      </label>
      <input type="text" class="ctrl-input" placeholder="Escribe el título..."
             [ngModel]="tituloBlog()"
             (input)="actualizarTitulo($any($event.target).value)"
             maxlength="40"/>
    </div>

    <div class="control-group">
      <label class="ctrl-label">
        <span class="ctrl-method">classList.toggle("dark")</span>
        Tema del sitio (evento click)
      </label>
      <button class="btn-neo" (click)="toggleTema()" [class.active]="temaOscuro()" type="button">
        {{ temaOscuro() ? '☀️ Modo claro' : '🌙 Modo oscuro' }}
      </button>
    </div>

    <div class="control-group">
      <label class="ctrl-label">
        <span class="ctrl-method">style.setProperty("--accent")</span>
        Color de acento (evento change)
      </label>
      <div class="color-row">
        <input type="color" [value]="colorAcento()"
               (change)="cambiarAcento($any($event.target).value)"
               class="color-picker"/>
        <code class="color-val">{{ colorAcento() }}</code>
      </div>
    </div>

    <div class="control-group">
      <label class="ctrl-label">
        <span class="ctrl-method">style.fontSize</span>
        Tamaño de fuente: {{ tamFuente() }}px
      </label>
      <input type="range" [value]="tamFuente()"
             (input)="cambiarFuente(+$any($event.target).value)"
             min="12" max="22" step="1" class="font-slider"/>
    </div>

  </div>

  <div class="actions-row">
    <input type="text" class="ctrl-input" placeholder="Título de nueva card..."
           [(ngModel)]="nuevoTitulo" (keydown)="onEnter($event)" maxlength="60"/>
    <select class="ctrl-select" [(ngModel)]="nuevaCategoria">
      <option value="gaming">Gaming</option>
      <option value="hardware">Hardware</option>
      <option value="reviews">Reviews</option>
      <option value="tutoriales">Tutoriales</option>
    </select>
    <button class="btn-neo btn-primary" (click)="agregarCard()" type="button">+ createElement</button>
    <button class="btn-neo" (click)="mostrarBanner()" type="button">prepend() banner</button>
  </div>

  <div class="posts-grid" id="postsGrid">
    @if (cards().length === 0) {
      <p class="empty-grid">Sin cards. Usa "+ createElement" para agregar.</p>
    }
    @for (card of cards(); track card.id) {
      <article class="post-card" [class.destacada]="card.destacada" [attr.data-id]="card.id">

        <span class="card-cat">{{ card.categoria }}</span>

        @if (cardEditandoId() === card.id) {
          <input type="text" class="card-edit-input"
                 [(ngModel)]="cardEditTexto"
                 (keydown)="onEditKeydown($event, card.id)"
                 (blur)="guardarEdicion(card.id)"
                 maxlength="60" autofocus/>
        } @else {
          <h3 class="card-title" title="Doble clic para editar" (dblclick)="iniciarEdicion(card)">
            {{ card.titulo }}
          </h3>
        }

        @if (card.destacada) { <span class="badge-dest">⭐ Destacada</span> }

        <div class="card-actions">
          <button class="btn-icon" (click)="toggleDestacada(card.id)" type="button" [title]="card.destacada ? 'Quitar' : 'Destacar'">{{ card.destacada ? '★' : '☆' }}</button>
          <button class="btn-icon" (click)="iniciarEdicion(card)"    type="button" title="Editar (replaceWith)">✏️</button>
          <button class="btn-icon" (click)="duplicarCard(card)"      type="button" title="Duplicar (cloneNode)">⧉</button>
          <button class="btn-icon btn-del" (click)="eliminarCard(card.id)" type="button" [attr.aria-label]="'Eliminar: ' + card.titulo">✕</button>
        </div>

      </article>
    }
  </div>

  <div class="log-panel">
    <div class="log-header">
      <span class="log-title">📋 Log de métodos DOM ejecutados</span>
      <button class="btn-sm" (click)="limpiarLog()" type="button">Limpiar</button>
    </div>
    <div class="log-body">
      @if (log().length === 0) {
        <p class="log-empty">Interactúa con los controles para ver los métodos registrados.</p>
      }
      @for (entry of log(); track entry.id) {
        <div class="log-line">
          <code class="log-method">{{ entry.metodo }}</code>
          <span class="log-detail">{{ entry.detalle }}</span>
          <span class="log-ts">{{ entry.ts }}</span>
        </div>
      }
    </div>
  </div>

</div>
```

### dom-demo.component.scss

```scss
.dom-page {
  max-width: 960px; margin: 0 auto; padding: 28px 24px;
  transition: background 0.3s ease, color 0.3s ease;

  &.dark {
    background: #1a1a2e; color: rgba(255,255,255,0.9); border-radius: 8px;
    .post-card   { background: #16213e; border-color: rgba(255,255,255,0.1); }
    .ctrl-input, .ctrl-select { background: #0f172a; border-color: rgba(255,255,255,0.15); color: #fff; }
    .btn-neo { background: #16213e; color: #fff; border-color: rgba(255,255,255,0.2); }
    .btn-neo.btn-primary { background: var(--accent, #e94560); border-color: var(--accent, #e94560); }
    .log-panel { background: #0f172a; border-color: rgba(255,255,255,0.1); }
    .log-method { color: var(--accent, #e94560); }
    .card-cat { background: var(--accent, #e94560); color: #fff; }
    .ctrl-label { color: rgba(255,255,255,0.6); }
    .ctrl-method { background: rgba(233,69,96,0.2); color: var(--accent, #e94560); }
    .control-group { background: #16213e; border-color: rgba(255,255,255,0.1); }
  }
}

.banner-top {
  background: var(--color-primary, #1b1b1b); color: white;
  padding: 12px 18px; border-radius: 4px; margin-bottom: 20px;
  display: flex; align-items: center; justify-content: space-between;
  font-size: 0.88rem; font-weight: 600;
  border-left: 5px solid var(--accent, #e94560);
  animation: slideDown 0.3s ease;
}
@keyframes slideDown { from { opacity: 0; transform: translateY(-12px); } to { opacity: 1; transform: translateY(0); } }
.banner-close { background: none; border: none; color: rgba(255,255,255,0.6); cursor: pointer; font-size: 1rem; padding: 2px 6px; &:hover { color: #fff; } }

.demo-header { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 24px; border-bottom: var(--border-width, 2px) solid var(--color-primary, #1b1b1b); padding-bottom: 12px; }
.blog-title { font-family: var(--font-heading); font-size: 1.5rem; font-weight: 700; color: var(--accent, #e94560); transition: color 0.2s; }
.card-count { font-size: 0.78rem; color: #888; font-weight: 600; }

.controls-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 20px; }
.control-group { background: white; border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b); border-radius: 4px; box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b); padding: 14px; }
.ctrl-label { display: block; font-size: 0.78rem; color: #555; margin-bottom: 8px; line-height: 1.4; }
.ctrl-method { display: inline-block; font-family: var(--font-heading); font-size: 0.72rem; background: #f0f0f0; padding: 1px 7px; border-radius: 3px; margin-right: 4px; }
.ctrl-input { width: 100%; padding: 8px 12px; border: 1.5px solid #ccc; border-radius: 4px; font-family: inherit; font-size: 0.88rem; outline: none; transition: border-color 0.15s; &:focus { border-color: var(--accent, #e94560); } }
.ctrl-select { width: 100%; padding: 8px 12px; border: 1.5px solid #ccc; border-radius: 4px; font-size: 0.88rem; background: white; cursor: pointer; }
.color-row { display: flex; align-items: center; gap: 10px; }
.color-picker { width: 44px; height: 36px; border: 1.5px solid #ccc; border-radius: 4px; cursor: pointer; padding: 2px; }
.color-val { font-family: var(--font-heading); font-size: 0.82rem; color: #555; }
.font-slider { width: 100%; accent-color: var(--accent, #e94560); cursor: pointer; }

.actions-row { display: flex; gap: 10px; margin-bottom: 20px; flex-wrap: wrap; .ctrl-input { max-width: 260px; } }

.btn-neo {
  padding: 9px 16px; border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b); border-radius: 4px;
  font-family: var(--font-heading); font-size: 0.82rem; font-weight: 700; background: white;
  box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b); cursor: pointer; transition: transform 0.1s, box-shadow 0.1s; white-space: nowrap;
  &:hover { transform: translate(-1px,-1px); box-shadow: 5px 5px 0 var(--color-primary, #1b1b1b); }
  &.btn-primary { background: var(--color-primary, #1b1b1b); color: white; }
  &.active { background: var(--accent, #e94560); border-color: var(--accent, #e94560); color: white; }
}

.posts-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 16px; margin-bottom: 28px; }
.empty-grid { grid-column: 1/-1; text-align: center; padding: 30px; color: #888; border: 2px dashed #ccc; border-radius: 4px; }

.post-card {
  background: white; border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b);
  border-radius: 4px; box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b); padding: 14px;
  &.destacada { border-color: var(--accent, #e94560); box-shadow: 4px 4px 0 var(--accent, #e94560); }
}
.card-cat { display: inline-block; background: var(--color-primary, #1b1b1b); color: white; font-size: 0.65rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.08em; padding: 2px 8px; border-radius: 3px; margin-bottom: 8px; }
.card-title { font-family: var(--font-heading); font-size: 0.9rem; font-weight: 700; margin-bottom: 8px; line-height: 1.35; min-height: 36px; cursor: default; }
.card-edit-input { width: 100%; padding: 4px 8px; border: 1.5px solid var(--accent, #e94560); border-radius: 3px; font-size: 0.88rem; font-family: inherit; margin-bottom: 8px; outline: none; }
.badge-dest { font-size: 0.68rem; font-weight: 700; color: var(--accent, #e94560); display: block; margin-bottom: 6px; }
.card-actions { display: flex; gap: 4px; border-top: 1px solid #eee; padding-top: 8px; margin-top: 4px; }
.btn-icon { background: none; border: 1.5px solid transparent; border-radius: 3px; padding: 3px 6px; font-size: 0.8rem; cursor: pointer; transition: all 0.12s; &:hover { border-color: #ccc; background: #f5f5f5; } &.btn-del:hover { border-color: #dc2626; color: #dc2626; background: #fef2f2; } }
.btn-sm { padding: 4px 12px; font-size: 0.75rem; font-weight: 700; border: 1.5px solid #ccc; border-radius: 3px; background: white; cursor: pointer; font-family: var(--font-heading); &:hover { border-color: var(--color-primary, #1b1b1b); } }

.log-panel { background: white; border: var(--border-width, 2px) solid var(--color-primary, #1b1b1b); border-radius: 4px; box-shadow: var(--shadow-neo, 4px 4px 0 #1b1b1b); overflow: hidden; }
.log-header { display: flex; align-items: center; justify-content: space-between; padding: 10px 14px; border-bottom: 1.5px solid #eee; background: #f8f8f8; }
.log-title { font-family: var(--font-heading); font-size: 0.82rem; font-weight: 700; }
.log-body { max-height: 200px; overflow-y: auto; }
.log-empty { padding: 16px; color: #888; font-size: 0.82rem; text-align: center; }
.log-line { display: flex; align-items: center; gap: 10px; padding: 7px 14px; border-bottom: 1px solid #f0f0f0; font-size: 0.78rem; animation: fadeIn 0.2s ease; &:last-child { border-bottom: none; } }
@keyframes fadeIn { from { opacity: 0; transform: translateX(-6px); } to { opacity: 1; transform: translateX(0); } }
.log-method { font-family: var(--font-heading); font-size: 0.72rem; background: #f0f0f0; padding: 2px 7px; border-radius: 3px; white-space: nowrap; flex-shrink: 0; }
.log-detail { flex: 1; color: #333; }
.log-ts { color: #aaa; flex-shrink: 0; font-size: 0.7rem; }

@media (max-width: 640px) {
  .controls-grid { grid-template-columns: 1fr; }
  .actions-row { flex-direction: column; .ctrl-input { max-width: 100%; } }
  .posts-grid { grid-template-columns: 1fr; }
}
```

---

## PASO 4 — Modificar app.routes.ts

Lee el archivo actual en `frontend/babywolf-frontend/src/app/app.routes.ts`.

Dentro del array `children` de la ruta `path: 'admin'`, agrega estas dos rutas **al final del array, antes del cierre `]`**:

```typescript
{
  path: 'tareas',
  loadComponent: () =>
    import('./features/admin/tareas/tareas.component')
      .then(m => m.TareasComponent),
},
{
  path: 'dom-demo',
  loadComponent: () =>
    import('./features/admin/dom-demo/dom-demo.component')
      .then(m => m.DomDemoComponent),
},
```

---

## PASO 5 — Modificar admin-layout.ts

Lee el archivo actual en `frontend/babywolf-frontend/src/app/features/admin/layout/admin-layout.ts`.

En el array `adminMenu`, agrega al final **antes del cierre `]`**:

```typescript
{ label: 'Tareas',   icon: '📋', routerLink: '/admin/tareas' },
{ label: 'DOM Demo', icon: '🧪', routerLink: '/admin/dom-demo' },
```

---

## PASO 6 — Commit final

Una vez creados y modificados todos los archivos, ejecuta:

```bash
git add docs/
git add frontend/babywolf-frontend/src/app/features/admin/tareas/
git add frontend/babywolf-frontend/src/app/features/admin/dom-demo/
git add frontend/babywolf-frontend/src/app/app.routes.ts
git add frontend/babywolf-frontend/src/app/features/admin/layout/admin-layout.ts
git commit -m "feat(labs): agregar practicas 1-6 con componentes Angular (tareas + dom-demo)"
```

---

## NOTAS IMPORTANTES

- No ejecutes `pnpm install` ni `npm install`, las dependencias ya existen
- No modifiques ningún archivo que no esté listado arriba
- Si algún archivo de destino ya existe con contenido diferente, reemplázalo completamente con el contenido indicado
- Los componentes usan `ChangeDetectionStrategy.OnPush` y Signals — no uses `ngZone` ni `markForCheck` salvo que veas que ya se usa en los archivos existentes del admin
- El `FormsModule` es necesario para `[(ngModel)]` en los dos componentes — ya está en las importaciones del `@Component`
