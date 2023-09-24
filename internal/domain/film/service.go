package film

type Service interface {
	GetAll() ([]Film, error)
	Get(id int64) (Film, error)
}
