package accrual

import (
	"errors"
	client "github.com/im-tollu/go-musthave-diploma-tpl/client/accrual"
	"github.com/im-tollu/go-musthave-diploma-tpl/model"
	storage "github.com/im-tollu/go-musthave-diploma-tpl/storage/accrual"
	"log"
	"time"
)

type Client interface {
	GetOrderAccruals(orderNr int64) (model.OrderAccrual, error)
}

type Service struct {
	client     *client.Client
	storage    storage.Storage
	orderQueue chan int64
}

const tickDuration = 12 * time.Second

func NewService(client *client.Client, storage storage.Storage) (*Service, error) {
	if storage == nil {
		return nil, errors.New("storage required")
	}

	j := Service{
		client:     client,
		storage:    storage,
		orderQueue: make(chan int64, 1),
	}

	//ticker := time.NewTicker(tickDuration)
	//
	//	for {
	//		select {
	//		case <-ticker.C:
	//			j.Run()
	//		case dur := <-j.pause:
	//			ticker.Stop()
	//			time.AfterFunc(dur, func() {
	//				resume <- struct{}{}
	//			})
	//		case <-resume:
	//			ticker.Reset(tickDuration)
	//		}
	//	}
	//

	return &j, nil
}

func (j *Service) pollDb() {
	for {
		orderID, errStorage := j.storage.NextOrder()
		if errors.Is(errStorage, storage.ErrNoOrders) {
			continue
		}
		if errStorage != nil {
			log.Printf("cannot get order for accrual because of DB: %s", errStorage.Error())
			continue
		}

		accrual, errClient := j.client.GetOrderAccruals(orderID)
		if errClient != nil {
			log.Printf("cannot get order for accrual because of service: %s", errClient.Error())
			continue
		}

		if errProcess := j.storage.ProcessOrder(accrual); errProcess != nil {
			log.Printf("cannot process apply accrual to order: %w", errProcess)
			continue
		}
	}
}

func (j *Service) Close() error {
	return nil
}