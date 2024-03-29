package server

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Options - Server 옵션 구조
type Options struct {
	Address             string           `hcl:"address" flagName:"address" flagSName:"a" flagDescribe:"IP address to listen" default:"0.0.0.0"`
	Port                string           `hcl:"port" flagName:"port" flagSName:"p" flagDescribe:"Port number to liten" default:"8080"`
	PermitWrite         bool             `hcl:"permit_write" flagName:"permit-write" flagSName:"w" flagDescribe:"Permit clients to write to the TTY (BE CAREFUL)" default:"false"`
	EnableBasicAuth     bool             `hcl:"enable_basic_auth" default:"false"`
	Credential          string           `hcl:"credential" flagName:"credential" flagSName:"c" flagDescribe:"Credential for Basic Authentication (ex: user:pass, default disabled)" default:""`
	EnableRandomURL     bool             `hcl:"enable_random_url" flagName:"random-url" flagSName:"r" flagDescribe:"Add a random string to the URL" default:"false"`
	RandomURLLength     int              `hcl:"random_url_length" flagName:"random-url-length" flagDescribe:"Random URL length" default:"8"`
	EnableTLS           bool             `hcl:"enable_tls" flagName:"tls" flagSName:"t" flagDescribe:"Enable TLS/SSL" default:"false"`
	TLSCrtFile          string           `hcl:"tls_crt_file" flagName:"tls-crt" flagDescribe:"TLS/SSL certificate file path" default:"~/.k3wt.crt"`
	TLSKeyFile          string           `hcl:"tls_key_file" flagName:"tls-key" flagDescribe:"TLS/SSL key file path" default:"~/.k3wt.key"`
	EnableTLSClientAuth bool             `hcl:"enable_tls_client_auth" default:"false"`
	TLSCACrtFile        string           `hcl:"tls_ca_crt_file" flagName:"tls-ca-crt" flagDescribe:"TLS/SSL CA certificate file for client certifications" default:"~/.k3wt.ca.crt"`
	IndexFile           string           `hcl:"index_file" flagName:"index" flagDescribe:"Custom index.html file" default:""`
	TitleFormat         string           `hcl:"title_format" flagName:"title-format" flagSName:"" flagDescribe:"Title format of browser window" default:"{{ .command }}@{{ .hostname }}"`
	EnableReconnect     bool             `hcl:"enable_reconnect" flagName:"reconnect" flagDescribe:"Enable reconnection" default:"false"`
	ReconnectTime       int              `hcl:"reconnect_time" flagName:"reconnect-time" flagDescribe:"Time to reconnect" default:"10"`
	MaxConnection       int              `hcl:"max_connection" flagName:"max-connection" flagDescribe:"Maximum connection to k3wt" default:"0"`
	Once                bool             `hcl:"once" flagName:"once" flagDescribe:"Accept only one client and exit on disconnection" default:"false"`
	Timeout             int              `hcl:"timeout" flagName:"timeout" flagDescribe:"Timeout seconds for waiting a client(0 to disable)" default:"0"`
	PermitArguments     bool             `hcl:"permit_arguments" flagName:"permit-arguments" flagDescribe:"Permit clients to send command line arguments in URL (e.g. http://example.com:8080/?arg=AAA&arg=BBB)" default:"true"`
	Preferences         *HtermPrefernces `hcl:"preferences"`
	Width               int              `hcl:"width" flagName:"width" flagDescribe:"Static width of the screen, 0(default) means dynamically resize" default:"0"`
	Height              int              `hcl:"height" flagName:"height" flagDescribe:"Static height of the screen, 0(default) means dynamically resize" default:"0"`
	WSOrigin            string           `hcl:"ws_origin" flagName:"ws-origin" flagDescribe:"A regular expression that matches origin URLs to be accepted by WebSocket. No cross origin requests are acceptable by default" default:""`
	Term                string           `hcl:"term" flagName:"term" flagDescribe:"Terminal name to use on the browser, one of xterm or hterm." default:"xterm"`

	TitleVariables map[string]interface{}
}

// Validate - 옵션 검증
func (options *Options) Validate() error {
	if options.EnableTLSClientAuth && !options.EnableTLS {
		return errors.New("TLS client authentication is enabled, but TLS is not enabled")
	}
	return nil
}

