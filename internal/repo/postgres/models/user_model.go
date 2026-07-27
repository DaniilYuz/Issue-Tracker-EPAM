package models

import "github.com/DaniilYuz/Issue-Tracker-EPAM/internal/domain"

type UserModel struct {
	ID           string `gorm:"primaryKey"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	EmailAddress string `gorm:"not null;uniqueindex"`
}

func (UserModel) TableName() string { return "users" }

func UserModelFromDomain(u *domain.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		EmailAddress: u.EmailAddress,
	}
}

func UserModelToDomain(u *UserModel) *domain.User {
	return &domain.User{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		EmailAddress: u.EmailAddress,
	}
}
