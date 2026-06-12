# Práctica de Laboratorio 4
## Buscador Interno del Sitio Web — Implementación y Pruebas
**Materia:** Desarrollo y Gestión de Software
**Proyecto:** BabyWolf Blog
**Autor:** Abraham

---

## 1. Arquitectura del buscador

El buscador es **client-side**: descarga todos los posts publicados desde Supabase y filtra en memoria.

```
Header input (input event)
    │
SearchStateService.setSearchQuery(query)  [Angular Signal]
    │
Navegación a /search?q=query
    │
SearchResultsPage.ngOnInit()
route.queryParams.subscribe(params)
    │
PostRepository.getPosts()  ← Supabase: SELECT * WHERE published=true
    │
posts.filter(p =>
  p.title.toLowerCase().includes(q)    ||
  p.excerpt.toLowerCase().includes(q)  ||
  p.content.toLowerCase().includes(q)  ||
  p.category.toLowerCase().includes(q)
)
    │
Renderizado con @for
```

## 2. Archivos involucrados

| Archivo | Rol |
|---------|-----|
| `core/services/search-state.service.ts` | Signal global del query |
| `core/layout/header/header.ts` | Captura input, llama `setSearchQuery()` |
| `features/search/search-results.ts` | Lee queryParams, filtra, renderiza |
| `features/posts/infrastructure/repositories/supabase-post.repository.ts` | `getPosts()` |

## 3. Casos de prueba

| ID | Query | Esperado | Campo | Resultado |
|----|-------|---------|-------|-----------|
| TC-01 | `rtx` | ≥ 1 resultado | título/contenido | ✅ PASA |
| TC-02 | `gaming` | ≥ 1 resultado | categoría/contenido | ✅ PASA |
| TC-03 | `ANGULAR` | ≥ 1 resultado | case-insensitive | ✅ PASA |
| TC-04 | `2026` | ≥ 1 resultado | título/excerpt | ✅ PASA |
| TC-05 | `hardware` | ≥ 1 resultado | categoría | ✅ PASA |
| TC-06 | ` ` (espacio) | 0 resultados | trim() vacío | ✅ PASA |
| TC-07 | `xkwqz` | 0 resultados | sin coincidencias | ✅ PASA |
| TC-08 | `a` | ≥ 1 resultado | carácter único | ✅ PASA |
| TC-09 | `mexico` | Condicional | sin normalización de acentos | ⚠️ |
| TC-10 | `borrador` | 0 resultados | RLS: published=false | ✅ PASA |

## 4. Limitaciones identificadas

| Limitación | Impacto |
|-----------|--------|
| Sin normalización de acentos ("mexico" ≠ "México") | Bajo-medio |
| Carga todos los posts en memoria | Medio-alto con 1000+ posts |
| Sin debounce en el input del header | Bajo |
| Sin paginación de resultados | Bajo |

## 5. Mejora recomendada (full-text search en Supabase)

```typescript
// En lugar de traer todos los posts y filtrar client-side:
this.sb.client
  .from('posts')
  .select('*')
  .textSearch('title', query, { type: 'websearch' })
  .eq('published', true)
```

---
*Práctica de Laboratorio 4 — Desarrollo y Gestión de Software*
