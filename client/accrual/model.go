package accrual

import "github.com/im-tollu/go-musthave-diploma-tpl/model"

type OrderAccrualJSON struct {
	OrderNr string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

func (j OrderAccrualJSON) ToOrderAccrual() (model.OrderAccrual, error) {
	result := model.OrderAccrual{}

	return result, nil
}
