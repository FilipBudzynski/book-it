package services

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username" form:"username"`
	Email    string `gorm:"unique;not null" json:"email" form:"email"`
	ID       uint   `gorm:"primaryKey" json:"id"`
}

func NewUserService(db *gorm.DB) *UserService {
	// TODO: user atlas as migration
	db.AutoMigrate(&User{})
	return &UserService{
		db: db,
	}
}

type UserService struct {
	db *gorm.DB
}

func (u *UserService) Create(user *User) error {
	return u.db.Create(user).Error
}

func (u *UserService) GetById(id uint) (*User, error) {
	var user User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) GetAll() ([]User, error) {
	var users []User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) Update(user *User) error {
	return u.db.Save(user).Error
}

func (u *UserService) Delete(user User) error {
	return u.db.Delete(user).Error
}
