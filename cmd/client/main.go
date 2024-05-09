package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const conStr = "amqp://guest:guest@localhost:5672/"

func main() {
	client, err := amqp.Dial(conStr)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	fmt.Println("Succesfully connected...")

	m, err := gamelogic.ClientWelcome()
	if err != nil {
		panic(err)
	}

	qName := fmt.Sprintf("pause.%s", m)
	_, _, err = pubsub.DeclareAndBind(client, routing.ExchangePerilDirect, qName, routing.PauseKey, 1)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Shutting down...")

	client.Close()
}
