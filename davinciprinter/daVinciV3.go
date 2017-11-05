package davinciprinter

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

	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

const (
	NewLineChar              = '\n'
	EndOfResponse            = "$"
	QueryResponseKeyValueSep = ":"
	QueryTypeAll             = "a"

	CommandTypeQuery  = "XYZv3/query"
	CommandTypeConfig = "XYZv3/config"
	CommandTypeAction = "XYZv3/action"
	CommandTypeUpload = "XYZv3/upload"
)

type DaVinciV3Relay struct {
	buffer      *bufio.ReadWriter
	queryFields QueryFieldsCache
	lock        sync.Mutex
}

func NewDaVinciV3Relay(w io.Writer, r io.Reader, qc QueryFieldsCache) *DaVinciV3Relay {
	d := DaVinciV3Relay{}
	d.queryFields = qc
	d.buffer = bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))

	return &d
}

func (d *DaVinciV3Relay) SendQueryResponse(strType string) error {
	if strType == QueryTypeAll {
		for k, v := range d.queryFields.GetAllFields() {
			fmt.Fprintf(d.buffer, "%s:%s\n", k, v)
		}
		fmt.Fprintf(d.buffer, "$\n")
	} else {
		fmt.Printf("sending [%s]", strType)
		if val, exists := d.queryFields.GetField(strType); exists {
			fmt.Fprintf(d.buffer, "%s:%s\n", strType, val)
		}
	}
	d.buffer.Flush()
	return nil
}

func (d *DaVinciV3Relay) RefreshStatus() (error, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	err := d.SendQueryRequest(QueryTypeAll)
	if err != nil {
		return err, false
	}

	return d.reciveQueryResponse()

}

func (d *DaVinciV3Relay) SendQueryRequest(queryType string) error {
	fmt.Fprintf(d.buffer, "%s=%s\n", CommandTypeQuery, queryType)
	return d.buffer.Flush()
}

func parseKeyValueLine(strLine string) (*string, *string) {
	sepIndex := strings.Index(strLine, QueryResponseKeyValueSep)
	if sepIndex == -1 {
		return nil, nil
	}

	key := strLine[:sepIndex]
	value := strLine[sepIndex+1:]

	return &key, &value
}

func (d *DaVinciV3Relay) reciveQueryResponse() (error, bool) {
	var hasChanged bool
	for {
		str, err := d.buffer.ReadString(NewLineChar)
		if strings.HasSuffix(str, "\n") {
			str = str[0 : len(str)-1]
		}
		fmt.Println(str, err)
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err, false
		}
		if str == EndOfResponse {
			break
		}

		key, value := parseKeyValueLine(str)
		if key == nil {
			return errors.New("unable to parse query"), false
		}

		if hc := d.queryFields.SetField(*key, *value); hc {
			hasChanged = true
		}
	}
	return nil, hasChanged
}
func (d *DaVinciV3Relay) checkIsOk() error {
	str, err := d.buffer.ReadString('\n')
	if err != nil {
		fmt.Errorf("error reading response %s\n", err.Error())
		return err
	}

	if !strings.HasPrefix(str, "ok") {
		return errors.New("expected ok got " + str)
	}
	return nil
}

func (d *DaVinciV3Relay) Upload(dataSrc UploadDataReader, okCallback UploadReadChunkComfirmedFunc, length int64) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	h := NewDaVinciV3Upload(nil, d.buffer, length)
	fmt.Fprintf(d.buffer, "XYZv3/upload=MyTest.gcode,%d\n", length)
	var err error
	err = d.buffer.Flush()
	if err != nil {
		fmt.Errorf("error flushing chunk %s\n", err.Error())
		return err
	}
	err = d.checkIsOk()
	if err != nil {
		fmt.Errorf("error response %s\n", err)
	} else {
		fmt.Println("got ok")
		//okCallback()
	}

	for h.HasNextChunk() {

		payload, err := dataSrc.GetNextChunk()
		if err != nil {
			fmt.Errorf("error reading chunk %s\n", err.Error())
			return err
		}

		err = h.sendChunk(payload)
		if err != nil {
			fmt.Errorf("error sending chunk %s\n", err.Error())
			return err
		}
		err = d.buffer.Flush()
		if err != nil {
			fmt.Errorf("error flushing chunk %s\n", err.Error())
			return err
		}
		err = d.checkIsOk()
		if err != nil {
			fmt.Errorf("error response %s\n", err)
		} else {
			fmt.Println("got ok")
			okCallback()
		}
		fmt.Println("chunk Sent")

	}
	fmt.Fprintf(d.buffer, "XYZv3/uploadDidFinish")
	err = d.buffer.Flush()
	if err != nil {
		fmt.Errorf("error flushing chunk %s\n", err.Error())
		return err
	} else {
		fmt.Println("all done!")
	}
	return nil
}

func (d *DaVinciV3Relay) SendRaw(str string) (string, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.buffer.WriteString(str)
	err := d.buffer.Flush()
	if err != nil {
		return "", err
	}
	return d.buffer.ReadString(NewLineChar)
}
