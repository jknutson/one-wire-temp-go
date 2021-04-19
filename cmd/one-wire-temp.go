package main

import (
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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

	// mqtt setup
	// TODO: get this from env
	mqHostname, ok := os.LookupEnv("HOSTNAME")
	if !ok {
		mqHostname = "raspberrypi" // default if HOSTNAME env var is not set
	}
	mqBroker, ok := os.LookupEnv("MQ_BROKER")
	if !ok {
		mqBroker = "tcp://192.168.2.6:1883"
	}
	mqOpts := MQTT.NewClientOptions().AddBroker(mqBroker)
	mqOpts.SetClientID(mqHostname)
	mqTopicBase := fmt.Sprintf("pi/%s/temperature", mqHostname)

	if version {
		log.Println(buildVersion)
		return
	}

	c := MQTT.NewClient(mqOpts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

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
			temperatureFahrenheit := (temperatureCelcius * 1.8) + 32
			token := c.Publish(fmt.Sprintf("%s/%s", mqTopicBase, device), 0, false, temperatureFahrenheit)
			token.Wait()
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
