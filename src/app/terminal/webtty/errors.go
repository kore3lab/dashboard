package webtty

import (
	"errors"
)

var (
	// ErrSlaveClosed - Slave (Local Command 처리 터미널)가 종료된 경우 오류
	ErrSlaveClosed = errors.New("slave closed")

	// ErrMasterClosed - Master (Web Socket Connection)가 종료된 경우 오류
	ErrMasterClosed = errors.New("master closed")

	// ErrConnectionLostPing - 지정한 시간 내에 ping이 없는 경우 오류 (Connection Lost)
	ErrConnectionLostPing = errors.New("connection lost ping")
)
