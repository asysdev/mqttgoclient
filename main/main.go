package main

import (
	"github.com/mqttgoclient/mqtt"
	"log"
)

func MessageArrived(client *mqtt.Client, topic string, message []byte) {
	log.Println("REceived Data: ", message)
}

func main() {

	var client mqtt.Client

	client.Broker = "tcp://test.mosquitto.org:1883"
	client.CleanSession = true
	client.ClientId = "goclienttest"
	client.MessageArrived = MessageArrived

	err := client.Init()

	if err != nil {
		log.Panic(err)
	}

	client.Connect()

	err = client.Subscribe("goclient_sub", 0)
	log.Println(err)

	client.Publish("goclient_pub", 0, false, "Hello")

}
