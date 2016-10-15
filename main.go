package main

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/facebookgo/inject"
	"github.com/koding/multiconfig"
)

//TODO make config general by using a map so we can get config from ENV,file or flag.
type ServerConfig struct {
	WebPort string `default:"80"`
	WebRoot string `default:"public/dist"`
}

type Startable interface {
	Start()
}

func main() {
	config := &ServerConfig{}
	m := loadMultiConfig()
	m.MustLoad(config)

	services := make([]interface{}, 0)
	ws := &webserver{}

	// Register the rest of the services
	services = append(
		services,
		&TempSensors{},
		ws,
	)

	err := inject.Populate(services...)
	if err != nil {
		panic(err)
	}

	saveConfigToFile(config)

	StartServices(services)

	select {}
}

func StartServices(services []interface{}) {
	for _, s := range services {
		if s, ok := s.(Startable); ok {
			s.Start()
		}
	}
}

func saveConfigToFile(config *ServerConfig) {
	configFile, err := os.Create("config.json")
	if err != nil {
		logrus.Error("creating config file", err.Error())
	}

	logrus.Info("Save config: ", config)
	var out bytes.Buffer
	b, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		logrus.Error("error marshal json", err)
	}
	json.Indent(&out, b, "", "\t")
	out.WriteTo(configFile)
}

func loadMultiConfig() *multiconfig.DefaultLoader {
	loaders := []multiconfig.Loader{}

	// Read default values defined via tag fields "default"
	loaders = append(loaders, &multiconfig.TagLoader{})

	if _, err := os.Stat("config.json"); err == nil {
		loaders = append(loaders, &multiconfig.JSONLoader{Path: "config.json"})
	}

	e := &multiconfig.EnvironmentLoader{}
	e.Prefix = "STAMPZILLA"
	f := &multiconfig.FlagLoader{}
	f.EnvPrefix = "STAMPZILLA"

	loaders = append(loaders, e, f)
	loader := multiconfig.MultiLoader(loaders...)

	d := &multiconfig.DefaultLoader{}
	d.Loader = loader
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})
	return d

}
