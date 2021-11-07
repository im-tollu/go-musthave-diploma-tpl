package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/model"
	storage "github.com/im-tollu/go-musthave-diploma-tpl/storage/accrual"
)

type AccrualStorage struct {
	*sql.DB
}

func NewAccrualStorage(db *sql.DB) (*AccrualStorage, error) {
	if db == nil {
		return nil, errors.New("db should not be nil")
	}
	return &AccrualStorage{db}, nil
}

func (s *AccrualStorage) NextOrder() (int64, error) {
	row := s.QueryRow(`
		select ORDERS_NR
		from ORDERS
		where ORDERS_STATUS in ('NEW', 'PROCESSING')
		order by ORDERS_UPLOADED_AT
		limit 1
	`)

	var orderID int64

	err := row.Scan(&orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return orderID, storage.ErrNoOrders
	}
	if err != nil {
		return orderID, fmt.Errorf("cannot map order ID: %w", err)
	}

	return orderID, nil
}

func (s *AccrualStorage) ProcessOrder(o model.OrderAccrual) error {
	result, errExec := s.Exec(`
		update ORDERS
		set ACCRUAL = $1
		where ORDERS_NR = $2
			and ORDERS_STATUS = $3
	`, o.Accrual*100)
	if errExec != nil {
		return fmt.Errorf("cannot update order: %w", errExec)
	}

	affected, errAffected := result.RowsAffected()
	if errAffected != nil {
		return fmt.Errorf("cannot get affected rows: %w", errExec)
	}
	if affected != 1 {
		return fmt.Errorf("order not updated; expected 1 row affected, got %d", affected)
	}

	return nil
}
