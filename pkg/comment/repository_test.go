package comment

import (
	"context"
	"fmt"
	"github.com/iBoBoTi/go-movie-api/internal/config"
	repo "github.com/iBoBoTi/go-movie-api/internal/database"
	"github.com/iBoBoTi/go-movie-api/pkg/comment/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"reflect"
	"testing"
)

var testRepo Repository
var testdb *pgxpool.Pool

func TestMain(m *testing.M) {
	conf, err := config.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	testdb, err = repo.ConnectPostgres(conf)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	testRepo = NewRespository(testdb)
	os.Exit(m.Run())
}

func Test_repository_AddComment(t *testing.T) {

	type args struct {
		c *types.Comment
	}
	tests := []struct {
		name string
		args args
		want *types.Comment
	}{
		{
			name: "Success",
			args: args{c: &types.Comment{
				MovieTitle: "A new movie",
				MovieID:    16,
				Author:     "172.18.0.1:59956",
				Content:    "An awesome movie",
			}},
			want: &types.Comment{
				MovieTitle: "A new movie",
				MovieID:    16,
				Author:     "172.18.0.1:59956",
				Content:    "An awesome movie",
			},
		},
	}
	defer clearTestDBTable("comments")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testRepo.AddComment(tt.args.c)
			tt.want.ID = got.ID
			tt.want.CreatedAt = got.CreatedAt
			require.Nil(t, err)
			require.Equal(t, got, tt.want)

		})
	}
}

func Test_repository_GetCommentsByMovieID(t *testing.T) {

	tests := []struct {
		name string
		args int
		want []types.Comment
	}{
		{
			name: "Success",
			args: 16,
			want: []types.Comment{
				types.Comment{
					MovieTitle: "A new movie",
					MovieID:    16,
					Author:     "172.18.0.1:59956",
					Content:    "An awesome movie",
				},
			},
		},
	}

	addCommentQuery := "INSERT INTO comments (movie_title, movie_id, author, content) VALUES ($1, $2, $3, $4) RETURNING *"
	defer clearTestDBTable("comments")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdb.Exec(context.Background(), addCommentQuery, tt.want[0].MovieTitle, tt.want[0].MovieID, tt.want[0].Author, tt.want[0].Content)
			got, err := testRepo.GetCommentsByMovieID(tt.args)
			require.Nil(t, err)
			require.Equal(t, reflect.TypeOf(got), reflect.TypeOf(tt.want))
			require.Equal(t, len(got), len(tt.want))
		})
	}
}

func Test_repository_GetCountForMovie(t *testing.T) {
	type args struct {
		movieTitle string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Success",
			args: args{
				movieTitle: "A new movie",
			},
			want: 1,
		},
	}

	comment := types.Comment{
		MovieTitle: "A new movie",
		MovieID:    16,
		Author:     "172.18.0.1:59956",
		Content:    "An awesome movie",
	}

	addCommentQuery := "INSERT INTO comments (movie_title, movie_id, author, content) VALUES ($1, $2, $3, $4) RETURNING *"
	testdb.Exec(context.Background(), addCommentQuery, comment.MovieTitle, comment.MovieID, comment.Author, comment.Content)
	defer clearTestDBTable("comments")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testRepo.GetCountForMovie(tt.args.movieTitle)
			require.Nil(t, err)
			if got != tt.want {
				t.Errorf("GetCountForMovie() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func clearTestDBTable(name string) {
	_, err := testdb.Exec(context.Background(), fmt.Sprintf("DELETE FROM %v", name))
	if err != nil {
		panic(err)
	}
}
