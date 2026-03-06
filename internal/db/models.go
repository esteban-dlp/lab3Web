package db

type Series struct {
	ID             int
	Name           string
	CurrentEpisode int
	TotalEpisodes  int
	Rating         int // -1 means "no rating"
}

type ListParams struct {
	Query    string
	Sort     string // name|current|total
	Dir      string // asc|desc
	Page     int
	PageSize int
}