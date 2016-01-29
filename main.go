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
	ccMessage             = 0xB0 // Table 2: Chan 1 Control/Mode Change in http://www.midi.org/techspecs/midimessages.php
	localControlCCNumber  = 0x7A // Table 3: Local Control On/Off (supported by MicroBrute)
	stepControlCCNumber   = 0x72 // Undefined by MIDI: MB uses it for Step On Clk | Gate
	localControlOnValue   = 127
	localControlOffValue  = 0
	stepControlGateValue  = 127
	stepControlClockValue = 0
)

func main() {

	var (
		ccValue  int64
		ccNumber int64
		deviceId portmidi.DeviceId = -1
		err      error
	)

	if len(os.Args) != 3 {
		fmt.Println("Not enough argument provided: (local|step) (on|off|clock|gate)")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "local":
		ccNumber = localControlCCNumber
	case "step":
		ccNumber = stepControlCCNumber
	default:
		fmt.Println("Operation must be (local|step)")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "on":
		ccValue = localControlOnValue
	case "off":
		ccValue = localControlOffValue
	case "clock":
		ccValue = stepControlClockValue
	case "gate":
		ccValue = stepControlGateValue
	default:
		fmt.Println("Value must be (on|off|clock|gate)")
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
