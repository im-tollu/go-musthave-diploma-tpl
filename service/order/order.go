package order

type Service interface {
	UploadOrder(pr ProcessRequest) error
	ListUserOrders(userID int64) ([]Order, error)
	GetUserBalance(userID int64) (Balance, error)
	Withdraw(wr WithdrawalRequest) error
	ListUserWithdrawals(userID int64) ([]Withdrawal, error)
}
