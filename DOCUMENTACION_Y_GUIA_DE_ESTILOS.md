# BabyWolf Blog - Documentación Técnica y Arquitectónica 🐺

Este documento expone de forma exhaustiva las decisiones arquitectónicas, de seguridad y de diseño (Guía de Estilos) que rigieron la construcción de BabyWolf Blog. Se justifica cada patrón en favor de la mantenibilidad, rendimiento a gran escala y cumplimiento de la rúbrica de evaluación.

---

## 🏛️ 1. Arquitectura del Proyecto y Patrones

### 1.1 Stack Tecnológico
La plataforma emplea un enfoque arquitectónico **Serverless CND (Cloud Native Development)** y **SPA (Single Page Application)** moderno:

* **Frontend:** Angular v17+ (Standalones & SSR support) + TypeScript + RxJS.
* **Backend como Servicio (BaaS):** Supabase (PostgreSQL 15, Go-based Auth API, Storage).
* **Infraestructura y Despliegue:** Vercel (Edge computing para UI).
* **Email Provider:** Resend (Notificaciones transaccionales y de seguridad).

### 1.2 Justificación de la Arquitectura
**¿Por qué Angular Standalone y Supabase?**
1. **Mantenimiento Cero de Infraestructura Backend:** Al delegar la BD, WebSockets y la API RESTful en Supabase, el coste de mantenimiento del servidor es cero. Supabase autogenera una API usando *PostgREST* documentada sobre Postgres, otorgando rendimiento de base de datos directa sin escribir un ORM intermedio.
2. **Escalabilidad Reactiva con Angular:** Se optó por Angular RxJS ante la necesidad de manejar flujos de datos asíncronos complejos (autenticación en tiempo real, cargas de archivos, lectura del blog sin recargar páginas). La nueva versión Standalone de Angular elimina el sobrepeso de `NgModules`, ofreciendo "Tree-Shaking" perfecto, lo que significa que el peso de descarga del sitio web es mínimo.
3. **Despliegue Global Automático:** La integración de GitHub a Vercel habilita *Continuous Deployment* (CD). Cada push a `main` orquesta un pipeline inmutable que asegura 0% de tiempo de inactividad (Zero-Downtime Deployment).

### 1.3 Patrones de Diseño de Software Implementados
- **Smart / Dumb Components (Container-Presenter):** Patrón utilizado implícitamente en vistas donde el contenedor (ej. `profile.page.ts`) posee inyección de dependencias pesada y lógica de negocio, derivando datos limpios a las directivas `@if` y `@for` en el template.
- **Dependency Injection (DI):** Servicios encapsulados (`SupabaseService`, `AuthService`, `StorageService`) se inyectan en componentes, asegurando el patrón Singleton para conexiones de base de datos.
- **Facade Pattern:** El servicio `SupabaseService` expone métodos simplificados, escondiendo la complejidad original de los imports y la validación de tokens de la libreta cliente oficial de Supabase.

---

## 🛡️ 2. Patrones de Seguridad y Control de Acceso

La seguridad es el pilar de un blog colaborativo. El sistema combina validación de lado-cliente (Front) con endurecimiento en BD (Back).

### 2.1 Políticas RLS (Row Level Security) - "Zero Trust Data"
Toda la base de datos PostgreSQL sigue el principio de "Confianza Cero" (*Zero Trust*). Las consultas, incluso poseyendo la URL de la base de datos pública, son filtradas estrictamente antes de devolverse vía Supabase, dictando qué puede hacer un JWT específico.
* **Usuarios:** Una política (`policy`) solo permite hacer UPDATE a la fila vinculada al `auth.uid()` actual. Nadie puede mutar perfiles de otros.
* **Posts (Artículos):** Todos los usuarios anónimos o no pueden enviar lecturas `SELECT` **solamente** si `published = true`. Los borradores están invisibles para todos salvo el creador, protegiendo contenido sin publicar.
* **Posts - Admins:** Solo un usuario cuyo id conste en la tabla users con `role = 'admin'` puede hacer `DELETE` indiscriminadamente. Esto blinda a la aplicación contra Inyección de Scripts u operaciones nocivas generadas por herramientas de pentesting hacia la API pública.

### 2.2 Route Guards en el Frontend (Angular)
Protegemos el acceso y la experiencia de usuario (UX) mediante `CanActivate` Guards basados en el estado asíncrono del sistema:
- `auth.guard.ts`: Evita que personas anónimas fuercen la barra de direcciones de URL para entrar a áreas privadas como `/profile` o `/editor`.
- `admin.guard.ts`: Valida además que la cuenta activa posea un rol administrativo. Redirige inmediatamente cualquier intento de vulnerar el panel y expulsa al infractor.

### 2.3 Seguridad de Segundo Factor (2FA) y Notificaciones
La plataforma va más allá de un log-in común, implementando seguridad bancaria (Autenticación Multi-Factor).
1. **Google Authenticator (TOTP RFC 6238):** Activado en `/profile`. Para operaciones sensibles posteriores (por ejemplo, inicio del siguiente log-in), se exige la coordenada de 6 dígitos que solo el celular del portador posee.
2. **Alertas Inteligentes:** Se configuraron alertas mediante la API de Resend para notificaciones push vía mail al correo del usuario cada vez que la plataforma detecta un nuevo inicio de sesión.

