package usecases

import (
	"context"
	"errors"
	"babywolf-blog-backend/internal/domain/entities"
	"babywolf-blog-backend/internal/domain/repositories"
)

type ManagePostsUseCase struct {
	postRepo repositories.PostRepository
}

func NewManagePostsUseCase(repo repositories.PostRepository) *ManagePostsUseCase {
	return &ManagePostsUseCase{postRepo: repo}
}

func (u *ManagePostsUseCase) CreatePost(ctx context.Context, post *entities.Post) error {
	if post.Title == "" || post.Content == "" {
		return errors.New("title and content are required")
	}
	// Slug generation would happen here ideally
	return u.postRepo.Create(ctx, post)
}

func (u *ManagePostsUseCase) GetPostBySlug(ctx context.Context, slug string) (*entities.Post, error) {
	return u.postRepo.GetBySlug(ctx, slug)
}

func (u *ManagePostsUseCase) ListPublishedPosts(ctx context.Context) ([]*entities.Post, error) {
	// False to includeUnpublished
	return u.postRepo.ListAll(ctx, false)
}

func (u *ManagePostsUseCase) DeletePost(ctx context.Context, id string, executorID string, executorRole string) error {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if executorRole != "admin" && post.AuthorID != executorID {
		return errors.New("unauthorized: only the author or an admin can delete a post")
	}

	return u.postRepo.Delete(ctx, id)
}
