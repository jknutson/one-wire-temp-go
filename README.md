# one-wire-temp

## Description

This program reads the temperature from [DS18B20 Digital Temperature Sensors](https://github.com/openenergymonitor/learn/blob/master/view/electricity-monitoring/temperature/DS18B20-temperature-sensing.md).
Temperature readings are emitted to MQTT topic(s).

## Usage

```sh
$ ./one-wire-temp_darwin -help
Usage: one-wire-temp [options]
Read temperature from one-wire sensores and emit to MQTT
Options:
  -count int
        count of times to poll/report, '-1' means continous (default -1)
  -devicesDir string
        path to one-wire devices (default "/sys/devices/w1_bus_master1/")
  -verbose
        verbose output
  -version
        show version

Environment Variables (will override command line flags):
  DEVICES_DIR - directory path containing one wire device directories
  HOSTNAME - hostname to interpolate in mq topic
  MQ_BROKER - mq broker, ex: "tcp://localhost:1833"
  POLL_INTERVAL - interval (in seconds) at which to poll for temperature

For more information, see https://github.com/jknutson/one-wire-temp-go
```
