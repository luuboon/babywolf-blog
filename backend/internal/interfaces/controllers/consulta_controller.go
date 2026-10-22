package controllers

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/infrastructure/security"
)

// Práctica 3 (Desarrollo móvil): ConsultaWeb.
//
// La app pide una URL preasignada (GET /api/consulta/url), hace la petición
// HTTP a esa URL (GET /api/consulta/recurso?exp=…&sig=…) y muestra el
// resumen del sitio. Si la URL se modifica, la firma no coincide y el
// recurso responde 403.

const (
	propositoConsulta = "consulta"
	ttlConsulta       = 5 * time.Minute
)

type ConsultaController struct {
	posts *usecases.ManagePostsUseCase
}

func NewConsultaController(posts *usecases.ManagePostsUseCase) *ConsultaController {
	return &ConsultaController{posts: posts}
}

// baseURL arma el origen público del backend. En Railway llega detrás de un
// proxy, así que X-Forwarded-Proto dice si la petición original era HTTPS.
// PUBLIC_BASE_URL permite fijarlo a mano si hiciera falta.
func baseURL(ctx *gin.Context) string {
	if b := strings.TrimRight(os.Getenv("PUBLIC_BASE_URL"), "/"); b != "" {
		return b
	}
	esquema := "http"
	if p := ctx.GetHeader("X-Forwarded-Proto"); p != "" {
		esquema = strings.Split(p, ",")[0]
	} else if ctx.Request.TLS != nil {
		esquema = "https"
	}
	return esquema + "://" + ctx.Request.Host
}

// URL devuelve una URL preasignada válida por 5 minutos.
func (c *ConsultaController) URL(ctx *gin.Context) {
	q, expira, err := security.Firmar(propositoConsulta, url.Values{"recurso": {"resumen"}}, ttlConsulta, time.Now())
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "El servicio de consulta no está disponible."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"url":          baseURL(ctx) + "/api/consulta/recurso?" + q.Encode(),
		"expira":       expira.Format(time.RFC3339),
		"ttl_segundos": int(ttlConsulta.Seconds()),
	})
}

type ultimaNoticia struct {
	Titulo    string `json:"titulo"`
	Categoria string `json:"categoria"`
	Fecha     string `json:"fecha"`
}

// Recurso es el contenido protegido por la URL preasignada: un resumen del
// blog generado en el momento.
func (c *ConsultaController) Recurso(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	if err := security.Verificar(propositoConsulta, q, time.Now()); err != nil {
		switch {
		case errors.Is(err, security.ErrExpirada):
			ctx.JSON(http.StatusGone, gin.H{"error": "La URL expiró. Solicita una nueva."})
		case errors.Is(err, security.ErrSinClave):
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "El servicio de consulta no está disponible."})
		default:
			ctx.JSON(http.StatusForbidden, gin.H{"error": "La URL no es válida o fue modificada."})
		}
		return
	}

	posts, err := c.posts.ListPublishedPosts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No fue posible generar el resumen."})
		return
	}

	porCategoria := map[string]int{}
	ultimas := []ultimaNoticia{}
	for i, p := range posts {
		cat := p.Category
		if cat == "" {
			cat = "general"
		}
		porCategoria[cat]++
		if i < 3 { // el repositorio ya las entrega de la más reciente a la más antigua
			ultimas = append(ultimas, ultimaNoticia{
				Titulo:    p.Title,
				Categoria: cat,
				Fecha:     p.CreatedAt.Format("2006-01-02"),
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sitio":          "BabyWolf Blog",
		"descripcion":    "Noticias de videojuegos retro y consolas",
		"generado":       time.Now().UTC().Format(time.RFC3339),
		"total_noticias": len(posts),
		"por_categoria":  porCategoria,
		"ultimas":        ultimas,
	})
}
