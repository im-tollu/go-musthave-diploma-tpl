package pkg

import (
	"database/sql/driver"
	"errors"
	"math/big"
)

type Money struct {
	*big.Rat
}

func NewMoney(sum *big.Rat) Money {
	return Money{sum}
}

func (m Money) Value() (driver.Value, error) {
	return m.Num().Int64(), nil
}

func (m *Money) Scan(value interface{}) error {
	if value == nil {
		*m = Money{big.NewRat(0, 100)}
		return nil
	}

	v, ok := value.(int64)
	if !ok {
		return errors.New("cannot scan value. cannot convert value to int64")
	}

	*m = Money{big.NewRat(v, 100)}

	return nil
}
