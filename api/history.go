package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type btcHistoryReq struct {
	StartDate time.Time `json:"startDatetime"`
	EndDate   time.Time `json:"endDatetime"`
}

func (h *Handler) GetBTCHistory(w http.ResponseWriter, r *http.Request) {
	btcHistoryReq := btcHistoryReq{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&btcHistoryReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(btcRes{
			ResultMessage: err.Error(),
		})
		return
	}
	defer r.Body.Close()

	btcHistory, err := h.store.GetBTCHistory(btcHistoryReq.StartDate, btcHistoryReq.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(btcRes{
			ResultMessage: err.Error(),
		})
		return
	}

	if len(btcHistory) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(btcHistory)
	return
}
