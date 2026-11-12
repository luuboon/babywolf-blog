package controllers

// Práctica 11 — Fase 3: caja blanca de los controladores HTTP.
//
// Se ejercitan los endpoints con httptest (sin red ni base de datos) para
// comprobar que el código HTTP corresponde al resultado real, que el JSON es
// consistente y que no se filtran datos internos. Un repositorio falso
// sustituye a la base.
//
// Ejecutar: go test ./internal/interfaces/controllers -race -cover -v

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"babywolf-blog-backend/internal/application/connect"
	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/domain/entities"
	"babywolf-blog-backend/internal/infrastructure/security"
)

func init() {
	gin.SetMode(gin.TestMode)
	security.ClaveFn = func() ([]byte, error) { return []byte("clave-de-prueba"), nil }
}

// repoFalso implementa repositories.PostRepository en memoria.
type repoFalso struct {
	posts []*entities.Post
	err   error
}

func (r *repoFalso) GetByID(_ context.Context, id string) (*entities.Post, error) {
	if r.err != nil {
		return nil, r.err
	}
	for _, p := range r.posts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("record not found")
}
func (r *repoFalso) GetBySlug(_ context.Context, slug string) (*entities.Post, error) {
	if r.err != nil {
		return nil, r.err
	}
	for _, p := range r.posts {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, errors.New("record not found")
}
func (r *repoFalso) Create(_ context.Context, p *entities.Post) error { return r.err }
func (r *repoFalso) Update(_ context.Context, p *entities.Post) error { return r.err }
func (r *repoFalso) Delete(_ context.Context, id string) error        { return r.err }
func (r *repoFalso) ListAll(_ context.Context, _ bool) ([]*entities.Post, error) {
	return r.posts, r.err
}

func (r *repoFalso) ListPublishedPage(_ context.Context, limit, offset int) ([]*entities.Post, int64, error) {
	if r.err != nil {
		return nil, 0, r.err
	}
	total := int64(len(r.posts))
	if offset >= len(r.posts) {
		return nil, total, nil
	}
	fin := offset + limit
	if fin > len(r.posts) {
		fin = len(r.posts)
	}
	return r.posts[offset:fin], total, nil
}

func hacer(r http.Handler, metodo, ruta, cuerpo string, cabeceras map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	if cuerpo != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range cabeceras {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func jsonDe(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("la respuesta no es un objeto JSON: %v\n%s", err, w.Body.String())
	}
	return m
}

func muestras() []*entities.Post {
	t0 := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	return []*entities.Post{
		{ID: "p1", Title: "Uno", Slug: "uno", Category: "retro", Published: true, CreatedAt: t0, AuthorID: "autor"},
		{ID: "p2", Title: "Dos", Slug: "dos", Category: "", Published: true, CreatedAt: t0.Add(-time.Hour), AuthorID: "autor"},
		{ID: "p3", Title: "Tres", Slug: "tres", Category: "retro", Published: true, CreatedAt: t0.Add(-2 * time.Hour), AuthorID: "autor"},
		{ID: "p4", Title: "Cuatro", Slug: "cuatro", Category: "tech", Published: true, CreatedAt: t0.Add(-3 * time.Hour), AuthorID: "autor"},
	}
}

// ── WEB-01/02/06 · ConsultaWeb ──────────────────────────────────

func routerConsulta(repo *repoFalso) *gin.Engine {
	c := NewConsultaController(usecases.NewManagePostsUseCase(repo))
	r := gin.New()
	r.GET("/api/consulta/url", c.URL)
	r.GET("/api/consulta/recurso", c.Recurso)
	return r
}

func TestWB_Consulta_URL_200_EstructuraYTTL(t *testing.T) {
	t.Setenv("PUBLIC_BASE_URL", "")
	r := routerConsulta(&repoFalso{})
	w := hacer(r, "GET", "/api/consulta/url", "", map[string]string{"X-Forwarded-Proto": "https"})
	if w.Code != http.StatusOK {
		t.Fatalf("código %d", w.Code)
	}
	m := jsonDe(t, w)
	for _, k := range []string{"url", "expira", "ttl_segundos"} {
		if _, ok := m[k]; !ok {
			t.Errorf("falta el campo %q: %v", k, m)
		}
	}
	if !strings.HasPrefix(m["url"].(string), "https://") {
		t.Errorf("con X-Forwarded-Proto=https la URL debe ser https: %v", m["url"])
	}
	if m["ttl_segundos"].(float64) != 300 {
		t.Errorf("ttl: %v", m["ttl_segundos"])
	}
}

func TestWB_Consulta_Recurso_Codigos(t *testing.T) {
	repo := &repoFalso{posts: muestras()}
	r := routerConsulta(repo)

	// Una URL firmada válida, una vencida y una alterada.
	firmada := func(ahora time.Time) string {
		q, _, _ := security.Firmar(propositoConsulta, url.Values{"recurso": {"resumen"}}, ttlConsulta, ahora)
		return "/api/consulta/recurso?" + q.Encode()
	}

	// 200 · éxito
	w := hacer(r, "GET", firmada(time.Now()), "", nil)
	if w.Code != 200 {
		t.Fatalf("válida: %d %s", w.Code, w.Body.String())
	}
	m := jsonDe(t, w)
	if m["total_noticias"].(float64) != 4 {
		t.Errorf("total_noticias: %v", m["total_noticias"])
	}
	cats := m["por_categoria"].(map[string]any)
	if cats["retro"].(float64) != 2 || cats["general"].(float64) != 1 || cats["tech"].(float64) != 1 {
		t.Errorf("por_categoria (la vacía cuenta como 'general'): %v", cats)
	}
	if n := len(m["ultimas"].([]any)); n != 3 {
		t.Errorf("sólo las 3 más recientes: %d", n)
	}

	// 410 · vencida (firmada hace 10 min con vida de 5)
	if w := hacer(r, "GET", firmada(time.Now().Add(-10*time.Minute)), "", nil); w.Code != http.StatusGone {
		t.Errorf("vencida: %d", w.Code)
	}
	// 403 · alterada (cambia el recurso) y sin firma
	alterada := strings.Replace(firmada(time.Now()), "resumen", "otro", 1)
	if w := hacer(r, "GET", alterada, "", nil); w.Code != http.StatusForbidden {
		t.Errorf("alterada: %d", w.Code)
	}
	if w := hacer(r, "GET", "/api/consulta/recurso", "", nil); w.Code != http.StatusForbidden {
		t.Errorf("sin firma: %d", w.Code)
	}
	// 500 · la base falla; el mensaje no filtra el error interno
	repo.err = errors.New("pq: password authentication failed for user postgres")
	w = hacer(r, "GET", firmada(time.Now()), "", nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("falla de base: %d", w.Code)
	}
	if strings.Contains(w.Body.String(), "pq:") || strings.Contains(w.Body.String(), "password") {
		t.Errorf("el 500 filtra el error interno: %s", w.Body.String())
	}
	repo.err = nil
	// 503 · sin clave de firma
	security.ClaveFn = func() ([]byte, error) { return nil, security.ErrSinClave }
	defer func() { security.ClaveFn = func() ([]byte, error) { return []byte("clave-de-prueba"), nil } }()
	if w := hacer(r, "GET", "/api/consulta/url", "", nil); w.Code != http.StatusServiceUnavailable {
		t.Errorf("sin clave /url: %d", w.Code)
	}
	if w := hacer(r, "GET", "/api/consulta/recurso?sig=x&exp=1", "", nil); w.Code != http.StatusServiceUnavailable {
		t.Errorf("sin clave /recurso: %d", w.Code)
	}
}

// ── WEB-01/02/03 · posts públicos ───────────────────────────────

func TestWB_Posts_Listar_Y_Slug(t *testing.T) {
	repo := &repoFalso{posts: muestras()}
	c := NewPostController(usecases.NewManagePostsUseCase(repo))
	r := gin.New()
	r.GET("/api/posts", c.ListPublicPosts)
	r.GET("/api/posts/:slug", c.GetPostBySlug)

	w := hacer(r, "GET", "/api/posts", "", nil)
	var lista []map[string]any
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &lista) != nil || len(lista) != 4 {
		t.Fatalf("listar: %d %s", w.Code, w.Body.String())
	}
	for _, k := range []string{"id", "title", "slug", "content", "published", "created_at"} {
		if _, ok := lista[0][k]; !ok {
			t.Errorf("el JSON de un post no tiene %q", k)
		}
	}
	if w := hacer(r, "GET", "/api/posts/uno", "", nil); w.Code != 200 {
		t.Errorf("slug existente: %d", w.Code)
	}
	if w := hacer(r, "GET", "/api/posts/no-existe", "", nil); w.Code != http.StatusNotFound {
		t.Errorf("slug inexistente: %d", w.Code)
	}
	repo.err = errors.New("conexión rechazada")
	if w := hacer(r, "GET", "/api/posts", "", nil); w.Code != http.StatusInternalServerError {
		t.Errorf("falla de base: %d", w.Code)
	}
}

// ── WEB-02/03/04 · escritura protegida ──────────────────────────

func TestWB_Posts_Crear_Codigos(t *testing.T) {
	repo := &repoFalso{}
	c := NewPostController(usecases.NewManagePostsUseCase(repo))
	conUsuario := func(id any) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			if id != nil {
				ctx.Set("user_id", id)
			}
		}
	}
	r := gin.New()
	r.POST("/ok", conUsuario("u1"), c.CreatePost)
	r.POST("/anonimo", c.CreatePost)

	// 201 · válido
	if w := hacer(r, "POST", "/ok", `{"title":"T","content":"C"}`, nil); w.Code != http.StatusCreated {
		t.Errorf("válido: %d %s", w.Code, w.Body.String())
	}
	// 400 · JSON mal formado
	if w := hacer(r, "POST", "/ok", `{"title":`, nil); w.Code != http.StatusBadRequest {
		t.Errorf("JSON roto: %d", w.Code)
	}
	// 401 · sin usuario en el contexto
	if w := hacer(r, "POST", "/anonimo", `{"title":"T","content":"C"}`, nil); w.Code != http.StatusUnauthorized {
		t.Errorf("sin usuario: %d", w.Code)
	}
	// Validación semántica (título vacío): HOY responde 500. Lo correcto sería 422.
	w := hacer(r, "POST", "/ok", `{"title":"","content":""}`, nil)
	t.Logf("DEFECTO D-05 · título y contenido vacíos → código %d (esperado 422)", w.Code)
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusUnprocessableEntity {
		t.Errorf("código inesperado: %d", w.Code)
	}
}

