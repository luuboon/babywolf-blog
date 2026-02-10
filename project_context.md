# BABY WOLF BLOG — Project Context

## 1. Objetivo del proyecto

BABY WOLF BLOG es un proyecto académico–profesional construido con:

- Frontend: Angular 21 (standalone)
- UI: Angular Material
- Package manager: pnpm
- Estilo visual: Retro / Neumorphic 90s
- Arquitectura: Clean Architecture (estricta)

El proyecto funciona como:

- Blog personal
- Organizador de ideas
- Plataforma futura para posts, comentarios y usuarios

El objetivo es aprender:

- Arquitectura limpia real
- Separación frontend/backend
- Componentización Angular moderna
- Buenas prácticas (sin antipatrón)
- Layout web completo según esquema clásico:
  - Header
  - Sidebar izquierda
  - Main content
  - Footer (pendiente)

---

## 2. Estado actual del proyecto

### Implementado

- Angular 21 standalone
- Angular Material configurado
- Header con navegación interna:
  - Home
  - Posts
  - Contact
- Sidebar izquierda con enlaces externos:
  - GitHub
  - LinkedIn
  - Docs
- Layout principal con grid:




- External links colocados en la parte inferior izquierda (sidebar)
- Router funcionando
- Rutas placeholder:
  - /
  - /posts
  - /contact

---

## 3. Estilo visual exacto (retro)

Inspirado en botones retro 90s:

### Colores base

- Fondo principal: #f6f2ec
- Cards: #e6e0d6
- Bordes: #1b1b1b

### Componentes retro

Todos los botones/enlaces usan:

```scss
border: 3px solid #1b1b1b;
box-shadow: 4px 4px #1b1b1b;
background: #e6e0d6;
border-radius: 6px;
font-weight: 800;
letter-spacing: 0.06em;


active state 

translate: 2px 2px;
box-shadow: 2px 2px #1b1b1b;


side bar card

border: 4px solid #1b1b1b;
box-shadow: 6px 6px #1b1b1b;
border-radius: 10px;
```

NO usar:

gradients modernos

glassmorphism

tailwind look

rounded excesivo

Todo debe verse sólido y físico.

frontend structure angular

src/app
│
├── app.ts (root standalone)
├── app.html
├── app.scss
├── app.routes.ts
│
└── shared
    ├── external-links
    │   └── external-links.(ts/html/scss)
    │
    └── layout
        └── header
            └── header.(ts/html/scss)

Angular es standalone:

Cada componente importa explícitamente lo que usa.

NO módulos tradicionales.

5. Reglas Angular obligatorias
Cada componente debe ser:
standalone: true

Cada componente importa sus dependencias:

imports: [...]

NO usar:

God components

Servicios gigantes

lógica en templates

acceso directo a APIs en componentes

6. Arquitectura objetivo completa

Backend (a implementar):

backend/
│
├── cmd/
│
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   └── repositories/
│   │
│   ├── application/
│   │   └── usecases/
│   │
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── http/
│   │   └── repositories/
│   │
│   └── interfaces/
│       └── controllers/
Frontend:

Solo habla con API Gateway

Nunca conoce base de datos

Nunca conoce ORM

7. Clean Architecture (reglas estrictas)

### Capas Generales:
- **Entities** (Enterprise Business Rules)
- **Use Cases** (Application Business Rules)
- **Interface Adapters** (Controllers, Gateways, Presenters)
- **Frameworks & Drivers** (Web, DB, UI, External Interfaces)

Dependencias SOLO hacia adentro:
UI → Controllers → UseCases → Entities
Nunca al revés.

### 7.1 Capas Frontend Específicas
1.  **Presentation Layer**: 
    - Componentes (Smart/Dumb) y Pages.
    - RESPONSABILIDAD: Solo lógica de UI y despacho de acciones. NO lógica de negocio.
    - Consume: State selectors o Use Cases.

2.  **Domain Layer**:
    - Modelos puros, Interfaces de Repositorios, Casos de Uso (Services de fachada).
    - REGLA: NO conoce HTTP, ni Angular HTTP Client, ni DTOs de backend.
    - Solo opera con entidades de dominio.

3.  **Infrastructure Layer**:
    - Implementación real de Repositorios, APIs HTTP.
    - Responsable de los MAPPERS (DTO <-> Domain).

8. Antipatrones prohibidos

🚫 Fat Controllers
🚫 Fat Components
🚫 lógica en views
🚫 repositorios con lógica
🚫 services que hacen de todo
🚫 singleton global state
🚫 models compartidos frontend/backend (DTO !== Domain Model)
🚫 llamadas HTTP directas desde componentes
🚫 carpetas "utils" genéricas
🚫 carpetas "helpers"
🚫 usar localStorage como DB
🚫 mezclar dominio con infraestructura

9. Estrategia de Gestión de Estado (Frontend)

- **Signals (Angular 21)**: Uso mandatorio de Signals para reactividad y fine-grained reactivity.
- **State Services**: El estado no reside en componentes. Se usan servicios de estado dedicados por Feature (ej. `PostStateService`).
- **Unidirectional Data Flow**: Action -> State Update -> Signal Change -> UI Update.

10. Mappers y Separación de Modelos

- **Prohibido usar DTOs en la Vista**: La API devuelve DTOs (Data Transfer Objects).
- **Mappers Obligatorios**: La capa de infraestructura DEBE mapear DTO -> Dominio antes de entregar datos a la capa de aplicación/ui.
- **Modelo de Dominio Puro**: Las interfaces de UI/Dominio no deben tener decoradores de serialización ni acoplamiento al backend.

11. Testing Pillars

- **Unit Tests**: Obligatorios para Mappers, Validadores y Lógica de Negocio (Domain/UseCases). Coverage alto.
- **Component Tests**: Shallow tests para asegurar que la UI renderiza estado y emite eventos. Lógica mínima.
- **E2E**: Flujos críticos (Smoke Tests).

12. Funcionalidad futura planeada
Core:

Posts

Comentarios

Usuarios

Admin panel

Roles:

Admin (1)

User (muchos)

Roadmap:

Auth

CRUD posts

CRUD comments

Moderación

Dashboard

13. External links (implementados)

Ubicados en sidebar izquierda inferior.

Buenas prácticas aplicadas:

target="_blank"

rel="noopener noreferrer"

Usados para:

GitHub (repositorio)

LinkedIn (perfil)

Docs (Angular)

14. Filosofía del proyecto

Este NO es un CRUD simple.

Este proyecto sirve para:

aprender arquitectura real

diseño frontend serio

separación de responsabilidades

patrones profesionales

evitar soluciones rápidas

Cada decisión debe ser:

explícita

razonada

mantenible

15. Regla final

Si algo parece fácil pero rompe arquitectura:

NO SE HACE.

Siempre preferir estructura > velocidad.
