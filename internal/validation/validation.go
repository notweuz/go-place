package validation

import (
	"errors"
	"fmt"
	"go-place/internal/errs"

	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

var validate = validator.New()

func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		log.Error().Err(err).Msg("Input validation failed")

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			var sb strings.Builder
			for i, e := range validationErrs {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(fmt.Sprintf("field '%s' failed on the '%s' tag", e.Field(), e.Tag()))
			}
			return errs.NewAppError(err, fiber.StatusBadRequest, sb.String())
		}
		return errs.NewAppError(err, fiber.StatusBadRequest, "validation failed")
	}
	return nil
}
