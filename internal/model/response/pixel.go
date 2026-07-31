package response

import "time"

type PixelFull struct {
	ID     uint64 `json:"id"`
	X      uint64 `json:"x"`
	Y      uint64 `json:"y"`
	UserID uint64 `json:"author_id"`
	Color  string `json:"color"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPixelFull(id, x, y, userID uint64, color string, createdAt, updatedAt time.Time) *PixelFull {
	return &PixelFull{
		ID:        id,
		X:         x,
		Y:         y,
		UserID:    userID,
		Color:     color,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type PixelsChangeFull struct {
	Pixels []PixelFull `json:"pixels"`
}

func NewPixelsChangeFull(pixels []PixelFull) *PixelsChangeFull {
	return &PixelsChangeFull{Pixels: pixels}
}
