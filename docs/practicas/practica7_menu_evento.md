# Práctica de Laboratorio 7
## Menú basado en evento de puntero
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Objetivo

Plantear, documentar e implementar un menú que se despliegue mediante un evento generado por el puntero, y que permita al usuario interactuar con él.

## 2. Planteamiento

El header del sitio público (`features/user/layout/user-layout.ts` / `.html`) ya tenía la navegación principal resuelta en el sidebar lateral. Se agregó un menú **"Categorías"** en el header superior que se despliega al hacer click (evento de puntero) sobre el botón, mostrando accesos directos a Gaming, Tech, Opinión y Retro, y que se cierra tanto al elegir una opción como al hacer click fuera de él.

## 3. Implementación

Se implementó como un menú propio (sin librería) para mantener el control total sobre el evento de apertura/cierre y demostrar el manejo directo de eventos del DOM:

```typescript
categoriesMenuOpen = false;

toggleCategoriesMenu(event: MouseEvent): void {
  event.stopPropagation();
  this.categoriesMenuOpen = !this.categoriesMenuOpen;
}

@HostListener('document:click', ['$event'])
onDocumentClick(event: MouseEvent): void {
  if (this.categoriesMenuOpen && !this.elementRef.nativeElement.contains(event.target as Node)) {
    this.categoriesMenuOpen = false;
  }
}
```

```html
<button class="header-btn categories-trigger" (click)="toggleCategoriesMenu($event)" aria-haspopup="true" [attr.aria-expanded]="categoriesMenuOpen">
  Categorías ▾
</button>
@if (categoriesMenuOpen) {
  <div class="categories-dropdown">
    @for (cat of categoriesMenu; track cat.label) {
      <a [routerLink]="cat.routerLink" class="categories-dropdown-item" (click)="closeCategoriesMenu()">{{ cat.label }}</a>
    }
  </div>
}
```

## 4. Comportamiento

| Evento de puntero | Resultado |
|---|---|
| `click` en el botón "Categorías" | Se despliega el menú (overlay flotante) |
| `click` en un ítem del menú | Navega a `/category/:name` y cierra el menú |
| `click` fuera del menú | El menú se cierra (comportamiento nativo de Material Overlay) |

## 5. Evidencias

> Anexar aquí capturas de pantalla mostrando:
> 1. El botón "Categorías" en el header, cerrado.
> 2. El menú desplegado tras el click, con las 5 categorías visibles.
> 3. La navegación resultante a `/category/gaming` (o la categoría elegida).

---
*Práctica de Laboratorio 7 — Desarrollo y Gestión de Software*
