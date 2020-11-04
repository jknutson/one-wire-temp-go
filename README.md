# one-wire-temp

## Description

This program reads the temperature (Celcius) from [DS18B20 Digital Temperature Sensors](https://github.com/openenergymonitor/learn/blob/master/view/electricity-monitoring/temperature/DS18B20-temperature-sensing.md).
Temperature readings are emitted as Custom Metrics to DataDog via the HTTP API. (Celcius)j

## Usage

```sh
$ ./one-wire-temp_darwin -help
Usage: one-wire-temp [options]
Read temperature from one-wire sensores and POST to DataDog
Options:
  -count int
        count of times to poll/report, '-1' means continous (default -1)
  -verbose
        verbose output
  -version
        show version

Environment Variables:
  DD_API_KEY - DataDog API Key (required)
  DEVICES_DIR - directory path containing one wire device directories
  POLL_INTERVAL - interval (in seconds) at which to poll for temperature

For more information, see https://github.com/jknutson/one-wire-temp-go
```
