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
)

type QueryFieldsCache interface {
	GetField(name string) (string, bool)
	SetField(name string, val string)
	GetAllFields() map[string]string
}

type QueryFieldsCacheMem struct {
	fields    map[string]string
	overrides map[string]string
}

func (qc *QueryFieldsCacheMem) GetField(name string) (string, bool) {
	k, v := qc.fields[name]

	return k, v
}

func (qc *QueryFieldsCacheMem) SetField(name string, val string) {
	if _, exists := qc.overrides[name]; exists {
		return
	}
	qc.fields[name] = val
}

func (qc *QueryFieldsCacheMem) GetAllFields() map[string]string {
	return qc.fields
}

func NewQueryFieldsCache(overrides map[string]string) *QueryFieldsCacheMem {
	qc := QueryFieldsCacheMem{}
	qc.fields = make(map[string]string, 10)
	qc.overrides = overrides
	for k, v := range overrides {
		qc.fields[k] = v
		fmt.Printf("set override for %s to %s\n", k, v)
	}
	return &qc
}
