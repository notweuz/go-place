package request

type ChangePixel struct {
	X     uint64 `json:"x" validate:"min=0"`
	Y     uint64 `json:"y" validate:"min=0"`
	Color string `json:"color" validate:"required"`
}

type ChangePixels struct {
	Pixels []ChangePixel `json:"pixels" validate:"dive"`
}
