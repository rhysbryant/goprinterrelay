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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type StatusQuery struct {
	PrinterStatus   string `json:"status"`
	Temperature     int    `json:"temperature"`
	PrintProgress   int    `json:"printProgress"`
	EstmatedMinutes int    `json:"estmatedTime"`
	ElapsedMinutes  int    `json:"elapsedTime"`
	Filament        struct {
		Totoal    int `json:"total"`
		Remaining int `json:"remaining"`
	} `json:"filament"`
	PrinterInfo struct {
		Type   string `json:"type"`
		Serial string `json:"serial"`
	} `json:"printerInfo"`
}

const (
	TemperatureId       = "t"
	ProgressInfoId      = "d"
	PrinterStatusId     = "j"
	FilamentLengthId    = "L"
	FilamentRemainingId = "f"
	ProductId           = "p"
	ProductSerialId     = "i"
)

func getStatusFromMap(values map[string]string) (*StatusQuery, error) {
	status := StatusQuery{}
	var err error

	if val, exists := values[TemperatureId]; exists {
		parts := strings.Split(val, ",")

		if len(parts) < 2 {
			return nil, errors.New("unexpected format")
		}
		status.Temperature, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	if val, exists := values[ProgressInfoId]; exists {
		parts := strings.Split(val, ",")

		if len(parts) < 3 {
			return nil, errors.New("unexpected format")
		}

		status.PrintProgress, err = strconv.Atoi(parts[0])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		status.ElapsedMinutes, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		status.EstmatedMinutes, err = strconv.Atoi(parts[2])
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	if val, exists := values[PrinterStatusId]; exists {
		statusCode, err := strconv.Atoi(val)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		status.PrinterStatus = getStatusText(statusCode)
		fmt.Println(statusCode)
	}

	if val, exists := values[FilamentLengthId]; exists {
		parts := strings.Split(val, ",")

		if len(parts) < 2 {
			return nil, errors.New("unexpected format")
		}

		status.Filament.Totoal, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)
			return nil, err
		}

	}

	if val, exists := values[FilamentRemainingId]; exists {
		parts := strings.Split(val, ",")

		if len(parts) < 2 {
			return nil, errors.New("unexpected format")
		}

		status.Filament.Remaining, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)
			return nil, err
		}

	}

	if val, exists := values[ProductId]; exists {
		status.PrinterInfo.Type = val
	}

	if val, exists := values[ProductSerialId]; exists {
		status.PrinterInfo.Serial = val
	}
	return &status, nil
}

func WriteJsonResponse(w http.ResponseWriter, object interface{}) {

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Server", "goprint")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(object)
}

func StatusRequest(w http.ResponseWriter, r *http.Request) {
	WriteJsonResponse(w, printerStatus)
}

func startHttpServer(listenerPath string) {

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/imagestream", imageStreamWebsocketsMgr.WsUpgradeHandler)
	router.HandleFunc("/status", StatusRequest).Methods("GET")
	router.HandleFunc("/toolsform", ToolsFormPostRequest).Methods("POST")
	router.HandleFunc("/tools", ToolsGetRequest).Methods("GET")
	router.HandleFunc("/applicationInfo", ApplicationInfoGetRequest).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("ui")))
	err := http.ListenAndServe(listenerPath, router)

	if err != nil {
		log.Fatal(err)
	}
}
