package dbrepo

import (
	"context"
	"time"

	"github.com/alvinahb/clavavin/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertWine inserts a wine into the database
func (m *postgresDBRepo) InsertWine(res models.Wine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into wines (name, domain, year, appellation,
		location, color, culture, varieties, robe, nose, taste, 
		dishes, season, created_at, updated_at) values ($1, $2, $3,
		$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err := m.DB.ExecContext(ctx, query,
		res.Name,
		res.Domain,
		res.Year,
		res.Appellation,
		res.Location,
		res.Color,
		res.Culture,
		res.Varieties,
		res.Robe,
		res.Nose,
		res.Taste,
		res.Dishes,
		res.Season,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// AllWinesSummary returns a slice of all wines in database
func (m *postgresDBRepo) AllWinesSummary() ([]models.Wine, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wines []models.Wine

	query := `select id, name, domain, year, appellation, location, color from wines`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return wines, err
	}

	for rows.Next() {
		var wine models.Wine
		err = rows.Scan(
			&wine.ID,
			&wine.Name,
			&wine.Domain,
			&wine.Year,
			&wine.Appellation,
			&wine.Location,
			&wine.Color,
		)
		if err != nil {
			return wines, err
		}

		wines = append(wines, wine)
	}

	if err = rows.Err(); err != nil {
		return wines, err
	}

	return wines, nil
}
