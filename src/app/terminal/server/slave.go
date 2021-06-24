package server

import (
	"github.com/acornsoftlab/dashboard/terminal/webtty"
)

// Slave - webtty를 통해 처리할 Socket 연계용 인터페이스
type Slave interface {
	webtty.Slave

	Close() error
}

// Factory - 각 요청에 대해 Slave 구성을 처리하는 인터페이스
type Factory interface {
	Name() string
	New(params map[string]string) (Slave, error)
}
