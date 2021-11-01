package dataset

import "btradoc/entities"

type Service interface {
	FetchByDialect(fullDialect string) ([]entities.Dataset, error)
	AddTranslatedIn(datasetID, fullDialect string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FetchByDialect(fullDialect string) ([]entities.Dataset, error) {
	return s.repository.FetchByDialect(fullDialect)
}

func (s *service) AddTranslatedIn(datasetID, fullDialect string) error {
	return s.repository.AddTranslatedIn(datasetID, fullDialect)
}
