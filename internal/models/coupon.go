package models

import (
	"time"

	"gorm.io/gorm"
)

// Coupon represents a coupon in the system
type Coupon struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`
	Amount    int            `gorm:"not null" json:"amount"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Claim represents a coupon claim by a user
type Claim struct {
	ID         uint           `gorm:"primaryKey" json:"-"`
	UserID     string         `gorm:"uniqueIndex:idx_user_coupon;not null" json:"user_id"`
	CouponName string         `gorm:"uniqueIndex:idx_user_coupon;not null;index" json:"coupon_name"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// CouponResponse represents the response structure for GetCouponDetails
type CouponResponse struct {
	Name           string   `json:"name"`
	Amount         int      `json:"amount"`
	RemainingAmount int     `json:"remaining_amount"`
	ClaimedBy      []string `json:"claimed_by"`
}

