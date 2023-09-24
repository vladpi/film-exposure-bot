package service

import "github.com/vladpi/film-exposure-bot/internal/domain/film"

type FilmService struct {
	repo film.Repository
}

func NewFilmService(repo film.Repository) *FilmService {
	return &FilmService{
		repo: repo,
	}
}

func (s *FilmService) GetAll() ([]film.Film, error) {
	return s.repo.GetAll()
}

func (s *FilmService) Get(id int64) (film.Film, error) {
	return s.repo.GetByID(id)
}
