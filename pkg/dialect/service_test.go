package dialect

import (
	"btradoc/storage/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	db := mongodb.NewMongoClient()
	dialectRepo := NewRepo(db)
	dialectService := NewService(dialectRepo)

	dialect := "auvernhat"
	subdialect := "brivadés"
	match, err := dialectService.Exists(dialect, subdialect)
	assert.Nil(t, err)
	t.Log(match)

	dialect = "ahuevratn"
	subdialect = "brivadés"
	match, err = dialectService.Exists(dialect, subdialect)
	assert.Nil(t, err)
	t.Log(match)
}

func TestFetchOccitan(t *testing.T) {
	db := mongodb.NewMongoClient()
	dialectRepo := NewRepo(db)
	dialectService := NewService(dialectRepo)

	translatorID := "6172fc761d23a6beb3a9666a"
	result, err := dialectService.FetchOccitan(translatorID)
	assert.Nil(t, err)
	t.Log(result)
}
