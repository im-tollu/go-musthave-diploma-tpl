package accrual

import (
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/model"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
)

type OrderAccrualJSON struct {
	OrderNr string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (j OrderAccrualJSON) ToOrderAccrual() (model.OrderAccrual, error) {
	nr, errNr := order.ParseOrderNr(j.OrderNr)
	if errNr != nil {
		return model.OrderAccrual{}, fmt.Errorf("envalid order nr [%s]: %w", j.OrderNr, errNr)
	}
	result := model.OrderAccrual{
		OrderNr: nr,
		Status:  j.Status,
		Accrual: int64(j.Accrual * 100),
	}

	if result.Status == model.StatusRegistered {
		result.Status = order.StatusProcessing
	} else if result.Status == model.StatusProcessing {
		result.Status = order.StatusProcessing
	} else if result.Status == model.StatusProcessed {
		result.Status = order.StatusProcessed
	} else {
		result.Status = order.StatusInvalid
	}

	return result, nil
}
