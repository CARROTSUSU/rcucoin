package abci

import (
	"encoding/json"
	"fmt"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type RCUCoinApp struct {
	abcitypes.BaseApplication
	balances map[string]int
}

func NewRCUCoinApp() *RCUCoinApp {
	return &RCUCoinApp{
		balances: map[string]int{
			"address1": 1000,
			"address2": 500,
		},
	}
}

func (app *RCUCoinApp) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{
		Data: "RCUCoin ABCI App",
	}
}

func (app *RCUCoinApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	var tx struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}

	err := json.Unmarshal(req.Tx, &tx)
	if err != nil {
		return abcitypes.ResponseDeliverTx{Code: 1, Log: "Invalid transaction"}
	}

	if app.balances[tx.From] < tx.Amount {
		return abcitypes.ResponseDeliverTx{Code: 1, Log: "Insufficient funds"}
	}

	app.balances[tx.From] -= tx.Amount
	app.balances[tx.To] += tx.Amount

	return abcitypes.ResponseDeliverTx{Code: 0, Log: "Transfer success"}
}
