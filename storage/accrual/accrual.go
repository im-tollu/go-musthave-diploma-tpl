package accrual

import (
	"errors"
	"github.com/im-tollu/go-musthave-diploma-tpl/model"
)

var ErrNoOrders = errors.New("no more orders to process")

type Storage interface {
	NextOrder() (int64, error)
	ApplyAccrual(o model.OrderAccrual) error
}
