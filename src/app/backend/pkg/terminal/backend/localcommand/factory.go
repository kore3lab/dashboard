package localcommand

import (
	"syscall"
	"time"

	"github.com/acornsoftlab/dashboard/pkg/terminal/server"
)

// Options - LocalCommand 운영에 필요한 옵션 정보 구조체
type Options struct {
	// CloseSignal - 터미널이 종료될 때 프로세스에 보낼 종료 시그널 (SIGHUB)
	CloseSignal int `hcl:"close_signal" flagName:"close-signal" flagSName:"" flagDescribe:"Signal sent to the command process when k3wt close it (default: SIGHUP)" default:"1"`
	// CloseTimout - 클라이언트 연결이 끊어진 이후에 프로세스 종료까지의 Timeout (초)
	CloseTimeout int `hcl:"close_timeout" flagName:"close-timeout" flagSName:"" flagDescribe:"Time in seconds to force kill process after client is disconnected (default: -1)" default:"-1"`
}

// Factory - LocalCommand 구성에 필요한 정보 구조체
type Factory struct {
	command string
	argv    []string
	options *Options
	opts    []Option
}

// NewFactory - 지정한 명령과 아규먼트 및 옵션들을 기준으로 LocalCommand 운영을 위한 Factory 인스턴스 생성
func NewFactory(command string, argv []string, options *Options) (*Factory, error) {
	opts := []Option{WithCloseSignal(syscall.Signal(options.CloseSignal))}
	if options.CloseTimeout >= 0 {
		opts = append(opts, WithCloseTimeout(time.Duration(options.CloseTimeout)*time.Second))
	}

	return &Factory{
		command: command,
		argv:    argv,
		options: options,
		opts:    opts,
	}, nil
}

// Name - Factory 명 반환
func (factory *Factory) Name() string {
	return "local command"
}

// New -  지정한 파라미터를 기준으로 터미널 Slave 인스턴스 생성
func (factory *Factory) New(params map[string]string) (server.Slave, error) {
	argv := make([]string, len(factory.argv))
	copy(argv, factory.argv)

	for key, val := range params {
		argv = append(argv, "--"+key)
		argv = append(argv, val)
	}

	// if params["arg"] != nil && len(params["arg"]) > 0 {
	// 	argv = append(argv, params["arg"]...)
	// 	/*
	// 		for _, arg := range params["arg"] {
	// 			argv = append(argv, arg)
	// 		}
	// 	*/
	// }

	return New(factory.command, argv, factory.opts...)
}
