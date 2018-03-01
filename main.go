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
	"fmt"
	"os"
	"sync"

	"github.com/rhysbryant/goprinterrelay/davinciprinter"
	"github.com/rhysbryant/goprinterrelay/httphandlers"
	"github.com/rhysbryant/goprinterrelay/transport"
)

const (
	CONN_TYPE = "tcp"
)

var printerConnection transport.PrinterConnection
var DaVinciPrinter *davinciprinter.DaVinciV3Relay
var queryCache *davinciprinter.QueryFieldsCacheMem
var printerStatus *StatusQuery
var config *Config
var serialConLock = sync.Mutex{}
var tools []Tool
var applicationInfo *ApplicationInfo
var imageStreamWebsocketsMgr = httphandlers.NewWsConnectionMgr()
var AppVersion string

func connectToPrinter() error {
	serialConLock.Lock()

	var err error
	r, w, err := printerConnection.Connect()
	if err != nil {
		serialConLock.Unlock()
		fmt.Println(err)
		return err
	}
	DaVinciPrinter = davinciprinter.NewDaVinciV3Relay(w, r, queryCache)
	return nil
}

func disconnectFromPrinter() {
	printerConnection.Disconnect()
	serialConLock.Unlock()
}

func getApplicationInfo(config *Config) *ApplicationInfo {
	appInfo := ApplicationInfo{}
	appInfo.FeatureConfig.Camera.AutoStart = config.ImageStream.AutoStart
	appInfo.Version = AppVersion

	return &appInfo
}

func startApplication() {
	var err error

	config, err = loadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	queryCache = davinciprinter.NewQueryFieldsCache(config.Printer.RelayQueryOverrides)

	applicationInfo = getApplicationInfo(config)

	printerConnection = transport.GetConnection(config.Printer.DevicePath)

	tools, err = loadTools("tools")
	if err != nil {
		fmt.Println("Error getting tools", err)
		os.Exit(3)
	}

	if config.ImageStream.ImageSourceCmd != "" {

		CreateImageStreamer(config.ImageStream.ImageSourceCmd, imageStreamWebsocketsMgr, config.ImageStream.EnableDebugLogging)
	}

	go startHttpServer(config.WebUI.Httplistener)
	go startDavinciTcpListener(config.Printer.RelayTCPListener)
	updateStatusLoop()
}

func main() {
	if AppVersion == "" {
		AppVersion = "dev-build"
	}

	if len(os.Args) > 1 {
		err := handleServiceCommand(os.Args[1])
		if err != nil {
			fmt.Printf("error %s", err.Error())
			os.Exit(5)
		}

	} else {
		startApplication()
	}
}
