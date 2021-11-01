package account

import "btradoc/entities"

type Service interface {
	Create(newTranslator entities.NewTranslator) error
	FetchAllSecretQuestions() ([]string, error)
	FetchSecretQuestionsAndResponses(token string) (*entities.TranslatorSecretQuestions, error)
	ResetPassword(email string) (*entities.TranslatorResetPassword, error)
	UpdatePassword(translatorID string, newHashedPassword string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(newTranslator entities.NewTranslator) error {
	return s.repository.InsertNewTranslator(newTranslator)
}

func (s *service) FetchAllSecretQuestions() ([]string, error) {
	return s.repository.GetAllSecretQuestions()
}

func (s *service) FetchSecretQuestionsAndResponses(token string) (*entities.TranslatorSecretQuestions, error) {
	return s.repository.GetTranslatorSecretQuestions(token)
}

func (s *service) ResetPassword(email string) (*entities.TranslatorResetPassword, error) {
	return s.repository.CreateResetToken(email)
}

func (s *service) UpdatePassword(translatorID string, newHashedPassword string) error {
	return s.repository.UpdatePassword(translatorID, newHashedPassword)
}
