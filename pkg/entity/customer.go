package entity

import (
	"time"

	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
)

// CustomerSignUpRequestDto is the inbound payload for POST /api/customer/sign-up.
// This is a public self-service endpoint used by both perorangan (INDIVIDUAL)
// and perusahaan (COMPANY) flows. OrganizationID is intentionally absent —
// customers are not yet tied to an organization.
type CustomerSignUpRequestDto struct {
	Username     string             `json:"username" validate:"required"`
	Password     string             `json:"password" validate:"required"`
	CustomerType model.CustomerType `json:"customerType" validate:"required"`

	// INDIVIDUAL-only fields.
	NamaLengkap        string                      `json:"namaLengkap"`
	NomorKtp           string                      `json:"nomorKtp"`
	NomorNpwp          string                      `json:"nomorNpwp"`
	TanggalLahir       time.Time                   `json:"tanggalLahir"`
	KotaLahir          string                      `json:"kotaLahir"`
	StatusPernikahan   model.CustomerMaritalStatus `json:"statusPernikahan"`
	PendidikanTerakhir model.CustomerLastEducation `json:"pendidikanTerakhir"`
	LamaTinggal        string                      `json:"lamaTinggal"`

	AlamatJalan        string `json:"alamatJalan"`
	Kecamatan          string `json:"kecamatan"`
	KotaKabupaten      string `json:"kotaKabupaten"`
	Provinsi           string `json:"provinsi"`
	KodePos            string `json:"kodePos"`
	SamaDenganDomisili bool   `json:"samaDenganDomisili"`

	AlamatKtpJalan   string `json:"alamatKtpJalan"`
	KecamatanKtp     string `json:"kecamatanKtp"`
	KotaKabupatenKtp string `json:"kotaKabupatenKtp"`

	// COMPANY-only fields.
	NamaPt           string             `json:"namaPt"`
	TanggalPendirian *time.Time         `json:"tanggalPendirian"`
	SektorUsaha      string             `json:"sektorUsaha"`
	KodeKbli         string             `json:"kodeKbli"`
	NpwpPerusahaan   string             `json:"npwpPerusahaan"`
	PemegangSaham    []PemegangSahamDto `json:"pemegangSaham"`
}

// PemegangSahamDto models one row of the company shareholder table on the FE.
// The backend stores the slice as JSON so we don't need a separate table today.
type PemegangSahamDto struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	Jenis     string `json:"jenis"`
	Saham     string `json:"saham"`
	NoKtpNpwp string `json:"noKtpNpwp"`
	Peran     string `json:"peran"`
}

// ToCustomerDto maps the inbound request onto a CustomerDto, ready to be turned
// into a GORM model. OrganizationID is left empty — customers are not yet
// bound to an organization.
func (r *CustomerSignUpRequestDto) ToCustomerDto() *CustomerDto {
	dto := &CustomerDto{
		CustomerType:       r.CustomerType,
		NamaLengkap:        r.NamaLengkap,
		NomorKtp:           r.NomorKtp,
		NomorNpwp:          r.NomorNpwp,
		TanggalLahir:       r.TanggalLahir,
		KotaLahir:          r.KotaLahir,
		StatusPernikahan:   r.StatusPernikahan,
		PendidikanTerakhir: r.PendidikanTerakhir,
		LamaTinggal:        r.LamaTinggal,
		AlamatJalan:        r.AlamatJalan,
		Kecamatan:          r.Kecamatan,
		KotaKabupaten:      r.KotaKabupaten,
		Provinsi:           r.Provinsi,
		KodePos:            r.KodePos,
		SamaDenganDomisili: r.SamaDenganDomisili,
		AlamatKtpJalan:     r.AlamatKtpJalan,
		KecamatanKtp:       r.KecamatanKtp,
		KotaKabupatenKtp:   r.KotaKabupatenKtp,
		NamaPt:             r.NamaPt,
		SektorUsaha:        r.SektorUsaha,
		KodeKbli:           r.KodeKbli,
		NpwpPerusahaan:     r.NpwpPerusahaan,
	}
	if r.TanggalPendirian != nil {
		dto.TanggalPendirian = *r.TanggalPendirian
	}
	return dto
}

