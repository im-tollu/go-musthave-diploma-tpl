package pg

import (
	"database/sql"
	"errors"
	"fmt"
	srv "github.com/im-tollu/go-musthave-diploma-tpl/service/order"
	"github.com/im-tollu/go-musthave-diploma-tpl/storage/pkg"
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
		insert into ORDERS (ORDERS_NR, USERS_ID, ORDERS_STATUS) 
		values($1, $2, $3) 
		returning ORDERS_NR, USERS_ID, ORDERS_STATUS, ORDERS_UPLOADED_AT
		`, pr.Nr, pr.UserID, srv.StatusNew)

	order := srv.Order{}

	err := mapOrder(&order, row)
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

func (s *OrderStorage) GetOrderByNr(nr int64) (srv.Order, error) {
	row := s.QueryRow(`
		select ORDERS_NR, USERS_ID, ORDERS_STATUS, ORDERS_UPLOADED_AT
		from ORDERS
		where ORDERS_NR = $1
		`, nr)
	order := srv.Order{}

	if err := mapOrder(&order, row); err != nil {
		return order, fmt.Errorf("cannot select order: %w", err)
	}

	return order, nil
}

func (s *OrderStorage) ListUserWithdrawals(userID int64) ([]srv.Withdrawal, error) {
	result := make([]srv.Withdrawal, 0)

	rows, err := s.Query(`
		select WITHDRAWALS_NR, USERS_ID, WITHDRAWALS_SUM, WITHDRAWALS_STATUS, WITHDRAWALS_REQUESTED_AT
		from WITHDRAWALS
		where USERS_ID = $1
		order by WITHDRAWALS_REQUESTED_AT
	`,
		userID)
	if err != nil {
		return result, fmt.Errorf("cannot select withdrawals for user [%d]: %w", userID, err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Cannot close result set: %s", err.Error())
		}
	}(rows)

	for rows.Next() {
		w := srv.Withdrawal{}
		s := Money{}
		if err := rows.Scan(&w.OrderNr, &w.UserID, &s, &w.Status, &w.RequestedAt); err != nil {
			return result, fmt.Errorf("cannot map all withdrawals from DB: %w", err)
		}

		w.Sum = s.Rat
		result = append(result, w)
	}
	if rows.Err() != nil {
		return result, fmt.Errorf("cannot iterate all results from DB: %w", rows.Err())
	}

	return result, nil

}

func (s *OrderStorage) ListUserOrders(userID int64) ([]srv.Order, error) {
	result := make([]srv.Order, 0)

	rows, err := s.Query(`
		select ORDERS_NR, USERS_ID, ORDERS_STATUS, ORDERS_UPLOADED_AT
		from ORDERS
		where USERS_ID = $1
		order by ORDERS_UPLOADED_AT
	`,
		userID)
	if err != nil {
		return result, fmt.Errorf("cannot select orders for user [%d]: %w", userID, err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Cannot close result set: %s", err.Error())
		}
	}(rows)

	for rows.Next() {
		order := srv.Order{}

		if err := mapOrder(&order, rows); err != nil {
			return result, fmt.Errorf("cannot map all orders from DB: %w", err)
		}

		result = append(result, order)
	}
	if rows.Err() != nil {
		return result, fmt.Errorf("cannot iterate all results from DB: %w", rows.Err())
	}

	return result, nil
}

func (s *OrderStorage) Withdraw(wr srv.WithdrawalRequest) error {
	log.Printf("LatestAccrual: %v\nLatestWithdrawal: %v\n", wr.LatestAccrual, wr.LatestWithdrawal)
	row := s.QueryRow(`
			with NEW_WITHDRAWAL as (
				select
					$1::bigint as WITHDRAWALS_NR,
					$2::bigint as USERS_ID,
					$3::bigint as WITHDRAWALS_SUM,
					$4::text as WITHDRAWALS_STATUS,
					coalesce((select ORDERS_NR from ORDERS where USERS_ID = $2 and ORDERS_STATUS = 'PROCESSED' order by ORDERS_UPLOADED_AT desc limit 1), -1) as LATEST_ACCRUAL,
					coalesce((select WITHDRAWALS_NR from WITHDRAWALS where USERS_ID = $2 and WITHDRAWALS_STATUS = 'PROCESSED' order by WITHDRAWALS_REQUESTED_AT desc limit 1), -1) as LATEST_WITHDRAWAL
			)
			insert into WITHDRAWALS
			select WITHDRAWALS_NR, USERS_ID, WITHDRAWALS_SUM, WITHDRAWALS_STATUS
			from NEW_WITHDRAWAL
			where
					LATEST_ACCRUAL = $5
			  and LATEST_WITHDRAWAL = $6
			returning WITHDRAWALS_NR, USERS_ID, WITHDRAWALS_SUM, WITHDRAWALS_STATUS;
		`, wr.OrderNr, wr.UserID, NewMoney(wr.Sum), srv.StatusNew, wr.LatestAccrual, wr.LatestWithdrawal)

	sum := Money{}
	w := srv.Withdrawal{}
	err := row.Scan(&w.OrderNr, &w.UserID, &sum, &w.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("withdrowal not accepted because of conflict: %w", err)
	}
	if err != nil {
		return fmt.Errorf("cannot select order: %w", err)
	}

	w.Sum = sum.Rat

	return nil
}

func mapOrder(o *srv.Order, row pkg.Scannable) error {
	errScan := row.Scan(&o.Nr, &o.UserID, &o.Status, &o.UploadedAt)
	if errScan == sql.ErrNoRows {
		return srv.ErrOrderNotFound
	}
	if errScan != nil {
		return fmt.Errorf("cannot scan order from DB results: %w", errScan)
	}

	return nil
}