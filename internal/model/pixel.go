package model

import "time"

type Pixel struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement"`
	X      uint64 `gorm:"not null;uniqueIndex:idx_pixel_coords"`
	Y      uint64 `gorm:"not null;uniqueIndex:idx_pixel_coords"`
	UserID uint64
	User   User `gorm:"foreignKey:UserID;references:ID"`
	Color  string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
