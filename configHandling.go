package main

/**
	This file is part of goPrinterRelay.

	goPrinterRelay - printer status page and protocol relay for daVinci jr 3d printers

    goPrinterRelay is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    goPrinterRelay is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with goPrinterRelay.  If not, see <http://www.gnu.org/licenses/>.

**/
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	WebUI struct {
		Httplistener string `json:"httpListener"`
	} `json:"webui"`

	Printer struct {
		DevicePath          string            `json:"devicePath"`
		RelayTCPListener    string            `json:"relayTcpListener"`
		RelayQueryOverrides map[string]string `json:"queryOverrides"`
	} `json:"printer"`
	ImageStream struct {
		ImageSourceCmd     string `json:"imageSourceCmd"`
		EnableDebugLogging bool   `json:"enableDebugLogging"`
		AutoStart          bool   `json:"autoStart"`
	} `json:"imagestream"`
}

func loadConfig(fileName string) (*Config, error) {
	config := Config{}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	e := json.Unmarshal(file, &config)
	if e != nil {
		return nil, e
	}
	for v, k := range config.Printer.RelayQueryOverrides {
		fmt.Println(v, string(k))
	}
	return &config, nil
}
