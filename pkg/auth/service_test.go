package auth

import (
	"btradoc/storage/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	db := mongodb.NewMongoClient()
	authRepo := NewRepo(db)
	authService := NewService(authRepo)

	username := "Dàrius"
	translator, err := authService.Login(username)

	assert.Nil(t, err)
	t.Logf("%+v\n", translator)
}

func TestCreateRefreshToken(t *testing.T) {
	db := mongodb.NewMongoClient()
	authRepo := NewRepo(db)
	authService := NewService(authRepo)

	translatorID := "6148c3f1ba78b40cdeb49289"
	refreshToken, err := authService.CreateRefreshToken(translatorID)
	assert.Nil(t, err)
	t.Logf("%s\n", refreshToken)
}

func TestFindByRefreshToken(t *testing.T) {
	db := mongodb.NewMongoClient()
	authRepo := NewRepo(db)
	authService := NewService(authRepo)

	refreshToken := "JfQt53reoJR"

	transl, err := authService.FindByRefreshToken(refreshToken)
	assert.Nil(t, err)
	t.Logf("%+v\n", transl)
}
