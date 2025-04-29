package main

import (
	"log"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/ceocarrotsusu/rcucoin/abci"
)

func main() {
	app := abci.NewRCUCoinApp()

	srv, err := server.NewSocketServer("tcp://0.0.0.0:26658", app)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	defer srv.Stop()
	select {}
}
