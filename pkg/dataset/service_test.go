package dataset

import (
	"btradoc/storage/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	db := mongodb.NewMongoClient()

	datasetRepo := NewRepo(db)
	datasetService := NewService(datasetRepo)

	datasetID := ""
	fullDialect := "test_test"

	_, _ = datasetService.FetchByDialect(fullDialect)
	_ = datasetService.AddTranslatedIn(datasetID, fullDialect)
}

func TestAddTransledIn(t *testing.T) {
	db := mongodb.NewMongoClient()
	datasetRepo := NewRepo(db)
	datasetService := NewService(datasetRepo)

	datasetID := "61422a33c632832bb589d914"
	fullDialect := "gascon_aranés"

	err := datasetService.AddTranslatedIn(datasetID, fullDialect)
	assert.Nil(t, err)
}
