package mongodb

import (
	"btradoc/data"
	"testing"
)

func TestConnection(t *testing.T) {
	mongodb := NewMongoClient()
	t.Log(mongodb.Name())
}

func TestCreateCollections(t *testing.T) {
	mongodb := NewMongoClient()
	createCollections(mongodb, mongo_COLLECTIONS[:])
}

func TestAddOccitanDialects(t *testing.T) {
	mongodb := NewMongoClient()
	addOccitanDialects(mongodb, data.OCCITAN)
}

func TestAddSecretQuestions(t *testing.T) {
	mongodb := NewMongoClient()
	addSecretQuestions(mongodb, data.SECRET_QUESTIONS)
}
