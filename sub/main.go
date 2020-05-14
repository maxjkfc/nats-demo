package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var (
	natsHost string
	n        *nats.Conn
	tid      = uuid.New().String()
)

func main() {
	natsHost = "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224"
	logrus.Infof("ConnectHost: %s", natsHost)

	nc, err := nats.Connect(natsHost, nats.DontRandomize())
	if err != nil {
		panic(err)
	}

	defer nc.Close()

	if nc.AuthRequired() {
		logrus.Info("need auth")
	}

	n = nc

	logrus.Info("connect success")
	time.Sleep(10 * time.Second)

	subj := "msg"
	//queueName := "TT123"
	i := 0
	reply := tid

	nc.Subscribe(subj, func(msg *nats.Msg) {
		i++
		logrus.Printf("[#%d] Received on [%s]: '%s' \n", i, msg.Subject, string(msg.Data))
		logrus.Printf("Reply data : %v", reply)

	})

	nc.Flush()

	fmt.Println("Listen Sub With: ", subj)

	if err := nc.LastError(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Printf("Listening on [%s]", subj)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logrus.Println()
	logrus.Printf("Draining...")
	nc.Drain()
	logrus.Fatalf("Exiting")

}