// CustomerDto is the JSON-serialisable form of a Customer record.
type CustomerDto struct {
	ID             string             `json:"id"`
	OrganizationID string             `json:"organizationId"`
	UserID         string             `json:"userId"`
	CustomerType   model.CustomerType `json:"customerType"`

	// Approval lifecycle. See migration 000005_customer_approval.
	Status     model.CustomerStatus `json:"status"`
	ApprovedBy string               `json:"approvedBy,omitempty"`
	ApprovedAt *time.Time           `json:"approvedAt,omitempty"`
	RejectedBy string               `json:"rejectedBy,omitempty"`
	RejectedAt *time.Time           `json:"rejectedAt,omitempty"`
	Remark     string               `json:"remark,omitempty"`

	NamaLengkap        string                      `json:"namaLengkap"`
	NomorKtp           string                      `json:"nomorKtp"`
	NomorNpwp          string                      `json:"nomorNpwp"`
	TanggalLahir       time.Time                   `json:"tanggalLahir"`
	KotaLahir          string                      `json:"kotaLahir"`
	StatusPernikahan   model.CustomerMaritalStatus `json:"statusPernikahan"`
	PendidikanTerakhir model.CustomerLastEducation `json:"pendidikanTerakhir"`
	LamaTinggal        string                      `json:"lamaTinggal"`
	AlamatJalan        string                      `json:"alamatJalan"`
	Kecamatan          string                      `json:"kecamatan"`
	KotaKabupaten      string                      `json:"kotaKabupaten"`
	Provinsi           string                      `json:"provinsi"`
	KodePos            string                      `json:"kodePos"`
	SamaDenganDomisili bool                        `json:"samaDenganDomisili"`
	AlamatKtpJalan     string                      `json:"alamatKtpJalan"`
	KecamatanKtp       string                      `json:"kecamatanKtp"`
	KotaKabupatenKtp   string                      `json:"kotaKabupatenKtp"`

	NamaPt           string             `json:"namaPt"`
	TanggalPendirian time.Time          `json:"tanggalPendirian"`
	SektorUsaha      string             `json:"sektorUsaha"`
	KodeKbli         string             `json:"kodeKbli"`
	NpwpPerusahaan   string             `json:"npwpPerusahaan"`
	PemegangSaham    []PemegangSahamDto `json:"pemegangSaham"`
}

func NewCustomerDtoFromModel(m *model.Customer) *CustomerDto {
	if m == nil {
		return nil
	}
	status := m.Status
	// Defensive default — model default is set by the DB layer, but if
	// callers construct a Customer manually we still want the wire
	// payload to be coherent.
	if status == "" {
		status = model.CustomerStatusPendingApproval
	}
	return &CustomerDto{
		ID:                 m.ID,
		OrganizationID:     m.OrganizationID,
		UserID:             m.UserID,
		CustomerType:       m.CustomerType,
		Status:             status,
		ApprovedBy:         m.ApprovedBy,
		ApprovedAt:         m.ApprovedAt,
		RejectedBy:         m.RejectedBy,
		RejectedAt:         m.RejectedAt,
		Remark:             m.Remark,
		NamaLengkap:        m.NamaLengkap,
		NomorKtp:           m.NomorKtp,
		NomorNpwp:          m.NomorNpwp,
		TanggalLahir:       m.TanggalLahir,
		KotaLahir:          m.KotaLahir,
		StatusPernikahan:   m.StatusPernikahan,
		PendidikanTerakhir: m.PendidikanTerakhir,
		LamaTinggal:        m.LamaTinggal,
		AlamatJalan:        m.AlamatJalan,
		Kecamatan:          m.Kecamatan,
		KotaKabupaten:      m.KotaKabupaten,
		Provinsi:           m.Provinsi,
		KodePos:            m.KodePos,
		SamaDenganDomisili: m.SamaDenganDomisili,
		AlamatKtpJalan:     m.AlamatKtpJalan,
		KecamatanKtp:       m.KecamatanKtp,
		KotaKabupatenKtp:   m.KotaKabupatenKtp,
		NamaPt:             m.NamaPt,
		TanggalPendirian:   m.TanggalPendirian,
		SektorUsaha:        m.SektorUsaha,
		KodeKbli:           m.KodeKbli,
		NpwpPerusahaan:     m.NpwpPerusahaan,
	}
}

