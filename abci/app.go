package abci

import (
	"encoding/json"
	"fmt"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

const MirocPerRCU = 1_000_000

type RcuCoinApp struct {
	abcitypes.BaseApplication
	balances map[string]int64 // disimpan dalam Miroc
}

func NewRcuCoinApp() *RcuCoinApp {
	return &RcuCoinApp{
		balances: map[string]int64{
			"genesis1": 100 * MirocPerRCU,
			"genesis2": 50 * MirocPerRCU,
		},
	}
}

// Struktur transaksi
type TransferTx struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"` // dalam Miroc
}

func (app *RcuCoinApp) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{
		Data: fmt.Sprintf("RCUCoin ABCI - Powered by Miroc"),
	}
}

func (app *RcuCoinApp) CheckTx(tx []byte) abcitypes.ResponseCheckTx {
	return app.validateTx(tx)
}

func (app *RcuCoinApp) DeliverTx(tx []byte) abcitypes.ResponseDeliverTx {
	return app.executeTx(tx)
}

func (app *RcuCoinApp) Query(req abcitypes.RequestQuery) abcitypes.ResponseQuery {
	address := string(req.Data)
	balance, ok := app.balances[address]
	if !ok {
		return abcitypes.ResponseQuery{Code: 1, Log: "Address not found"}
	}

	return abcitypes.ResponseQuery{
		Code:  0,
		Log:   "Balance retrieved",
		Value: []byte(fmt.Sprintf("%d", balance)),
	}
}

func (app *RcuCoinApp) validateTx(tx []byte) abcitypes.ResponseCheckTx {
	var t TransferTx
	if err := json.Unmarshal(tx, &t); err != nil {
		return abcitypes.ResponseCheckTx{Code: 1, Log: fmt.Sprintf("Invalid JSON: %v", err)}
	}

	if t.Amount <= 0 {
		return abcitypes.ResponseCheckTx{Code: 1, Log: "Amount must be positive"}
	}

	if app.balances[t.From] < t.Amount {
		return abcitypes.ResponseCheckTx{Code: 1, Log: "Insufficient funds"}
	}

	return abcitypes.ResponseCheckTx{Code: 0}
}

func (app *RcuCoinApp) executeTx(tx []byte) abcitypes.ResponseDeliverTx {
	var t TransferTx
	if err := json.Unmarshal(tx, &t); err != nil {
		return abcitypes.ResponseDeliverTx{Code: 1, Log: fmt.Sprintf("Invalid JSON: %v", err)}
	}

	app.balances[t.From] -= t.Amount
	app.balances[t.To] += t.Amount

	log := fmt.Sprintf("Transferred %d Miroc (%.6f RCU) from %s to %s",
		t.Amount, float64(t.Amount)/float64(MirocPerRCU), t.From, t.To)

	return abcitypes.ResponseDeliverTx{Code: 0, Log: log}
}
