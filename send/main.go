package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/axolotlteam/thunder/logger"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var (
	natsHost string
)

func main() {
	flag.StringVar(&natsHost, "n", "", "set nats host")

	flag.Parse()
	if !flag.Parsed() {
		fmt.Println("Not parse")
	}

	natsHost = "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224"

	logrus.Infof("ConnectHost: %s", natsHost)
	fmt.Println(natsHost)

	nc, err := nats.Connect(natsHost, nats.DontRandomize())
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	if nc.AuthRequired() {
		logger.Info("need auth")
	}

	logger.Info("connect success")
	time.Sleep(10 * time.Second)

	subj := "test"

	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		data := input.Bytes()
		msg, err := nc.Request(subj, data, 2*time.Second)
		if err != nil {
			if nc.LastError() != nil {
				logrus.Errorf("last error : %v for request", nc.LastError())
			}
			logrus.Errorf(" xx : %v for request", err)
		} else {

			t := time.Now()
			logrus.WithTime(t).Printf("Published [%s] : %s", subj, data)
			logrus.WithTime(t).Printf("Received  [%s] : %s[%s]", msg.Subject, string(msg.Data), msg.Reply)
		}
	}

	logrus.Info("exit")
}
