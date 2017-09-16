package httphandlers

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
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ClientDisconnectedFunc func()
type ClientConnectedFunc func()

type WsConnectionMgr struct {
	clients                   map[*websocket.Conn]bool
	clientCount               int
	clientConnectedHandler    ClientConnectedFunc
	clientDisconnectedHandler ClientDisconnectedFunc
}

func NewWsConnectionMgr() *WsConnectionMgr {
	wsm := WsConnectionMgr{}
	wsm.clients = make(map[*websocket.Conn]bool, 1)

	return &wsm
}

func (w *WsConnectionMgr) addConnection(conn *websocket.Conn) {
	w.clients[conn] = true
	w.clientCount++
	if w.clientConnectedHandler != nil {
		w.clientConnectedHandler()
	}
}

func (w *WsConnectionMgr) removeConnection(conn *websocket.Conn) {
	delete(w.clients, conn)
	w.clientCount--
	if w.clientDisconnectedHandler != nil {
		w.clientDisconnectedHandler()
	}
}

func (m *WsConnectionMgr) WsUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		//http.Error(w, "Origin not allowed", 403)
		//return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	m.addConnection(conn)
}

func (w *WsConnectionMgr) SendToAll(msgType int, payload []byte) error {
	pm, err := websocket.NewPreparedMessage(msgType, payload)
	if err != nil {
		return err
	}

	for c, _ := range w.clients {
		err := c.WritePreparedMessage(pm)
		if err != nil {
			w.removeConnection(c)
		}
	}

	return nil

}

func (w *WsConnectionMgr) ClientCount() int {
	return w.clientCount
}

func (w *WsConnectionMgr) SetClientDisconnectHandler(f ClientDisconnectedFunc) {
	w.clientDisconnectedHandler = f
}
func (w *WsConnectionMgr) SetClientConnectHandler(f ClientConnectedFunc) {
	w.clientConnectedHandler = f
}

func (w *WsConnectionMgr) CloseAll(codeCode int, message string) {
	cm := websocket.FormatCloseMessage(codeCode, message)

	for c, _ := range w.clients {
		c.WriteControl(websocket.CloseMessage, cm, time.Now().Add(1*time.Second))
		w.removeConnection(c)
	}
}
