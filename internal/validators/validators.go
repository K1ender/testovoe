package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	monthYearRegex = regexp.MustCompile(`^(0[1-9]|1[0-2])-\d{4}$`)
)

func MonthYearValidator(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return monthYearRegex.MatchString(val)
}
