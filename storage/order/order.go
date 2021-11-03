package order

import (
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
)

type Storage interface {
	AddOrder(pr order.ProcessRequest) error
	GetOrderByNr(nr int64) (order.Order, error)
	ListUserOrders(userID int64) ([]order.Order, error)
	Withdraw(wr order.WithdrawalRequest) error
	ListUserWithdrawals(userID int64) ([]order.Withdrawal, error)
}
