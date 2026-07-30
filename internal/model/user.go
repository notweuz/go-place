package model

import "time"

type User struct {
	ID       uint64  `gorm:"primaryKey;autoIncrement"`
	Username string  `gorm:"not null;uniqueIndex;size:50"`
	Password string  `gorm:"not null;size:24"`
	Pixels   []Pixel `gorm:"foreignKey:UserID"`
	Charges  uint    `gorm:"not null;default:0"`

	LastChargeAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}

func (u *User) CurrentCharges(maxCharges uint, regenRate time.Duration) uint {
	if u.Charges >= maxCharges {
		return maxCharges
	}
	timePassed := time.Since(u.LastChargeAt)

	gainedCharges := uint(timePassed / regenRate)
	totalCharges := u.Charges + gainedCharges
	if totalCharges > maxCharges {
		return maxCharges
	}

	return totalCharges
}
