package main

import (
	"fmt"
	"os"

	"github.com/rakyll/portmidi"
)

func main() {

	var val int64
	var selectedDeviceId portmidi.DeviceId = -1

	if len(os.Args) < 2 {
		fmt.Println("must provide on or off")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "on":
		val = 127
	case "off":
		val = 0
	default:
		val = 127
	}

	portmidi.Initialize()
	defer portmidi.Terminate()
	total := portmidi.CountDevices()

	for id := 0; id < total; id++ {
		id := portmidi.DeviceId(id)
		info := portmidi.GetDeviceInfo(id)
		if info.Name == "MicroBrute" && info.IsOutputAvailable {
			selectedDeviceId = id
			fmt.Println("Found Microbrute output at device id", selectedDeviceId)
		}
	}

	if selectedDeviceId == -1 {
		fmt.Println("Microbrute not found")
		os.Exit(0)
	}

	out, err := portmidi.NewOutputStream(selectedDeviceId, 1024, 0)
	if err != nil {
		// fmt.Println()
		fmt.Println("Error opening device: ", err)
		os.Exit(0)
	}

	// 0x80     Note Off
	// 0x90     Note On
	// 0xA0     Aftertouch
	// 0xB0     Continuous controller
	// 0xC0     Patch change
	// 0xD0     Channel Pressure
	// 0xE0     Pitch bend
	// 0xF0     (non-musical commands)

	// CC 122 is Local ON/OFF according to MicroBurte Connection Manual
	// https://dl.dropboxusercontent.com/u/976344/MICROBRUTE_Connection-Manual_v1.0.pdf
	err = out.WriteShort(0xB0, 122, val)
	if err != nil {
		fmt.Println(err)
	}

}
