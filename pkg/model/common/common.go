package concern

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CommonWithIDs struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (u *CommonWithIDs) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	return
}

func (u *CommonWithIDs) AfterUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

type CommonWithID struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (u *CommonWithID) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

func (u *CommonWithID) AfterUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}
