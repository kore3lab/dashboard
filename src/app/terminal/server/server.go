package server

import (
	"log"
	"net/http"
	"regexp"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/kore3lab/dashboard/terminal/cache/token"
	"github.com/kore3lab/dashboard/terminal/webtty"
	"github.com/pkg/errors"
)

// Server - API 및 Webtty 제공을 위한 서버 구조체
type Server struct {
	// Slave 생성을 위한 Factory
	Factory Factory
	Options *Options

	Upgrader *websocket.Upgrader
	Cache    token.Cache
}

// New - 서버 인스턴스 생성
// 각 요청을 처리하기 위해 factory의 New() 메서드 사용
func New(factory Factory, options *Options, redisOptions *RedisOptions) (*Server, error) {
	// Request 정보에서 Web socket 관련 orgin 정보 검증
	var originChecker func(r *http.Request) bool
	if options.WSOrigin != "" {
		matcher, err := regexp.Compile(options.WSOrigin)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compile regular expression of Websocket Origin: %s", options.WSOrigin)
		}
		originChecker = func(r *http.Request) bool {
			return matcher.MatchString(r.Header.Get("Origin"))
		}
	} else {
		// Don't check origin if ws-origin is not set
		originChecker = func(r *http.Request) bool {
			return true
		}
	}

	// Token Cache 처리
	var cache token.Cache
	if redisOptions.UseRedisTokenCache == "true" {
		log.Println("use redis store token")
		client := redis.NewClient(redisOptions.Convert())
		cache = token.NewRedisCache(client, "k3webterminal-")
	} else {
		cache = token.NewMemCache()
	}

	// 서버 정보를 설정한 구조체 반환
	return &Server{
		Factory: factory,
		Options: options,

		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  webtty.MaxBufferSize,
			WriteBufferSize: webtty.MaxBufferSize,
			Subprotocols:    webtty.Protocols,
			CheckOrigin:     originChecker,
		},
		Cache: cache,
	}, nil
}