// HtermPrefernces - HTerm 속성 구조
type HtermPrefernces struct {
	AltGrMode                     *string                      `hcl:"alt_gr_mode" json:"alt-gr-mode,omitempty"`
	AltBackspaceIsMetaBackspace   bool                         `hcl:"alt_backspace_is_meta_backspace" json:"alt-backspace-is-meta-backspace,omitempty"`
	AltIsMeta                     bool                         `hcl:"alt_is_meta" json:"alt-is-meta,omitempty"`
	AltSendsWhat                  string                       `hcl:"alt_sends_what" json:"alt-sends-what,omitempty"`
	AudibleBellSound              string                       `hcl:"audible_bell_sound" json:"audible-bell-sound,omitempty"`
	DesktopNotificationBell       bool                         `hcl:"desktop_notification_bell" json:"desktop-notification-bell,omitempty"`
	BackgroundColor               string                       `hcl:"background_color" json:"background-color,omitempty"`
	BackgroundImage               string                       `hcl:"background_image" json:"background-image,omitempty"`
	BackgroundSize                string                       `hcl:"background_size" json:"background-size,omitempty"`
	BackgroundPosition            string                       `hcl:"background_position" json:"background-position,omitempty"`
	BackspaceSendsBackspace       bool                         `hcl:"backspace_sends_backspace" json:"backspace-sends-backspace,omitempty"`
	CharacterMapOverrides         map[string]map[string]string `hcl:"character_map_overrides" json:"character-map-overrides,omitempty"`
	CloseOnExit                   bool                         `hcl:"close_on_exit" json:"close-on-exit,omitempty"`
	CursorBlink                   bool                         `hcl:"cursor_blink" json:"cursor-blink,omitempty"`
	CursorBlinkCycle              [2]int                       `hcl:"cursor_blink_cycle" json:"cursor-blink-cycle,omitempty"`
	CursorColor                   string                       `hcl:"cursor_color" json:"cursor-color,omitempty"`
	ColorPaletteOverrides         []*string                    `hcl:"color_palette_overrides" json:"color-palette-overrides,omitempty"`
	CopyOnSelect                  bool                         `hcl:"copy_on_select" json:"copy-on-select,omitempty"`
	UseDefaultWindowCopy          bool                         `hcl:"use_default_window_copy" json:"use-default-window-copy,omitempty"`
	ClearSelectionAfterCopy       bool                         `hcl:"clear_selection_after_copy" json:"clear-selection-after-copy,omitempty"`
	CtrlPlusMinusZeroZoom         bool                         `hcl:"ctrl_plus_minus_zero_zoom" json:"ctrl-plus-minus-zero-zoom,omitempty"`
	CtrlCCopy                     bool                         `hcl:"ctrl_c_copy" json:"ctrl-c-copy,omitempty"`
	CtrlVPaste                    bool                         `hcl:"ctrl_v_paste" json:"ctrl-v-paste,omitempty"`
	EastAsianAmbiguousAsTwoColumn bool                         `hcl:"east_asian_ambiguous_as_two_column" json:"east-asian-ambiguous-as-two-column,omitempty"`
	Enable8BitControl             *bool                        `hcl:"enable_8_bit_control" json:"enable-8-bit-control,omitempty"`
	EnableBold                    *bool                        `hcl:"enable_bold" json:"enable-bold,omitempty"`
	EnableBoldAsBright            bool                         `hcl:"enable_bold_as_bright" json:"enable-bold-as-bright,omitempty"`
	EnableClipboardNotice         bool                         `hcl:"enable_clipboard_notice" json:"enable-clipboard-notice,omitempty"`
	EnableClipboardWrite          bool                         `hcl:"enable_clipboard_write" json:"enable-clipboard-write,omitempty"`
	EnableDec12                   bool                         `hcl:"enable_dec12" json:"enable-dec12,omitempty"`
	Environment                   map[string]string            `hcl:"environment" json:"environment,omitempty"`
	FontFamily                    string                       `hcl:"font_family" json:"font-family,omitempty"`
	FontSize                      int                          `hcl:"font_size" json:"font-size,omitempty"`
	FontSmoothing                 string                       `hcl:"font_smoothing" json:"font-smoothing,omitempty"`
	ForegroundColor               string                       `hcl:"foreground_color" json:"foreground-color,omitempty"`
	HomeKeysScroll                bool                         `hcl:"home_keys_scroll" json:"home-keys-scroll,omitempty"`
	Keybindings                   map[string]string            `hcl:"keybindings" json:"keybindings,omitempty"`
	MaxStringSequence             int                          `hcl:"max_string_sequence" json:"max-string-sequence,omitempty"`
	MediaKeysAreFkeys             bool                         `hcl:"media_keys_are_fkeys" json:"media-keys-are-fkeys,omitempty"`
	MetaSendsEscape               bool                         `hcl:"meta_sends_escape" json:"meta-sends-escape,omitempty"`
	MousePasteButton              *int                         `hcl:"mouse_paste_button" json:"mouse-paste-button,omitempty"`
	PageKeysScroll                bool                         `hcl:"page_keys_scroll" json:"page-keys-scroll,omitempty"`
	PassAltNumber                 *bool                        `hcl:"pass_alt_number" json:"pass-alt-number,omitempty"`
	PassCtrlNumber                *bool                        `hcl:"pass_ctrl_number" json:"pass-ctrl-number,omitempty"`
	PassMetaNumber                *bool                        `hcl:"pass_meta_number" json:"pass-meta-number,omitempty"`
	PassMetaV                     bool                         `hcl:"pass_meta_v" json:"pass-meta-v,omitempty"`
	ReceiveEncoding               string                       `hcl:"receive_encoding" json:"receive-encoding,omitempty"`
	ScrollOnKeystroke             bool                         `hcl:"scroll_on_keystroke" json:"scroll-on-keystroke,omitempty"`
	ScrollOnOutput                bool                         `hcl:"scroll_on_output" json:"scroll-on-output,omitempty"`
	ScrollbarVisible              bool                         `hcl:"scrollbar_visible" json:"scrollbar-visible,omitempty"`
	ScrollWheelMoveMultiplier     int                          `hcl:"scroll_wheel_move_multiplier" json:"scroll-wheel-move-multiplier,omitempty"`
	SendEncoding                  string                       `hcl:"send_encoding" json:"send-encoding,omitempty"`
	ShiftInsertPaste              bool                         `hcl:"shift_insert_paste" json:"shift-insert-paste,omitempty"`
	UserCSS                       string                       `hcl:"user_css" json:"user-css,omitempty"`
}

