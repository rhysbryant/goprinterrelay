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
func GetStatusText(code int) string {
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

	case 30:
		return "Printer Busy"

	case 9601:
		return "Print Job Paused"

	case 9602:
		return "Print Job Cancelled"

	case 9500:
		return "Initializing..."

	case 9501:
		return "Heating"

	case 9502:
		return "Printing"

	case 9503:
		return "Calibrating"

	case 9504:
		return "Calibrating Done"

	case 9505:
		return "Printing In Progress"

	case 9506:
		return "Cooling Finished"

	case 9507:
		return "Cooling Finished."

	case 9508:
		return "Print Process Ending."

	case 9509:
		return "Print Process Ending."

	case 9510:
		return "Print Job Finished"

	case 9511:
		return "Ready"

	case 9512:
		return "Print Stopped"

	case 9513:
		return "Loading Filament."

	case 9514:
		return "Unloading filament."

	case 9515:
		return "Auto Calibration"

	case 9516:
		return "Job Mode."

	case 9517:
		return "Print Error."

	case 9520:
		return "Print file check"

	case 9530:
		return "Loading Filament.."

	case 9531:
		return "Unloading Filament.."

	case 9532:
		return "Job Mode.."

	case 9533:
		return "Print Error.."

	case 9534:
		return "Homing"

	case 9535:
		return "Calibrating."

	case 9536:
		return "Cleaning Nozzle"

	case 9537:
		return "Get SD File"

	case 9538:
		return "Print SD File"

	case 9539:
		return "Print Engrave Place Object"

	case 9540:
		return "Adjusting Z-Offset"

	case 9700:
		return "Busy"

	case 9800:
		return "Scanner Idle"

	case 9801:
		return "Scanner Running"
	}

	return "<Unknown>"
}
