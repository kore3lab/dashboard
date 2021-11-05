package app

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type Status struct {
	Message string `json:"message"`
}

func Response(resp *resty.Response) Status {
	if resp.StatusCode() != http.StatusOK {
		e := Status{}
		json.Unmarshal(resp.Body(), &e)
		return e
	}
	return Status{Message: ""}
}
func Error(err error) Status {
	return Status{Message: err.Error()}
}

type Gin struct {
	C *gin.Context
}

func (g *Gin) SendMessage(httpCode int, msg string, err error) {
	if err != nil {
		if err.Error() == msg {
			log.Errorln(msg)
		} else {
			log.Errorf("%s (cause=%v)", msg, err.Error())
		}

	}
	g.C.JSON(httpCode, Status{Message: msg})
	return
}
func (g *Gin) SendError(err error) {

	msg := ""
	if err != nil {
		msg = err.Error()
		log.Errorln(msg)
	}

	g.C.JSON(http.StatusInternalServerError, Status{Message: msg})
	return
}

func (g *Gin) Send(httpCode int, json interface{}) {
	g.C.JSON(httpCode, json)
	return
}
func (g *Gin) SendOK() {
	g.C.JSON(http.StatusOK, Status{Message: ""})
	return
}

//url 검사
func (g *Gin) ValidateUrl(params []string) error {

	valid := validation.Validation{}

	for _, name := range params {
		valid.Required(g.C.Param(name), name)
	}

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(fmt.Sprintf("[%s]%s", err.Key, err.Error()))
		}
	}
	return nil

}

// parse querystrings
func (g *Gin) ParseQuery() (url.Values, error) {

	u, err := url.Parse(g.C.Request.RequestURI)
	if err != nil {
		return nil, err
	}
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	return query, nil
}
