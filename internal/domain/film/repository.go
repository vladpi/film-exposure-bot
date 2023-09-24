package film

type Repository interface {
	GetAll() ([]Film, error)
}
