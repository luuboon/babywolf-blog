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
