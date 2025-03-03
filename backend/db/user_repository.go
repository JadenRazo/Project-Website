package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UserRepository provides specialized methods for user data
type UserRepository struct {
	Repository[User]
}

// FindByEmail finds a user by email address
func (r *UserRepository) FindByEmail(email string) (*User, error) {
	return r.FindBy("email", email)
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*User, error) {
	return r.FindBy("username", username)
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(id uint) error {
	result := GetDB().Model(&User{}).
		Where("id = ?", id).
		UpdateColumn("last_login_at", time.Now())
	
	if result.Error != nil {
		return fmt.Errorf("failed to update last login: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// FindByRole finds all users with a specific role
func (r *UserRepository) FindByRole(role string, page, pageSize int) ([]User, *Pagination, error) {
	return r.Paginate(page, pageSize, 
		WithOrder("username", SortAscending),
		func(db *gorm.DB) *gorm.DB {
			return db.Where("role = ?", role)
		})
}

// Search searches for users by username or email
func (r *UserRepository) Search(query string, page, pageSize int) ([]User, *Pagination, error) {
	return r.Paginate(page, pageSize,
		WithOrder("username", SortAscending),
		func(db *gorm.DB) *gorm.DB {
			if query != "" {
				return db.Where("username LIKE ? OR email LIKE ?", 
					"%"+query+"%", "%"+query+"%")
			}
			return db
		})
}

// ChangePassword updates a user's password
func (r *UserRepository) ChangePassword(userID uint, newPassword string) error {
	user, err := r.Find(userID)
	if err != nil {
		return err
	}
	
	user.Password = newPassword
	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	return r.Update(user)
}

// SetActiveStatus enables or disables a user account
func (r *UserRepository) SetActiveStatus(userID uint, active bool) error {
	result := GetDB().Model(&User{}).
		Where("id = ?", userID).
		Update("is_active", active)
		
	if result.Error != nil {
		return fmt.Errorf("failed to update active status: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}
