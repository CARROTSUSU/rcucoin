package main

import (
    "log"
    "github.com/tendermint/tendermint/abci/server"
    "github.com/tendermint/tendermint/libs/cli"
    "github.com/CARROTSUSU/rcucoin/abci"
)

func main() {
    app := abci.NewRCUCoinApp()
    srv := server.NewSocketServer("tcp://0.0.0.0:26658", app)
    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
    defer srv.Stop()
    cli.TrapSignal(nil, func() {
        srv.Stop()
    })
    select {}
}
