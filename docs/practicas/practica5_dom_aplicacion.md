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
