package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/alvinahb/clavavin/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// InsertWine inserts a wine into the database
func (m *postgresDBRepo) InsertWine(res models.Wine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into wines (name, domain, year, appellation_type,
		appellation_name, location, color, culture, varieties, robe, nose,
		taste, dishes, season, created_at, updated_at) values ($1, $2, $3,
		$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err := m.DB.ExecContext(ctx, query,
		res.Name,
		res.Domain,
		res.Year,
		res.AppellationType,
		res.AppellationName,
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

	query := `select id, name, domain, year, appellation_type, appellation_name, location, color from wines`

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
			&wine.AppellationType,
			&wine.AppellationName,
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

// WineByID returns the information about a wine given its ID
func (m *postgresDBRepo) WineByID(id int) (models.Wine, error) {
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wine models.Wine

	query := `select id, name, domain, year, appellation_type, appellation_name,
		location, color, culture, varieties, robe, nose, taste, dishes, season
		from wines where id=$1`

	row := m.DB.QueryRow(query, id)
	err := row.Scan(
		&wine.ID,
		&wine.Name,
		&wine.Domain,
		&wine.Year,
		&wine.AppellationType,
		&wine.AppellationName,
		&wine.Location,
		&wine.Color,
		&wine.Culture,
		&wine.Varieties,
		&wine.Robe,
		&wine.Nose,
		&wine.Taste,
		&wine.Dishes,
		&wine.Season,
	)
	if err != nil {
		return wine, err
	}

	if err = row.Err(); err != nil {
		return wine, err
	}

	return wine, nil
}

// UserByID returns the information about a user given its ID
func (m *postgresDBRepo) UserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	query := `select id, first_name, last_name, email, password, access_level,
		created_at, updated_at from users where id=$1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	if err = row.Err(); err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3,
		access_level=$4, updated_at=$5`

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates the user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `select id, password from users where email=$1`

	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
