package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/nats-io/stan.go"

	"l0/order"
)

type appEnv struct {
	messageDelay time.Duration
	clusterID    string
	clientID     string
	natsURL      string
	stanConn     stan.Conn
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("stan-publisher", flag.ContinueOnError)
	fl.DurationVar(&app.messageDelay, "d", 1*time.Second, "Message delay")
	fl.StringVar(&app.clusterID, "crid", "test-cluster", "Nats streaming cluster ID")
	fl.StringVar(&app.clientID, "clid", "client-123", "Nats streaming client ID")
	fl.StringVar(&app.natsURL, "u", "nats://127.0.0.1:4222", "Nats streaming URL")

	if err := fl.Parse(args); err != nil {
		return err
	}

	sc, err := stan.Connect(app.clusterID, app.clientID, stan.NatsURL(app.natsURL))
	if err != nil {
		return err
	}

	app.stanConn = sc

	err = faker.SetRandomMapAndSliceMinSize(1)
	if err != nil {
		return err
	}
	err = faker.SetRandomMapAndSliceSize(5)
	if err != nil {
		return err
	}

	return nil
}

func (app *appEnv) publishMessage() {
	o := new(order.Order)
	err := faker.FakeData(o)
	if err != nil {
		log.Printf("can't create fake data: %s", err.Error())
	}

	b, err := json.Marshal(o)
	if err != nil {
		log.Printf("error marshaling message %s", err.Error())
	}

	err = app.stanConn.Publish("orders", b)
	if err != nil {
		log.Printf("error publishing message %s", err.Error())
	}
	log.Printf("%s order is sent", o.OrderUID)
}

func (app *appEnv) run() {
	defer app.stanConn.Close()

	ticker := time.NewTicker(app.messageDelay)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			app.publishMessage()
		case <-signalChan:
			log.Printf("\nreceived an interrupt signal\n")
			ticker.Stop()
			return
		}
	}

	return
}

func cli(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		log.Println(err.Error())
		return 2
	}

	app.run()

	return 0
}

func main() {
	os.Exit(cli(os.Args[1:]))
}
