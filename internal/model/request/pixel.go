package request

type ChangePixel struct {
	X     uint64 `json:"x" validate:"required"`
	Y     uint64 `json:"y" validate:"required"`
	Color string `json:"color" validate:"required"`
}
