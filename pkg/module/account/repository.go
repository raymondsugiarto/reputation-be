package account

import (
	"context"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	FindByID(ctx context.Context, id string) (*entity.AccountDto, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]entity.AccountDto, error)
	FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.AccountDto], error)
	Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	m := dto.ToModel()
	err := r.db.WithContext(ctx).Create(m).Error
	if err != nil {
		return nil, err
	}
	dto.ID = m.ID
	return dto, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.AccountDto, error) {
	var m model.Account
	err := r.db.WithContext(ctx).Preload("Customer").First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return entity.NewAccountDtoFromModel(&m), nil
}

func (r *repository) FindByCustomerID(ctx context.Context, customerID string) ([]entity.AccountDto, error) {
	var accounts []model.Account
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	result := make([]entity.AccountDto, len(accounts))
	for i, m := range accounts {
		result[i] = *entity.NewAccountDtoFromModel(&m)
	}
	return result, nil
}

func (r *repository) FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.AccountDto], error) {
	info, paginationResult, err := pagination.NewTable[*entity.AccountFilterDto, *model.Account, entity.AccountDto]().
		Paginate(ctx, func(req *entity.AccountFilterDto) *gorm.DB {
			query := r.db.WithContext(ctx).Model(&model.Account{}).Preload("Customer")
			return query
		}, &pagination.TableRequest[*entity.AccountFilterDto]{
			Request:       req.(*entity.AccountFilterDto),
			QueryField:    []string{"account_name"},
			AllowedFields: []string{"customer_id", "account_type"},
		})
	if err != nil {
		return nil, err
	}
	result := make([]entity.AccountDto, len(paginationResult))
	for i, m := range paginationResult {
		result[i] = *entity.NewAccountDtoFromModel(m)
	}
	info.Data = result
	return info, nil
}

func (r *repository) Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	err := r.db.WithContext(ctx).Save(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Account{}).Error
	return err
}
