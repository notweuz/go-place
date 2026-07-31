package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Option func(*gorm.DB) *gorm.DB

func WithPreload(relation string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(relation, args...)
	}
}

func WithLock() Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
}

func WithTx(tx *gorm.DB) Option {
	return func(db *gorm.DB) *gorm.DB {
		return tx
	}
}
