// File: "conf.go"

package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config file name and file mode
const CONFIG = "showlog.json"
const CONFIG_MODE = 0644

// Configuration scructure
type Conf struct {
	File   string  `json:"file,omitempty"`   // log file
	Center bool    `json:"center,omitempty"` // center on screen
	Full   bool    `json:"full,omitempty"`   // full screen
	W      float32 `json:"w,omitempty"`      // width
	H      float32 `json:"h,omitempty"`      // height
}

// Make default config
func NewConf() *Conf {
	return &Conf{
		Center: true,
		W:      640,
		H:      480,
	}
}

// Read config from file
func (c *Conf) Read() {
	data, err := os.ReadFile(CONFIG)
	if err != nil {
		log.Print(err)
		return
	}

	var C Conf
	err = json.Unmarshal(data, &C)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("read conf: File='%s' Center=%v Full=%v W=%f H=%f",
		C.File, C.Center, C.Full, C.W, C.H)

	if C.Center {
		c.Center = true
	}
	if C.Full {
		c.Full = true
	}
	if C.W > 0 {
		c.W = C.W
	}
	if C.H > 0 {
		c.H = C.H
	}
}

// Write config to file
func (c *Conf) Write() {
	data, err := json.Marshal(c)
	if err != nil {
		log.Print(err)
		return
	}

	err = os.WriteFile(CONFIG, data, CONFIG_MODE)
	if err != nil {
		log.Print(err)
	}

	log.Printf("write conf: File='%s' Center=%v Full=%v W=%f H=%f",
		c.File, c.Center, c.Full, c.W, c.H)
}

// EOF: "conf.go"
