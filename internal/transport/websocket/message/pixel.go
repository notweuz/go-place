package message

type PixelChanged struct {
	ID     uint64 `json:"id"`
	X      uint64 `json:"x"`
	Y      uint64 `json:"y"`
	UserID uint64 `json:"user_id"`
	Color  string `json:"color"`
}
