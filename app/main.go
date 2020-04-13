package main

import (
	"bytes"
	"flag"
	"github.com/pedidosya/minesweeper-API/app/config"
	"github.com/pedidosya/minesweeper-API/app/server"
	"github.com/pedidosya/minesweeper-API/utils"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func main() {
	StartApp()
}

func StartApp() {
	initLog()
	readConfiguration()
	startWebServer()
}

func initLog() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func startWebServer() {

	s := server.New(&server.Config{
		Port: viper.GetInt("server.port"),
	})

	Routes(s)

	utils.LogInfo("starting http listener ...")
	go func() {
		s.ListenAndServe()
	}()

	// Wait for terminate signal to shut down server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	utils.LogInfo("services started and listening in port %d ...", viper.GetInt("server.port"))
	<-c

}

func readConfiguration() {

	//Util for several environments
	env := flag.String("E", "live", "Execution environment")
	flag.Parse()
	utils.LogInfo("Starting minesweeper in %s environment ...", *env)

	viper.SetConfigType("json")
	var err error

	switch *env {
	case "live":
		err = viper.ReadConfig(bytes.NewBuffer(config.Live))
	}

	if err != nil {
		log.Fatalf("error reading configuration from viper: %v", err)
	}
}
