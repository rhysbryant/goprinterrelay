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
func getStatusText(code int) string {
	switch code {
	case 0:

		return "Initializing..."

	case 1:

		return "Heating"

	case 2:

		return "Printing"

	case 3:

		return "Calibrating"

	case 4:

		return "Calibrating Done"

	case 5:

		return "Cooling Finished"

	case 6:

		return "Cooling Finished."

	case 7:

		return "Print Process Ending"

	case 8:

		return "Print Process Ending."

	case 9:

		return "Print Job Finished"

	case 10:

		return "Ready"

	case 11:
		return "Preparing Print"

	case 12:

		return "Print Stopped"

	case 13:

		return "Loading Filament"

	case 14:

		return "Unloading filament"

	case 15:

		return "Auto Calibration"

	case 16:

		return "Job Mode"

	case 17:

		return "Print Error"

	case 9601:
		return "Print Job Paused"

	case 9602:
		return "Print Job Cancelled"

	}
	return ""
}
