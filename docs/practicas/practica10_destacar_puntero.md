# Práctica de Laboratorio 10
## Destacar contenido mediante el puntero
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Objetivo

Aplicar un evento tipo "destacar contenido mediante el puntero" sobre un elemento del sitio, con grabación de pantalla como evidencia.

## 2. Implementación

Se instrumentó `PostCard` (`features/posts/presentation/components/post-card/`), la tarjeta que representa cada publicación del blog, con eventos DOM explícitos `mouseenter` / `mouseleave`:

```typescript
export class PostCard {
    @Input({ required: true }) post!: Post;

    highlighted = false;

    onPointerEnter(): void {
        this.highlighted = true;
    }

    onPointerLeave(): void {
        this.highlighted = false;
    }
}
```

```html
<article
    class="post-card"
    [class.highlighted]="highlighted"
    [routerLink]="['/posts', post.slug]"
    (mouseenter)="onPointerEnter()"
    (mouseleave)="onPointerLeave()"
>
```

```scss
&.highlighted {
    outline: 3px solid #e94560;
    outline-offset: 2px;
}
```

## 3. Comportamiento

Al pasar el puntero sobre una tarjeta de post en `/posts` o en `/` (home), el evento `mouseenter` agrega la clase `highlighted`, que dibuja un contorno de color sobre la tarjeta; al retirar el puntero (`mouseleave`) el contorno desaparece. Se suma al efecto de elevación (`transform` + `box-shadow`) que la tarjeta ya tenía vía `:hover` en CSS.

## 4. Evidencias

> Anexar aquí una grabación de pantalla breve (10-15s) mostrando:
> 1. El listado de posts en estado normal.
> 2. El puntero recorriendo varias tarjetas, mostrando cómo cada una se destaca al pasar el mouse y vuelve a su estado normal al salir.

---
*Práctica de Laboratorio 10 — Desarrollo y Gestión de Software*
