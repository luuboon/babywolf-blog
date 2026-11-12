package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/domain/entities"
)

type PostController struct {
	managePostsUseCase *usecases.ManagePostsUseCase
}

func NewPostController(useCase *usecases.ManagePostsUseCase) *PostController {
	return &PostController{managePostsUseCase: useCase}
}

// Paginación de GET /api/posts (Práctica 12). Sin parámetros se devuelve la
// primera página de PostsLimiteDefecto, así el JSON sigue siendo un arreglo y
// los clientes existentes no cambian; X-Total-Count dice cuántos hay en total.
const (
	PostsLimiteDefecto = 50
	PostsLimiteMax     = 100
)

// parsePagina valida limit y offset. Devuelve un mensaje si son inválidos.
func parsePagina(limitStr, offsetStr string) (limit, offset int, errMsg string) {
	limit, offset = PostsLimiteDefecto, 0
	if limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err != nil || n < 1 || n > PostsLimiteMax {
			return 0, 0, "limit debe ser un entero entre 1 y " + strconv.Itoa(PostsLimiteMax)
		}
		limit = n
	}
	if offsetStr != "" {
		n, err := strconv.Atoi(offsetStr)
		if err != nil || n < 0 {
			return 0, 0, "offset debe ser un entero mayor o igual a 0"
		}
		offset = n
	}
	return limit, offset, ""
}

func (c *PostController) ListPublicPosts(ctx *gin.Context) {
	limit, offset, msg := parsePagina(ctx.Query("limit"), ctx.Query("offset"))
	if msg != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	posts, total, err := c.managePostsUseCase.ListPublishedPostsPage(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list posts"})
		return
	}
	if posts == nil {
		posts = []*entities.Post{} // [] y no null
	}
	ctx.Header("X-Total-Count", strconv.FormatInt(total, 10))
	ctx.JSON(http.StatusOK, posts)
}

func (c *PostController) CreatePost(ctx *gin.Context) {
	var post entities.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming auth middleware set this
	if authorId, exists := ctx.Get("user_id"); exists {
		// Needs mapping from Supabase Auth ID to Database UUID in real scenario
		post.AuthorID = authorId.(string) 
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.managePostsUseCase.CreatePost(ctx.Request.Context(), &post); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

func (c *PostController) GetPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	post, err := c.managePostsUseCase.GetPostBySlug(ctx.Request.Context(), slug)
	
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	
	ctx.JSON(http.StatusOK, post)
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("id")
	userId := ctx.GetString("user_id")
	
	// App metadata injected by auth middleware
	var role string = "user"
	if meta, exists := ctx.Get("app_metadata"); exists {
		if metaMap, ok := meta.(map[string]interface{}); ok {
			role, _ = metaMap["role"].(string)
		}
	}

	if err := c.managePostsUseCase.DeletePost(ctx.Request.Context(), postId, userId, role); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
