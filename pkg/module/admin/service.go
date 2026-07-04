package admin

import (
	"context"
	"time"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
)

// ServiceName is the lookup key for the admin service in the DI container.
const ServiceName = "adminService"

// Service exposes platform-wide operations for internal admins only.
//
// `GetStats` is intentionally simple — counters against the customer /
// organization / user tables. Kept separate from the customer module so
// internal-admin endpoints can grow without polluting the customer API.
type Service interface {
	GetProfile(ctx context.Context, session *entity.UserSessionDto) *entity.AdminProfileDto
	GetStats(ctx context.Context) (*entity.AdminStatsDto, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetProfile maps the existing JWT session onto the public AdminProfileDto
// shape. Returned by GET /api/admin/me so the FE can render the sidebar /
// topbar without a separate session call.
func (s *service) GetProfile(_ context.Context, session *entity.UserSessionDto) *entity.AdminProfileDto {
	if session == nil || session.UserCredential == nil || session.UserCredential.User == nil {
		return nil
	}
	cred := session.UserCredential
	return &entity.AdminProfileDto{
		ID:             cred.ID,
		UserID:         cred.UserID,
		OrganizationID: cred.OrganizationID,
		Username:       cred.Username,
		UserType:       string(cred.User.UserType),
	}
}

func (s *service) GetStats(ctx context.Context) (*entity.AdminStatsDto, error) {
	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.GeneratedAt = time.Now()
	return stats, nil
}
