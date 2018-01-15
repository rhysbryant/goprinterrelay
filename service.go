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
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

const (
	CmdRunAsService     = "-service"
	CmdInstallService   = "-installSvc"
	CmdUnInstallService = "-uninstallSvc"
)

type svcHandler struct {
}

func (this *svcHandler) Start(s service.Service) error {
	os.Chdir(filepath.Dir(os.Args[0]))
	go startApplication()
	return nil
}

func (this *svcHandler) Stop(s service.Service) error {
	return nil
}

func logStartupError(s service.Service, err error) {
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		logger.Error(err)
	}
}

func handleServiceCommand(cmd string) error {

	svcConfig := &service.Config{
		Name:        "goPrint",
		DisplayName: "Go Printer Relay",
		Description: "Proxy and status interface for davinci 3d printers",
		Arguments:   []string{CmdRunAsService},
	}
	var err error
	s, err := service.New(&svcHandler{}, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case CmdRunAsService:
		err = s.Run()
		if err != nil {
			logStartupError(s, err)
		}
	case CmdInstallService:
		err = s.Install()
	case CmdUnInstallService:
		err = s.Uninstall()
	default:
		err = errors.New("unknown command")
	}
	return err
}
