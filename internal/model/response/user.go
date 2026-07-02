package response

import "time"

type UserPublic struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func NewUserPublic(id uint64, username string, createdAt time.Time) *UserPublic {
	return &UserPublic{
		ID:        id,
		Username:  username,
		CreatedAt: createdAt,
	}
}
