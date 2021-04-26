package localcommand

import (
	"syscall"
	"time"
)

// Option - LocalCommand에 Option을 설정하는 함수 시그니처
type Option func(*LocalCommand)

// WithCloseSignal - 지정한 Signal을 LocalCommand에 설정
func WithCloseSignal(signal syscall.Signal) Option {
	return func(lcmd *LocalCommand) {
		lcmd.closeSignal = signal
	}
}

// WithCloseTimeout - 지정한 시간을 LocalCommand에 설정
func WithCloseTimeout(timeout time.Duration) Option {
	return func(lcmd *LocalCommand) {
		lcmd.closeTimeout = timeout
	}
}
