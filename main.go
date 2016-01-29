// Brutetoggle is a command line utility for switching
// MIDI Local Control on or off on a MicroBrute synthesizer
// Usage:
// microbrute [on|off]
package main

import (
	"fmt"
	"os"

	"github.com/rakyll/portmidi"
)

const (
	ccMessage             = 0xB0 // Chan 1 Control/Mode Change (CC message on MIDI channel 1): https://www.midi.org/specifications/item/table-2-expanded-messages-list-status-bytes
	localControlCCNumber  = 0x7A // Local Control On/Off, according to Standard MIDI (https://www.midi.org/specifications/item/table-3-control-change-messages-data-bytes-2) and MicroBrute Connection Manual
	stepControlCCNumber   = 0x72 // Undefined by MIDI: MicroBrute uses it for Step On Clock or Gate (MicroBrute Connection Manual)
	localControlOnValue   = 127
	localControlOffValue  = 0
	stepControlGateValue  = 127
	stepControlClockValue = 0
)

func main() {

	var (
		ccValue  int64
		ccNumber int64
		deviceId portmidi.DeviceId
		err      error
	)

	ccNumber, ccValue, err = parseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	portmidi.Initialize()
	defer portmidi.Terminate()

	deviceId, err = getDeviceId()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := portmidi.NewOutputStream(deviceId, 1024, 0)
	if err != nil {
		fmt.Println("Error opening MicroBrute output port:", err)
		os.Exit(1)
	}

	err = out.WriteShort(ccMessage, ccNumber, ccValue)
	if err != nil {
		fmt.Println("Error sending message to MicroBrute:", err)
		os.Exit(1)
	}

}

func parseArgs(args []string) (ccNumber, ccValue int64, err error) {
	if len(args) != 3 {
		return 0, 0, fmt.Errorf("Not enough argument provided:\n%s (local|step) (on|off|clock|gate)", args)
	}

	switch args[1] {
	case "local":
		ccNumber = localControlCCNumber
		switch os.Args[2] {
		case "on":
			ccValue = localControlOnValue
		case "off":
			ccValue = localControlOffValue
		default:
			return 0, 0, fmt.Errorf("Value for 'local' must be (on|off)")
		}
	case "step":
		ccNumber = stepControlCCNumber
		switch os.Args[2] {
		case "clock":
			ccValue = stepControlClockValue
		case "gate":
			ccValue = stepControlGateValue
		default:
			return 0, 0, fmt.Errorf("Value for 'step' must be (clock|gate)")
		}
	default:
		fmt.Println("Operation must be (local|step)")
	}

	return
}

func getDeviceId() (portmidi.DeviceId, error) {
	total := portmidi.CountDevices()

	for id := 0; id < total; id++ {
		id := portmidi.DeviceId(id)
		info := portmidi.GetDeviceInfo(id)
		if info.Name == "MicroBrute" && info.IsOutputAvailable {
			return id, nil
		}
	}

	return -1, fmt.Errorf("MicroBrute not found")
}
