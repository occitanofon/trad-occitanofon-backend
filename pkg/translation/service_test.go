package translation

import (
	"btradoc/entities"
	"btradoc/storage/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTranslations(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	translatorID := ""
	translations := []entities.Translation{
		{
			Oc:        "Es de Canadà mai naissèt en Anglatèrra",
			Fr:        "Je suis au Canada mais je suis né en Angleterre",
			En:        "I'm in Canada but I was born in England",
			DatasetID: "9969",
			Occitan:   "auvernhat_estandard",
		},
	}

	err := translationService.AddTranslations(translatorID, translations)
	assert.Nil(t, err)
}

func TestFetchTotalOnGoingTranslations(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	translatorID := "6148c3f1ba78b40cdeb49288"
	fullDialect := "auvernhat_estandard"

	counter, err := translationService.FetchTotalOnGoingTranslations(fullDialect, translatorID)
	assert.Nil(t, err)
	t.Logf("total: %d", counter)
}

func TestAddOnGoingTranslations(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	fullDialect := "auvernhat_estandard"

	datasets := []entities.Dataset{
		{
			ID:       "",
			Sentence: "",
		},
	}

	translatorID := "6148c3f1ba78b40cdeb49288"

	err := translationService.AddOnGoingTranslations(fullDialect, translatorID, datasets)
	assert.Nil(t, err)
}

func TestFetchPathnameFiles(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	translationsFiles, err := translationService.FetchPathnameFiles()
	assert.Nil(t, err)
	t.Logf("%+v\n", translationsFiles)
}
