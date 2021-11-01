package translation

import "btradoc/entities"

type Service interface {
	AddTranslations(translatorID string, translations []entities.Translation) error
	FetchTotalOnGoingTranslations(fullDialect, translatorID string) (int, error)
	AddOnGoingTranslations(fullDialect, translatorID string, datasets []entities.Dataset) error
	RemoveOnGoingTranslations(translations []entities.Translation) error
	FetchPathnameFiles() ([]entities.TranslationFile, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddTranslations(translatorID string, translations []entities.Translation) error {
	return s.repository.InsertTranslations(translatorID, translations)
}

func (s *service) FetchTotalOnGoingTranslations(fullDialect, translatorID string) (int, error) {
	return s.repository.GetTotalOnGoingTranslation(fullDialect, translatorID)
}

func (s *service) AddOnGoingTranslations(fullDialect, translatorID string, datasets []entities.Dataset) error {
	return s.repository.InsertDatasetsOnGoingTranslations(fullDialect, translatorID, datasets)
}

func (s *service) RemoveOnGoingTranslations(translations []entities.Translation) error {
	return s.repository.RemoveDatasetsOnGoingTranslations(translations)
}

func (s *service) FetchPathnameFiles() ([]entities.TranslationFile, error) {
	return s.repository.GetTranslationsFiles()
}
