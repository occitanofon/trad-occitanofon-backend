package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsernameValidity(t *testing.T) {
	usernames := []string{
		"Loís",
		"Loìs",
		"William",
		"koli",
		"Gregòri",
		"Timotèu",
		"Emanuèl",
	}
	for _, username := range usernames {
		err := UsernameValidity(username)
		assert.NotNil(t, err)
	}
}

func TestIsEmailValid(t *testing.T) {
	type MockEmail struct {
		Email    string
		Expected bool
	}

	mockEmails := []MockEmail{
		{
			Email:    "gille-dubois@orange.fr",
			Expected: true,
		},
		{
			Email:    "poirier_chloe@free.fr",
			Expected: true,
		},
		{
			Email:    "bernard-dupont@test.vf",
			Expected: false,
		},
	}

	for _, mockEmail := range mockEmails {
		ok := IsEmailValid(mockEmail.Email)
		assert.Equal(t, mockEmail.Expected, ok, mockEmail.Email)
	}
}

func TestNormalize(t *testing.T) {
	type mockWord struct {
		Word     string
		Expected string
	}

	mockWords := []mockWord{
		{
			Word:     "niçard",
			Expected: "nicard",
		},
		{
			Word:     "aranés",
			Expected: "aranes",
		},
		{
			Word:     "gavòt",
			Expected: "gavot",
		},
	}

	for _, mk := range mockWords {
		wordResult, err := Normalize(mk.Word)
		assert.Nil(t, err)
		assert.Equal(t, mk.Expected, wordResult)
		t.Log(mk.Word, wordResult, fmt.Sprintf("%s_%s", "auv", wordResult[:3]))
	}
}

func TestAccent(t *testing.T) {
	str := "niçard"
	s := []rune(str)
	fmt.Println(len(s), s[2], string(s[2]), string(s[:3]))
	fmt.Println(len(str), str[2], string(str[2]), str[:3])
	fmt.Println()
}
