package main

import (
    "encoding/json"
    "fmt"
    "log"

    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/tendermint/tendermint/abci/server"
)

type RcpuCoinApp struct {
    abci.BaseApplication
    balances map[string]int
}

func NewRcpuCoinApp() *RcpuCoinApp {
    app := &RcpuCoinApp{
        balances: make(map[string]int),
    }
    app.balances["address1"] = 1000
    app.balances["address2"] = 500
    return app
}

func (app *RcpuCoinApp) DeliverTx(tx abci.RequestDeliverTx) abci.ResponseDeliverTx {
    var transfer struct {
        From   string `json:"from"`
        To     string `json:"to"`
        Amount int    `json:"amount"`
    }

    if err := json.Unmarshal(tx.Tx, &transfer); err != nil {
        return abci.ResponseDeliverTx{
            Code: 1,
            Log:  fmt.Sprintf("Error decoding tx: %v", err),
        }
    }

    if app.balances[transfer.From] < transfer.Amount {
        return abci.ResponseDeliverTx{
            Code: 1,
            Log:  "Insufficient balance",
        }
    }

    app.balances[transfer.From] -= transfer.Amount
    app.balances[transfer.To] += transfer.Amount

    fmt.Printf("Transferred %d RCU from %s to %s\n", transfer.Amount, transfer.From, transfer.To)

    return abci.ResponseDeliverTx{Code: 0, Log: "Transfer successful"}
}

func (app *RcpuCoinApp) CheckTx(tx abci.RequestCheckTx) abci.ResponseCheckTx {
    var transfer struct {
        From   string `json:"from"`
        To     string `json:"to"`
        Amount int    `json:"amount"`
    }

    if err := json.Unmarshal(tx.Tx, &transfer); err != nil {
        return abci.ResponseCheckTx{
            Code: 1,
            Log:  fmt.Sprintf("Error decoding tx: %v", err),
        }
    }

    if app.balances[transfer.From] < transfer.Amount {
        return abci.ResponseCheckTx{
            Code: 1,
            Log:  "Insufficient balance",
        }
    }

    return abci.ResponseCheckTx{Code: 0}
}

func (app *RcpuCoinApp) Info(req types.RequestInfo) types.ResponseInfo {
    // Menyediakan maklumat aplikasi, termasuk hash blok terkini
    blockInfo := types.ResponseInfo{
        Data:            "RCUCOIN Blockchain",
        Validators:      []types.Validator{},
        LatestBlockHash: []byte("hash-blok-sebelum"),
        LatestAppHash:   []byte("hash-aplikasi-sebelum"),
    }
    // Hash blok terkini yang akan ditetapkan selepas transaksi berlaku
    blockInfo.LatestBlockHash = app.latestBlockHash  // Dapatkan nilai daripada keadaan terkini
    return blockInfo
}

func main() {
    app := NewRcpuCoinApp()

// Set up your validators and their addresses
validator1 := "ADDRESS1_HEX"
validator2 := "ADDRESS2_HEX"

// Other code to initialize and use the validators...
    srv := server.NewSocketServer(":26658", app)

    if err := srv.Start(); err != nil {
        log.Fatalf("Failed to start ABCI server: %v", err)
    }

    defer srv.Stop()
    log.Println("ABCI server running on port 26658")
    select {}
}
