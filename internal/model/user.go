package model

import "time"

type User struct {
	ID       uint64  `gorm:"primaryKey;autoIncrement"`
	Username string  `gorm:"not null;uniqueIndex;size:50"`
	Password string  `gorm:"not null;size:24"`
	Charges  uint    `gorm:"not null;default:0"`
	Pixels   []Pixel `gorm:"foreignKey:UserID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}
