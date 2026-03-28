package normalize

import (
	"errors"
	"math"
	"strings"
)

var ErrInvalidAmount = errors.New("monto inválido")
var ErrAmountMustBePositive = errors.New("el monto debe ser mayor a cero")

func Spaces(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func Optional(value string) string {
	return Spaces(strings.TrimSpace(value))
}

func LowerIdentifier(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func LowerKey(value string) string {
	return strings.ToLower(Optional(value))
}

func Money(value float64) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, ErrInvalidAmount
	}

	rounded := math.Round(value*100) / 100
	if rounded <= 0 {
		return 0, ErrAmountMustBePositive
	}

	return rounded, nil
}
