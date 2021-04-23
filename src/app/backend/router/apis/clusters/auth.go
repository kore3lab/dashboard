package apis

import (
	"net/http"
	"os"

	"github.com/acornsoftlab/dashboard/model/v1alpha1"
	"github.com/acornsoftlab/dashboard/pkg/app"
	"github.com/acornsoftlab/dashboard/pkg/config"
	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	g := app.Gin{C: c}

	req := map[string]string{}

	if g.C.BindJSON(&req) != nil {
		g.SendMessage(http.StatusInternalServerError, "Unable to bind request body", nil)
	} else {
		if os.Getenv("APP_ENV") != "production" || req["secret"] == config.Value.SecretToken {
			if err := model.CreateToken(c); err != nil {
				g.SendMessage(http.StatusUnauthorized, err.Error(), err)
			} else {
				g.SendOK()
			}
		} else {
			g.SendMessage(http.StatusUnauthorized, "invalidate access token", nil)
		}
	}

}

func Logout(c *gin.Context) {
	g := app.Gin{C: c}

	model.CleanToken(c)
	g.SendOK()

}
