package apiModel

import (
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
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
