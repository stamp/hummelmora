package main

import (
	"fmt"
	"time"

	"github.com/yryz/ds18b20"
)

type TempSensors struct {
	sensors map[string]float64
}

func (self *TempSensors) Start() {
	self.sensors = make(map[string]float64)
	go self.Worker()
}

func (self *TempSensors) Worker() {
	for {
		self.ReadSensors()
		<-time.After(time.Second * 60)
	}
}

func (self *TempSensors) Get() map[string]float64 {
	return self.sensors
}

func (self *TempSensors) ReadSensors() {
	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	fmt.Printf("sensor IDs: %v\n", sensors)

	for _, sensor := range sensors {
		t, err := ds18b20.Temperature(sensor)
		if err == nil {
			self.sensors[sensor] = t
			//			fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, t)
		}
	}
}
