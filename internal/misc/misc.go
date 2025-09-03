package misc

import (
	"math"
)

func CalculatePercentage(base float64, percentage float64) float64 {
	return Round(base * percentage / 100)
}

func Round(f float64) float64 {
	return math.Round(f*100) / 100
}
