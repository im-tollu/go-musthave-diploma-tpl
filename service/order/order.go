package order

type Service interface {
	ScheduleOrder(pr ProcessRequest) error
}
