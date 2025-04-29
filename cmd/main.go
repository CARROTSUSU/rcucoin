package main

import (
	"log"

	"github.com/tendermint/tendermint/abci/server"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/ceocarrotsusu/rcucoin/abci"
)

func main() {
	app := abci.NewRCUCoinApp()

	srv, err := server.NewSocketServer("tcp://0.0.0.0:26658", app)
	if err != nil {
		log.Fatal(err)
	}

	srv.Start()
	defer srv.Stop()
	select {}
}
