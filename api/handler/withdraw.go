package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/apiModel"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
	"log"
	"net/http"
)

func (h *LoyaltyHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		http.Error(w, "Expected application/json content type", http.StatusBadRequest)
		return
	}

	j := apimodel.WithdrawalRequestJSON{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&j); err != nil {
		msg := fmt.Sprintf("Invalid withdrawal request json: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	uID := userID(r)
	wr, errReq := apimodel.NewWithdrawalRequest(j, uID)
	if errors.Is(errReq, order.ErrInvalidOrderNr) {
		msg := fmt.Sprintf("Invalid order nr: %s", errReq.Error())
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}
	if errReq != nil {
		log.Printf("Cannot withdraw: %s", errReq.Error())
		http.Error(w, "Cannot withdraw", http.StatusInternalServerError)
		return
	}

	errWithdraw := h.OrderSrv.Withdraw(wr)
	if errors.Is(errWithdraw, order.ErrInsufficientBalance) {
		http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
		return
	}
	if errWithdraw != nil {
		log.Printf("Cannot withdraw for user [%d]: %s", uID, errWithdraw.Error())
		http.Error(w, "Cannot get balance for user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
