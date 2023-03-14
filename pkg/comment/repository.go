package comment

import (
	"context"
	"github.com/iBoBoTi/go-movie-api/pkg/comment/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	AddComment(comment *types.Comment) (*types.Comment, error)
	GetCountForMovie(movieTitle string) (int, error)
	GetCommentsByMovieID(movieID int) ([]types.Comment, error)
}

type repository struct {
	Db *pgxpool.Pool
}

func NewRespository(Db *pgxpool.Pool) Repository {
	return &repository{Db: Db}
}

func (r *repository) AddComment(c *types.Comment) (*types.Comment, error) {
	queryString := `INSERT INTO comments (movie_title, movie_id, author, content) VALUES ($1, $2, $3, $4) RETURNING *`
	result := &types.Comment{}
	row := r.Db.QueryRow(context.Background(), queryString, c.MovieTitle, c.MovieID, c.Author, c.Content)
	err := row.Scan(&result.ID, &result.MovieTitle, &result.MovieID, &result.Author, &result.Content, &result.CreatedAt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) GetCountForMovie(movieTitle string) (int, error) {
	var count int
	queryString := "SELECT COUNT(*) FROM comments WHERE movie_title=$1"
	row := r.Db.QueryRow(context.Background(), queryString, movieTitle)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetCommentsByMovieID(movieID int) ([]types.Comment, error) {
	comments := make([]types.Comment, 0)
	queryString := `SELECT * FROM comments WHERE movie_id = $1 ORDER BY created_at DESC`
	rows, err := r.Db.Query(context.Background(), queryString, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := types.Comment{}
		err := rows.Scan(&c.ID, &c.MovieTitle, &c.MovieID, &c.Author, &c.Content, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
