package localcommand

import (
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"github.com/pkg/errors"
)

const (
	// DefaultCloseSignal - 기본으로 사용할 LocalCommand 종료 Signal
	DefaultCloseSignal = syscall.SIGINT
	// DefaultCloseTimeout - 기본으로 사용할 Timeout (초)
	DefaultCloseTimeout = 10 * time.Second
)

// LocalCommand - 터미널을 운영하기 위한 Local Command 정보 구조체
type LocalCommand struct {
	command string
	argv    []string

	closeSignal  syscall.Signal
	closeTimeout time.Duration

	cmd       *exec.Cmd
	pty       *os.File
	ptyClosed chan struct{}
}

// New - 지정한 명령을 지정한 아규먼트와 옵션들을 기준으로 실행할 Local Command (Slave) 인스턴스 생성
func New(command string, argv []string, options ...Option) (*LocalCommand, error) {
	cmd := exec.Command(command, argv...)
	pty, err := pty.Start(cmd)
	if err != nil {
		// todo close cmd?
		return nil, errors.Wrapf(err, "failed to start command `%s`", command)
	}
	ptyClosed := make(chan struct{})

	lcmd := &LocalCommand{
		command: command,
		argv:    argv,

		closeSignal:  DefaultCloseSignal,
		closeTimeout: DefaultCloseTimeout,

		cmd:       cmd,
		pty:       pty,
		ptyClosed: ptyClosed,
	}

	for _, option := range options {
		option(lcmd)
	}

	// 사용자에 의해서 프로세스가 종료되면 터미널 종료
	go func() {
		defer func() {
			lcmd.pty.Close()
			close(lcmd.ptyClosed)
		}()
		lcmd.cmd.Wait()
	}()

	return lcmd, nil
}

// Read - 터미널에서 정보 읽기
func (lcmd *LocalCommand) Read(p []byte) (n int, err error) {
	return lcmd.pty.Read(p)
}

// Write - 터미널로 정보 출력
func (lcmd *LocalCommand) Write(p []byte) (n int, err error) {
	return lcmd.pty.Write(p)
}

// Close - 터미널 종료 (연계된 프로세스 종료)
func (lcmd *LocalCommand) Close() error {
	if lcmd.cmd != nil && lcmd.cmd.Process != nil {
		//자식프로세스 까지 종료시그널 전달을 위해서 재정의 함
		//lcmd.cmd.Process.Signal(lcmd.closeSignal)
		syscall.Kill(-lcmd.cmd.Process.Pid, syscall.SIGINT)
	}
	for {
		select {
		case <-lcmd.ptyClosed:
			return nil
		case <-lcmd.closeTimeoutC():
			//lcmd.cmd.Process.Signal(syscall.SIGKILL)
			syscall.Kill(-lcmd.cmd.Process.Pid, syscall.SIGKILL)
		}
	}
}

// WindowTitleVariables - 현재 실행되고 터미널 정보 반환
func (lcmd *LocalCommand) WindowTitleVariables() map[string]interface{} {
	return map[string]interface{}{
		"command": lcmd.command,
		"argv":    lcmd.argv,
		"pid":     lcmd.cmd.Process.Pid,
	}
}

// ResizeTerminal - 지정한 크기로 터미널 크기 조정
func (lcmd *LocalCommand) ResizeTerminal(width int, height int) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(height),
		uint16(width),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		lcmd.pty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

// closeTimeoutC - Timer에 의해 Close Timeout 채널 설정
func (lcmd *LocalCommand) closeTimeoutC() <-chan time.Time {
	if lcmd.closeTimeout >= 0 {
		return time.After(lcmd.closeTimeout)
	}

	return make(chan time.Time)
}
