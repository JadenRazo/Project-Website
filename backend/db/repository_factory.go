// Package db provides database access patterns and implementations
package db

// Factory functions for creating repository instances
// These provide a convenient way to create properly initialized repositories

// NewURLRepository creates a specialized URL repository
func NewURLRepository() *URLRepository {
	base := NewRepository[URL](URL{})
	return &URLRepository{Repository: base}
}

// NewUserRepository creates a specialized User repository
func NewUserRepository() *UserRepository {
	base := NewRepository[User](User{})
	return &UserRepository{Repository: base}
}

// NewAnalyticsRepository creates a specialized analytics repository
func NewAnalyticsRepository() *AnalyticsRepository {
	base := NewRepository[ClickAnalytics](ClickAnalytics{})
	return &AnalyticsRepository{Repository: base}
}

