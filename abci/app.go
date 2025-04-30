package abci

import (
    "encoding/json"
    "fmt"
    abcitypes "github.com/tendermint/tendermint/abci/types"
)

type RCUCoinApp struct {
    abcitypes.BaseApplication
    balances map[string]int64 // simpan dalam Miroc
}

func NewRCUCoinApp() *RCUCoinApp {
    return &RCUCoinApp{
        balances: map[string]int64{
            "address1": 5_000_000_000, // 5,000 RCU = 5B Miroc
            "address2": 2_000_000_000,
        },
    }
}

func (app *RCUCoinApp) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
    return abcitypes.ResponseInfo{Data: "RCUCoin (Miroc-based)"}
}

func (app *RCUCoinApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
    var tx struct {
        From  string `json:"from"`
        To    string `json:"to"`
        Miroc int64  `json:"miroc"` // dalam Miroc
    }

    if err := json.Unmarshal(req.Tx, &tx); err != nil {
        return abcitypes.ResponseDeliverTx{Code: 1, Log: "Invalid JSON"}
    }

    if tx.Miroc <= 0 || app.balances[tx.From] < tx.Miroc {
        return abcitypes.ResponseDeliverTx{Code: 1, Log: "Invalid or insufficient balance"}
    }

    app.balances[tx.From] -= tx.Miroc
    app.balances[tx.To] += tx.Miroc

    fmt.Printf("Transferred %d Miroc from %s to %s\n", tx.Miroc, tx.From, tx.To)

    return abcitypes.ResponseDeliverTx{Code: 0, Log: "Success"}
}
