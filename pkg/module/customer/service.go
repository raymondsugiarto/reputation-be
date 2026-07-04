package customer

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/module/user"
	usercredential "github.com/raymondsugiarto/reputation-be/pkg/module/user-credential"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/transaction"
)

// ErrCustomerNotFound is the sentinel returned by FindAll when a search
// yields zero rows. The handler maps this to HTTP 404 so the FE can
// render its dedicated "result not found" UI instead of treating an
// empty list as a server error.
var ErrCustomerNotFound = errors.New("customerNotFound")

const ServiceName = "customerService"

type Service interface {
	SignUp(ctx context.Context, dto *entity.CustomerSignUpRequestDto) (*entity.CustomerSignUpResponseDto, error)
	Create(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error)
	FindByID(ctx context.Context, id string) (*entity.CustomerDto, error)
	FindByUserID(ctx context.Context, userID string) (*entity.CustomerDto, error)
	FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error)
	Update(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error)
	Delete(ctx context.Context, id string) error

	// Approval workflow. Implemented in approval.go.
	Approve(ctx context.Context, customerID, adminUserID, remark string) (*ApprovalResultDto, error)
	Reject(ctx context.Context, customerID, adminUserID, remark string) (*ApprovalResultDto, error)
	FindPendingApprovals(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error)
	FindApprovalHistory(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error)
	GetApprovalStats(ctx context.Context) (*entity.CustomerApprovalStatsDto, error)
}

type service struct {
	txManager             transaction.Manager
	repo                  Repository
	userService           user.Service
	userCredentialService usercredential.Service
}

func NewService(
	txManager transaction.Manager,
	repo Repository,
	userService user.Service,
	userCredentialService usercredential.Service,
) Service {
	return &service{
		txManager:             txManager,
		repo:                  repo,
		userService:           userService,
		userCredentialService: userCredentialService,
	}
}

// SignUp is the public self-service endpoint for onboarding a customer.
// Customers are not tied to an organization yet, so OrganizationID is left
// empty on every record created here. CustomerType (INDIVIDUAL or COMPANY)
// is required. The flow is:
//  1. Create a UserCredential (username/password).
//  2. Create a User of type CUSTOMER.
//  3. Create the Customer profile pointing at the User.
//
// All three are wrapped in a single transaction.
func (s *service) SignUp(ctx context.Context, dto *entity.CustomerSignUpRequestDto) (*entity.CustomerSignUpResponseDto, error) {
	if dto == nil {
		return nil, errors.New("errorInvalidRequest")
	}
	if dto.CustomerType == "" {
		return nil, errors.New("errorMissingCustomerType")
	}

	var (
		userCred *entity.UserCredentialDto
		userDto  *entity.UserDto
		cust     *entity.CustomerDto
	)

	err := s.txManager.Execute(ctx, func(txCtx context.Context) error {
		cred, err := s.userCredentialService.Create(txCtx, &entity.UserCredentialDto{
			Username: dto.Username,
			Password: dto.Password,
		})
		if err != nil {
			return errors.New("errorCreatingUserCredential")
		}
		userCred = cred

		createdUser, err := s.userService.Create(txCtx, &entity.UserDto{
			UserType: model.CUSTOMER,
		})
		if err != nil {
			return errors.New("errorCreatingUser")
		}
		userDto = createdUser

		custDto := dto.ToCustomerDto()
		custDto.UserID = createdUser.ID

		createdCust, err := s.repo.Create(txCtx, custDto)
		if err != nil {
			log.WithContext(ctx).Errorf("error creating customer", err)
			return errors.New("errorCreatingCustomer")
		}
		cust = createdCust
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &entity.CustomerSignUpResponseDto{
		UserID:           userDto.ID,
		UserCredentialID: userCred.ID,
		CustomerID:       cust.ID,
		CustomerType:     cust.CustomerType,
		Status:           model.CustomerStatusPendingApproval,
	}, nil
}

func (s *service) Create(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error) {
	var created *entity.CustomerDto
	err := s.txManager.Execute(ctx, func(txCtx context.Context) error {
		res, err := s.repo.Create(txCtx, dto)
		if err != nil {
			return err
		}
		created = res
		return nil
	})
	return created, err
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.CustomerDto, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindByUserID(ctx context.Context, userID string) (*entity.CustomerDto, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *service) FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	result, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}
	// The FE treats 404 from this endpoint as "no results" and shows the
	// result-not-found illustration. We only 404 when the caller is
	// actually searching (query != "") — an empty list when no filter is
	// given is a legitimate "no data" response, not a not-found error.
	if result.Count == 0 && req.GetQuery() != "" {
		return nil, ErrCustomerNotFound
	}
	return result, nil
}

func (s *service) Update(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error) {
	if _, err := s.FindByID(ctx, dto.ID); err != nil {
		return nil, err
	}
	err := s.txManager.Execute(ctx, func(txCtx context.Context) error {
		_, err := s.repo.Update(txCtx, dto)
		return err
	})
	return dto, err
}

func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := s.FindByID(ctx, id); err != nil {
		return err
	}
	return s.txManager.Execute(ctx, func(txCtx context.Context) error {
		return s.repo.Delete(txCtx, id)
	})
}
