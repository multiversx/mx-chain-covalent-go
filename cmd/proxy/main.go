package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

const (
	logFilePrefix = "covalent-proxy"
	tomlFile      = "./config.toml"
)

func main() {
	app := cli.NewApp()
	app.Name = "Covalent proxy indexer tool"
	app.Usage = "This is the entry point for covalent proxy indexer tool. It acts as a proxy to fetch hyperblocks from Elrond  It converts hyperblocks data and provides API calls for covalent in their desired format"
	app.Flags = getFlags()
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Action = startProxy
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
		return
	}
}

func startProxy(ctx *cli.Context) error {
	flagsConfig := getFlagsLogConfig(ctx)
	errLogger := attachFileLogger(log, logFilePrefix, flagsConfig)
	if errLogger != nil {
		return errLogger
	}

	log.Info("starting server")

	cfg, err := config.LoadConfig(tomlFile)
	if err != nil {
		return err
	}
	server := createServer(cfg)

	go func() {
		err = server.ListenAndServe()
		log.LogIfError(err)
	}()

	waitForServerShutdown(server)
	return nil
}

func createServer(cfg *config.Config) api.HTTPServer {
	hyperBlockFacade := api.NewHyperBlockFacade(80, cfg.ElrondProxyUrl)

	//processor, err := factory.CreateDataProcessor(&factory.ArgsDataProcessor{
	//	PubKeyConvertor:  nil,
	//	Accounts:         nil,
	//	Hasher:           nil,
	//	Marshaller:       nil,
	//	ShardCoordinator: nil,
	//})
	hyperBlockProxy := api.NewHyperBlockProxy(hyperBlockFacade, &utility.AvroMarshaller{}, &mock.DataHandlerStub{})

	router := gin.Default()
	router.GET(fmt.Sprintf("%s/by-nonce/:nonce", cfg.HyperBlockPath), hyperBlockProxy.GetHyperBlockByNonce)

	return &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}
}

func waitForServerShutdown(httpServer api.HTTPServer) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	shutdownContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := httpServer.Shutdown(shutdownContext)
	log.LogIfError(err)
	err = httpServer.Close()
	log.LogIfError(err)
}
