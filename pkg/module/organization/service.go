package organization

import (
	"context"
	"errors"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/module/user"
	usercredential "github.com/raymondsugiarto/reputation-be/pkg/module/user-credential"
)

const ServiceName = "organizationService"

type Service interface {
	SignUp(ctx context.Context, dto *entity.SignUpRequestDto) (*entity.SignUpResponseDto, error)
}

type service struct {
	organizationRepo      Repository
	userService           user.Service
	userCredentialService usercredential.Service
}

func NewService(
	organizationRepo Repository,
	userService user.Service,
	userCredentialService usercredential.Service,
) Service {
	return &service{
		organizationRepo:      organizationRepo,
		userService:           userService,
		userCredentialService: userCredentialService,
	}
}

func (s *service) SignUp(ctx context.Context, dto *entity.SignUpRequestDto) (*entity.SignUpResponseDto, error) {
	organization := &model.Organization{
		Code:   dto.Name,
		Name:   dto.Name,
		Origin: "",
	}
	createdOrg, err := s.organizationRepo.Create(ctx, organization)
	if err != nil {
		return nil, errors.New("errorCreatingOrganization")
	}

	userDto := &entity.UserDto{
		OrganizationID: createdOrg.ID,
		UserType:       model.ADMIN,
	}
	createdUser, err := s.userService.Create(ctx, userDto)
	if err != nil {
		return nil, errors.New("errorCreatingUser")
	}

	userCredDto := &entity.UserCredentialDto{
		OrganizationID: createdOrg.ID,
		UserID:         createdUser.ID,
		Username:       dto.Username,
		Password:       dto.Password,
	}
	_, err = s.userCredentialService.Create(ctx, userCredDto)
	if err != nil {
		return nil, errors.New("errorCreatingUserCredential")
	}

	return &entity.SignUpResponseDto{
		OrganizationID: createdOrg.ID,
		UserID:         createdUser.ID,
	}, nil
}
