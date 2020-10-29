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

			/*
				// build metric payload
				tags := []string{fmt.Sprintf("device:%s", device)}
				metricPoints := []string{fmt.Sprintf("%v", time.Now().Unix()), fmt.Sprintf("%f", temperatureCelcius)}
				var series []metricSeries
				series = append(series, metricSeries{MetricName: "w1_temperature.celcius.gauge", Points: [][]string{metricPoints}, Tags: tags, MetricType: "gauge"})
				payload := metricData{Series: series}
				jsonPayload, err := json.Marshal(payload)
				check(err)
				// log.Printf("JSON payload: %s\n", jsonPayload)

				// make http request
				req, err := http.NewRequest("POST", datadogAPIUrl, bytes.NewBuffer(jsonPayload))
				check(err)
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				// TODO: make this configurable
				client.Timeout = time.Second * 5
				resp, err := client.Do(req)
				check(err)
				if resp.StatusCode != 202 {
					body, _ := ioutil.ReadAll(resp.Body)
					log.Printf("problem submitting metric to DataDog: %s", string(body))
				}
				resp.Body.Close()
				log.Println("response Status:", resp.Status)
				body, _ := ioutil.ReadAll(resp.Body)
				log.Println("response Body:", string(body))
			*/
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
