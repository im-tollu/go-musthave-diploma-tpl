package pkg

type Scannable interface {
	Scan(dest ...interface{}) error
}
