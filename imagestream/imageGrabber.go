package imagestream

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
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/mgutz/str"
)

type ImageReceivedFunc func(img image.Image)

type ImageGrapper struct {
	cmd                *exec.Cmd
	stdout             *io.Reader
	imageReceivedFunc  ImageReceivedFunc
	isRunning          bool
	enableDebugLogging bool
	strCmd             string
	cmdArgs            []string
}

func NewImageGrabber(cmd string, enableDebugLogging bool, imageReceivedHandler ImageReceivedFunc) *ImageGrapper {
	img := ImageGrapper{}
	img.imageReceivedFunc = imageReceivedHandler
	img.enableDebugLogging = enableDebugLogging
	tmp := str.ToArgv(cmd)

	if len(tmp) > 1 {
		img.cmdArgs = tmp[1:]
	}
	img.strCmd = tmp[0]

	return &img
}

func (c *ImageGrapper) handleReadError(err error) {
	log.Println("stopping image grabber error ", err)
	if c.Running() {
		err := c.Stop()
		if err != nil {
			log.Println("error stopping image grabber error ", err)
		}
	}
	c.isRunning = false
}

func (c *ImageGrapper) readImageStream(strm io.Reader) {
	for {
		img, err := jpeg.Decode(strm)
		if err != nil {
			c.handleReadError(err)
			return
		}
		c.imageReceivedFunc(img)
	}
}

func (c *ImageGrapper) Start() error {

	c.cmd = exec.Command(c.strCmd, c.cmdArgs...)
	o, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if c.enableDebugLogging {
		c.cmd.Stderr = os.Stdout
	}

	err = c.cmd.Start()
	if err != nil {
		return err
	}
	c.isRunning = true
	go c.readImageStream(o)
	return nil
}

func (c *ImageGrapper) Stop() error {
	err := c.cmd.Process.Kill()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = c.cmd.Wait()
	c.isRunning = false
	return err
}

func (c *ImageGrapper) Running() bool {
	return c.isRunning
}
