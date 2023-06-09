package repository

import "github.com/alvinahb/clavavin/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertWine(res models.Wine) error
}
