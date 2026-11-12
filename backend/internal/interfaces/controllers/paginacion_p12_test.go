package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/domain/entities"
)

// Práctica 12 · pruebas de la paginación de GET /api/posts.

func routerPaginado(n int) *gin.Engine {
	posts := make([]*entities.Post, n)
	t0 := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	for i := range posts {
		posts[i] = &entities.Post{ID: fmt.Sprintf("p%03d", i), Slug: fmt.Sprintf("s%03d", i), Published: true, CreatedAt: t0.Add(-time.Duration(i) * time.Minute)}
	}
	c := NewPostController(usecases.NewManagePostsUseCase(&repoFalso{posts: posts}))
	r := gin.New()
	r.GET("/api/posts", c.ListPublicPosts)
	return r
}

func TestP12_ParsePagina_ValoresLimite(t *testing.T) {
	casos := []struct {
		limit, offset string
		wl, wo        int
		invalido      bool
	}{
		{"", "", PostsLimiteDefecto, 0, false},
		{"1", "0", 1, 0, false},
		{"100", "5", 100, 5, false},
		{"0", "", 0, 0, true},   // bajo el mínimo
		{"101", "", 0, 0, true}, // sobre el máximo
		{"-3", "", 0, 0, true},
		{"abc", "", 0, 0, true},
		{"", "-1", 0, 0, true},
		{"", "x", 0, 0, true},
	}
	for _, c := range casos {
		l, o, msg := parsePagina(c.limit, c.offset)
		if c.invalido {
			if msg == "" {
				t.Errorf("limit=%q offset=%q debía ser inválido", c.limit, c.offset)
			}
			continue
		}
		if msg != "" || l != c.wl || o != c.wo {
			t.Errorf("limit=%q offset=%q -> %d,%d,%q; esperado %d,%d", c.limit, c.offset, l, o, msg, c.wl, c.wo)
		}
	}
}

func TestP12_Listar_PaginaYTotal(t *testing.T) {
	r := routerPaginado(120)
	w := hacer(r, "GET", "/api/posts?limit=10&offset=5", "", nil)
	var lista []map[string]any
	if w.Code != http.StatusOK || json.Unmarshal(w.Body.Bytes(), &lista) != nil || len(lista) != 10 {
		t.Fatalf("página: %d %s", w.Code, w.Body.String())
	}
	if lista[0]["id"] != "p005" {
		t.Errorf("la página no empieza en el offset: %v", lista[0]["id"])
	}
	if got := w.Header().Get("X-Total-Count"); got != "120" {
		t.Errorf("X-Total-Count = %q, esperado 120", got)
	}
}

func TestP12_Listar_SinParametros_UsaLimiteDefecto(t *testing.T) {
	w := hacer(routerPaginado(120), "GET", "/api/posts", "", nil)
	var lista []map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &lista)
	if len(lista) != PostsLimiteDefecto {
		t.Errorf("sin parámetros devolvió %d, esperado %d", len(lista), PostsLimiteDefecto)
	}
}

func TestP12_Listar_PocosPosts_SigueSiendoArreglo(t *testing.T) {
	// Compatibilidad: con menos posts que el límite el cliente recibe todos.
	w := hacer(routerPaginado(8), "GET", "/api/posts", "", nil)
	var lista []map[string]any
	if json.Unmarshal(w.Body.Bytes(), &lista) != nil || len(lista) != 8 {
		t.Fatalf("esperaba 8, %s", w.Body.String())
	}
}

func TestP12_Listar_OffsetFueraDeRango_ArregloVacio(t *testing.T) {
	w := hacer(routerPaginado(5), "GET", "/api/posts?offset=99", "", nil)
	if w.Code != http.StatusOK || w.Body.String() != "[]" {
		t.Errorf("fuera de rango: %d %q (debe ser 200 y [])", w.Code, w.Body.String())
	}
}

func TestP12_Listar_ParametroInvalido_400(t *testing.T) {
	for _, q := range []string{"?limit=0", "?limit=101", "?limit=x", "?offset=-1"} {
		if w := hacer(routerPaginado(5), "GET", "/api/posts"+q, "", nil); w.Code != http.StatusBadRequest {
			t.Errorf("%s -> %d, esperado 400", q, w.Code)
		}
	}
}
