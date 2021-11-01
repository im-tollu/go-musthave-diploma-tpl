package pg

import (
	"database/sql"
	"errors"
	"fmt"
	srv "github.com/im-tollu/go-musthave-diploma-tpl/service/order"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"log"
)

type OrderStorage struct {
	*sql.DB
}

func NewOrderStorage(db *sql.DB) (*OrderStorage, error) {
	if db == nil {
		return nil, errors.New("db should not be nil")
	}
	return &OrderStorage{db}, nil
}

func (s *OrderStorage) AddOrder(pr srv.ProcessRequest) error {
	row := s.QueryRow(`
		insert into ORDERS (ORDERS_NR, USERS_ID) 
		values($1, $2) 
		returning ORDERS_NR, USERS_ID
		`, pr.Nr, pr.UserID)
	order := srv.ProcessRequest{}

	err := row.Scan(&order.Nr, &order.UserID)
	var dbErr *pgconn.PgError
	if errors.As(err, &dbErr) && dbErr.Code == pgerrcode.UniqueViolation {
		log.Printf("Duplicate order [%d]", pr.Nr)
		err = srv.ErrDuplicateOrder
	}
	if err != nil {
		return fmt.Errorf("cannot insert order: %w", err)
	}

	return nil
}

func (s *OrderStorage) GetOrderByNr(nr int64) (srv.ProcessRequest, error) {
	row := s.QueryRow(`
		select ORDERS_NR, USERS_ID 
		from ORDERS
		where ORDERS_NR = $1
		`, nr)
	order := srv.ProcessRequest{}

	if err := row.Scan(&order.Nr, &order.UserID); err != nil {
		return order, fmt.Errorf("cannot select order: %w", err)
	}

	return order, nil
}
