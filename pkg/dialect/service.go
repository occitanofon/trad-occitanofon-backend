package dialect

import "btradoc/entities"

type Service interface {
	Exists(dialect, subdialect string) (bool, error)
	FetchOccitan(translatorID string) ([]entities.Occitan, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FetchOccitan(translatorID string) ([]entities.Occitan, error) {
	return s.repository.GetOccitanWithFurtherInfo(translatorID)
}

func (s *service) Exists(dialect, subdialect string) (bool, error) {
	return s.repository.IsItExists(dialect, subdialect)
}
