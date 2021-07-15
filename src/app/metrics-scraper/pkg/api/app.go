package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type App struct {
	Writer http.ResponseWriter
}

func (self *App) Send(v interface{}) {

	self.Writer.Header().Add("Content-Type", "application/json")

	j, err := json.Marshal(v)
	if err != nil {
		self.SendMessage(http.StatusInternalServerError, err.Error())
	} else {
		_, err = self.Writer.Write(j)
		if err != nil {
			log.Errorf("Error cannot write response: %v", err)
		}
	}

}

func (self *App) SendMessage(httpCode int, msg string) {

	self.Writer.Header().Add("Content-Type", "application/json")
	self.Writer.WriteHeader(httpCode)
	_, err := self.Writer.Write([]byte(fmt.Sprintf("{'message':'%s'}", msg)))
	if err != nil {
		log.Errorf("Error cannot write response: %v", err)
	}

}
