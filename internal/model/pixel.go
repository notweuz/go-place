package model

import "time"

type Pixel struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement"`
	X      uint64 `gorm:"not null"`
	Y      uint64 `gorm:"not null"`
	UserID uint64
	User   User `gorm:"foreignKey:UserID;references:ID"`
	Color  string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
