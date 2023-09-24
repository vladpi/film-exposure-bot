package film

type Repository interface {
	GetAll() ([]Film, error)
	GetByID(id int64) (Film, error)
}
