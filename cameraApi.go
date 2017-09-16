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
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"

	"github.com/gorilla/websocket"
	"github.com/rhysbryant/goprinterrelay/httphandlers"
	"github.com/rhysbryant/goprinterrelay/imagestream"
)

type StatusMessage struct {
	StreamStatus string `json:"status"`
	Message      string `json:"message"`
}

func sendErrorStatus(wsconMgr *httphandlers.WsConnectionMgr, message string) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(&StatusMessage{"Error", message})
	wsconMgr.SendToAll(websocket.TextMessage, buf.Bytes())
	wsconMgr.CloseAll(websocket.CloseTryAgainLater, message)
}

func CreateImageStreamer(cmd string, wsconMgr *httphandlers.WsConnectionMgr, streamerDebugLogging bool) *imagestream.ImageGrapper {
	var buf bytes.Buffer

	ig := imagestream.NewImageGrabber(cmd, streamerDebugLogging, func(img image.Image) {
		buf.Reset()
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			log.Println(err)
			sendErrorStatus(wsconMgr, "streaming process error "+err.Error())
			return
		}
		wsconMgr.SendToAll(websocket.BinaryMessage, buf.Bytes())
	})

	wsconMgr.SetClientConnectHandler(func() {
		fmt.Println(ig.Running())
		if !ig.Running() {
			err := ig.Start()
			if err != nil {
				log.Printf("error starting streaming process %s\n", err.Error())
				sendErrorStatus(wsconMgr, "error starting streaming process "+err.Error())
			}
		}
	})

	wsconMgr.SetClientDisconnectHandler(func() {
		if wsconMgr.ClientCount() <= 0 {
			if ig.Running() {
				ig.Stop()
			}
		}
	})

	return ig
}
