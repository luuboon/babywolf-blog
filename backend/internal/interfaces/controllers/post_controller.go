package controllers

import (
	"net/http"

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

func (c *PostController) ListPublicPosts(ctx *gin.Context) {
	posts, err := c.managePostsUseCase.ListPublishedPosts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list posts"})
		return
	}
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