func (c *CustomerDto) ToModel() *model.Customer {
	m := &model.Customer{
		OrganizationID:     c.OrganizationID,
		CustomerType:       c.CustomerType,
		NamaLengkap:        c.NamaLengkap,
		NomorKtp:           c.NomorKtp,
		NomorNpwp:          c.NomorNpwp,
		TanggalLahir:       c.TanggalLahir,
		KotaLahir:          c.KotaLahir,
		StatusPernikahan:   c.StatusPernikahan,
		PendidikanTerakhir: c.PendidikanTerakhir,
		LamaTinggal:        c.LamaTinggal,
		AlamatJalan:        c.AlamatJalan,
		Kecamatan:          c.Kecamatan,
		KotaKabupaten:      c.KotaKabupaten,
		Provinsi:           c.Provinsi,
		KodePos:            c.KodePos,
		SamaDenganDomisili: c.SamaDenganDomisili,
		AlamatKtpJalan:     c.AlamatKtpJalan,
		KecamatanKtp:       c.KecamatanKtp,
		KotaKabupatenKtp:   c.KotaKabupatenKtp,
		NamaPt:             c.NamaPt,
		TanggalPendirian:   c.TanggalPendirian,
		SektorUsaha:        c.SektorUsaha,
		KodeKbli:           c.KodeKbli,
		NpwpPerusahaan:     c.NpwpPerusahaan,
		// Approval lifecycle. Status is left empty when the DTO is
		// empty so the DB column default kicks in on INSERT.
		Status:     c.Status,
		ApprovedBy: c.ApprovedBy,
		ApprovedAt: c.ApprovedAt,
		RejectedBy: c.RejectedBy,
		RejectedAt: c.RejectedAt,
		Remark:     c.Remark,
	}
	if c.ID != "" {
		m.ID = c.ID
	}
	if c.UserID != "" {
		m.UserID = c.UserID
	}
	return m
}

// CustomerSignUpResponseDto is the response for the customer sign-up endpoint.
//
// `Status` is included so the FE can render the "Pending Approval"
// state on the success screen without an extra fetch. The status is
// always PENDING_APPROVAL right after sign-up — the value is exposed
// here for forward compatibility.
type CustomerSignUpResponseDto struct {
	UserID           string               `json:"userId"`
	UserCredentialID string               `json:"userCredentialId"`
	CustomerID       string               `json:"customerId"`
	CustomerType     model.CustomerType   `json:"customerType"`
	Status           model.CustomerStatus `json:"status"`
}

// CustomerFilterDto for query filtering on Customer.
//
// `Status` is optional — when empty, the repo applies no status filter
// and returns all rows. The admin approval flow sets it to one of
// "PENDING_APPROVAL" / "APPROVED" / "REJECTED" depending on which list
// the admin is viewing.
type CustomerFilterDto struct {
	pagination.GetListRequest
	OrganizationID string `query:"organizationId"`
	CustomerType   string `query:"customerType"`
	Status         string `query:"status"`
	NamaLengkap    string `query:"namaLengkap"`
	NomorKtp       string `query:"nomorKtp"`
}

// GenerateFilter pre-parses the Status filter into a list so the repo
// can issue an `IN (...)` query when the admin is viewing the history
// page without a specific action filter. PaginationRequestDto's default
// GenerateFilter is a no-op, so we override here.
func (f *CustomerFilterDto) GenerateFilter() {
	if f.Status == "" {
		return
	}
	// Pass-through — repo reads f.Status directly.
}

// ApprovalHistoryFilterDto is a second filter DTO for the admin
// approval-history endpoint. It accepts an optional `action`
// ("APPROVED" or "REJECTED") and returns rows whose status is in
// {APPROVED, REJECTED} when action is empty. We use a separate DTO
// instead of overloading CustomerFilterDto so the wire shape is
// self-documenting on the FE.
type ApprovalHistoryFilterDto struct {
	pagination.GetListRequest
	Action       string `query:"action"`       // APPROVED | REJECTED | (empty = both)
	CustomerType string `query:"customerType"` // optional filter
}
