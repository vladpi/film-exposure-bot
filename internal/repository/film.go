package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/vladpi/film-exposure-bot/internal/domain/film"
)

type SQLFilm struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	ISO  int64  `db:"iso"`
}

func (f *SQLFilm) ToEntity() film.Film {
	return film.Film{
		ID:   f.ID,
		Name: f.Name,
		ISO:  f.ISO,
	}
}

type SQLFilmRepository struct {
	db *sqlx.DB
}

func NewSQLFilmRepository(db *sqlx.DB) *SQLFilmRepository {
	return &SQLFilmRepository{db: db}
}

func (r *SQLFilmRepository) GetAll() ([]film.Film, error) {

	rows, err := r.db.Queryx("SELECT * FROM films")
	if err != nil {
		return []film.Film{}, err
	}

	var films []film.Film
	for rows.Next() {
		var f SQLFilm
		err = rows.StructScan(&f)
		if err != nil {
			return []film.Film{}, err
		}
		films = append(films, f.ToEntity())
	}

	return films, nil
}
