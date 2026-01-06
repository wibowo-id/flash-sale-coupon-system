package handlers

import (
	"errors"
	"net/http"
	"strings"

	"ubersnap/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CouponHandler struct {
	db *gorm.DB
}

func NewCouponHandler(db *gorm.DB) *CouponHandler {
	return &CouponHandler{db: db}
}

// CreateCouponRequest represents the request body for creating a coupon
type CreateCouponRequest struct {
	Name   string `json:"name" binding:"required"`
	Amount int    `json:"amount" binding:"required,min=0"`
}

// ClaimCouponRequest represents the request body for claiming a coupon
type ClaimCouponRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	CouponName string `json:"coupon_name" binding:"required"`
}

// CreateCoupon handles POST /api/coupons
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon := models.Coupon{
		Name:   req.Name,
		Amount: req.Amount,
	}

	if err := h.db.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create coupon: " + err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// ClaimCoupon handles POST /api/coupons/claim
func (h *CouponHandler) ClaimCoupon(c *gin.Context) {
	var req ClaimCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use transaction to ensure atomicity
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// Check if coupon exists
		var coupon models.Coupon
		if err := tx.Where("name = ?", req.CouponName).First(&coupon).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &ClaimError{Message: "Coupon not found", StatusCode: http.StatusNotFound}
			}
			return err
		}

		// Check if user has already claimed this coupon
		var existingClaim models.Claim
		if err := tx.Where("user_id = ? AND coupon_name = ?", req.UserID, req.CouponName).First(&existingClaim).Error; err == nil {
			return &ClaimError{Message: "User has already claimed this coupon", StatusCode: http.StatusConflict}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// Count existing claims for this coupon
		var claimCount int64
		if err := tx.Model(&models.Claim{}).Where("coupon_name = ?", req.CouponName).Count(&claimCount).Error; err != nil {
			return err
		}

		// Check stock availability
		if int(claimCount) >= coupon.Amount {
			return &ClaimError{Message: "Coupon stock exhausted", StatusCode: http.StatusBadRequest}
		}

		// Create claim record
		claim := models.Claim{
			UserID:     req.UserID,
			CouponName: req.CouponName,
		}

		if err := tx.Create(&claim).Error; err != nil {
			// Check if it's a unique constraint violation (race condition)
			if isUniqueConstraintError(err) {
				return &ClaimError{Message: "User has already claimed this coupon", StatusCode: http.StatusConflict}
			}
			return err
		}

		return nil
	})

	if err != nil {
		if claimErr, ok := err.(*ClaimError); ok {
			c.JSON(claimErr.StatusCode, gin.H{"error": claimErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to claim coupon: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetCouponDetails handles GET /api/coupons/:name
func (h *CouponHandler) GetCouponDetails(c *gin.Context) {
	couponName := c.Param("name")

	var coupon models.Coupon
	if err := h.db.Where("name = ?", couponName).First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupon"})
		return
	}

	// Get all claims for this coupon
	var claims []models.Claim
	if err := h.db.Where("coupon_name = ?", couponName).Find(&claims).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch claims"})
		return
	}

	// Build claimed_by list
	claimedBy := make([]string, 0, len(claims))
	for _, claim := range claims {
		claimedBy = append(claimedBy, claim.UserID)
	}

	// Calculate remaining amount
	remainingAmount := coupon.Amount - len(claims)
	if remainingAmount < 0 {
		remainingAmount = 0
	}

	response := models.CouponResponse{
		Name:            coupon.Name,
		Amount:          coupon.Amount,
		RemainingAmount: remainingAmount,
		ClaimedBy:       claimedBy,
	}

	c.JSON(http.StatusOK, response)
}

// ClaimError is a custom error type for claim operations
type ClaimError struct {
	Message    string
	StatusCode int
}

func (e *ClaimError) Error() string {
	return e.Message
}

// isUniqueConstraintError checks if the error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	// Check for PostgreSQL unique constraint violation (error code 23505)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // unique_violation
	}

	// Fallback: check error message
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "unique constraint") ||
		strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "idx_user_coupon")
}

