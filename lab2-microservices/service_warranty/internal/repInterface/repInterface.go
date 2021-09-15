package repInterface

import (
	"services/internal/models"

	uuid "github.com/satori/go.uuid"
)

type RepInterface interface {
	StartWarrantyPeriod(itemUID uuid.UUID) (err error)
	GetWarranty(itemUID uuid.UUID) (warranty models.Warranty, err error)
	UpdateWarranty(warranty models.Warranty) (err error)
	CloseWarranty(itemUID uuid.UUID) (err error)
}
