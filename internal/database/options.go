package database

import "gorm.io/gorm"

type Option func(*gorm.DB) *gorm.DB

func WithPreload(relation string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(relation, args...)
	}
}
