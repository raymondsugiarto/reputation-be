package model

import (
	"time"

	concern "github.com/raymondsugiarto/reputation-be/pkg/model/common"
)

type CustomerMaritalStatus string

const (
	CustomerMaritalSingle    CustomerMaritalStatus = "SINGLE"
	CustomerMaritalMarried   CustomerMaritalStatus = "MARRIED"
	CustomerMaritalDivorced  CustomerMaritalStatus = "DIVORCED"
	CustomerMaritalWidowed   CustomerMaritalStatus = "WIDOWED"
	CustomerMaritalSeparated CustomerMaritalStatus = "SEPARATED"
	CustomerMaritalOther     CustomerMaritalStatus = "OTHER"
)

type CustomerLastEducation string

const (
	CustomerEducationNone       CustomerLastEducation = "NONE"
	CustomerEducationElementary CustomerLastEducation = "ELEMENTARY"
	CustomerEducationJuniorHigh CustomerLastEducation = "JUNIOR_HIGH"
	CustomerEducationSeniorHigh CustomerLastEducation = "SENIOR_HIGH"
	CustomerEducationDiploma    CustomerLastEducation = "DIPLOMA"
	CustomerEducationBachelor   CustomerLastEducation = "BACHELOR"
	CustomerEducationMaster     CustomerLastEducation = "MASTER"
	CustomerEducationDoctorate  CustomerLastEducation = "DOCTORATE"
	CustomerEducationOther      CustomerLastEducation = "OTHER"
)

// CustomerType distinguishes the two flavours of customer: perorangan
// (individual / personal) and perusahaan (company / corporate). Stored as a
// plain varchar(50) — the values below are the only ones the application
// writes today.
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "INDIVIDUAL"
	CustomerTypeCompany    CustomerType = "COMPANY"
)

// CustomerStatus is the approval lifecycle for a customer account. All
// customer self-service sign-ups land in PENDING_APPROVAL and must be
// approved by an internal admin before the customer can sign in. The
// `rejected_at`/`remark` columns are populated when the customer is
// moved to REJECTED.
//
// Storage column: varchar(50). Stored as plain strings (no native enum)
// because the application treats them as constants, not a closed DB enum.
type CustomerStatus string

const (
	// CustomerStatusPendingApproval — initial state for every new
	// customer sign-up. Customer cannot sign in.
	CustomerStatusPendingApproval CustomerStatus = "PENDING_APPROVAL"
	// CustomerStatusApproved — internal admin approved the customer.
	// Customer can sign in.
	CustomerStatusApproved CustomerStatus = "APPROVED"
	// CustomerStatusRejected — internal admin rejected the customer.
	// Customer cannot sign in.
	CustomerStatusRejected CustomerStatus = "REJECTED"
)

type Customer struct {
	concern.CommonWithIDs
	OrganizationID string
	UserID         string
	User           *User

	CustomerType CustomerType

	NamaLengkap        string
	NomorKtp           string
	NomorNpwp          string
	TanggalLahir       time.Time
	KotaLahir          string
	StatusPernikahan   CustomerMaritalStatus
	PendidikanTerakhir CustomerLastEducation
	LamaTinggal        string

	AlamatJalan        string
	Kecamatan          string
	KotaKabupaten      string
	Provinsi           string
	KodePos            string
	SamaDenganDomisili bool

	AlamatKtpJalan   string
	KecamatanKtp     string
	KotaKabupatenKtp string

	// Company-only fields (populated when CustomerType == CustomerTypeCompany).
	// Column tags point at the English snake_case columns introduced in
	// migration 000004_customer_company_fields. The Go names stay Indonesian
	// so the rest of the codebase (entities, services, handlers) reads naturally.
	NamaPt           string    `gorm:"column:company_name"`
	TanggalPendirian time.Time `gorm:"column:establishment_date"`
	SektorUsaha      string    `gorm:"column:business_sector"`
	KodeKbli         string    `gorm:"column:kbli_code"`
	NpwpPerusahaan   string    `gorm:"column:company_tax_id"`

	// PemegangSaham is modelled as JSON so we don't need a separate table for
	// shareholder structure today. Empty slice for individual customers.
	// Stored as a jsonb column on the customer table.
	PemegangSaham string `gorm:"column:shareholders;type:jsonb"`

	// Approval lifecycle. See migration 000005_customer_approval. All
	// four timestamps are nullable so we can distinguish "never set"
	// from a zero time.Time. Default status is PENDING_APPROVAL via the
	// migration's column DEFAULT — application code does NOT set it on
	// insert.
	Status     CustomerStatus `gorm:"column:status;type:varchar(50);default:PENDING_APPROVAL"`
	ApprovedBy string         `gorm:"column:approved_by;type:varchar(255)"`
	ApprovedAt *time.Time     `gorm:"column:approved_at"`
	RejectedBy string         `gorm:"column:rejected_by;type:varchar(255)"`
	RejectedAt *time.Time     `gorm:"column:rejected_at"`
	Remark     string         `gorm:"column:remark;type:text"`

	// User                  *User `gorm:"foreignKey:UserID;references:ID"`
}
