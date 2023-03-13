package respository

import (
	"context"
	"github.com/iBoBoTi/go-movie-api/internal/domain"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CommentRepository interface {
	AddComment(comment *domain.Comment) (*domain.Comment, error)
	GetCommentCountForMovie(movieTitle string) (int, error)
	GetCommentsByMovieByID(movieID int) ([]domain.Comment, error)
}

type commentRepository struct {
	Db *pgxpool.Pool
}

func NewCommentRespository(Db *pgxpool.Pool) CommentRepository {
	return &commentRepository{Db: Db}
}

func (m *commentRepository) AddComment(comment *domain.Comment) (*domain.Comment, error) {
	queryString := `INSERT INTO comments (movie_title, movie_id, author, content) VALUES ($1, $2, $3, $4) RETURNING *`
	result := &domain.Comment{}
	row := m.Db.QueryRow(context.Background(), queryString, comment.MovieTitle, comment.MovieID, comment.Author, comment.Content)
	err := row.Scan(&result.ID, &result.MovieTitle, &result.MovieID, &result.Author, &result.Content, &result.CreatedAt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *commentRepository) GetCommentCountForMovie(movieTitle string) (int, error) {
	var count int
	queryString := "SELECT COUNT(*) FROM comments WHERE movie_title=$1"
	row := m.Db.QueryRow(context.Background(), queryString, movieTitle)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *commentRepository) GetCommentsByMovieByID(movieID int) ([]domain.Comment, error) {
	comments := make([]domain.Comment, 0)
	queryString := `SELECT * FROM comments WHERE movie_id = $1 ORDER BY created_at DESC`
	rows, err := m.Db.Query(context.Background(), queryString, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		comment := domain.Comment{}
		err := rows.Scan(&comment.ID, &comment.MovieTitle, &comment.MovieID, &comment.Author, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
