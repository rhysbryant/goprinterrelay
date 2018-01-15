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
	"time"
)

func updateStatus() {
	err := connectToPrinter()
	if err != nil {
		fmt.Println("unable to connect to printer ", err)
		time.Sleep(time.Second * 30)
		return
	}

	err, valuesChanged := DaVinciPrinter.RefreshStatus()
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
	if valuesChanged {
		SendStatusUpdate()
	}

	disconnectFromPrinter()
	time.Sleep(time.Second * 10)
}

func updateStatusLoop() {
	for {
		updateStatus()
	}
}
