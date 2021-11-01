package dataset

import (
	"btradoc/entities"
	"btradoc/pkg"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FetchByDialect(fullDialect string) ([]entities.Dataset, error)
	AddTranslatedIn(datasetID, fullDialect string) error
}

type repository struct {
	MongoDB *mongo.Database
}

func NewRepo(mongoDB *mongo.Database) Repository {
	return &repository{
		MongoDB: mongoDB,
	}
}

func (r *repository) FetchByDialect(fullDialect string) ([]entities.Dataset, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "translatedIn", Value: bson.D{{Key: "$nin", Value: []interface{}{fullDialect}}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "OnGoingTranslations"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "dataset"}, {Key: "as", Value: "join"}}}}
	match2Stage := bson.D{{Key: "$match", Value: bson.D{{Key: "join", Value: bson.D{{Key: "$size", Value: 0}}}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}, {Key: "sentence", Value: 1}}}}
	sampleStage := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 5}}}}

	datasetsColl := r.MongoDB.Collection("Datasets")
	ctx := context.Background()
	cursor, err := datasetsColl.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, match2Stage, projectStage, sampleStage})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var datasets []entities.Dataset
	if err = cursor.All(ctx, &datasets); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return datasets, nil
}

func (r *repository) AddTranslatedIn(datasetID, fullDialect string) error {
	datasetsColl := r.MongoDB.Collection("Datasets")

	datasetObjectID, err := primitive.ObjectIDFromHex(datasetID)
	if err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	if _, err := datasetsColl.UpdateOne(context.Background(), bson.M{"_id": datasetObjectID}, bson.D{{Key: "$push", Value: bson.D{{Key: "translatedIn", Value: fullDialect}}}}); err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}
