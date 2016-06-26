package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type webserver struct {
}

func (self *webserver) Start() {
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
		m.Broadcast(msg)
	})

	go func() {
		http.ListenAndServe(":80", r)
	}()
}
