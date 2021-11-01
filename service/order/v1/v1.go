package v1

import (
	"errors"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/service/order"
	orderStorage "github.com/im-tollu/go-musthave-diploma-tpl/storage/order"
)

type Service struct {
	storage orderStorage.Storage
}

func NewService(storage orderStorage.Storage) (*Service, error) {
	if storage == nil {
		return nil, errors.New("storage required")
	}

	return &Service{storage}, nil
}

func (s *Service) ScheduleOrder(pr order.ProcessRequest) error {
	if errAdd := s.storage.AddOrder(pr); errAdd != nil {
		if errors.Is(errAdd, order.ErrDuplicateOrder) {
			dupO, errGet := s.storage.GetOrderByNr(pr.Nr)
			if errGet != nil {
				return fmt.Errorf("cannot get details of a duplicate order: %w", errGet)
			}

			if dupO.UserID == pr.UserID {
				return order.ErrDuplicateOrderForUser
			} else {
				return order.ErrDuplicateOrderForAnotherUser
			}
		}

		return fmt.Errorf("cannot schedule order for processing: %w", errAdd)
	}

	return nil
}
