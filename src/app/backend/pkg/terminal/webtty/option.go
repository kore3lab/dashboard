package webtty

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Option - WebTTY 옵션 설정 함수 시그니쳐
type Option func(*WebTTY) error

// WithPermitWrite - Slaves로 부터 입력을 받을 수 있도록 설정
func WithPermitWrite() Option {
	return func(wt *WebTTY) error {
		wt.permitWrite = true
		return nil
	}
}

// WithFixedColumns - 터미널의 Width 설정
func WithFixedColumns(columns int) Option {
	return func(wt *WebTTY) error {
		wt.columns = columns
		return nil
	}
}

// WithFixedRows - 터미널의 Height 설정
func WithFixedRows(rows int) Option {
	return func(wt *WebTTY) error {
		wt.rows = rows
		return nil
	}
}

// WithWindowTitle - 터미널 타이틀 설정
func WithWindowTitle(windowTitle []byte) Option {
	return func(wt *WebTTY) error {
		wt.windowTitle = windowTitle
		return nil
	}
}

// WithReconnect - 재 연결 가능 설정
func WithReconnect(timeInSeconds int) Option {
	return func(wt *WebTTY) error {
		wt.reconnect = timeInSeconds
		return nil
	}
}

// WithMasterPreferences - Web Socket 연결 옵션 설정
func WithMasterPreferences(preferences interface{}) Option {
	return func(wt *WebTTY) error {
		prefs, err := json.Marshal(preferences)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal preferences as JSON")
		}
		wt.masterPrefs = prefs
		return nil
	}
}
