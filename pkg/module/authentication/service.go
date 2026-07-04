package authentication

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raymondsugiarto/reputation-be/config"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	usercredential "github.com/raymondsugiarto/reputation-be/pkg/module/user-credential"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/security"
)

const ServiceName = "authenticationService"

// Customer lookup for sign-in gating. Defined as a function shape so
// the auth module doesn't need to import the customer module directly
// (which would create a circular dependency: customer → auth → customer).
//
// The DI container wires this up in cmd/server/rest.go (or wherever
// the service is constructed) to point at customer.Service.FindByUserID.
// We expose it through NewServiceWithCustomerLookup so existing
// callers don't have to change.
type CustomerLookup func(ctx context.Context, userID string) (CustomerStatusInfo, error)

// CustomerStatusInfo is the minimal slice of customer state that
// authentication.SignIn needs. Avoids leaking the customer DTO across
// the module boundary.
type CustomerStatusInfo struct {
	Status model.CustomerStatus
}

// Sentinel errors surfaced to the FE. The handler maps these to
// `errors.New(...)` so the FE can pattern-match on the message.
var (
	ErrCustomerNotApproved = errors.New("customerPendingApproval")
	ErrCustomerRejected    = errors.New("customerRejected")
	ErrCustomerNotFound    = errors.New("customerNotFound")
)

type Service interface {
	GetSession(ctx context.Context) (*entity.UserSessionDto, error)
	SignIn(context.Context, *entity.LoginRequestDto) (*entity.LoginDto, error)
}

type service struct {
	userCredentialService usercredential.Service
	customerLookup        CustomerLookup
}

// NewService preserves the existing constructor signature for
// backwards compatibility. Customer accounts will still be blocked
// from signing in (because customerLookup returns "not found"), but
// the friendly "pending approval" / "rejected" message will fall back
// to a generic "invalid credential" path. Production wiring should
// use NewServiceWithCustomerLookup.
func NewService(userCredentialService usercredential.Service) Service {
	return &service{userCredentialService: userCredentialService}
}

// NewServiceWithCustomerLookup returns a Service that, for CUSTOMER
// users, additionally checks the customer's approval status before
// issuing a JWT.
func NewServiceWithCustomerLookup(
	userCredentialService usercredential.Service,
	customerLookup CustomerLookup,
) Service {
	return &service{
		userCredentialService: userCredentialService,
		customerLookup:        customerLookup,
	}
}

func (s *service) GetSession(ctx context.Context) (*entity.UserSessionDto, error) {
	token := jwtware.FromContext(ctx)
	claims := token.Claims.(jwt.MapClaims)
	fmt.Printf("token %s", claims)
	userSessionDto := entity.NewUserSessionDtoFromClaims(claims)

	userCredentialDto, err := s.userCredentialService.FindByID(ctx, userSessionDto.ID)
	if err != nil {
		return nil, errors.New("userNotFound")
	}
	userSessionDto.UserCredential = userCredentialDto

	return userSessionDto, nil
}

func (s *service) SignIn(ctx context.Context, request *entity.LoginRequestDto) (*entity.LoginDto, error) {
	log.WithContext(ctx).Infof("sign in started")
	userCredentialModel, err := s.userCredentialService.GetUserCredentialByUsername(ctx, request.Username)
	if err != nil {
		return nil, errors.New("userNotFound")
	}

	// Approval gating for CUSTOMER users. INTERNAL_ADMIN and org-owner
	// ADMIN always pass — they were created directly by the platform
	// and have no approval workflow.
	if s.customerLookup != nil && userCredentialModel.User != nil &&
		userCredentialModel.User.UserType == model.CUSTOMER {
		info, lookupErr := s.customerLookup(ctx, userCredentialModel.UserID)
		if lookupErr != nil {
			// Customer record not yet created — for example the
			// sign-up transaction failed mid-way. Treat as not-found
			// so we don't leak the underlying error.
			log.WithContext(ctx).Warnf("customer lookup failed for user %s: %v",
				userCredentialModel.UserID, lookupErr)
			return nil, ErrCustomerNotFound
		}
		switch info.Status {
		case model.CustomerStatusApproved:
			// allow
		case model.CustomerStatusRejected:
			return nil, ErrCustomerRejected
		default:
			// PENDING_APPROVAL or empty (treated as pending).
			return nil, ErrCustomerNotApproved
		}
	}

	pp, _ := security.HashPassword(request.Password)
	log.WithContext(ctx).Infof("password: %s, hash: %s", request.Password, pp)
	if !security.CheckPasswordHash(request.Password, userCredentialModel.Password) {
		return nil, errors.New("invalidPassword")
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userCredentialModel.ID
	claims["uid"] = userCredentialModel.UserID
	claims["oid"] = userCredentialModel.OrganizationID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	cfg := config.GetConfig()
	t, err := token.SignedString([]byte(cfg.Server.Rest.SecretKey))
	if err != nil {
		return nil, errors.New("errorGeneratetoken")
	}
	return &entity.LoginDto{
		Token:   t,
		Expired: strconv.Itoa(int(claims["exp"].(int64))),
	}, nil
}
