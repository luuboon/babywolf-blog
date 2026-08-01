# Práctica de Laboratorio 9
## Documentación del evento JavaScript de actualización de temporada (versionado en Git)
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Objetivo

Documentar, dentro del versionado de Git, el evento JavaScript pensado para la actualización de temporada del sitio (el mismo evento calendarizado implementado en la [Práctica 8](practica8_evento_calendarizado.md)).

## 2. Qué cambia y por qué

El evento es `SeasonalThemeService.apply()` (`src/app/core/services/seasonal-theme.service.ts`), invocado una vez al arrancar la app desde `App` (`app.ts`). Compara la fecha del sistema contra rangos de meses fijos y aplica una clase CSS al `<body>` que cambia el acento visual del sitio (Halloween, Navidad, Verano, o el estado por defecto).

Es un evento "basado en periodo": no depende de que el usuario haga click ni de datos remotos, sino del reloj del sistema en el momento en que la SPA se inicializa.

## 3. Archivos que forman parte del cambio

| Archivo | Rol |
|---|---|
| `src/app/core/services/seasonal-theme.service.ts` | Lógica del evento (nuevo) |
| `src/app.ts` | Dispara el evento al arrancar la app |
| `src/styles.scss` | Reglas visuales por clase de temporada |

## 4. Mensaje de commit sugerido

Siguiendo Conventional Commits (convención ya usada en el repositorio, ver `git log`):

```
feat(theme): agregar evento calendarizado de temporada para banner del sitio

Introduce SeasonalThemeService, que aplica una clase CSS al body
según el mes actual (Halloween, Navidad, Verano, default) al
iniciar la aplicación, cambiando el acento visual del sitio sin
intervención del usuario.
```

## 5. Evidencia de versionado

> Anexar aquí:
> 1. `git log --oneline -1` mostrando el commit del cambio.
> 2. `git show <hash> --stat` mostrando los archivos modificados.
> 3. Captura del diff (`git show <hash>`) del servicio `seasonal-theme.service.ts`.

---
*Práctica de Laboratorio 9 — Desarrollo y Gestión de Software*
