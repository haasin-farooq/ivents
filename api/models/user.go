package main

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);unique_index" json:"email"`
	FirstName string `gorm:"size:100;not null" json:"firstname"`
	LastName string `gorm:"size:100;not null" json:"lastname"`
	Password string `gorm:"size:100;not null" json:"password"`
	ProfileImage string `gorm:"size:255" json:"profileimage"`
}

// HashPassword hashes password from user input
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

// CheckPasswordHash checks password hash and password from user input if they match
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

// BeforeSave hashes user password
func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// Prepare strips user input of any white spaces
func (u *User) Prepare() {
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.ProfileImage = strings.TrimSpace(u.ProfileImage)
}

// Validate user input
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		return nil
	default:
		if u.FirstName == "" {
			return errors.New("firstname is required")
		}
		if u.LastName == "" {
			return errors.New("lastname is required")
		}
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}

// SaveUser adds a user to the database
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// GetUser returns a user based on email
func (u *User) GetUser(db *gorm.DB) (*User, error) {
	user := &User{}
	if err := db.Debug().Table("users").Where("email = ?", u.Email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers returns a list of all the user
func GetAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, nil
}