package token

import "time"

const (
	//DefaultExpiration never expire
	DefaultExpiration = 0
)

//TtyParameter kubectl tty param
type TtyParameter struct {
	//	Arg []string
	Arg map[string]string
}

//interface that defines token cache behavior
type Cache interface {
	Get(token string) *TtyParameter
	Delete(token string) error
	Add(token string, param *TtyParameter, d time.Duration) error
}
