package user

import (
	"context"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserDto, error)
	FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.UserDto], error)
	Update(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	dto.ID = m.ID
	return dto, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.UserDto, error) {
	var m entity.UserDto
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.UserDto], error) {
	info, paginationResult, err := pagination.NewTable[*entity.UserFilterDto, *model.User, entity.UserDto]().
		Paginate(ctx, func(req *entity.UserFilterDto) *gorm.DB {
			query := r.db.WithContext(ctx).Model(&model.User{})
			return query
		}, &pagination.TableRequest[*entity.UserFilterDto]{})
	if err != nil {
		return nil, err
	}
	result := make([]entity.UserDto, len(paginationResult))
	for i, m := range paginationResult {
		result[i] = new(entity.UserDto).FromModel(m)
	}
	info.Data = result
	return info, nil
}

func (r *repository) Update(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error) { // Implementation of updating a funder in the database
	return nil, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	// Implementation of deleting a funder from the database
	return nil
}

// // FindByReferralCode is a function to find user by referral code
// func (r *repository) FindByReferralCode(referralCode string) (*entity.CreateUser, error) {
// 	var user model.User
// 	if err := r.db.Joins("Customer").
// 		Where("customer.referral_code = ?", referralCode).
// 		First(&user).Error; err != nil {
// 		return nil, err
// 	}
// 	return &entity.CreateUser{}, nil
// }

// // CreateUser is a function to create user
// func (r *repository) CreateUser(createUser *entity.CreateUser) (*entity.CreateUser, error) {
// 	user := new(model.User)
// 	user.OrganizationID = createUser.OrganizationData.ID

// 	password, _ := utils.HashPassword(createUser.Password)
// 	user.UserCredential = []model.UserCredential{
// 		{
// 			OrganizationID: createUser.OrganizationData.ID,
// 			Username:       createUser.Username,
// 			Password:       password,
// 		},
// 	}

// 	if err := r.db.Create(user).Error; err != nil {
// 		return nil, err
// 	}
// 	createUser.UserID = user.ID
// 	return createUser, nil
// }
