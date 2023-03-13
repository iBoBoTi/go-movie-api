package usecase

import (
	"github.com/iBoBoTi/go-movie-api/internal/domain"
	repo "github.com/iBoBoTi/go-movie-api/internal/respository"
)

type CommentService interface {
	AddComment(comment *domain.Comment) (*domain.Comment, error)
	GetCommentCountForMovie(movieTitle string) (int, error)
	GetCommentsByMovieByID(movieID int) ([]domain.Comment, error)
}

type commentService struct {
	commentRepository repo.CommentRepository
}

func NewCommentService(repo repo.CommentRepository) CommentService {
	return &commentService{commentRepository: repo}
}

func (m *commentService) GetCommentCountForMovie(movieTitle string) (int, error) {
	return m.commentRepository.GetCommentCountForMovie(movieTitle)
}

func (m *commentService) GetCommentsByMovieByID(movieID int) ([]domain.Comment, error) {
	return m.commentRepository.GetCommentsByMovieByID(movieID)
}

func (m *commentService) AddComment(comment *domain.Comment) (*domain.Comment, error) {
	return m.commentRepository.AddComment(comment)
}
