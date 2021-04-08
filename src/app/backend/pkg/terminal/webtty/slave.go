package webtty

import (
	"io"
)

// Slave - Local Command를 실행하는 터미널
type Slave interface {
	io.ReadWriter

	// WindowTitleVariables - 터미널 타이틀 구성 정보
	WindowTitleVariables() map[string]interface{}

	// ResizeTerminal - 터미널 크기 정보
	ResizeTerminal(columns int, rows int) error
}
