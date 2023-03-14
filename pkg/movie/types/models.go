package types

type Movie struct {
	Title         string `json:"title"`
	OpeningCrawl  string `json:"opening_crawl"`
	ReleaseDate   string `json:"release_date"`
	CommentsCount int    `json:"comments_count"`
}
