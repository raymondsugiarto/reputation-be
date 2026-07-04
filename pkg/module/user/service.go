package user

import (
	"context"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	usercredential "github.com/raymondsugiarto/reputation-be/pkg/module/user-credential"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
)

const ServiceName = "userService"

type Service interface {
	Create(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserDto, error)
	FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.UserDto], error)
	Update(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo                  Repository
	userCredentialService usercredential.Service
}

func NewService(repository Repository, userCredentialService usercredential.Service) Service {
	return &service{
		repo:                  repository,
		userCredentialService: userCredentialService,
	}
}

func (s *service) Create(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error) {
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.UserDto, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.UserDto], error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) Update(ctx context.Context, dto *entity.UserDto) (*entity.UserDto, error) {
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// // FindByUserID is a function to find user by user id
// func (s *service) FindByUserID(ctx context.Context, userID string) (*entity.UserDto, error) {
// 	userDto, err := s.repo.FindByID(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return userDto, nil
// }

// CreateUser is a function to create user
// func (s *service) CreateUser(createUser *entity.CreateUser) (*entity.CreateUser, error) {
// 	userCredential := &entity.UserCredential{
// 		OrganizationData: createUser.OrganizationData,
// 		Username:         createUser.Username,
// 	}
// 	_, err := s.userCredentialService.FindByUsername(userCredential)
// 	if err == nil {
// 		return nil, errors.New("errorAccountCodeAlreadyExist")
// 	}

// 	_, err = s.userCredentialService.FindByEmail(userCredential)
// 	if err == nil {
// 		return nil, errors.New("errorEmailAlreadyExist")
// 	}

// 	return s.repository.CreateUser(createUser)
// }
