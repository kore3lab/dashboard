package webtty

import (
	"io"
)

// Master - 터미널 (Slave)과 연계되는 Socket Connection
type Master io.ReadWriter
