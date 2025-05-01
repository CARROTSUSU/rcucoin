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
    balances         map[string]int
    latestAppHash    []byte
    latestBlockHash  []byte
}

// Fungsi untuk cipta aplikasi baru
func NewRcpuCoinApp() *RcpuCoinApp {
    app := &RcpuCoinApp{
        balances:         make(map[string]int),
        latestAppHash:    []byte("genesis-app-hash"),
        latestBlockHash:  []byte("genesis-block-hash"),
    }

    // Baki permulaan untuk dua alamat (anda boleh ganti dengan alamat sebenar)
    app.balances["rcpuaddress1"] = 1000
    app.balances["rcpuaddress2"] = 500

    return app
}

// Inisialisasi rantai (jika perlu baca validator atau maklumat permulaan lain)
func (app *RcpuCoinApp) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
    fmt.Println("Chain initialized with genesis")
    return abci.ResponseInitChain{}
}

// Semak transaksi sebelum ia dihantar ke blok
func (app *RcpuCoinApp) CheckTx(req abci.RequestCheckTx) abci.ResponseCheckTx {
    var tx struct {
        From   string `json:"from"`
        To     string `json:"to"`
        Amount int    `json:"amount"`
    }

    if err := json.Unmarshal(req.Tx, &tx); err != nil {
        return abci.ResponseCheckTx{Code: 1, Log: fmt.Sprintf("Invalid tx format: %v", err)}
    }

    if app.balances[tx.From] < tx.Amount {
        return abci.ResponseCheckTx{Code: 1, Log: "Insufficient balance"}
    }

    return abci.ResponseCheckTx{Code: 0}
}

// Laksanakan transaksi sebenar dan kemas kini baki
func (app *RcpuCoinApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
    var tx struct {
        From   string `json:"from"`
        To     string `json:"to"`
        Amount int    `json:"amount"`
    }

    if err := json.Unmarshal(req.Tx, &tx); err != nil {
        return abci.ResponseDeliverTx{Code: 1, Log: fmt.Sprintf("Invalid tx format: %v", err)}
    }

    if app.balances[tx.From] < tx.Amount {
        return abci.ResponseDeliverTx{Code: 1, Log: "Insufficient balance"}
    }

    app.balances[tx.From] -= tx.Amount
    app.balances[tx.To] += tx.Amount

    // Simpan app hash
    app.latestAppHash = []byte(fmt.Sprintf("apphash-%s-%s-%d", tx.From, tx.To, tx.Amount))

    fmt.Printf("Transferred %d RCU from %s to %s\n", tx.Amount, tx.From, tx.To)
    return abci.ResponseDeliverTx{Code: 0, Log: "Transfer successful"}
}

// Maklumat blok semasa & aplikasi
func (app *RcpuCoinApp) Info(req abci.RequestInfo) abci.ResponseInfo {
    return abci.ResponseInfo{
        Data:              "RCUCOIN Blockchain",
        LastBlockHeight:   1,
        LastBlockAppHash:  app.latestAppHash,
    }
}

// Hantar hash terkini setiap kali blok selesai
func (app *RcpuCoinApp) Commit() abci.ResponseCommit {
    app.latestBlockHash = []byte(fmt.Sprintf("blockhash-%x", app.latestAppHash))
    return abci.ResponseCommit{Data: app.latestBlockHash}
}

// Fungsi utama jalankan ABCI server
func main() {
    app := NewRcpuCoinApp()
    srv := server.NewSocketServer(":26658", app)

    if err := srv.Start(); err != nil {
        log.Fatalf("Failed to start ABCI server: %v", err)
    }

    defer srv.Stop()
    log.Println("RCUCOIN ABCI server running on port 26658")
    select {}
}
