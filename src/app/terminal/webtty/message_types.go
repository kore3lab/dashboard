package webtty

// Protocols - 터미널 처리를 위한 프로토콜 리스트
// webtty는 Websocket Streams의 하위 프로토콜
var Protocols = []string{"webtty"}

const (
	// UnknownInput - 알수 없는 입력
	UnknownInput = '0'
	// Input - 사용자 입력
	Input = '1'
	// Ping - ping to server
	Ping = '2'
	// ResizeTerminal - 브라우저의 터미널 크기 변경
	ResizeTerminal = '3'
)

const (
	// UnknownOutput - 알수 없는 출력
	UnknownOutput = '0'
	// Output - 일반 출력
	Output = '1'
	// Pong - pont to client
	Pong = '2'
	// SetWindowTitle - 터미널 타이틀 출력
	SetWindowTitle = '3'
	// SetPreferences - 터미널 속성 출력
	SetPreferences = '4'
	// SetReconnect - 터미널 재 연결결 출력
	SetReconnect = '5'
)
