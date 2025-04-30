package main

import (
	"encoding/json"
	"fmt"
	"log"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/abci/server"
)

const MirocPerRCU = 1_000_000

type RcuCoinApp struct {
	abci.BaseApplication
	balances map[string]int64 // simpan dalam unit Miroc
}

func NewRcuCoinApp() *RcuCoinApp {
	return &RcuCoinApp{
		balances: map[string]int64{
			"address1": 10 * MirocPerRCU, // 10 RCUCoin
			"address2": 5 * MirocPerRCU,  // 5 RCUCoin
		},
	}
}

type TransferTx struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"` // dalam Miroc
}

func (app *RcuCoinApp) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	var t TransferTx
	err := json.Unmarshal(tx, &t)
	if err != nil {
		return abci.ResponseDeliverTx{Code: 1, Log: fmt.Sprintf("Decode error: %v", err)}
	}

	if t.Amount <= 0 {
		return abci.ResponseDeliverTx{Code: 1, Log: "Invalid amount"}
	}

	if app.balances[t.From] < t.Amount {
		return abci.ResponseDeliverTx{Code: 1, Log: "Insufficient balance"}
	}

	app.balances[t.From] -= t.Amount
	app.balances[t.To] += t.Amount

	logStr := fmt.Sprintf("Transfer %d Miroc (%.6f RCU) from %s to %s",
		t.Amount, float64(t.Amount)/float64(MirocPerRCU), t.From, t.To)

	return abci.ResponseDeliverTx{Code: 0, Log: logStr}
}

func (app *RcuCoinApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{Data: "RCUCoin Miroc Blockchain"}
}

func main() {
	app := NewRcuCoinApp()
	srv := server.NewSocketServer("tcp://127.0.0.1:26658", app)

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start ABCI server: %v", err)
	}
	defer srv.Stop()

	select {} // keep running
}
