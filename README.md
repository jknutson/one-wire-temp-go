# one-wire-temp

## Usage

```sh
$ ./one-wire-temp -help
Usage: one-wire-temp [options]
Read temperature from one-wire sensores and POST to DataDog
Options:
  -count int
        count of times to poll/report, '-1' means continous (default -1)
  -version
        show version

Environment Variables:
  DD_API_KEY - DataDog API Key (required)
  DEVICES_DIR - directory path containing one wire device directories
  POLL_INTERVAL - interval (in seconds) at which to poll for temperature

For more information, see https://github.com/jknutson/one-wire-temp-go
```
