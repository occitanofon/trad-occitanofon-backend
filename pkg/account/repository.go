package account

import (
	"btradoc/entities"
	"btradoc/helpers"
	"btradoc/pkg"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	InsertNewTranslator(newTranslator entities.NewTranslator) error
	GetAllSecretQuestions() ([]string, error)
	GetTranslatorSecretQuestions(token string) (*entities.TranslatorSecretQuestions, error)
	CreateResetToken(email string) (*entities.TranslatorResetPassword, error)
	UpdatePassword(translatorID string, newHashedPassword string) error
}

type repository struct {
	MongoDB *mongo.Database
}

func NewRepo(mongoDB *mongo.Database) Repository {
	return &repository{
		MongoDB: mongoDB,
	}
}

func (r *repository) InsertNewTranslator(newTranslator entities.NewTranslator) error {
	translatorsColl := r.MongoDB.Collection("Translators")

	doc := bson.D{
		{Key: "username", Value: newTranslator.Username},
		{Key: "email", Value: newTranslator.Email},
		{Key: "hpwd", Value: newTranslator.Hpwd},
		{Key: "confirmed", Value: false},
		{Key: "suspended", Value: false},
		{Key: "secretQuestions", Value: []bson.D{
			{{Key: "question", Value: newTranslator.SecretQuestionsAndResponses[0].Question}, {Key: "response", Value: newTranslator.SecretQuestionsAndResponses[0].Response}},
			{{Key: "question", Value: newTranslator.SecretQuestionsAndResponses[1].Question}, {Key: "response", Value: newTranslator.SecretQuestionsAndResponses[1].Response}},
		}},
		{Key: "permissions", Value: []string{}},
		{Key: "createdAt", Value: primitive.NewDateTimeFromTime(time.Now())},
		{Key: "refreshToken", Value: ""},
	}

	if _, err := translatorsColl.InsertOne(context.Background(), doc); err != nil {
		if ok := mongo.IsDuplicateKeyError(err); ok {
			if strings.Contains(err.Error(), "username_1") {
				return &pkg.DBError{
					Code:    409,
					Message: pkg.ErrPseudoUsed,
					Wrapped: err,
				}
			} else if strings.Contains(err.Error(), "email_1") {
				return &pkg.DBError{
					Code:    409,
					Message: pkg.ErrEmailUsed,
					Wrapped: err,
				}
			}
		}
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}

func (r *repository) GetAllSecretQuestions() ([]string, error) {
	secretQuestionsColl := r.MongoDB.Collection("SecretQuestions")
	ctx := context.Background()

	cursor, err := secretQuestionsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var secretQuestionsDocs []bson.M
	if err = cursor.All(ctx, &secretQuestionsDocs); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	var secretQuestions []string
	for _, value := range secretQuestionsDocs {
		sq := value["question"].(string)

		secretQuestions = append(secretQuestions, sq)
	}

	return secretQuestions, nil
}

func (r *repository) GetTranslatorSecretQuestions(token string) (*entities.TranslatorSecretQuestions, error) {
	temporaryTokensColl := r.MongoDB.Collection("TemporaryTokens")

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "token", Value: token}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "Translators"}, {Key: "localField", Value: "translator"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "translator"}}}}
	project1Stage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}, {Key: "translator", Value: 1}, {Key: "issuedAt", Value: "$issuedAt"}, {Key: "secretQuestions", Value: "$translator.secretQuestions"}}}}
	project2Stage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}, {Key: "translator", Value: 1}, {Key: "issuedAt", Value: "$issuedAt"}, {Key: "secretQuestions", Value: bson.D{{Key: "$first", Value: "$secretQuestions"}}}}}}

	ctx := context.Background()
	cursor, err := temporaryTokensColl.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, project1Stage, project2Stage})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	} else if len(result) == 0 {
		return nil, &pkg.DBError{
			Code:    404,
			Message: pkg.ErrSecretQuestionsNotFound,
			Wrapped: fmt.Errorf("no secret questions found"),
		}
	}

	translatorObjectID := result[0]["translator"].(primitive.A)[0].(primitive.M)["_id"].(primitive.ObjectID)
	timestamp := result[0]["issuedAt"].(primitive.DateTime)

	if time.Now().After(timestamp.Time().Add(12 * time.Hour)) {
		return nil, &pkg.DBError{
			Code:    403,
			Message: pkg.ErrResetPasswordTokenExpired,
			Wrapped: errors.New("reset token password is expired"),
		}
	}

	mapSecretQuestions := make(map[string]string, len(result[0]["secretQuestions"].(primitive.A)))
	for _, mapElems := range result[0]["secretQuestions"].(primitive.A) {
		maps := mapElems.(primitive.M)
		question := maps["question"].(string)
		response := maps["response"].(string)
		mapSecretQuestions[question] = response
	}

	translatorSecretQuestions := entities.TranslatorSecretQuestions{
		TranslatorID:                translatorObjectID,
		SecretQuestionsAndResponses: mapSecretQuestions,
	}

	return &translatorSecretQuestions, nil
}

func (r *repository) CreateResetToken(email string) (*entities.TranslatorResetPassword, error) {
	translatorColl := r.MongoDB.Collection("Translators")
	temporaryTokensColl := r.MongoDB.Collection("TemporaryTokens")

	opts := options.FindOne().SetSort(bson.D{{Key: "email", Value: 1}})
	var result bson.M
	if err := translatorColl.FindOne(context.Background(), bson.D{{Key: "email", Value: email}, {Key: "confirmed", Value: true}, {Key: "suspended", Value: false}}, opts).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &pkg.DBError{
				Code:    404,
				Message: pkg.ErrTranslatorNotFound,
				Wrapped: err,
			}
		}
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	translatorID := result["_id"].(primitive.ObjectID)
	usernameRes := result["username"].(string)

	translatorObjectID, err := primitive.ObjectIDFromHex(translatorID.Hex())
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ctx := context.Background()

	// remove last stored temporary token before adding a new one
	if _, err := temporaryTokensColl.DeleteOne(ctx, bson.D{{Key: "translator", Value: translatorObjectID}}); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	token := helpers.GenerateID(20)
	doc := bson.D{
		{Key: "type", Value: "PASSWORD_RESET"},
		{Key: "token", Value: token},
		{Key: "translator", Value: translatorObjectID},
		{Key: "issuedAt", Value: primitive.NewDateTimeFromTime(time.Now())},
	}

	if _, err := temporaryTokensColl.InsertOne(ctx, doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &pkg.DBError{
				Code:    404,
				Message: pkg.ErrTranslatorNotFound,
				Wrapped: err,
			}
		}
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	translatorResetPassword := entities.TranslatorResetPassword{
		Email:    email,
		Username: usernameRes,
		Token:    token,
	}

	return &translatorResetPassword, nil
}

func (r *repository) UpdatePassword(translatorID string, newHashedPassword string) error {
	translatorColl := r.MongoDB.Collection("Translators")

	translatorObjectID, err := primitive.ObjectIDFromHex(translatorID)
	if err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	if _, err = translatorColl.UpdateOne(context.Background(), bson.M{"_id": translatorObjectID}, bson.D{{Key: "$set", Value: bson.D{{Key: "hpwd", Value: newHashedPassword}}}}); err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}
