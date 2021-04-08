package server

import (
	"github.com/gorilla/websocket"
)

// wsWrapper - Web Socket 운영을 위한 Wrapper 구조체
type WsWrapper struct {
	*websocket.Conn
}

// Write - Web Socket Stream에 출력
func (wsw *WsWrapper) Write(p []byte) (n int, err error) {
	writer, err := wsw.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
}

// Read - Web Socket Stream에서 입력
func (wsw *WsWrapper) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := wsw.Conn.NextReader()
		if err != nil {
			return 0, err
		}

		if msgType != websocket.TextMessage {
			continue
		}

		return reader.Read(p)
	}
}