//RedisOptions - Redis 사용 여부 및 Redis 서버 연결 정보 구조체
type RedisOptions struct {
	UseRedisTokenCache string `hcl:"use_redis_token_cache" flagName:"use_redis_token_cache" flagDescribe:"if true,will use redis cache token;if false,will use memory cache token." default:"false"`
	// host:port address.
	Addr string `hcl:"redis_addr" flagName:"redis_addr" flagDescribe:"redis conntect host:port address." default:""`
	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string `hcl:"redis_user" flagName:"redis_user" flagDescribe:"redis conntect user name." default:""`
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string `hcl:"redis_vpassword" flagName:"redis_password" flagDescribe:"redis conntect user password." default:""`

	// Database to be selected after connecting to the server.
	DB int `hcl:"redis_db" flagName:"redis_db" flagDescribe:"redis database to be selected after connecting to the server.." default:"0"`

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `hcl:"redis_pool_size" flagName:"redis_pool_size" flagDescribe:"redis timeout for socket writes." default:"1"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `hcl:"redis_min_idle_conns" flagName:"redis_min_idle_conns" flagDescribe:"redis minimum number of idle connections which is useful when establishing." default:"1"`
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration `hcl:"redis_max_conn_age" flagName:"redis_max_conn_age" flagDescribe:"redis connection age at which client retires (closes) the connection" default:"0"`
}

// Convert - Redis 옵션을 redis-go 옵션으로 변경
func (ro *RedisOptions) Convert() *redis.Options {
	return &redis.Options{
		Addr:         ro.Addr,
		Username:     ro.Username,
		Password:     ro.Password,
		DB:           ro.DB,
		PoolSize:     ro.PoolSize,
		MinIdleConns: ro.MinIdleConns,
		MaxConnAge:   ro.MaxConnAge,
	}
}

// Validate - Redis 옵션 검증
func (ro *RedisOptions) Validate() error {
	if ro.UseRedisTokenCache == "true" && ro.Addr == "" {
		return errors.New("redis addr must not be null")
	}
	return nil
}
