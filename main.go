package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/senseyman/bitcoin-handshake/client"
	"github.com/senseyman/bitcoin-handshake/core"
	"github.com/senseyman/bitcoin-handshake/service"
)

var (
	nodeHostFlag = flag.String("node.host", "127.0.0.1", "Host of blockchain node")
	nodePortFlag = flag.Int("node.port", 18333, "Port of blockchain node")
)

func main() {
	flag.Parse()

	// setup logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	log.Info("App starting...")

	// create all services
	log.Info("Initializing all services...")
	readSrv := service.NewDecodeService()
	writeSrv := service.NewEncodeService()
	msgGenerator := service.NewMessageGenerator()

	// create node client
	log.Infof("Connecting to bitcoin node, host %s, port %d...", *nodeHostFlag, *nodePortFlag)
	btcnCli, err := client.NewBitcoinClient(
		*nodeHostFlag, *nodePortFlag,
		func(host string, port int) (client.Connection, error) {
			return net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// create main core logic service
	coreSystem := core.New(readSrv, writeSrv, msgGenerator, btcnCli)

	// general context
	globalCtx, globalCtxCancel := context.WithCancel(context.Background())

	// handshake context. If nothing work in 1 minutes - stop the app by timeout
	handshakeCtx, cancel := context.WithTimeout(globalCtx, time.Minute)
	defer cancel()

	// init graceful shutdown if got some terminations from outside
	setupGracefulShutdown(globalCtxCancel)

	// start listening messages from node
	log.Info("starting reading incoming messages from node")
	coreSystem.ReceiveMessages(handshakeCtx)

	// start main task flow
	log.Info("starting handshake")

	execTimeMs, err := coreSystem.Handshake(handshakeCtx)
	if err != nil {
		log.Fatalf("error while doing main flow: %v", err)
	}

	log.Info("All necessary messages for connection are received.")
	log.Infof("Handshake took %d ms.", execTimeMs)
	log.Info("Stopping the App...")
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Warnf("Got Interrupt signal: %v", sig.String())
		stop()
	}()
}
