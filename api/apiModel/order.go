package apiModel

import (
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
	"math/big"
	"time"
)

type OrderView struct {
	Number     int64     `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrderView(o order.Order) OrderView {
	return OrderView{
		Number:     o.Nr,
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

type WithdrawalRequestJson struct {
	OrderNr string `json:"order"`
	Sum     int64  `json:"sum"`
}

func NewWithdrawalRequest(j WithdrawalRequestJson, userID int64) (order.WithdrawalRequest, error) {
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
