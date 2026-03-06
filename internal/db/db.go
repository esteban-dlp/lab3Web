package db

import (
  "database/sql"
  "errors"
  "fmt"
  "strings"

  _ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// Verify connection
	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, err
	}
	return database, nil
}

func InitSchema(database *sql.DB) error {
	_, err := database.Exec(`
CREATE TABLE IF NOT EXISTS series (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  current_episode INTEGER NOT NULL,
  total_episodes INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS series_rating (
  series_id INTEGER PRIMARY KEY,
  rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 10),
  FOREIGN KEY(series_id) REFERENCES series(id) ON DELETE CASCADE
);
`)
	return err
}

func InsertSeries(database *sql.DB, name string, current, total int) error {
	_, err := database.Exec(
		`INSERT INTO series (name, current_episode, total_episodes) VALUES (?, ?, ?)`,
		name, current, total,
	)
	return err
}

func UpdateSeries(database *sql.DB, id int, name string, current, total int) error {
	res, err := database.Exec(
		`UPDATE series SET name = ?, current_episode = ?, total_episodes = ? WHERE id = ?`,
		name, current, total, id,
	)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("serie no encontrada")
	}
	return nil
}

func DeleteSeries(database *sql.DB, id int) error {
	res, err := database.Exec(`DELETE FROM series WHERE id = ?`, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("serie no encontrada")
	}
	return nil
}

func IncrementEpisode(database *sql.DB, id int) (Series, error) {
	_, err := database.Exec(`
UPDATE series
SET current_episode = current_episode + 1
WHERE id = ? AND current_episode < total_episodes
`, id)
	if err != nil {
		return Series{}, err
	}

	var s Series
	err = database.QueryRow(`
SELECT s.id, s.name, s.current_episode, s.total_episodes, COALESCE(r.rating, -1)
FROM series s
LEFT JOIN series_rating r ON r.series_id = s.id
WHERE s.id = ?
`, id).Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.Rating)
	if err != nil {
		return Series{}, err
	}
	return s, nil
}

func SetRating(database *sql.DB, id int, rating int) error {
	// Upsert rating
	_, err := database.Exec(`
INSERT INTO series_rating(series_id, rating) VALUES(?, ?)
ON CONFLICT(series_id) DO UPDATE SET rating = excluded.rating
`, id, rating)
	return err
}

func GetSeriesByID(database *sql.DB, id int) (Series, error) {
	var s Series
	err := database.QueryRow(`
SELECT s.id, s.name, s.current_episode, s.total_episodes, COALESCE(r.rating, -1)
FROM series s
LEFT JOIN series_rating r ON r.series_id = s.id
WHERE s.id = ?
`, id).Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.Rating)
	return s, err
}

func ListSeries(database *sql.DB, p ListParams) ([]Series, int, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 50 {
		p.PageSize = 50
	}

	sortCol := "s.name"
	switch p.Sort {
	case "name":
		sortCol = "s.name"
	case "current":
		sortCol = "s.current_episode"
	case "total":
		sortCol = "s.total_episodes"
	default:
		p.Sort = "name"
		sortCol = "s.name"
	}

	dir := "ASC"
	if strings.ToLower(p.Dir) == "desc" {
		dir = "DESC"
		p.Dir = "desc"
	} else {
		p.Dir = "asc"
	}

	q := strings.TrimSpace(p.Query)
	like := "%"
	if q != "" {
		like = "%" + q + "%"
	} else {
		like = "%"
	}

	// Count for pagination
	var totalCount int
	err := database.QueryRow(`
SELECT COUNT(*)
FROM series s
WHERE s.name LIKE ?
`, like).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	offset := (p.Page - 1) * p.PageSize

	query := fmt.Sprintf(`
SELECT s.id, s.name, s.current_episode, s.total_episodes, COALESCE(r.rating, -1)
FROM series s
LEFT JOIN series_rating r ON r.series_id = s.id
WHERE s.name LIKE ?
ORDER BY %s %s
LIMIT ? OFFSET ?
`, sortCol, dir)

	rows, err := database.Query(query, like, p.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []Series
	for rows.Next() {
		var s Series
		if err := rows.Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.Rating); err != nil {
			return nil, 0, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return out, totalCount, nil
}