func TestWB_Posts_Eliminar_Codigos(t *testing.T) {
	repo := &repoFalso{posts: muestras()}
	c := NewPostController(usecases.NewManagePostsUseCase(repo))
	como := func(id, rol string) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			ctx.Set("user_id", id)
			ctx.Set("app_metadata", map[string]interface{}{"role": rol})
		}
	}
	r := gin.New()
	r.DELETE("/autor/:id", como("autor", "user"), c.DeletePost)
	r.DELETE("/ajeno/:id", como("otro", "user"), c.DeletePost)
	r.DELETE("/admin/:id", como("otro", "admin"), c.DeletePost)

	if w := hacer(r, "DELETE", "/autor/p1", "", nil); w.Code != 200 {
		t.Errorf("el autor borra lo suyo: %d", w.Code)
	}
	if w := hacer(r, "DELETE", "/admin/p1", "", nil); w.Code != 200 {
		t.Errorf("el admin borra lo ajeno: %d", w.Code)
	}
	if w := hacer(r, "DELETE", "/ajeno/p1", "", nil); w.Code != http.StatusForbidden {
		t.Errorf("un usuario ajeno no puede borrar: %d", w.Code)
	}
	// Inexistente: HOY devuelve 403 porque el controlador traduce CUALQUIER
	// error a Forbidden. Lo correcto sería 404.
	w := hacer(r, "DELETE", "/admin/no-existe", "", nil)
	t.Logf("DEFECTO D-04 · borrar un post inexistente → código %d (esperado 404)", w.Code)
	if w.Code != http.StatusForbidden && w.Code != http.StatusNotFound {
		t.Errorf("código inesperado: %d", w.Code)
	}
	// Falla de la base: también sale 403 en lugar de 500.
	repo.err = errors.New("timeout")
	w = hacer(r, "DELETE", "/admin/p2", "", nil)
	t.Logf("DEFECTO D-04 · falla de la base al borrar → código %d (esperado 500)", w.Code)
}

