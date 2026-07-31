package message

type PixelChanged struct {
	ID     uint64 `json:"id"`
	X      uint64 `json:"x"`
	Y      uint64 `json:"y"`
	UserID uint64 `json:"user_id"`
	Color  string `json:"color"`
}

func NewPixelChanged(id, x, y, userID uint64, color string) *PixelChanged {
	return &PixelChanged{
		ID:     id,
		X:      x,
		Y:      y,
		UserID: userID,
		Color:  color,
	}
}

type PixelsChanged struct {
	Pixels []PixelChanged `json:"pixels"`
}

func NewPixelsChanged(pixels []PixelChanged) *PixelsChanged {
	return &PixelsChanged{
		Pixels: pixels,
	}
}
