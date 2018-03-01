package transport

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
	"io"

	"github.com/tarm/serial"
)

type serialCon struct {
	connectionStr string
	port          *serial.Port
}

func (s *serialCon) Connect() (io.Reader, io.Writer, error) {
	port, err := serial.OpenPort(&serial.Config{Name: s.connectionStr, Baud: 115200})
	if err != nil {
		return nil, nil, err
	}
	s.port = port

	return port, port, nil
}

func (s *serialCon) Disconnect() {
	s.port.Close()
}
