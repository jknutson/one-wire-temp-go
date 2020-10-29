package onewire

import "testing"

func TestGetDevices(t *testing.T) {
	devicesPath := "./test/devices"
	devices, err := GetDevices(devicesPath)
	if err != nil {
		t.Fatalf("could not read devices from path:  %s", devicesPath)
	}

	numDevices := 2
	if len(devices) != numDevices {
		t.Fatalf("unexpected number of devices found. expected: %v, got: %v", numDevices, len(devices))
	}
}

func TestReadDevice(t *testing.T) {
	device := "./test/devices/28-0516a42628ff/w1_slave"
	temperature, err := ReadDevice(device)
	if err != nil {
		t.Fatalf("error reading temperature from device %s: %s", device, err)
	}

	expectedTemp := 24.937
	if temperature != expectedTemp {
		t.Fatalf("incorrect temperature value read. expected %v, got: %v", temperature, expectedTemp)
	}
}
