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
	"testing"
)

func TestParseKeyValue(t *testing.T) {
	key, value := parseKeyValueLine("key:value")
	if *key != "key" {
		t.Errorf("expected key got %s", *key)
		t.Fail()
	}
	if *value != "value" {
		t.Errorf("expected value got %s", *value)
		t.Fail()
	}
}

func TestParseKeyValueError(t *testing.T) {
	key, value := parseKeyValueLine("key")
	if key != nil {
		t.Errorf("expected nil for key")
		t.Fail()
	}
	if value != nil {
		t.Errorf("expected nil for value got %s", *value)
		t.Fail()
	}
}
