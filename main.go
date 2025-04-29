package main

import (
    "log"

    "github.com/ceocarrotsusu/rcucoin/abci"
    "github.com/tendermint/tendermint/abci/server"
    abcitypes "github.com/tendermint/tendermint/abci/types"
)

func main() {
    app := abci.NewRCUCoinApp()
    s := server.NewSocketServer("tcp://0.0.0.0:26658", app)
    if err := s.Start(); err != nil {
        log.Fatalf("Failed to start ABCI server: %v", err)
    }
    defer s.Stop()
    select {}
}
