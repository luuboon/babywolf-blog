# Práctica de Laboratorio 13
## Administración de cuentas de usuario en el sistema operativo (macOS Tahoe)
**Materia:** Desarrollo y Gestión de Software
**Fecha límite:** Jueves 23, 23:00 hrs
**Autor:** Abraham

---

## 1. Objetivo

Comprender la administración de cuentas de usuario en macOS Tahoe: creación, modificación, habilitación, deshabilitación y verificación de estados, y analizar sus implicaciones de seguridad y administración.

## 2. Marco conceptual

- **Habilitada:** el usuario puede autenticarse normalmente.
- **Deshabilitada:** la cuenta y su carpeta `/Users/<usuario>` permanecen en el sistema, pero no puede iniciar sesión.
- **Eliminada:** se pierde la cuenta y, salvo respaldo, su configuración y archivos — mala práctica para auditoría cuando el usuario podría volver.

## 3. Creación de las tres cuentas

Comandos (Terminal, requieren `sudo`):

```bash
sudo sysadminctl -addUser alumno1 -fullName "Alumno Uno" -password 'ClaveSegura#2026' -admin
sudo sysadminctl -addUser alumno2 -fullName "Alumno Dos" -password 'ClaveSegura#2026' -admin
sudo sysadminctl -addUser invitado -fullName "Usuario Invitado" -password 'ClaveSegura#2026'
```

`alumno1` y `alumno2` se crean con el mismo nivel (estándar o admin, a elección); `invitado` se crea sin `-admin` (cuenta estándar, rol de invitado dentro del alcance de esta práctica — macOS Tahoe eliminó la "Guest Account" clásica de versiones previas).

> Anexar aquí: captura de los tres comandos ejecutados y de **System Settings → Users & Groups** mostrando las tres cuentas creadas.

## 4. Identificación de cada usuario

Para cada cuenta, obtener los datos con:

```bash
dscl . -read /Users/alumno1 RealName
dscl . -read /Users/alumno1 UniqueID
id alumno1                     # grupos a los que pertenece
sudo dscl . -read /Users/alumno1 AuthenticationAuthority   # estado / política de password
last alumno1 | head -1         # último inicio de sesión
stat -f "%Sc" /Users/alumno1   # fecha de creación aproximada (ctime de la carpeta home)
```

| Campo | alumno1 | alumno2 | invitado |
|---|---|---|---|
| Nombre completo | Alumno Uno | Alumno Dos | Usuario Invitado |
| Estado de la cuenta | Habilitada | Habilitada | Habilitada |
| Fecha de creación | _(completar)_ | _(completar)_ | _(completar)_ |
| Último inicio de sesión | _(completar, `last`)_ | _(completar)_ | _(completar)_ |
| Grupos | staff, admin* | staff, admin* | staff |
| ¿Requiere cambio de contraseña? | No (a definir con `pwpolicy`) | No | No |

\* si se crearon con `-admin`.

> Anexar captura de la salida de `id alumno1` / `id alumno2` / `id invitado`.

## 5. Tabla resumen

| Usuario | Habilitado | Grupo | Cambio de contraseña requerido |
|---|---|---|---|
| alumno1 | Sí | staff (+admin si aplica) | No |
| alumno2 | Sí | staff (+admin si aplica) | No |
| invitado | Sí | staff | No |

> Anexar evidencia (captura de Users & Groups o de la terminal).

## 6. Deshabilitar a `alumno2`

macOS no tiene un flag único "activo/inactivo" como Windows; la forma soportada por línea de comandos es marcar la autoridad de autenticación como deshabilitada:

```bash
sudo dscl . -append /Users/alumno2 AuthenticationAuthority ";DisabledUser;"
```

(Alternativa por GUI: **System Settings → Users & Groups → alumno2 → Login Options**, desactivar el acceso.)

Verificación:

```bash
sudo dscl . -read /Users/alumno2 AuthenticationAuthority   # debe listar ;DisabledUser;
```

- **¿Aparece como deshabilitado?** Sí, `AuthenticationAuthority` incluye `;DisabledUser;` y en Users & Groups aparece con el acceso restringido.
- **¿Puede iniciar sesión?** No — macOS rechaza el login (pantalla de bloqueo) aunque la contraseña sea correcta.

> Anexar captura del comando `dscl` y de la pantalla de login rechazando a `alumno2`.

## 7. Volver a habilitar a `alumno2`

