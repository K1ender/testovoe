package validators

import (
    "testing"

    "github.com/go-playground/validator/v10"
)

type payload struct {
    Date string `validate:"mm_yyyy"`
}

func TestMonthYearValidator(t *testing.T) {
    v := validator.New(validator.WithRequiredStructEnabled())
    v.RegisterValidation("mm_yyyy", MonthYearValidator)

    cases := []struct {
        in   string
        want bool
    }{
        {"01-2006", true},
        {"12-1999", true},
        {"00-2006", false},
        {"13-2006", false},
        {"1-2006", false},
        {"01-06", false},
        {"2006-01", false},
        {"", false},
    }

    for _, tc := range cases {
        err := v.Struct(payload{Date: tc.in})
        got := err == nil
        if got != tc.want {
            t.Fatalf("input %q: expected %v, got %v (err=%v)", tc.in, tc.want, got, err)
        }
    }
}

