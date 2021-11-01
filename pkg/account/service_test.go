package account

import (
	"btradoc/entities"
	"btradoc/storage/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db := mongodb.NewMongoClient()

	accountRepo := NewRepo(db)
	accountService := NewService(accountRepo)

	newTranslator := entities.NewTranslator{
		Email:    "lois@sfr.fr",
		Username: "loisssvfvfvfvf",
		Hpwd:     "$argon2id$v=19$m=16,t=2,p=1$RXhFUTNTWGFueUI5Umx6Qg$wtoGpTHQ43Vt5kFiR342Xg",
		SecretQuestionsAndResponses: []entities.SecretQuestionAndResponse{
			{
				Question: "Qué siguèt ton faus-nom quand ères un enfant ?",
				Response: "$argon2id$v=19$m=16,t=2,p=1$NkFLZFFPeEllS2VLeGpjWg$HlmRZGzfHlpD8ItbTdWABA",
			},
			{
				Question: "Qué siguèt lo premier filme que veguères au cinemà ?",
				Response: "$argon2id$v=19$m=16,t=2,p=1$dW54UThKWjJQTU9WM05Tdw$7KTR8EuJOF9QZ+jqtMDAPw",
			},
		},
	}

	err := accountService.Create(newTranslator)
	assert.Nil(t, err)
}

func TestFetchAllSecretQuestions(t *testing.T) {
	db := mongodb.NewMongoClient()
	accountRepo := NewRepo(db)
	accountService := NewService(accountRepo)

	allSecretQuestions, err := accountService.FetchAllSecretQuestions()
	assert.Nil(t, err)
	t.Logf("%s", allSecretQuestions)
}

func TestFetchSecretQuestionsAndResponses(t *testing.T) {
	db := mongodb.NewMongoClient()
	accountRepo := NewRepo(db)
	accountService := NewService(accountRepo)

	token := "v4hKKGd4ttyXJyuiXbc0"

	sq, err := accountService.FetchSecretQuestionsAndResponses(token)
	assert.Nil(t, err)

	if err == nil {
		t.Logf("%s\n", sq)
	}
}

func TestResetPassword(t *testing.T) {
	db := mongodb.NewMongoClient()
	accountRepo := NewRepo(db)
	accountService := NewService(accountRepo)

	email := "dupont-lois@gmail.com"
	transl, err := accountService.ResetPassword(email)
	assert.Nil(t, err)

	t.Log(transl)
}
