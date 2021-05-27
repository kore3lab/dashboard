package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

/**
<Authenticator>
	CookieAuthenticator
	LocalAuthenticator
	BasicAuthAuthenticator (미사용)
*/
type Authenticator struct {
	Realm          string
	HandlerFunc    func() gin.HandlerFunc
	Validate       func(map[string]string) error
	LoginHandler   func(body map[string]string) (interface{}, error)
	RefreshHandler func(body map[string]string) (interface{}, error)
	LogoutHandler  func(body map[string]string)
}

func CreateAuthenticator(conf *AuthConfig, c *rest.Config) (*Authenticator, error) {

	//validator
	validator, err := getValidateFunc(conf, c)
	if err != nil {
		return nil, err
	}

	var authenticator *Authenticator
	//authentocator
	if conf.Strategy == StrategyCookie {
		authenticator = CookieAuthenticator(validator)
	} else if conf.Strategy == StrategyLocal {
		authenticator = LocalAuthenticator(conf.AccessKey, conf.RefreshKey, validator)
	} else {
		return nil, fmt.Errorf("not supported '%s' strategy yet", conf.Strategy)
	}

	return authenticator, nil

}

func DummyAuthenticator() *Authenticator {

	h := &Authenticator{}
	h.Validate = func(map[string]string) error {
		return nil
	}

	h.HandlerFunc = func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return h

}

func CookieAuthenticator(validateFunc ValidateFunc) *Authenticator {

	h := &Authenticator{}
	h.Validate = validateFunc
	h.HandlerFunc = func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return h

}

func LocalAuthenticator(accessKey string, refreshKey string, validateFunc ValidateFunc) *Authenticator {

	h := &Authenticator{}
	h.Validate = validateFunc
	h.HandlerFunc = func() gin.HandlerFunc {

		return func(c *gin.Context) {

			postfix, _ := c.Cookie("auth.strategy")
			accessToken, err := c.Cookie(fmt.Sprintf("auth._token.%s", postfix))
			if err != nil {
				log.Warnf("prasing token (%s) failed  (cause=%s)", accessToken, err.Error())
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			if expired, err := ValidateSessionToken(accessKey, accessToken); err != nil {
				log.Warnf("validate token (%s) failed  (cause=%s, expired=%v)", accessToken, err.Error(), expired)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			} else {
				if expired {
					log.Warnf("expired=%s", err.Error(), expired)
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			}
			c.Next()
		}
	}

	//login, refresh, logout callback
	h.LoginHandler = func(params map[string]string) (interface{}, error) {
		return newJWTToken(accessKey, refreshKey)
	}
	h.RefreshHandler = func(params map[string]string) (interface{}, error) {
		// validating refresh-token
		if expired, err := ValidateSessionToken(refreshKey, params["refreshToken"]); err != nil {
			return nil, fmt.Errorf("invalid refresh token (cause=%s)", err.Error())
		} else if expired {
			return nil, errors.New("refresh token expired")
		} else {
			// new access, refresh token
			return newJWTToken(refreshKey, refreshKey)
		}
	}

	return h

}

func newJWTToken(accessSecret string, refreshSecret string) (map[string]string, error) {

	token, err := GenerateSessionToken(accessSecret, 60*15)
	if err != nil {
		return nil, errors.New("can't genrated a access-token")
	}
	refreshToken, err := GenerateSessionToken(refreshSecret, 60*60*24*7)
	if err != nil {
		return nil, errors.New("can't genrated a refresh-token")
	}
	return map[string]string{
		"token":        token,
		"refreshToken": refreshToken,
	}, nil

}

func BasicAuthAuthenticator(filename string, validateFunc ValidateFunc) *Authenticator {

	h := &Authenticator{}
	h.Validate = validateFunc
	h.HandlerFunc = func() gin.HandlerFunc {

		return func(c *gin.Context) {
			// Get the Basic Authentication credentials
			user, password, ok := c.Request.BasicAuth()
			if ok {
				err := validateFunc(map[string]string{"username": user, "password": password})
				ok = (err == nil)
			}
			if !ok {
				c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Next()
		}
	}

	return h

}

func getValidateFunc(conf *AuthConfig, c *rest.Config) (ValidateFunc, error) {

	var secret SecretProvider

	// choice provider
	ty := conf.Secret
	if ty == SecretBasicAuth {
		secret = UserFileSecretProvider(conf.Data["dir"])
	} else if ty == SecretStaticUser {
		secret = StaticUserSecretProvider(conf.Data["username"], conf.Data["password"])
	} else if ty == SecretStaticToken {
		secret = StaticTokenSecretProvider(conf.Data["token"])
	} else if ty == SecretServiceAccountToken {
		if c == nil {
			return nil, fmt.Errorf("can't initialized becuase connection cluster is empty")
		} else {
			secret = ServiceAccountTokenSecretProvider(c)
		}

	} else if ty == "" {
	} else {
		return nil, fmt.Errorf("cannot found '%s' secret provider", ty)
	}

	schema := conf.GetSchema()
	if schema == "user" {
		//username, password
		return func(params map[string]string) error {
			if params["username"] == "" {
				return errors.New("username is empty")
			}
			if secret(params["username"], Realm) != params["password"] {
				return errors.New("invalid password")
			} else {
				return nil
			}
		}, nil

	} else if schema == "token" {
		//token
		return func(params map[string]string) error {
			if params["token"] == "" {
				return errors.New("token is empty")
			}
			if secret(params["token"], Realm) != params["token"] {
				return errors.New("invalid token")
			} else {
				return nil
			}
		}, nil

	} else {
		return func(params map[string]string) error {
			return nil
		}, nil

	}

}
