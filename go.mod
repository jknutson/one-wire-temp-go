module github.com/jknutson/one-wire-temp

go 1.14

require (
	github.com/DataDog/datadog-go v4.1.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.3.3
	github.com/jknutson/one-wire-temp-go v0.0.0-00010101000000-000000000000
)

replace github.com/jknutson/one-wire-temp-go => ./
