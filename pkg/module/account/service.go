package account

import (
	"context"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/transaction"
)

const ServiceName = "accountService"

type Service interface {
	Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	FindByID(ctx context.Context, id string) (*entity.AccountDto, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]entity.AccountDto, error)
	FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.AccountDto], error)
	Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	txManager transaction.Manager
	repo      Repository
}

func NewService(
	txManager transaction.Manager,
	repo Repository,
) Service {
	return &service{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *service) Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	var accountDtoResult *entity.AccountDto
	err := s.txManager.Execute(ctx, func(txCtx context.Context) error {
		res, err := s.repo.Create(txCtx, dto)
		if err != nil {
			return err
		}
		accountDtoResult = res
		return nil
	})
	return accountDtoResult, err
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.AccountDto, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindByCustomerID(ctx context.Context, customerID string) ([]entity.AccountDto, error) {
	return s.repo.FindByCustomerID(ctx, customerID)
}

func (s *service) FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.AccountDto], error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	_, err := s.FindByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	err = s.txManager.Execute(ctx, func(txCtx context.Context) error {
		_, err := s.repo.Update(txCtx, dto)
		return err
	})
	return dto, err
}

func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.txManager.Execute(ctx, func(txCtx context.Context) error {
		return s.repo.Delete(txCtx, id)
	})
	return err
}
