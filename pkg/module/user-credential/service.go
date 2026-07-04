package usercredential

import (
	"context"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/security"
)

const ServiceName = "userCredentialService"

type Service interface {
	FindByUsername(ctx context.Context, userCredential *entity.UserCredentialDto) (*entity.UserCredentialDto, error)
	FindByEmail(ctx context.Context, userCredential *entity.UserCredentialDto) (*entity.UserCredentialDto, error)
	GetUserCredentialByUsername(ctx context.Context, username string) (*model.UserCredential, error)
	ChangePassword(ctx context.Context, userID, password string) error

	Create(ctx context.Context, dto *entity.UserCredentialDto) (*entity.UserCredentialDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserCredentialDto, error)
	FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.UserCredentialDto], error)
	Update(ctx context.Context, dto *entity.UserCredentialDto) (*entity.UserCredentialDto, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return &service{
		repo: repository,
	}
}

// FindByUsername is a function to find user credential by username
func (s *service) FindByUsername(ctx context.Context, userCredential *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return s.repo.FindByUsername(ctx, userCredential)
}

// FindByEmail is a function to find user credential by username
func (s *service) FindByEmail(ctx context.Context, userCredential *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return s.repo.FindByEmail(ctx, userCredential)
}

func (s *service) GetUserCredentialByUsername(ctx context.Context, username string) (*model.UserCredential, error) {
	return s.repo.GetUserCredentialByUsername(ctx, username)
}

func (s *service) ChangePassword(ctx context.Context, userID, password string) error {
	encryptedPass, _ := security.HashPassword(password)
	return s.repo.ChangePassword(ctx, userID, encryptedPass)
}

func (s *service) Create(ctx context.Context, dto *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.UserCredentialDto, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.UserCredentialDto], error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) Update(ctx context.Context, dto *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
