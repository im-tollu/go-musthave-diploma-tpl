package apimodel

import (
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
	"math/big"
	"strconv"
	"time"
)

type OrderView struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrderView(o order.Order) OrderView {
	return OrderView{
		Number:     fmt.Sprintf("%d", o.Nr),
		Status:     o.Status,
		Accrual:    0,
		UploadedAt: time.Time{},
	}
}

type BalanceView struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewBalanceView(b order.Balance) BalanceView {
	current, _ := b.Current.Float64()
	withdrawn, _ := b.Withdrawn.Float64()
	return BalanceView{
		Current:   current,
		Withdrawn: withdrawn,
	}
}

type WithdrawalRequestJSON struct {
	OrderNr string `json:"order"`
	Sum     int64  `json:"sum"`
}

func NewWithdrawalRequest(j WithdrawalRequestJSON, userID int64) (order.WithdrawalRequest, error) {
	wr := order.WithdrawalRequest{}

	orderNr, err := order.ParseOrderNr(j.OrderNr)
	if err != nil {
		return wr, fmt.Errorf("cannot make withdrawal request: %s", err)
	}

	wr.OrderNr = orderNr
	wr.Sum = big.NewRat(j.Sum*100, 100)
	wr.UserID = userID

	return wr, nil
}

type WithdrawalView struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func NewWithdrawalView(w order.Withdrawal) WithdrawalView {
	s, _ := w.Sum.Float64()
	return WithdrawalView{
		Order:       strconv.FormatInt(w.OrderNr, 10),
		Sum:         s,
		ProcessedAt: w.RequestedAt,
	}
}
