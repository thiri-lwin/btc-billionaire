package api

import (
	"encoding/json"
	"net/http"
	"time"

	"btc_billionaire/postgres"

	"github.com/shopspring/decimal"
)

type sendBTCReq struct {
	Date   time.Time       `json:"datetime"`
	Amount decimal.Decimal `json:"amount"`
}

type btcRes struct {
	ResultMessage string `json:"ResultMessage"`
}

func (h *Handler) SendBTC(w http.ResponseWriter, r *http.Request) {
	btcInfoReq := sendBTCReq{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&btcInfoReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(btcRes{
			ResultMessage: err.Error(),
		})
		return
	}
	defer r.Body.Close()

	// validate input data
	if btcInfoReq.Amount.BigInt().Int64() <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(btcRes{
			ResultMessage: "Invalid amount",
		})
		return
	}

	_, offset := btcInfoReq.Date.Zone()
	// insert btc info into db
	_, err := h.store.InsertBTCInfo(postgres.BTCInfo{
		Amount:  btcInfoReq.Amount,
		Created: btcInfoReq.Date,
		Offset:  offset,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(btcRes{
			ResultMessage: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(btcRes{
		ResultMessage: "BTC was sent successfully",
	})
	return
}
