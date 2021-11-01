package auth

import (
	"btradoc/entities"
	"btradoc/helpers"
	"btradoc/pkg"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	FindByUsername(translatorUsername string) (*entities.Translator, error)
	InsertRefreshToken(translatorID string) (string, error)
	GetTranslatorByRefreshToken(refreshToken string) (*entities.Translator, error)
}

type repository struct {
	MongoDB *mongo.Database
}

func NewRepo(mongoDB *mongo.Database) Repository {
	return &repository{
		MongoDB: mongoDB,
	}
}

func (r *repository) FindByUsername(translatorUsername string) (*entities.Translator, error) {
	translatorsColl := r.MongoDB.Collection("Translators")

	opts := options.FindOne().SetSort(bson.D{{Key: "username", Value: 1}})
	var result bson.M
	err := translatorsColl.FindOne(context.Background(), bson.D{{Key: "username", Value: translatorUsername}}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &pkg.DBError{
				Code:    404,
				Message: pkg.ErrBadCredentials,
				Wrapped: err,
			}
		}
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ID := result["_id"].(primitive.ObjectID)

	translator := entities.Translator{
		ID:        ID.Hex(),
		Email:     result["email"].(string),
		Username:  result["username"].(string),
		Hpwd:      result["hpwd"].(string),
		Confirmed: result["confirmed"].(bool),
		Suspended: result["suspended"].(bool),
	}

	permissions := result["permissions"].(primitive.A)
	for _, permission := range permissions {
		perm := permission.(string)
		translator.Permissions = append(translator.Permissions, perm)
	}

	return &translator, nil
}

func (r *repository) InsertRefreshToken(translatorID string) (string, error) {
	// retry if refreshToken generated is already set
	for i := 0; i < 50; i++ {
		refreshToken := helpers.GenerateID(12)

		translatorObjectID, err := primitive.ObjectIDFromHex(translatorID)
		if err != nil {
			return "", &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}

		translatorsColl := r.MongoDB.Collection("Translators")
		filter := bson.M{"_id": translatorObjectID}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "refreshToken", Value: refreshToken}}}}
		if _, err = translatorsColl.UpdateOne(context.Background(), filter, update); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				continue
			}

			return "", &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}

		return refreshToken, nil
	}

	return "", &pkg.DBError{
		Code:    500,
		Message: pkg.ErrDefault,
		Wrapped: errors.New("cannot update refresh token"),
	}
}

func (r *repository) GetTranslatorByRefreshToken(refreshToken string) (*entities.Translator, error) {
	translatorsColl := r.MongoDB.Collection("Translators")

	opts := options.FindOne().SetSort(bson.D{{Key: "refreshToken", Value: 1}})
	var result bson.M
	err := translatorsColl.FindOne(context.Background(), bson.D{{Key: "refreshToken", Value: refreshToken}}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &pkg.DBError{
				Code:    404,
				Message: pkg.ErrRefreshTokenNotFound,
				Wrapped: err,
			}
		}
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ID := result["_id"].(primitive.ObjectID)

	translator := entities.Translator{
		ID: ID.Hex(),
	}

	permissions := result["permissions"].(primitive.A)
	for _, permission := range permissions {
		perm := permission.(string)
		translator.Permissions = append(translator.Permissions, perm)
	}

	return &translator, nil
}
