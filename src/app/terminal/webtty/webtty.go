package webtty

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// WebTTY - 터미널(Slave)와 Master(Web Socket Connection) 연결하는 정보 구조체
// 텍스트 스트림과 터미널 크기 조정과 같은 부수적인 명령 처리를 담당하며 기본 프로토콜을 사용한다.
type WebTTY struct {
	// PTY Master, which probably a connection to browser
	masterConn Master
	// PTY Slave
	slave Slave

	windowTitle []byte
	permitWrite bool
	columns     int
	rows        int
	reconnect   int // in seconds
	masterPrefs []byte

	bufferSize   int
	writeMutex   sync.Mutex
	lastPingTime time.Time
}

const (
	// MaxBufferSize - 연결에 사용할 Buffer Size
	MaxBufferSize = 1024 * 1024 * 1
)

// New - 지정한 Web Socket 연결과 Local Command 정보를 기준으로 WebTTY 인스턴스 생성
func New(masterConn Master, slave Slave, options ...Option) (*WebTTY, error) {
	wt := &WebTTY{
		masterConn: masterConn,
		slave:      slave,

		permitWrite: false,
		columns:     0,
		rows:        0,

		bufferSize:   MaxBufferSize,
		lastPingTime: time.Now(),
	}

	for _, option := range options {
		option(wt)
	}

	return wt, nil
}

// Run - WetbTTY 프로세스 시작
// 해당 메서드는 지정한 Context가 취소될 떄까지 유지되며, Master와 Slave도 유지된다. 따라서 호출한 곳에서 종료처리해야 한다.
// Master나 Slave 중에 하나가 종료되면 ErrSlaveClosed 또는 ErrMasterClosed 오류를 반환한다.
func (wt *WebTTY) Run(ctx context.Context) error {
	// 터미널 초기화 메시지 전송
	err := wt.sendInitializeMessage()
	if err != nil {
		return errors.Wrapf(err, "failed to send initializing message")
	}

	errs := make(chan error, 3)

	// Slave와 연동 처리
	slaveBuffer := make([]byte, wt.bufferSize)
	go func() {
		errs <- func() error {
			defer func() {
				if e := recover(); e != nil {
				}
			}()
			for {
				if slaveBuffer == nil {
					return ErrSlaveClosed
				}
				n, err := wt.slave.Read(slaveBuffer)
				if err != nil {
					return ErrSlaveClosed
				}
				err = wt.handleSlaveReadEvent(slaveBuffer[:n])
				if err != nil {
					return err
				}
			}
		}()
	}()

	// Master와 연동 처리
	masterBuffer := make([]byte, wt.bufferSize)
	go func() {
		errs <- func() error {
			defer func() {
				if e := recover(); e != nil {
				}
			}()
			for {
				if masterBuffer == nil {
					return ErrMasterClosed
				}
				n, err := wt.masterConn.Read(masterBuffer)
				if err != nil {
					return ErrMasterClosed
				}
				err = wt.handleMasterReadEvent(masterBuffer[:n])
				if err != nil {
					return err
				}
			}
		}()
	}()

	//
	go func() {
		errs <- func() error {
			lostPingTimeout := time.Duration(180) * time.Second
			seconds, _err := strconv.Atoi(os.Getenv("LOST_PING_TIMEOUT_SECONDS"))
			if _err != nil && seconds > 30 {
				lostPingTimeout = time.Duration(seconds) * time.Second
			}
			for {
				time.Sleep(time.Duration(30) * time.Second)
				if err != nil {
					return err
				}
				if time.Now().After(wt.lastPingTime.Add(lostPingTimeout)) {
					return ErrConnectionLostPing
				}
			}
		}()
	}()

	defer func() {
		slaveBuffer = nil
		masterBuffer = nil
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errs:
	}

	return err
}

// sendInitializeMessage - Socket으로 초기화 메시지 전송
func (wt *WebTTY) sendInitializeMessage() error {
	err := wt.masterWrite(append([]byte{SetWindowTitle}, wt.windowTitle...))
	if err != nil {
		return errors.Wrapf(err, "failed to send window title")
	}

	// 재 연결 여부 및 재연결 출력
	if wt.reconnect > 0 {
		reconnect, _ := json.Marshal(wt.reconnect)
		err := wt.masterWrite(append([]byte{SetReconnect}, reconnect...))
		if err != nil {
			return errors.Wrapf(err, "failed to set reconnect")
		}
	}

	// 터미널 속성 검증 및 출력
	if wt.masterPrefs != nil {
		err := wt.masterWrite(append([]byte{SetPreferences}, wt.masterPrefs...))
		if err != nil {
			return errors.Wrapf(err, "failed to set preferences")
		}
	}

	return nil
}

// handleSlaveReadEvent - Slave에서 수신된 데이터 처리
func (wt *WebTTY) handleSlaveReadEvent(data []byte) error {
	safeMessage := base64.StdEncoding.EncodeToString(data)
	err := wt.masterWrite(append([]byte{Output}, []byte(safeMessage)...))
	if err != nil {
		return errors.Wrapf(err, "failed to send message to master")
	}

	return nil
}

// masterWrite - 데이터를 Master로 출력
func (wt *WebTTY) masterWrite(data []byte) error {
	wt.writeMutex.Lock()
	defer wt.writeMutex.Unlock()

	_, err := wt.masterConn.Write(data)
	if err != nil {
		return errors.Wrapf(err, "failed to write to master")
	}

	return nil
}

// handleMasterReadEvent - Master로 수신된 데이터 처리
func (wt *WebTTY) handleMasterReadEvent(data []byte) error {
	if len(data) == 0 {
		return errors.New("unexpected zero length read from master")
	}

	switch data[0] {
	case Input:
		if !wt.permitWrite {
			return nil
		}

		if len(data) <= 1 {
			return nil
		}
		_, err := wt.slave.Write(data[1:])
		if err != nil {
			return errors.Wrapf(err, "failed to write received data to slave")
		}

	case Ping:
		err := wt.masterWrite([]byte{Pong})
		wt.lastPingTime = time.Now()
		if err != nil {
			return errors.Wrapf(err, "failed to return Pong message to master")
		}

	case ResizeTerminal:
		if wt.columns != 0 && wt.rows != 0 {
			break
		}

		if len(data) <= 1 {
			return errors.New("received malformed remote command for terminal resize: empty payload")
		}

		var args argResizeTerminal
		err := json.Unmarshal(data[1:], &args)
		if err != nil {
			return errors.Wrapf(err, "received malformed data for terminal resize")
		}
		rows := wt.rows
		if rows == 0 {
			rows = int(args.Rows)
		}

		columns := wt.columns
		if columns == 0 {
			columns = int(args.Columns)
		}

		return wt.slave.ResizeTerminal(columns, rows)
	default:
		return errors.Errorf("unknown message type `%c`", data[0])
	}

	return nil
}

// argResizeTerminal - 터미널 크기 변경 정보 구조체
type argResizeTerminal struct {
	Columns float64
	Rows    float64
}
