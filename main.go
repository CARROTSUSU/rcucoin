package main

import (
	"fmt"
	"log"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/abci/server"
)

// Struktur aplikasi RcpuCoin
type RcpuCoinApp struct {
	types.BaseApplication
	balances map[string]int // Penyimpanan dalam ingatan untuk baki akaun
}

// Fungsi untuk mencipta aplikasi RcpuCoin baru
func NewRcpuCoinApp() *RcpuCoinApp {
	app := &RcpuCoinApp{
		balances: make(map[string]int),
	}

	// Menambah beberapa akaun permulaan dengan baki
	app.balances["address1"] = 1000
	app.balances["address2"] = 500

	return app
}

// Menerima dan memproses transaksi
func (app *RcpuCoinApp) DeliverTx(tx types.RequestDeliverTx) types.ResponseDeliverTx {
	var transfer struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}

	// Mendekod transaksi
	err := json.Unmarshal(tx.Tx, &transfer)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: 1,
			Log:  fmt.Sprintf("Error decoding transaction: %v", err),
		}
	}

	// Validasi jika penghantar ada baki yang mencukupi
	if app.balances[transfer.From] < transfer.Amount {
		return types.ResponseDeliverTx{
			Code: 1,
			Log:  "Insufficient balance",
		}
	}

	// Mengemas kini baki akaun
	app.balances[transfer.From] -= transfer.Amount
	app.balances[transfer.To] += transfer.Amount

	// Log transaksi
	fmt.Printf("Transfer %d RCPU Coin from %s to %s\n", transfer.Amount, transfer.From, transfer.To)

	// Menyediakan respons untuk transaksi yang berjaya
	return types.ResponseDeliverTx{
		Code: 0,
		Log:  "Transfer successful",
	}
}

// Fungsi untuk mendapatkan informasi aplikasi
func (app *RcpuCoinApp) Info(req types.RequestInfo) types.ResponseInfo {
	// Menyediakan maklumat aplikasi
	return types.ResponseInfo{
		Data:             "RCUCoin Blockchain",
		Validators:       []types.Validator{},
		LatestBlockHash:  []byte("latest-block-hash"),
		LatestAppHash:    []byte("latest-app-hash"),
	}
}

// Fungsi untuk memeriksa konsistensi dan keadaan blockchain
func (app *RcpuCoinApp) CheckTx(tx types.RequestCheckTx) types.ResponseCheckTx {
	var transfer struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}

	// Mendekod transaksi
	err := json.Unmarshal(tx.Tx, &transfer)
	if err != nil {
		return types.ResponseCheckTx{
			Code: 1,
			Log:  fmt.Sprintf("Error decoding transaction: %v", err),
		}
	}

	// Validasi jika penghantar ada baki yang mencukupi
	if app.balances[transfer.From] < transfer.Amount {
		return types.ResponseCheckTx{
			Code: 1,
			Log:  "Insufficient balance",
		}
	}

	// Validasi transaksi
	return types.ResponseCheckTx{
		Code: 0,
	}
}

// Fungsi utama untuk menjalankan aplikasi
func main() {
	// Cipta aplikasi baru
	app := NewRcpuCoinApp()

	// Mulakan server ABCI untuk aplikasi ini pada port :26658
	server := types.NewSocketServer(":26658", app)
	log.Fatal(server.Start())
}