```bash
sudo dscl . -delete /Users/alumno2 AuthenticationAuthority
sudo dscl . -passwd /Users/alumno2 'ClaveSegura#2026'   # regenera AuthenticationAuthority válida
```

- **¿Puede iniciar sesión?** Sí, de inmediato, con la misma contraseña y sin perder ningún dato de su carpeta home.
- **¿Qué propiedad cambió?** `AuthenticationAuthority` (se removió el flag `;DisabledUser;`); ninguna otra propiedad (UID, grupos, home, archivos) se modificó.

> Anexar captura confirmando el login exitoso de `alumno2`.

## 8. Pruebas adicionales de deshabilitación por comando

> Repetir el procedimiento de la sección 6 sobre `invitado` (o de nuevo sobre `alumno2`) y anexar evidencia de al menos una segunda prueba, incluyendo el intento de inicio de sesión fallido.

## 9. Deshabilitar vs. eliminar

| Acción | ¿Se conserva la información? | ¿Puede iniciar sesión? | ¿Puede recuperarse fácilmente? |
|---|---|---|---|
| Deshabilitar usuario | Sí, cuenta y home intactos | No | Sí, revirtiendo `AuthenticationAuthority` |
| Eliminar usuario | No (salvo respaldo explícito de `/Users/<user>`) | No (la cuenta ya no existe) | No — hay que recrear la cuenta desde cero y restaurar backups si existen |

## 10. Casos prácticos

| Empleado | Situación | ¿Eliminar? | ¿Deshabilitar? | Justificación |
|---|---|---|---|---|
| Ana | Vacaciones por 6 meses | No | **Sí** | Va a volver; eliminar perdería su configuración y forzaría recrear la cuenta. Deshabilitar bloquea el acceso sin perder nada. |
| Luis | Renunció definitivamente | **Sí** (tras respaldo) | No basta solo con deshabilitar | No va a volver; por política de datos y licencias conviene depurar cuentas que ya no representan a nadie en la organización — pero primero se hace un backup y, mientras dure el proceso de baja administrativa, puede pasar por un estado "deshabilitada" antes de la eliminación definitiva. |
| Pedro | Suspendido durante una investigación | No | **Sí** | Es clave conservar su cuenta, historial y archivos intactos como posible evidencia; sólo se bloquea el acceso mientras dura la investigación. |
| María | Cambio temporal de departamento | No | No (o solo ajustar grupos/permisos) | No perdió su relación laboral ni dejó de necesitar acceso — lo que corresponde es reasignar sus grupos/permisos al nuevo departamento, no deshabilitar ni eliminar la cuenta. |

## 11. Preguntas finales

**¿Cuál es la diferencia entre deshabilitar y eliminar un usuario?**
Deshabilitar bloquea la autenticación pero conserva la cuenta, su UID, sus grupos y su carpeta home tal cual estaban, de forma reversible. Eliminar borra la cuenta del sistema (y opcionalmente su carpeta home); es una operación destructiva y, sin respaldo, irreversible.

**¿Qué ventajas ofrece deshabilitar una cuenta?**
Permite reaccionar rápido ante bajas temporales, sospechas de seguridad o investigaciones sin perder historial ni configuración; conserva la trazabilidad para auditoría (quién es el dueño de cada archivo/acción) y evita el trabajo de recrear la cuenta si la persona regresa.

**¿Qué riesgos existen si nunca se deshabilitan cuentas antiguas?**
Cuentas de ex-empleados o de personal inactivo quedan como vector de ataque (credenciales filtradas, cuentas huérfanas con privilegios elevados), inflan la superficie de auditoría, y dificultan saber quién debería tener acceso real a qué recursos en un momento dado.

**¿Qué sucede con los archivos personales de un usuario deshabilitado?**
Permanecen intactos en su carpeta `/Users/<usuario>`, con los mismos permisos y propietario; sólo se le impide iniciar sesión para acceder a ellos interactivamente (un administrador con privilegios sigue pudiendo leerlos si es necesario).

**¿En qué escenarios empresariales se recomienda esta práctica?**
Licencias o incapacidades prolongadas, suspensiones o investigaciones internas, transiciones de rol donde el acceso se re-evalúa antes de reactivarse, periodos de prueba antes de decidir una baja definitiva, y como paso intermedio obligatorio antes de cualquier eliminación de cuenta (para dar margen de reversión ante errores administrativos).

---
*Práctica de Laboratorio 13 — Desarrollo y Gestión de Software*
