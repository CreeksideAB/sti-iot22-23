package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func main() {
	keepAlive := make(chan os.Signal)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	broker := os.Getenv("BROKER_ADDRESS")
	port := os.Getenv("BROKER_PORT")
	opts := mqtt.NewClientOptions()
	url := fmt.Sprintf("tcp://%s:%s", broker, port)
	log.Println(url)
	opts.AddBroker(url)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	var token mqtt.Token
	log.Println("Wait for mqtt")
	for i := 0; i < 30; i++ {
		if token = client.Connect(); token.Wait() && token.Error() != nil {
			log.Print(".")
			time.Sleep(time.Second)
			continue
		}
		if i > 0 {
			log.Print("\n")
		}
		break
	}
	if token.Error() != nil {
		panic(token.Error())
	} else {
		log.Print("mqtt connected succesful\n")
	}

	client.Subscribe("hej/på/+", 1, messagePubHandler)

	<-keepAlive
}
