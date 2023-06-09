package repository

import "github.com/alvinahb/clavavin/internal/models"

type DatabaseRepo interface {
	InsertWine(res models.Wine) error
	AllWinesSummary() ([]models.Wine, error)
	WineByID(id int) (models.Wine, error)
}
