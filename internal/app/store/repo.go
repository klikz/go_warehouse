package store

import "warehouse/internal/app/models"

type Repo interface {
	Create(*models.User) error
	FindByEmail(*models.User) (*models.User, error)
	GetAllComponents() (interface{}, error)
}
