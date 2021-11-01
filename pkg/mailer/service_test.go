package mailer

import (
	"btradoc/entities"
	"btradoc/storage/mongodb"
	"testing"
)

func TestMailer(t *testing.T) {
	db := mongodb.NewMongoClient()
	mailerService := NewService(db, nil)

	translTest := &entities.TranslatorResetPassword{
		Email:    "test@test.com",
		Username: "test",
		Token:    "4a5z11vf24e5",
	}

	for i := 10_000; i > 0; i-- {
		mailerService.SendResetPasswordLink(translTest)
	}
}
