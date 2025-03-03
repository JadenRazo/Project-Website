package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DomainRepository provides specialized methods for custom domains
type DomainRepository struct {
	Repository[CustomDomain]
}

// NewDomainRepository creates a new domain repository
func NewDomainRepository() *DomainRepository {
	base := NewRepository[CustomDomain](CustomDomain{})
	return &DomainRepository{Repository: base}
}

// FindByDomain finds a custom domain by its domain name
func (r *DomainRepository) FindByDomain(domain string) (*CustomDomain, error) {
	return r.FindBy("domain", domain)
}

// FindUserDomains gets all domains for a specific user
func (r *DomainRepository) FindUserDomains(userID uint) ([]CustomDomain, error) {
	return r.FindAllBy("user_id", userID, 
		WithOrder("domain", SortAscending))
}

// SetVerified marks a domain as verified or unverified
func (r *DomainRepository) SetVerified(domainID uint, verified bool) error {
	updates := map[string]interface{}{
		"is_verified": verified,
	}
	
	if verified {
		now := time.Now()
		updates["verified_at"] = now
	}
	
	result := GetDB().Model(&CustomDomain{}).
		Where("id = ?", domainID).
		Updates(updates)
		
	if result.Error != nil {
		return fmt.Errorf("failed to update verification status: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// VerifyByToken verifies a domain using a token
func (r *DomainRepository) VerifyByToken(domain string, token string) error {
	// This is a placeholder for actual verification logic
	// In a real implementation, you might check a DNS record or a file on the domain
	
	// Find domain first
	domainRecord, err := r.FindByDomain(domain)
	if err != nil {
		return err
	}
	
	// For demonstration purposes, we'll just verify it
	return r.SetVerified(domainRecord.ID, true)
}

// CountVerifiedDomains counts verified domains for a user
func (r *DomainRepository) CountVerifiedDomains(userID uint) (int64, error) {
	return r.Count(func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ? AND is_verified = ?", userID, true)
	})
}
