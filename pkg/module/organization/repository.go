package organization

import (
	"context"

	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, m *model.Organization) (*model.Organization, error)
	FindByID(ctx context.Context, id string) (*model.Organization, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, m *model.Organization) (*model.Organization, error) {
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*model.Organization, error) {
	var m model.Organization
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
