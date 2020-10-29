package main

import (
	"flag"
	"fmt"
	"github.com/jknutson/one-wire-temp-go"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	buildVersion string
	version      bool
	count        int
)

func usage() {
	println(`Usage: one-wire-temp [options]
Read temperature from one-wire sensores and POST to DataDog
Options:`)
	flag.PrintDefaults()
	println(`
Environment Variables:
  DD_API_KEY - DataDog API Key (required)
  DEVICES_DIR - directory path containing one wire device directories
  POLL_INTERVAL - interval (in seconds) at which to poll for temperature
`)
	println(`For more information, see https://github.com/jknutson/one-wire-temp-go`)
}

func initFlags() {
	flag.BoolVar(&version, "version", false, "show version")
	flag.IntVar(&count, "count", -1, "count of times to poll/report, '-1' means continous")

	flag.Usage = usage
	flag.Parse()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	initFlags()

	if version {
		log.Println(buildVersion)
		return
	}

	datadogAPIKey := os.Getenv("DD_API_KEY")
	datadogAPIUrl := fmt.Sprintf("https://api.datadoghq.com/api/v1/series?api_key=%s", datadogAPIKey)

	devicesDir := "/sys/devices/w1_bus_master1/"
	if os.Getenv("DEVICES_DIR") != "" {
		devicesDir = os.Getenv("DEVICES_DIR")
	}
	devices, err := onewire.GetDevices(devicesDir)
	check(err)

	pollInterval := int64(30)
	if os.Getenv("POLL_INTERVAL") != "" {
		pollInterval, err = strconv.ParseInt(os.Getenv("POLL_INTERVAL"), 10, 32)
		check(err)
	}

	pollCount := 0
	for {
		for _, device := range devices {
			deviceFile := path.Join(devicesDir, device, "w1_slave")
			temperatureCelcius, err := onewire.ReadDevice(deviceFile)
			check(err)
			log.Printf("device: %s, temperature (celcius): %f", device, temperatureCelcius)

			metricPayload, err := onewire.BuildMetric(device, temperatureCelcius)
			check(err)
			err = onewire.PostMetric(datadogAPIUrl, metricPayload)
			check(err)
		}
		if count != -1 {
			pollCount = pollCount + 1
			if pollCount >= count {
				log.Printf("exiting after %v polls", pollCount)
				os.Exit(0)
			}
		}
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
