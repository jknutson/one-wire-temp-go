package main

import (
	"flag"
	"fmt"
	"github.com/jknutson/one-wire-temp-go"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"
)

var (
	buildVersion string
	count        int
	verbose      bool
	version      bool
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
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
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

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Ctrl+C pressed, exiting.")
		os.Exit(0)
	}()
}

func main() {
	initFlags()
	setupCloseHandler()

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

	log.Println("devices dir: ", devicesDir)
	log.Println("poll interval: ", pollInterval)
	log.Println("devices found: ", devices)

	pollCount := 0
	for {
		for _, device := range devices {
			deviceFile := path.Join(devicesDir, device, "w1_slave")
			temperatureCelcius, err := onewire.ReadDevice(deviceFile)
			check(err)
			if verbose {
				log.Printf("device: %s, temperature (celcius): %f", device, temperatureCelcius)
			}

			metricPayload, err := onewire.BuildMetric(device, temperatureCelcius)
			check(err)
			if verbose {
				log.Printf("payload: %s", metricPayload)
			}
			err = onewire.PostMetric(datadogAPIUrl, metricPayload)
			if err != nil {
				log.Printf("error posting metric: %v", err)
			}
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
