package domain

import "time"

type Comment struct {
	ID         int       `json:"id"`
	MovieTitle string    `json:"movie_title"`
	MovieID    int       `json:"movie_id"`
	Author     string    `json:"author"`                                   // IP ADDRESS OF COMMENTER
	Content    string    `json:"content" binding:"required,min=1,max=500"` // Set limit to 500 characters
	CreatedAt  time.Time `json:"created_at"`
}
