package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	userDB UserDB
)

func SetUserDB(d UserDB) {
	userDB = d
}

type UserDB interface {
	SaveUser(*User) error
	GetUser(*User) error
	DeleteUser(*User) error
}

type User struct {
	ID          uuid.UUID  `json:"id"`
	DisplayName string     `json:"display_name"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	BornDate    *time.Time `json:"born_date"`
	Address     *string    `json:"address"`
}

type UserConfig func(*User)

func WithUsername(username string) UserConfig {
	return func(u *User) {
		u.Username = username
	}
}

func WithPassword(password string) UserConfig {
	return func(u *User) {
		u.Password = password
	}
}

func WithDisplayName(name string) UserConfig {
	return func(u *User) {
		u.DisplayName = name
	}
}

func NewUser(configs ...UserConfig) *User {
	u := &User{}
	for _, config := range configs {
		config(u)
	}
	return u
}

func (u *User) Login() error {
	if err := userDB.GetUser(u); err != nil {
		return err
	}
	return nil
}

func (u *User) Logout() error {
	return nil
}

func (u *User) Update() error {
	return nil
}

func (u *User) Save() error {
	return userDB.SaveUser(u)
}
