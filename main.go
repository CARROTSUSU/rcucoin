package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "os"

    abci "github.com/tendermint/tendermint/abci/types"
    "github.com/tendermint/tendermint/abci/server"
)

type AppState struct {
    Balances        map[string]int `json:"balances"`
    LatestAppHash   string         `json:"latest_app_hash"`
    LatestBlockHash string         `json:"latest_block_hash"`
    LastBlockHeight int64          `json:"last_block_height"`
}

type RcpuCoinApp struct {
    abci.BaseApplication
    balances        map[string]int
    latestAppHash   []byte
    latestBlockHash []byte
    lastBlockHeight int64
}

func NewRcpuCoinApp() *RcpuCoinApp {
    app := &RcpuCoinApp{
        balances:        make(map[string]int),
        latestAppHash:   []byte{},
        latestBlockHash: []byte{},
        lastBlockHeight: 0,
    }
    app.LoadState()
    return app
}

func (app *RcpuCoinApp) LoadState() {
    data, err := os.ReadFile("state.json")
    if err != nil {
        log.Println("Tiada state.json, mula dari kosong")
        app.balances["9758A0E9A531642AE9E781BBDCE8F1298501298501BFFB"] = 1000
        app.balances["FC8EA6FB4A04F93845D2F8E8E6ED63F698965499"] = 500
        return
    }

    var state AppState
    if err := json.Unmarshal(data, &state); err != nil {
        log.Println("Gagal parse state.json:", err)
        return
    }

    app.balances = state.Balances
    app.latestAppHash, _ = hex.DecodeString(state.LatestAppHash)
    app.latestBlockHash, _ = hex.DecodeString(state.LatestBlockHash)
    app.lastBlockHeight = state.LastBlockHeight
    log.Println("Berjaya load state.json")
}

func (app *RcpuCoinApp) saveState() {
    state := AppState{
        Balances:        app.balances,
        LatestAppHash:   hex.EncodeToString(app.latestAppHash),
        LatestBlockHash: hex.EncodeToString(app.latestBlockHash),
        LastBlockHeight: app.lastBlockHeight,
    }

    data, err := json.MarshalIndent(state, "", "  ")
    if err != nil {
        log.Println("Gagal serialize state:", err)
        return
    }

    _ = os.WriteFile("state.json", data, 0644)
    log.Println("State disimpan ke state.json")
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
            Log:  fmt.Sprintf("Gagal decode tx: %v", err),
        }
    }

    if app.balances[transfer.From] < transfer.Amount {
        return abci.ResponseDeliverTx{
            Code: 1,
            Log:  "Baki tidak mencukupi",
        }
    }

    app.balances[transfer.From] -= transfer.Amount
    app.balances[transfer.To] += transfer.Amount

    log.Printf("Pindahan %d RCU dari %s ke %s", transfer.Amount, transfer.From, transfer.To)

    return abci.ResponseDeliverTx{Code: 0, Log: "Pindahan berjaya"}
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
            Log:  fmt.Sprintf("Gagal decode tx: %v", err),
        }
    }

    if app.balances[transfer.From] < transfer.Amount {
        return abci.ResponseCheckTx{
            Code: 1,
            Log:  "Baki tidak mencukupi",
        }
    }

    return abci.ResponseCheckTx{Code: 0}
}

func (app *RcpuCoinApp) Commit() abci.ResponseCommit {
    balancesBytes, _ := json.Marshal(app.balances)
    hash := sha256.Sum256(balancesBytes)
    app.latestAppHash = hash[:]
    app.latestBlockHash = hash[:]
    app.lastBlockHeight++

    app.saveState()

    return abci.ResponseCommit{Data: app.latestAppHash}
}

func (app *RcpuCoinApp) Info(req abci.RequestInfo) abci.ResponseInfo {
    return abci.ResponseInfo{
        Data:             "RCUCOIN Blockchain ABCI",
        LastBlockHeight:  app.lastBlockHeight,
        LastBlockAppHash: app.latestAppHash,
    }
}

func main() {
    app := NewRcpuCoinApp()
    srv := server.NewSocketServer(":26658", app)

    if err := srv.Start(); err != nil {
        log.Fatalf("Gagal mulakan ABCI server: %v", err)
    }

    defer srv.Stop()
    log.Println("ABCI server berjalan di port 26658")
    select {}
}
