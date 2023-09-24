package film

type Service interface {
	GetAll() ([]Film, error)
}
