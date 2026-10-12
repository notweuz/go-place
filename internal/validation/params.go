package validation

import (
	"fmt"
	"go-place/internal/errs"
	"strconv"
)

func ParseUint64(raw, name string, allowZero bool) (uint64, error) {
	if raw == "" {
		return 0, errs.BadRequest(nil, fmt.Sprintf("%s is required", name))
	}
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || (!allowZero && n == 0) {
		return 0, errs.BadRequest(err, fmt.Sprintf("invalid %s", name))
	}
	return n, nil
}
