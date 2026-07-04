package admin

import (
	"context"
	"time"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetStats(ctx context.Context) (*entity.AdminStatsDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetStats runs all counts in a single function so the FE only needs one
// round-trip to populate the dashboard. All queries are platform-wide —
// internal admins are not scoped to a tenant.
func (r *repository) GetStats(ctx context.Context) (*entity.AdminStatsDto, error) {
	stats := &entity.AdminStatsDto{}

	if err := r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Count(&stats.TotalCustomers).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Where("customer_type = ?", model.CustomerTypeIndividual).
		Count(&stats.IndividualCustomers).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Where("customer_type = ?", model.CustomerTypeCompany).
		Count(&stats.CompanyCustomers).Error; err != nil {
		return nil, err
	}

	// "This month" uses calendar-month boundaries. created_at is the
	// canonical sign-up timestamp on the customer table.
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if err := r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Where("created_at >= ?", monthStart).
		Count(&stats.CustomersThisMonth).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Model(&model.Organization{}).
		Count(&stats.TotalOrganizations).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("user_type = ?", model.INTERNAL_ADMIN).
		Count(&stats.TotalInternalAdmins).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
