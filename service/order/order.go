package order

type Service interface {
	UploadOrder(pr ProcessRequest) error
	ListUserOrders(userID int64) ([]Order, error)
}
