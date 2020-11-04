package onewire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// BuildMetric ...
func BuildMetric(device string, temperature float64) ([]byte, error) {
	tags := []string{fmt.Sprintf("device:%s", device)}
	metricPoints := []string{fmt.Sprintf("%v", time.Now().Unix()), fmt.Sprintf("%f", temperature)}
	var series []metricSeries
	series = append(series, metricSeries{MetricName: "w1_temperature.celcius.gauge", Points: [][]string{metricPoints}, Tags: tags, MetricType: "gauge"})
	payload := metricData{Series: series}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return jsonPayload, nil
}

// PostMetric ...
func PostMetric(url string, payload []byte) error {
	// make http request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	// TODO: make this configurable
	client.Timeout = time.Second * 5
	resp, err := client.Do(req)
	// TODO: handle errors/retry
	// context deadline exceeded
	// dial tcp: lookup api.datadoghq.com on 8.8.4.4:53: dial udp 8.8.4.4:53: connect: network is unreachable
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("problem submitting metric to DataDog: %s", string(body))
	}
	resp.Body.Close()
	/*
		log.Println("response Status:", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body:", string(body))
	*/
	return nil
}
