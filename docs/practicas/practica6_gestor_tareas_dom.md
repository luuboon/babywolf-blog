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
