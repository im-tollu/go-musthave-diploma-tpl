package handler

import (
	"encoding/json"
	"github.com/im-tollu/go-musthave-diploma-tpl/api/apimodel"
	"log"
	"net/http"
)

func (h *LoyaltyHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		http.Error(w, "Content not expected", http.StatusBadRequest)
		return
	}

	uID := userID(r)

	orders, errList := h.OrderSrv.ListUserOrders(uID)
	if errList != nil {
		log.Printf("Cannot list orders for user [%d], %s", uID, errList.Error())
		http.Error(w, "Cannot list orders for user", http.StatusInternalServerError)
		return
	}

	view := make([]apimodel.OrderView, 0, len(orders))
	for _, order := range orders {
		view = append(view, apimodel.NewOrderView(order))
	}

	w.Header().Set("Content-Type", "application/json")
	if len(view) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	log.Printf("Listing user orders: %v", view)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(view); err != nil {
		log.Printf("Cannot write response: %s", err.Error())
	}
}
