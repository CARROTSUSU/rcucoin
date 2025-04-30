package abci

import (
	"encoding/json"
	"fmt"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type RCUCoinApp struct {
	abcitypes.BaseApplication
	balances map[string]int64
}

func NewRCUCoinApp() *RCUCoinApp {
	return &RCUCoinApp{
		balances: map[string]int64{
			"rcu1": 1000000,
			"rcu2": 500000,
		},
	}
}

type Tx struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"` // dalam unit Miroc (1 RCU = 1,000,000 Miroc)
}

func (app *RCUCoinApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	var tx Tx
	if err := json.Unmarshal(req.Tx, &tx); err != nil {
		return abcitypes.ResponseDeliverTx{Code: 1, Log: "Invalid TX format"}
	}

	if app.balances[tx.From] < tx.Amount {
		return abcitypes.ResponseDeliverTx{Code: 1, Log: "Insufficient funds"}
	}

	app.balances[tx.From] -= tx.Amount
	app.balances[tx.To] += tx.Amount

	return abcitypes.ResponseDeliverTx{Code: 0, Log: "Transaction successful"}
}

func (app *RCUCoinApp) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{
		Data: "RCUCoin v1.0",
	}
}
