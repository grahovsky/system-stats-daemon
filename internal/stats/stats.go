package stats

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidData = errors.New("incorrect output")

func SafeParseFloat(sVal string) float64 {
	prep := strings.ReplaceAll(sVal, ",", ".")
	fVal, err := strconv.ParseFloat(prep, 64)
	if err != nil {
		return 0.0
	}
	return fVal
}
