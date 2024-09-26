package entity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func NewUser(id, email, password string) (*User, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if email == "" {
		log.Error("Email is required", log.Fstring("email", email))
		return nil, fmt.Errorf("email is required")
	}
	if password == "" {
		log.Error("Password is required", log.Fstring("password", password))
		return nil, fmt.Errorf("password is required")
	}
	return &User{
		ID:       id,
		Email:    email,
		Password: password,
	}, nil
}

func (u *User) CompareHashAndPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
