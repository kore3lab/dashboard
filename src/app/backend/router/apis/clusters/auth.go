package apis

import (
	"net/http"
	"time"

	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/pkg/auth"
	log "github.com/sirupsen/logrus"

	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/gin-gonic/gin"
)

type user struct {
	Username string `json:"username"`
}

// get auth info
func GetAuth(c *gin.Context) {
	g := app.Gin{C: c}

	g.Send(http.StatusOK, map[string]string{
		"strategy": config.Value.AuthConfig.Strategy,
		"scheme":   config.Value.AuthConfig.GetScheme(),
		"provider": config.Value.AuthConfig.Secret["type"],
	})

}

// sing-in validation
func Login(c *gin.Context) {
	g := app.Gin{C: c}

	body := make(map[string]string)

	// parse body params
	if g.C.BindJSON(&body) != nil {
		g.SendMessage(http.StatusBadRequest, "Unable to bind request body", nil)
		return
	}

	//validation
	if err := config.Value.Authenticator.Validate(body); err != nil {
		g.SendMessage(http.StatusUnauthorized, err.Error(), err)
	} else {
		// login evnet 실행
		if config.Value.Authenticator.LoginHandler != nil {
			if resp, err := config.Value.Authenticator.LoginHandler(body); err != nil {
				g.SendMessage(http.StatusUnauthorized, err.Error(), err)
			} else {
				g.Send(http.StatusOK, resp)
			}
		} else {
			g.SendOK()
		}
	}

}

// return issuanced refresh-token
func RefreshToken(c *gin.Context) {
	g := app.Gin{C: c}

	if config.Value.AuthConfig.Strategy != auth.StrategyLocal {
		g.Send(http.StatusNotFound, nil)
		return
	}

	body := make(map[string]string)

	// parse body params
	if g.C.BindJSON(&body) != nil {
		g.SendMessage(http.StatusBadRequest, "Unable to bind request body", nil)
		return
	}

	// on refresh
	if config.Value.Authenticator.RefreshHandler != nil {
		if resp, err := config.Value.Authenticator.RefreshHandler(body); err != nil {
			g.SendMessage(http.StatusInternalServerError, err.Error(), err)
		} else {
			g.Send(http.StatusOK, resp)
		}
	} else {
		g.SendOK()
	}

}

// return user info
func GetUser(c *gin.Context) {
	g := app.Gin{C: c}

	scheme := config.Value.AuthConfig.GetScheme()

	user := &user{}

	if scheme == "user" {
		user.Username = config.Value.AuthConfig.Secret["username"]
	} else {
		user.Username = "admin"
	}

	g.Send(http.StatusOK, map[string]interface{}{
		"user": user,
	})

}

// sign-out
func Logout(c *gin.Context) {
	g := app.Gin{C: c}

	log.Infof("logout at %s", time.Now().Format("2006-01-02 15:04:05"))
	g.SendOK()

}