---

## 🎨 3. Guía de Estilos y Design System (UI/UX)

La plataforma emplea una estética intencional denominada **"Dark Gaming / Cyberpunk Minimal"** para conectar estéticamente con audiencias amantes de los videojuegos, hardware, y cultura pop general.

### 3.1 Fundamentos UX y Psicología del Color
El uso del "Modo Oscuro como Defecto" no solo prolonga el tiempo de batería de dispositivos móviles OLED, sino que la psicología del color determina que los entornos nocturnos disminuyen drásticamente el estrés ocular provocado por exceso de luma durante la lectura de artículos y tutoriales largos.

| Concepto | Código CSS (HEX) | Justificación de su Uso |
| :--- | :--- | :--- |
| **Space Void (Body)** | `#1a1a2e` | _Negro violáceo de bajo contraste. Abandona el tradicional gris carbón por un matiz frío para separar visualmente contenido vs fondo._ |
| **Elevación (Cards / Nav)** | `#16213e` | _Materiales que parecen sobrepuestos respecto al fondo "Void". Permite agrupar la información (Tarjetas de post, cuadros de login)._ |
| **Action Red (Primario)** | `#e94560` | _Tono fucsia o "rojo neón". En la teoría del color, este derivado del rojo evoca adrenalina y velocidad (gaming). Usamos sus derivaciones al 20% para glows focales y al 100% para Call To Actions._ |
| **Text White Alpha** | `rgba(255,255,255,0.9)` | _Evitamos el blanco puro (`#fff`) en grandes bloques de texto ya que en pantallas ultra brillantes desgasta la retina frente a un fondo negro (`#1a1a2e`)._ |

### 3.2 Tipografía Analítica
Para mantener consistencia interplataforma y minimizar el tiempo de render (First Contentful Paint), prescindimos de cargar mega-librerías webfont pesadas, apostando por las fuentes del Sistema Operativo de la persona que lee:
```css
font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
```
- **Inter y SF Pro (Apple):** Abarcan el 80% de dispositivos móviles y laptops hoy, tienen altísima neutralidad y no distraen la lectura del Blog.
- **Alturas y Anchos:** Títulos robustos (Weight +700), párrafos con respiro visual (`line-height: 1.6`) asegurando que los artículos encajen fluidamente en Mobile y Desktop sin verse saturados.

### 3.3 Patrones Constructivos (Componentes UI)
- **Botones con Neumorfismo Oscuro Invertido:** En vez de sombras de impacto planas, los botones primarios `#e94560` gozan de un sutil resplandor box-shadow color neón al pasar el cursor (Hover State), confirmando que hay interactividad lista (`box-shadow: 0 4px 15px rgba(233,69,96,0.3)`).
- **Entradas Inmersivas (Inputs):** En vez de utilizar cajas de input convencionales blancas de los 90s, se usan rectángulos transparentes o más oscuros (`#0f172a`) que solo "despiertan" la frontera exterior cuando el usuario hace focus/click sobre ellos. Evitan la "ansiedad de llenado".
- **Glassmorphism (Efecto Cristal):** Se emplea en el componente Menú Header (`app-header`) a través de un simple `backdrop-filter: blur(10px)`. Da sensación inmersiva dejando vislumbrar el texto que hay detrás de un menú anclado en la cabecera sin robar visibilidad.
- **Skeleton Views:** Los clásicos "spinners" y "iconos de reloj de arena" fueron catalogados en UX moderno como ansiogénicos. En `/profile` y en el Dashboard de admin se sustituyen por rectángulos base con animación "Pulsar" (Skeletons). Le anticipan al cerebro la estructura de los datos venideros un par de milissegundos antes, dando percepción global de extrema velocidad.

---

## 💾 4. Almacenamiento y Estructura Relacional (DBMS)

La Rúbrica exige operaciones integrales CRUD manejando contenido y media. Se estructuró de la siguiente forma:

- **Estructura Cíclica:** Cada nueva vinculación al sistema Auth detona un Trigger en Postgres que genera una tabla `users` automáticamente, amarrando el "Owner" del ID a la entidad relacional.
- **Manejo de Imágenes en Nube:** Supabase Storage distribuye la media alojada mediante una CDN local y políticas RLS:
  - Solo usuarios `authenticated` pueden "SUBIR" al bucket `blog_assets`.
  - Archivos generados con `Date.now()` evitan bloqueos DNS por repetición de sintaxis en nombres de usuario.

### 4.1 Jerarquía Operativa
1. Rol de Guest (Lector Diario): Explora `posts` con estado true (Puntos 1, 6).
2. Rol Registered User (Subscriptor/Blogger): Posee acceso a su perfil y AuthGuard superado (Punto 2, 5). Postea en editor HTML Rico.
3. Rol Admin (Root Manager): CRUD total sin restricciones, expulsiones de miembros (Punto 4, 7).

*(Fin de esta iteración documental de la infraestructura de BabyWolf Blog)*
