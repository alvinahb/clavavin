package dbrepo

import (
	"github.com/alvinahb/clavavin/internal/models"
)

// InsertWine inserts a wine into the database
func (m *testDBRepo) InsertWine(res models.Wine) error {
	return nil
}

// AllWinesSummary returns a slice of all wines in database
func (m *testDBRepo) AllWinesSummary() ([]models.Wine, error) {
	var wines []models.Wine

	return wines, nil
}

// WineByID returns the information about a wine given its ID
func (m *testDBRepo) WineByID(id int) (models.Wine, error) {
	var wine models.Wine

	return wine, nil
}

// UserByID returns the information about a user given its ID
func (m *testDBRepo) UserByID(id int) (models.User, error) {
	var user models.User

	return user, nil
}

// UpdateUser updates a user in the database
func (m *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

// Authenticate authenticates the user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
