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
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rhysbryant/goprinterrelay/davinciprinter"
	"github.com/rhysbryant/goprinterrelay/httphandlers"
	"github.com/tarm/serial"
)

const (
	CONN_TYPE = "tcp"
)

var port *serial.Port
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
	fmt.Printf("%+v", config)
	var err error
	port, err = serial.OpenPort(&serial.Config{Name: config.Printer.DevicePath, Baud: 115200})
	if err != nil {
		fmt.Println(err)
		return err
	}
	DaVinciPrinter = davinciprinter.NewDaVinciV3Relay(port, port, queryCache)
	return nil
}

func disconnectFromPrinter() {
	port.Close()
	serialConLock.Unlock()
}

func getApplicationInfo(config *Config) *ApplicationInfo {
	appInfo := ApplicationInfo{}
	appInfo.FeatureConfig.Camera.AutoStart = config.ImageStream.AutoStart
	appInfo.Version = AppVersion

	return &appInfo
}

func main() {
	if AppVersion == "" {
		AppVersion = "dev-build"
	}

	var err error
	config, err = loadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	queryCache = davinciprinter.NewQueryFieldsCache(config.Printer.RelayQueryOverrides)
	l, err := net.Listen(CONN_TYPE, config.Printer.RelayTCPListener)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(2)
	}

	applicationInfo = getApplicationInfo(config)

	tools, err = loadTools("tools")
	if err != nil {
		fmt.Println("Error getting tools", err)
		os.Exit(3)
	}

	if config.ImageStream.ImageSourceCmd != "" {

		CreateImageStreamer(config.ImageStream.ImageSourceCmd, imageStreamWebsocketsMgr, config.ImageStream.EnableDebugLogging)
	}

	go func() {
		for {
			err = connectToPrinter()
			if err != nil {
				fmt.Println("unable to connect to printer ", err)
				time.Sleep(time.Second * 30)
				continue
			}

			err := DaVinciPrinter.RefreshStatus()
			if err != nil {
				fmt.Printf("getStatus error %s\n", err.Error())
			}

			status, err := getStatusFromMap(queryCache.GetAllFields())
			if err != nil {
				fmt.Println(err)

			} else {
				printerStatus = status
				fmt.Printf("%+v", status)
			}

			disconnectFromPrinter()
			time.Sleep(time.Second * 15)
		}
	}()

	go startHttpServer(config.WebUI.Httplistener)

	defer l.Close()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buffer := bufio.NewReader(conn)
	daVinciTcpIpCon := davinciprinter.NewDaVinciV3Relay(conn, conn, queryCache)

	for {
		str, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println(err)
			return
		}
		fmt.Println(str, err)
		if strings.HasPrefix(str, davinciprinter.CommandTypeQuery) {
			queryType := str[len(davinciprinter.CommandTypeQuery)+1 : len(davinciprinter.CommandTypeQuery)+2]

			daVinciTcpIpCon.SendQueryResponse(queryType)

		} else if strings.HasPrefix(str, davinciprinter.CommandTypeUpload) {

			strUpload := str[len(davinciprinter.CommandTypeUpload)+1 : len(str)-2]
			fields := strings.Split(strUpload, ",")
			fmt.Println(strUpload, fields)

			length, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Println(err)
				return
			}

			fileName := fields[0]
			fmt.Println(length, err, fileName)
			fmt.Fprintln(conn, "ok")
			err = connectToPrinter()
			if err != nil {
				fmt.Println("unable to connect to printer ", err)
				return
			}

			uploadHandler := davinciprinter.NewDaVinciV3Upload(conn, nil, int64(length))

			err = DaVinciPrinter.Upload(uploadHandler, func() { fmt.Fprintln(conn, "ok") }, int64(length))
			if err != nil {
				fmt.Printf("Upload returned error [%s]", err.Error())
			}

			disconnectFromPrinter()

		} else if strings.HasPrefix(str, davinciprinter.CommandTypeConfig) || strings.HasPrefix(str, davinciprinter.CommandTypeAction) {
			err = connectToPrinter()
			if err != nil {
				fmt.Println("unable to connect to printer ", err)
				return
			}
			fmt.Printf("got: [%s]", str)
			res, err := DaVinciPrinter.SendRaw(str)
			if err != nil {
				log.Println(err)
				break

			}
			fmt.Fprintf(conn, "%s", res)
			fmt.Println(res)
			disconnectFromPrinter()
		}

	}

}
