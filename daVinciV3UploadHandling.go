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
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

type DaVinciV3UploadHandler struct {
	src          io.Reader
	dst          io.Writer
	dstComfirm   io.ByteReader
	length       int64
	remaining    int64
	chunkNumSent int32
}

type DaVinciUploadChunk struct {
	offset     int32
	chunkSize  int32
	Payload    []byte
	remainding int32
}

type UploadReadChunkComfirmedFunc func()

type UploadDataReader interface {
	GetNextChunk() ([]byte, error)
}

func (d *DaVinciV3UploadHandler) GetNextChunk() ([]byte, error) {
	chunk, err := d.readChunk()
	if err != nil {
		return nil, err
	}
	fmt.Printf("read chunk:%d size %d other %d\n", chunk.offset, chunk.chunkSize, chunk.remainding)
	return chunk.Payload, nil
}

func (d *DaVinciV3UploadHandler) readChunk() (*DaVinciUploadChunk, error) {

	var l int
	var err error

	/**
	* format
	*
	* chunk number int32
	* chunksize int32
	* payload[chunksize]
	* unknown int32 always 0
	 */

	uc := DaVinciUploadChunk{}

	err = binary.Read(d.src, binary.BigEndian, &uc.offset)
	if err != nil {
		fmt.Printf("Error reading chunk num\n")
		return nil, err
	}

	err = binary.Read(d.src, binary.BigEndian, &uc.chunkSize)
	if err != nil {
		fmt.Printf("Error reading chunk size\n")
		return nil, err
	}
	uc.Payload = make([]byte, uc.chunkSize)

	offset := int32(0)
	start := time.Now()
	end := start.Add(time.Second * 10)
	for offset < uc.chunkSize && start.Before(end) {
		l, err = d.src.Read(uc.Payload[offset:])
		if err != nil {
			return nil, err
		}
		offset += int32(l)
	}

	if offset != uc.chunkSize {
		fmt.Printf("got %d bytes expected %d bytes", l, uc.chunkSize)
		return nil, errors.New("expected more data")
	}

	d.remaining -= int64(uc.chunkSize)

	err = binary.Read(d.src, binary.BigEndian, &uc.remainding)
	if err != nil {
		return nil, err
	}

	return &uc, nil
}

func (d *DaVinciV3UploadHandler) sendChunk(payload []byte) error {
	/**
	* format
	*
	* chunk number int32
	* chunksize int32
	* payload[chunksize]
	* unknown int32 always 0
	 */
	pLen := int32(len(payload))
	var err error
	var l int
	err = binary.Write(d.dst, binary.BigEndian, d.chunkNumSent)
	if err != nil {
		return err
	}
	err = binary.Write(d.dst, binary.BigEndian, pLen)
	if err != nil {
		return err
	}
	l, err = d.dst.Write(payload)
	if err != nil {
		return err
	} else if int32(l) != pLen {
		return errors.New("some data not sent")
	}

	err = binary.Write(d.dst, binary.BigEndian, int32(0))
	if err != nil {
		return err
	}
	fmt.Printf("sent chunk:%d size %d other %d\n", d.chunkNumSent, pLen, 0)
	d.chunkNumSent++
	d.remaining -= int64(pLen)

	return nil
}

func (d *DaVinciV3UploadHandler) HasNextChunk() bool {

	return d.remaining > 0
}

func NewDaVinciV3Upload(src io.Reader, dst io.Writer, length int64) *DaVinciV3UploadHandler {
	d := DaVinciV3UploadHandler{src: src, dst: dst, length: length, remaining: length}
	return &d
}
