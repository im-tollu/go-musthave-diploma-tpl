package order

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrDuplicateOrder = errors.New("duplicate order")

var ErrDuplicateOrderForUser = errors.New("order already posted by the user")

var ErrDuplicateOrderForAnotherUser = errors.New("order already posted by another user")

var ErrInvalidOrderNr = errors.New("invalid order number")

func ParseOrderNr(s string) (int64, error) {
	nr, errConv := strconv.ParseInt(s, 10, 64)
	if errConv != nil {
		return nr, fmt.Errorf("cannot parse order number: %w", errConv)
	}

	n := len(s)
	checksum := 0

	for i := 1; i <= len(s); i++ {
		d, err := strconv.Atoi(string(s[n-i]))
		if err != nil {
			return nr, fmt.Errorf("contains non-digit character: %w", ErrInvalidOrderNr)
		}

		if i%2 == 0 {
			s := 2 * d
			if s > 9 {
				s -= 9
			}
			checksum += s
		} else {
			checksum += d
		}
	}

	if checksum%10 != 0 {
		return nr, fmt.Errorf("incorrect checksum: %w", ErrInvalidOrderNr)
	}

	return nr, nil
}

type ProcessRequest struct {
	Nr     int64
	UserID int64
}
