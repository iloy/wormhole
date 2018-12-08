package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"

	"github.com/iloy/wormhole/server"
)

func main() {
	flagversion := flag.Bool("version", false, "output version information and exit")
	flagnocolor := flag.Bool("nocolor", false, "log messages without color")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   *flagnocolor,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	log.SetOutput(colorable.NewColorableStdout())

	if *flagversion {
		fmt.Println(version())
		return
	}

	sServer := server.StreamingServer{}

	err := sServer.Serve()
	if err != nil {
		panic(err)
	}
	log.Info("server started")

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Infoln("got signal:", sig)

	sServer.Stop()

	log.Info("it is done")
}
