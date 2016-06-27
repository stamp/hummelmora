package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/stianeikeland/go-rpio"
)

type webserver struct {
}

func (self *webserver) Start() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}

	heat1 := rpio.Pin(22)
	heat1.Output()

	heat2 := rpio.Pin(23)
	heat2.Output()

	lights := rpio.Pin(24)
	lights.Output()

	r := gin.Default()
	m := melody.New()

	// allow all origin hosts
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	r.Use(static.Serve("/", static.LocalFile("gui/web", false)))
	r.StaticFile("/", "gui/web/index.html")

	r.GET("/socket", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		logrus.Info("WS connect")
	})
	m.HandleDisconnect(func(s *melody.Session) {
		logrus.Info("WS disconnect")
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		type Msg map[string]interface{}
		data := make(Msg)
		err := json.Unmarshal(msg, &data)
		if err != nil {
			logrus.Error("Failed to decode JSON: ", err)
			logrus.Debug(string(msg))
		}

		for key, line := range data {
			switch key {
			case "heat1":
				switch val := line.(type) {
				case bool:
					logrus.Info("Heat1=", val)
					if !val {
						heat1.High()
					} else {
						heat1.Low()
					}
				default:
					logrus.Error("Heat1")
				}
			case "heat2":
				switch val := line.(type) {
				case bool:
					logrus.Info("Heat2=", val)
					if !val {
						heat2.High()
					} else {
						heat2.Low()
					}
				default:
					logrus.Error("Heat2")
				}
			case "lights":
				logrus.Info("Pulse lights")
				lights.High()
				<-time.After(time.Millisecond * 200)
				lights.Low()
			default:
				logrus.Error("key", key)
			}
		}

		//m.Broadcast(msg)

	})

	go func() {
		http.ListenAndServe(":80", r)
	}()
}
