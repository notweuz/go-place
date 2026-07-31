package model

import "time"

type Setting struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement"`
	CanvasWidth      uint64 `gorm:"not null;default:1000"`
	CanvasHeight     uint64 `gorm:"not null;default:1000"`
	CooldownSeconds  uint64 `gorm:"not null;default:30"`
	MaxCharges       uint   `gorm:"not null;default:5"`
	RegistrationOpen bool   `gorm:"not null;default:true"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewDefaultSetting() *Setting {
	return &Setting{
		CanvasWidth:      1000,
		CanvasHeight:     1000,
		CooldownSeconds:  30,
		MaxCharges:       5,
		RegistrationOpen: true,
	}
}
