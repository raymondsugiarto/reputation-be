package authentication_test

import (
	"context"
	"errors"
	"testing"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/module/authentication"
	usercredential "github.com/raymondsugiarto/reputation-be/pkg/module/user-credential"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
)

// fakeUserCredService satisfies usercredential.Service for the SignIn
// gating tests. Only the methods exercised by SignIn are non-nil; the
// rest are stubs that keep the file self-contained.
type fakeUserCredService struct {
	usernameToModel map[string]*model.UserCredential
}

func (f *fakeUserCredService) GetUserCredentialByUsername(_ context.Context, username string) (*model.UserCredential, error) {
	m, ok := f.usernameToModel[username]
	if !ok {
		return nil, errors.New("userNotFound")
	}
	return m, nil
}

// Stubs — SignIn doesn't call these.
func (f *fakeUserCredService) FindByUsername(_ context.Context, _ *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return nil, nil
}
func (f *fakeUserCredService) FindByEmail(_ context.Context, _ *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return nil, nil
}
func (f *fakeUserCredService) ChangePassword(_ context.Context, _, _ string) error { return nil }
func (f *fakeUserCredService) Create(_ context.Context, _ *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return nil, nil
}
func (f *fakeUserCredService) FindByID(_ context.Context, _ string) (*entity.UserCredentialDto, error) {
	return nil, nil
}
func (f *fakeUserCredService) FindAll(_ context.Context, _ pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.UserCredentialDto], error) {
	return nil, nil
}
func (f *fakeUserCredService) Update(_ context.Context, _ *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return nil, nil
}
func (f *fakeUserCredService) Delete(_ context.Context, _ string) error { return nil }

// Compile-time interface guard.
var _ usercredential.Service = (*fakeUserCredService)(nil)

// bcrypt hash of "secret123" — generated once via bcrypt.GenerateFromPassword
// with cost 7 to match the production hashing parameters.
const passwordHash = "$2a$07$jFUQR8gMKK47XHtRbLO1cuJ1Slm4MmyHpi57gK/esBLYz2bg9agRa"

func newCred(username, userID string, userType model.UserType) *model.UserCredential {
	return &model.UserCredential{
		OrganizationID: "org-1",
		UserID:         userID,
		Username:       username,
		Password:       passwordHash,
		User: &model.User{
			OrganizationID: "org-1",
			UserType:       userType,
		},
	}
}

// ---------------------------------------------------------------------------
// SignIn approval gating
// ---------------------------------------------------------------------------

func TestSignIn_PendingCustomerIsBlocked(t *testing.T) {
	svc := authentication.NewServiceWithCustomerLookup(
		&fakeUserCredService{usernameToModel: map[string]*model.UserCredential{
			"alice": newCred("alice", "u-1", model.CUSTOMER),
		}},
		func(_ context.Context, _ string) (authentication.CustomerStatusInfo, error) {
			return authentication.CustomerStatusInfo{Status: model.CustomerStatusPendingApproval}, nil
		},
	)

	_, err := svc.SignIn(context.Background(), &entity.LoginRequestDto{
		Username: "alice",
		Password: "secret123",
	})
	if !errors.Is(err, authentication.ErrCustomerNotApproved) {
		t.Fatalf("expected ErrCustomerNotApproved, got %v", err)
	}
}

func TestSignIn_RejectedCustomerIsBlocked(t *testing.T) {
	svc := authentication.NewServiceWithCustomerLookup(
		&fakeUserCredService{usernameToModel: map[string]*model.UserCredential{
			"alice": newCred("alice", "u-1", model.CUSTOMER),
		}},
		func(_ context.Context, _ string) (authentication.CustomerStatusInfo, error) {
			return authentication.CustomerStatusInfo{Status: model.CustomerStatusRejected}, nil
		},
	)

	_, err := svc.SignIn(context.Background(), &entity.LoginRequestDto{
		Username: "alice",
		Password: "secret123",
	})
	if !errors.Is(err, authentication.ErrCustomerRejected) {
		t.Fatalf("expected ErrCustomerRejected, got %v", err)
	}
}

func TestSignIn_ApprovedCustomerSucceeds(t *testing.T) {
	svc := authentication.NewServiceWithCustomerLookup(
		&fakeUserCredService{usernameToModel: map[string]*model.UserCredential{
			"alice": newCred("alice", "u-1", model.CUSTOMER),
		}},
		func(_ context.Context, _ string) (authentication.CustomerStatusInfo, error) {
			return authentication.CustomerStatusInfo{Status: model.CustomerStatusApproved}, nil
		},
	)

	result, err := svc.SignIn(context.Background(), &entity.LoginRequestDto{
		Username: "alice",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected success for APPROVED customer, got %v", err)
	}
	if result == nil || result.Token == "" {
		t.Fatalf("expected non-empty token, got %+v", result)
	}
}

func TestSignIn_AdminIsNotGatedByCustomerLookup(t *testing.T) {
	// Admin users (ADMIN or INTERNAL_ADMIN) must not trigger the
	// customerLookup gating. If they did, the spy would call t.Fatalf.
	for _, userType := range []model.UserType{model.ADMIN, model.INTERNAL_ADMIN} {
		svc := authentication.NewServiceWithCustomerLookup(
			&fakeUserCredService{usernameToModel: map[string]*model.UserCredential{
				"boss": newCred("boss", "u-1", userType),
			}},
			func(_ context.Context, _ string) (authentication.CustomerStatusInfo, error) {
				t.Fatalf("customerLookup should not be called for %s", userType)
				return authentication.CustomerStatusInfo{}, nil
			},
		)

		result, err := svc.SignIn(context.Background(), &entity.LoginRequestDto{
			Username: "boss",
			Password: "secret123",
		})
		if err != nil {
			t.Fatalf("expected success for %s, got %v", userType, err)
		}
		if result == nil || result.Token == "" {
			t.Fatalf("expected non-empty token for %s", userType)
		}
	}
}
