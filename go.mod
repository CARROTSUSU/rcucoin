package main

import (
    "log"

    "github.com/ceocarrotsusu/rcucoin/abci"
    "github.com/tendermint/tendermint/abci/server"
)

func main() {
    app := abci.NewRCUCoinApp()
    srv, err := server.NewSocketServer("tcp://0.0.0.0:26658", app)
    if err != nil {
        log.Fatalf("Failed to start ABCI server: %v", err)
    }

    err = srv.Start()
    if err != nil {
        log.Fatalf("Failed to run ABCI server: %v", err)
    }

    defer srv.Stop()

    select {}
}
