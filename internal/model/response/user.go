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

type UserPrivateDetailed struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Charges  uint   `json:"charges"`

	LastChargedAt time.Time `json:"last_charged_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewUserPrivateDetailed(id uint64, username string, charges uint, lastChargedAt, createdAt, updatedAt time.Time) *UserPrivateDetailed {
	return &UserPrivateDetailed{
		ID:            id,
		Username:      username,
		Charges:       charges,
		LastChargedAt: lastChargedAt,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
