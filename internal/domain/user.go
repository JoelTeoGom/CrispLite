package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

type UserSummary struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (u *User) Validate() error {
	switch {
	case u.Username == "":
		return ErrUsernameEmpty
	case len(u.Username) < 3:
		return ErrUsernameTooShort
	case len(u.Username) > 50:
		return ErrUsernameTooLong
	case u.Password == "":
		return ErrPasswordEmpty
	case len(u.Password) < 8:
		return ErrPasswordTooShort
	case len(u.Password) > 72:
		return ErrPasswordTooLong
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashed, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
