package handler

import (
	"encoding/json"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/apiModel"
	"log"
	"net/http"
)

func (h *LoyaltyHandler) Balance(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		http.Error(w, "Content not expected", http.StatusBadRequest)
		return
	}

	uID := userID(r)

	balance, errBalance := h.OrderSrv.GetUserBalance(uID)
	if errBalance != nil {
		log.Printf("Cannot get balance for user [%d]: %e", uID, errBalance.Error())
		http.Error(w, "Cannot get balance for user", http.StatusInternalServerError)
		return
	}

	view := apiModel.NewBalanceView(balance)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(view); err != nil {
		log.Printf("Cannot write response: %s", err.Error())
	}
}
