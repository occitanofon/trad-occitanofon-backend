package auth

import "btradoc/entities"

type Service interface {
	Login(translatorUsername string) (*entities.Translator, error)
	CreateRefreshToken(translatorID string) (string, error)
	FindByRefreshToken(refreshToken string) (*entities.Translator, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Login(translatorUsername string) (*entities.Translator, error) {
	return s.repository.FindByUsername(translatorUsername)
}

func (s *service) CreateRefreshToken(translatorID string) (string, error) {
	return s.repository.InsertRefreshToken(translatorID)
}

func (s *service) FindByRefreshToken(refreshToken string) (*entities.Translator, error) {
	return s.repository.GetTranslatorByRefreshToken(refreshToken)
}
