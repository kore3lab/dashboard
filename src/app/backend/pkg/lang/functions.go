package lang

import "math"

// NVL is null value logic
func NVL(str string, def string) string {
	if len(str) == 0 {
		return def
	}
	return str
}
func Divide(a int64, b int64) float64 {

	if b == 0 {
		return 0
	} else {
		return float64(a) / float64(b)
	}

}
func DivideRound(a int64, b int64, decimal int) float64 {

	if b == 0 {
		return 0
	} else {
		pow := math.Pow10(decimal)
		return math.Round((float64(a)/float64(b))*pow) / pow
	}

}
