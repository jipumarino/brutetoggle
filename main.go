package main

import (
	"fmt"
	"os"

	"github.com/rakyll/portmidi"
)

const (
	ccMessage      = 0xB0 // Table 2: Chan 1 Control/Mode Change in http://www.midi.org/techspecs/midimessages.php
	localControlCC = 0x7A // Table 3: Local Control On/Off (supported by MicroBrute)
	onMessage      = 127
	offMessage     = 0
)

func main() {

	var (
		localControlVal  int64
		selectedDeviceId portmidi.DeviceId = -1
	)

	if len(os.Args) != 2 {
		fmt.Println("No argument provided")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "on":
		localControlVal = onMessage
	case "off":
		localControlVal = offMessage
	default:
		fmt.Println("Argument must be 'on' or 'off'")
		os.Exit(1)
	}

	portmidi.Initialize()
	defer portmidi.Terminate()
	total := portmidi.CountDevices()

	for id := 0; id < total; id++ {
		id := portmidi.DeviceId(id)
		info := portmidi.GetDeviceInfo(id)
		if info.Name == "MicroBrute" && info.IsOutputAvailable {
			selectedDeviceId = id
			fmt.Println("Found MicroBrute output port at device id", selectedDeviceId)
			break // Only apply command to first MicroBrute found
		}
	}

	if selectedDeviceId == -1 {
		fmt.Println("MicroBrute not found")
		os.Exit(1)
	}

	out, err := portmidi.NewOutputStream(selectedDeviceId, 1024, 0)
	if err != nil {
		fmt.Println("Error opening MicroBrute output port:", err)
		os.Exit(1)
	}

	err = out.WriteShort(ccMessage, localControlCC, localControlVal)
	if err != nil {
		fmt.Println("Error sending message to MicroBrute:", err)
		os.Exit(1)
	}

}
