package accrual

import (
	"errors"
	"fmt"
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
	resultChan chan error
	resumeChan chan struct{}
	stopChan   chan struct{}
}

func NewService(client *client.Client, storage storage.Storage) (*Service, error) {
	if storage == nil {
		return nil, errors.New("storage required")
	}

	srv := Service{
		client:     client,
		storage:    storage,
		resultChan: make(chan error, 1),
		resumeChan: make(chan struct{}, 1),
		stopChan:   make(chan struct{}, 1),
	}

	go srv.run()

	return &srv, nil
}

func (s *Service) run() {
	s.processOrder()
	for {
		select {
		case res := <-s.resultChan:
			tooMayRequests := &model.ErrTooManyRequests{}
			if errors.Is(res, storage.ErrNoOrders) {
				log.Println("No orders to process, waiting...")
				time.AfterFunc(10*time.Second, s.resume)
			} else if errors.As(res, &tooMayRequests) {
				log.Println("Too many requests, waiting...")
				time.AfterFunc(tooMayRequests.RetryAfter, s.resume)
			} else if res != nil {
				log.Printf("Could not process order: %s", res.Error())
			} else {
				s.processOrder()
			}
		case <-s.resumeChan:
			s.processOrder()
		case <-s.stopChan:
			log.Println("Stopping Accruals service...")
			close(s.stopChan)
			return
		}
	}
}

func (s *Service) resume() {
	s.resumeChan <- struct{}{}
}

func (s *Service) processOrder() {
	orderID, errStorage := s.storage.NextOrder()
	if errStorage != nil {
		err := fmt.Errorf("cannot get order for accrual because of DB: %w", errStorage)
		s.resultChan <- err
		return
	}

	log.Printf("Processing order [%d]", orderID)
	accrual, errClient := s.client.GetOrderAccruals(orderID)
	if errClient != nil {
		err := fmt.Errorf("cannot get order for accrual because of service: %w", errClient)
		s.resultChan <- err
		return
	}
	log.Printf("Got accrual: %v", accrual)

	if errApply := s.storage.ApplyAccrual(accrual); errApply != nil {
		err := fmt.Errorf("cannot process apply accrual to order: %w", errApply)
		s.resultChan <- err
		return
	}
}

func (s *Service) Stop() {
	s.stopChan <- struct{}{}
	<-s.stopChan
}
