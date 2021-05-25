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
	Realm                     = "Kore-Board"
	StrategyCookie            = "cookie"
	StrategyLocal             = "local"
	SecretStaticUser          = "static-user"
	SecretBasicAuth           = "basic-auth"
	SecretStaticToken         = "static-token"
	SecretServiceAccountToken = "service-account-token"
)

type AuthConfig struct {
	Strategy   string //nuxt-auth strategy
	Secret     string //static-user, static-token , service-account-token
	AccessKey  string //local access-token-key
	RefreshKey string //local refresh-token-key
	Data       map[string]string
}

// auth scheme (user, token)
func (me *AuthConfig) GetSchema() string {

	schema := "user"
	if me.Secret == "" {
		schema = ""
	} else if strings.Contains(me.Secret, "token") {
		schema = "token"
	}

	return schema

}

type SecretProvider func(user, realm string) string
type ValidateFunc func(map[string]string) error
