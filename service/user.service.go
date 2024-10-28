package service

import (
	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"not null" json:"username" form:"username"`
	Email    string `gorm:"unique;not null" json:"email" form:"email"`
	ID       uint   `gorm:"primaryKey" json:"id"`
}

func NewUserService(db *gorm.DB) *UserService {
	db.AutoMigrate(&User{})

	return &UserService{
		db: db,
	}
}

type UserService struct {
	db   *gorm.DB
	User User
}

func (u *UserService) GetUserById(id uint) (*User, error) {
	var user User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) GetAllUsers() ([]User, error) {
	var users []User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) CreateUser(user *User) error {
	return u.db.Create(user).Error
}

func (u *UserService) UpdateUser(user *User) error {
	return u.db.Save(user).Error
}

func (u *UserService) DeleteUser(id uint) error {
	return u.db.Delete(&User{}, id).Error
}
