# Práctica de Laboratorio 8
## Evento basado en periodo (banner/fondo estacional)
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Objetivo

Seleccionar un evento basado en periodo (fecha calendarizada) que cambie la interfaz del sitio — en este caso, un acento visual en el fondo/banner — en función de la temporada del año.

## 2. Planteamiento

Se creó `core/services/seasonal-theme.service.ts`, un servicio que se ejecuta una sola vez al iniciar la aplicación (`app.ts`) y evalúa el mes actual (`Date.getMonth()`) contra rangos fijos:

| Temporada | Meses | Clase CSS aplicada | Color de acento |
|---|---|---|---|
| Halloween | Octubre | `season-halloween` | Naranja `#ff7518` |
| Navidad | Diciembre | `season-navidad` | Verde `#2ecc71` |
| Verano | Junio–Agosto | `season-verano` | Amarillo `#f4c542` |
| Resto del año | — | `season-default` | Sin acento |

## 3. Implementación

```typescript
apply(date: Date = new Date()): string {
  const month = date.getMonth() + 1;
  const theme = THEMES.find(t => t.months.includes(month));
  const className = theme?.className ?? DEFAULT_THEME;

  document.body.classList.remove(...THEMES.map(t => t.className), DEFAULT_THEME);
  document.body.classList.add(className);
  return className;
}
```

El cambio visual se resuelve en `styles.scss` mediante clases sobre `body`:

```scss
body.season-halloween { border-top-color: #ff7518; }
body.season-navidad { border-top-color: #2ecc71; }
body.season-verano { border-top-color: #f4c542; }
```

Se aplica automáticamente al arrancar el sitio, sin intervención del usuario — el evento "disparador" es el propio ciclo de vida de la aplicación combinado con la fecha del sistema (`new Date()`).

## 4. Evidencias

> Anexar capturas de pantalla del sitio en al menos dos "temporadas" distintas. Para verificar sin esperar al mes real, se puede forzar temporalmente desde la consola del navegador:
> ```js
> // Sólo para pruebas/evidencia — no forma parte del código de producción
> document.body.classList.add('season-navidad')
> ```
> 1. Captura con el acento por defecto (fuera de temporada).
> 2. Captura con `season-navidad` o `season-halloween` forzado, mostrando el borde superior con el color correspondiente.

---
*Práctica de Laboratorio 8 — Desarrollo y Gestión de Software*
