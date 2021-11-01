package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionDB int

const (
	Datasets collectionDB = iota
	Occitan
	OnGoingTranslations
	SecretQuestions
	TemporaryTokens
	Translations
	TranslationsFiles
	Translators
)

func (c collectionDB) ColletionName() string {
	switch c {
	case Datasets:
		return "Datasets"
	case Occitan:
		return "Occitan"
	case OnGoingTranslations:
		return "OnGoingTranslations"
	case SecretQuestions:
		return "SecretQuestions"
	case TemporaryTokens:
		return "TemporaryTokens"
	case Translations:
		return "Translations"
	case TranslationsFiles:
		return "TranslationsFiles"
	case Translators:
		return "Translators"
	}
	return "unknown"
}

var (
	ttlOnGoingTranslation int32  = 64800  // 18 hours
	ttlTemporaryTokens    int32  = 259200 // 3 days
	uniqueness            bool   = true
	noneLanguage          string = "none"
)

var mongo_COLLECTIONS = [...]mongoCollection{
	{
		Name:    Datasets.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name:    Occitan.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name: OnGoingTranslations.ColletionName(),
		Indexes: []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "createdAt", Value: 1}},
				Options: &options.IndexOptions{ExpireAfterSeconds: &ttlOnGoingTranslation},
			},
		},
	},
	{
		Name:    SecretQuestions.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name: TemporaryTokens.ColletionName(),
		Indexes: []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "issuedAt", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, ExpireAfterSeconds: &ttlTemporaryTokens},
			},
		},
	},
	{
		Name:    Translations.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name:    TranslationsFiles.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name: Translators.ColletionName(),
		Indexes: []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, DefaultLanguage: &noneLanguage},
			},
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, DefaultLanguage: &noneLanguage},
			},
			{
				Keys:    bson.D{{Key: "refreshToken", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, DefaultLanguage: &noneLanguage},
			},
		},
	},
}