// ── WEB-03 · MobileConnect ──────────────────────────────────────

func TestWB_Connect_Codigos(t *testing.T) {
	c := NewConnectController(connect.NewStore())
	r := gin.New()
	r.GET("/s/:sala/m", c.Listar)
	r.POST("/s/:sala/m", c.Enviar)
	r.DELETE("/s/:sala/m", c.Limpiar)

	if w := hacer(r, "POST", "/s/LAB-1/m", `{"emisor":"A","texto":"hola"}`, nil); w.Code != http.StatusCreated {
		t.Errorf("enviar: %d %s", w.Code, w.Body.String())
	}
	if w := hacer(r, "POST", "/s/LAB-1/m", `no es json`, nil); w.Code != http.StatusBadRequest {
		t.Errorf("cuerpo inválido: %d", w.Code)
	}
	if w := hacer(r, "POST", "/s/LAB-1/m", `{"emisor":"A","texto":"  "}`, nil); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("mensaje vacío: %d", w.Code)
	}
	if w := hacer(r, "POST", "/s/x/m", `{"emisor":"A","texto":"hola"}`, nil); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("sala inválida: %d", w.Code)
	}
	w := hacer(r, "GET", "/s/LAB-1/m", "", nil)
	m := jsonDe(t, w)
	if w.Code != 200 || len(m["mensajes"].([]any)) != 1 || m["conteo"].(map[string]any)["A"].(float64) != 1 {
		t.Errorf("listar: %d %v", w.Code, m)
	}
	if w := hacer(r, "GET", "/s/x/m", "", nil); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("listar sala inválida: %d", w.Code)
	}
	if w := hacer(r, "DELETE", "/s/LAB-1/m", "", nil); w.Code != http.StatusNoContent || w.Body.Len() != 0 {
		t.Errorf("limpiar debe ser 204 sin cuerpo: %d %q", w.Code, w.Body.String())
	}
	if w := hacer(r, "DELETE", "/s/x/m", "", nil); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("limpiar sala inválida: %d", w.Code)
	}
}
