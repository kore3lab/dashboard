package lang

import (
	"math"
	"math/rand"
)

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

func ArrayContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
