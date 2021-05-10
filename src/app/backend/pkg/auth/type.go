package auth

import (
	"strings"
)

/**
auth-info USE-CASES
	{"strategy":"cookie"}
	{"strategy":"cookie",	"secret": {"type": "static-user",			"username": "admin", "password": "acornsoft"} }
	{"strategy":"cookie",	"secret": {"type": "static-token",			"token": "acornsoft"} }
	{"strategy":"cookie",	"secret": {"type": "basic-auth",			"dir": "/var/tmp"} }
	{"strategy":"cookie",	"secret": {"type": "service-account-token"} }
	{"strategy":"local",	"key": {"access": "whdmstkddk", "refresh":"hsthvmxm"} }
	{"strategy":"local",	"key": {"access": "whdmstkddk", "refresh":"hsthvmxm"},	"secret": {"type": "static-user",			"username": "admin", "password": "acornsoft"} }
	{"strategy":"local",	"key": {"access": "whdmstkddk", "refresh":"hsthvmxm"},	"secret": {"type": "static-token"	,		"token": "acornsoft"} }
	{"strategy":"local",	"key": {"access": "whdmstkddk", "refresh":"hsthvmxm"},	"secret": {"type": "basic-auth",			"dir": "/var/tmp"} }
	{"strategy":"local",	"key": {"access": "whdmstkddk", "refresh":"hsthvmxm"},	"secret": {"type": "service-account-token"} }
*/
const (
	Realm                = "Kore-Board"
	DefaultAuthenticator = `{"strategy":"cookie", "secret": {"type": "static-token", "token": "acornsoft"} }`

	StrategyCookie                  = "cookie"
	StrategyLocal                   = "local"
	ProviderTypeStaticUser          = "static-user"
	ProviderTypeBasicAuth           = "basic-auth"
	ProviderTypeStaticToken         = "static-token"
	ProviderTypeServiceAccountToken = "service-account-token"
)

type AuthConfig struct {
	Strategy string            `json:"strategy"` //nuxt-auth strategy
	Key      map[string]string `json:"key"`      //if auth='local' , access&refresh token secret
	Secret   map[string]string `json:"secret"`   //user validation secret provider
}

// auth scheme (user, token)
func (me *AuthConfig) GetScheme() string {

	scheme := "user"

	if me.Secret == nil {
		scheme = ""
	} else if strings.Contains(me.Secret["type"], "token") {
		scheme = "token"
	}

	return scheme

}

type SecretProvider func(user, realm string) string
type ValidateFunc func(map[string]string) error
