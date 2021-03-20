package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `json:"name" gorm:"type:varchar(255) unique"`
	Password string `json:"password,omitempty" gorm:"type:varchar(255)"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("Please provide username")
	}
	if u.Password == "" {
		return fmt.Errorf("Please provide password")
	}

	return nil
}

func (u *User) Create() error {
	enc, err := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	if err != nil {
		return err
	}

	u.Password = string(enc)

	return db.Create(&u).Error
}

func (u *User) Find() (*User, error) {
	var (
		tmp User
		err error
	)

	if err = db.Where("name = ?", u.Name).First(&tmp).Error; err != nil {
		log.Printf("User not found - %s\n", err)
		return nil, fmt.Errorf("wrong username or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(u.Password)); err != nil {
		log.Printf("bcrypt error - %s\n", err)
		return nil, fmt.Errorf("wrong username or password")
	}

	return &tmp, nil
}

type JWTToken struct {
	Token string `json:"token"`
}

func (u *User) GenerateJWT() (JWTToken, error) {
	// TODO, move secret to ENV
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Minute * 60).Unix(),
		"id":   u.ID,
		"name": u.Name,
	})

	tokenString, err := token.SignedString(signingKey)

	return JWTToken{tokenString}, err
}
