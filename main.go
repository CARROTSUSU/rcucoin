package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/abci/server"
)

// Struktur aplikasi RCUCoin
type RCUCoinApp struct {
	types.BaseApplication
	balances map[string]int64 // Gunakan int64 untuk unit mikro (Miroc)
}

// Unit minimum
const Unit = 1_000_000 // 1 RCUCoin = 1,000,000 Miroc

// Cipta aplikasi
func NewRCUCoinApp() *RCUCoinApp {
	return &RCUCoinApp{
		balances: map[string]int64{
			"address1": 10 * Unit, // 10 RCU
			"address2": 5 * Unit,  // 5 RCU
		},
	}
}

// Terima transaksi
func (app *RCUCoinApp) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	var tx struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int64  `json:"amount"` // dalam Miroc
	}
	err := json.Unmarshal(req.Tx, &tx)
	if err != nil {
		return types.ResponseDeliverTx{Code: 1, Log: "Invalid transaction format"}
	}

	if app.balances[tx.From] < tx.Amount {
		return types.ResponseDeliverTx{Code: 2, Log: "Insufficient balance"}
	}

	app.balances[tx.From] -= tx.Amount
	app.balances[tx.To] += tx.Amount

	fmt.Printf("Transferred %d Miroc from %s to %s\n", tx.Amount, tx.From, tx.To)
	return types.ResponseDeliverTx{Code: 0, Log: "Transaction successful"}
}

// Info aplikasi
func (app *RCUCoinApp) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{
		Data:             "RCUCoin ABCI",
		AppVersion:       1,
		LastBlockHeight:  0,
		LastBlockAppHash: []byte{},
	}
}

// Fungsi utama
func main() {
	app := NewRCUCoinApp()

	srv, err := server.NewSocketServer("tcp://127.0.0.1:26658", app)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	srv.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	// Signal trap (ganti cli.TrapSignal)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("Shutting down ABCI server...")
		if err := srv.Stop(); err != nil {
			log.Fatalf("Error stopping server: %v", err)
		}
	}()

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting ABCI server: %v", err)
	}

	// Biarkan terus jalan
	select {}
}
