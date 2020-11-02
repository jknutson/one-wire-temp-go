package onewire

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

// GetDevices ...
func GetDevices(devicesDir string) ([]string, error) {
	var devices []string

	files, err := ioutil.ReadDir(devicesDir)
	if err != nil {
		return nil, err
	}

	deviceRegexp := regexp.MustCompile(`^28.*`)

	for _, f := range files {
		matched := deviceRegexp.MatchString(f.Name())
		// TODO: make this support symlinks?
		if f.IsDir() && matched {
			devices = append(devices, f.Name())
		}
	}
	return devices, nil
}

// ReadDevice ...
func ReadDevice(deviceFile string) (float64, error) {
	dat, err := ioutil.ReadFile(deviceFile)
	if err != nil {
		return 0, err
	}

	temperatureRegexp := regexp.MustCompile(`(?s)^.*t\=(\d+)\n$`)
	temperatureCelciusMatch := temperatureRegexp.FindSubmatch(dat)
	if temperatureCelciusMatch == nil {
		return 0, fmt.Errorf("could not parse temperature from file: %s\ncontents: %s", deviceFile, string(dat))
	}
	temperatureCelcius, err := strconv.ParseFloat(string(temperatureCelciusMatch[1]), 32)
	if err != nil {
		return 0, fmt.Errorf("could not parse temperature to integer: %v", temperatureCelcius)
	}
	temperatureCelcius = temperatureCelcius / 1000
	return temperatureCelcius, nil
}
