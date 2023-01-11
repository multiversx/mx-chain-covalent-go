package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-covalent-go/api"
	"github.com/multiversx/mx-chain-covalent-go/cmd/proxy/config"
	"github.com/multiversx/mx-chain-covalent-go/facade"
	"github.com/multiversx/mx-chain-covalent-go/process/factory"
	"github.com/multiversx/mx-chain-covalent-go/process/utility"
	"github.com/urfave/cli"
)

const (
	logFilePrefix = "covalent-proxy"
	tomlFile      = "./config.toml"
)

func main() {
	app := cli.NewApp()
	app.Name = "Covalent proxy indexer tool"
	app.Usage = "This is the entry point for covalent proxy indexer tool. It acts as a proxy to fetch hyperblocks from Multiversx  It converts hyperblocks data and provides API calls for covalent in their desired format"
	app.Flags = getFlags()
	app.Authors = []cli.Author{
		{
			Name:  "The Multiversx Team",
			Email: "contact@multiversx.com",
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

	server, err := createServer(cfg)
	if err != nil {
		return err
	}

	go func() {
		err = server.ListenAndServe()
		log.LogIfError(err)
	}()

	waitForServerShutdown(server)
	return nil
}

func createServer(cfg *config.Config) (api.HTTPServer, error) {
	httpClient := api.NewDefaultHttpClient(cfg.RequestTimeOutSec)
	multiversxHyperBlockEndpointHandler, err := api.NewMultiversxHyperBlockEndPoint(httpClient)
	if err != nil {
		return nil, err
	}

	hyperBlockProcessor, err := factory.CreateHyperBlockProcessor()
	if err != nil {
		return nil, err
	}

	avroEncoder := &utility.AvroMarshaller{}
	hyperBlockFacade, err := facade.NewHyperBlockFacade(
		cfg.MultiversxProxyUrl,
		avroEncoder,
		multiversxHyperBlockEndpointHandler,
		hyperBlockProcessor,
	)
	if err != nil {
		return nil, err
	}

	hyperBlockProxy, err := api.NewHyperBlockProxy(hyperBlockFacade, *cfg)
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	router.GET(fmt.Sprintf("%s", cfg.HyperBlocksPath), hyperBlockProxy.GetHyperBlocksByInterval)
	router.GET(fmt.Sprintf("%s/by-nonce/:nonce", cfg.HyperBlockPath), hyperBlockProxy.GetHyperBlockByNonce)
	router.GET(fmt.Sprintf("%s/by-hash/:hash", cfg.HyperBlockPath), hyperBlockProxy.GetHyperBlockByHash)

	return &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", cfg.Port),
	}, nil
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
