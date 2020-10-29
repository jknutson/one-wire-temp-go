package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

type metricSeries struct {
	MetricName string     `json:"metric"`
	Points     [][]string `json:"points"`
	Tags       []string   `json:"tags"`
	MetricType string     `json:"type"`
}

type metricData struct {
	Series []metricSeries `json:"series"`
}

var (
	buildVersion string
	version      bool
	count        int
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func usage() {
	println(`Usage: one-wire-temp [options]
Read temperature from one-wire sensores and POST to DataDog
Options:`)
	flag.PrintDefaults()
	println(`For more information, see https://github.com/jwilder/docker-gen`)
}

func initFlags() {
	flag.BoolVar(&version, "version", false, "show version")
	flag.IntVar(&count, "count", -1, "count of times to poll/report, '-1' means continous")

	flag.Usage = usage
	flag.Parse()
}
func main() {
	initFlags()

	if version {
		log.Println(buildVersion)
		return
	}

	datadogAPIKey := os.Getenv("DD_API_KEY")
	datadogAPIUrl := fmt.Sprintf("https://api.datadoghq.com/api/v1/series?api_key=%s", datadogAPIKey)

	devicesDir := "/sys/bus/w1/devices/"
	if os.Getenv("DEVICES_DIR") != "" {
		devicesDir = os.Getenv("DEVICES_DIR")
	}

	var devices []string

	files, err := ioutil.ReadDir(devicesDir)
	check(err)

	deviceRegexp := regexp.MustCompile(`^28.*`)

	for _, f := range files {
		matched := deviceRegexp.MatchString(f.Name())
		// TODO: make this support symlinks?
		if f.IsDir() && matched {
			devices = append(devices, f.Name())
		}
	}
	log.Printf("devices found: %q", devices)

	temperatureRegexp := regexp.MustCompile(`(?s)^.*t\=(\d+)\n$`)

	pollInterval := int64(30)
	if os.Getenv("POLL_INTERVAL") != "" {
		pollInterval, err = strconv.ParseInt(os.Getenv("POLL_INTERVAL"), 10, 32)
		check(err)
	}

	pollCount := 0
	for {
		for _, device := range devices {
			deviceFile := path.Join(devicesDir, device, "w1_slave")
			dat, err := ioutil.ReadFile(deviceFile)
			check(err)
			temperatureCelciusMatch := temperatureRegexp.FindSubmatch(dat)
			if temperatureCelciusMatch == nil {
				log.Fatalf("could not parse temperature from file: %s\ncontents: %s", deviceFile, string(dat))
			}
			temperatureCelcius, err := strconv.ParseFloat(string(temperatureCelciusMatch[1]), 32)
			check(err)
			temperatureCelcius = temperatureCelcius / 1000
			log.Printf("device: %s, temperature (celcius): %f", device, temperatureCelcius)

			tags := []string{fmt.Sprintf("device:%s", device)}

			metricPoints := []string{fmt.Sprintf("%v", time.Now().Unix()), fmt.Sprintf("%f", temperatureCelcius)}

			var series []metricSeries
			series = append(series, metricSeries{MetricName: "w1_temperature.celcius.gauge", Points: [][]string{metricPoints}, Tags: tags, MetricType: "gauge"})

			payload := metricData{Series: series}
			jsonPayload, err := json.Marshal(payload)
			check(err)
			// log.Printf("JSON payload: %s\n", jsonPayload)
			req, err := http.NewRequest("POST", datadogAPIUrl, bytes.NewBuffer(jsonPayload))
			check(err)
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			// TODO: make this configurable
			client.Timeout = time.Second * 5
			_, err = client.Do(req)
			check(err)
			// resp, err := client.Do(req)
			// TODO: should we `defer`?
			// defer resp.Body.Close()
			/*
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
