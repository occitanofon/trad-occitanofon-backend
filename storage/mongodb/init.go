package mongodb

import (
	"btradoc/data"
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCollection struct {
	Name    string
	Indexes []mongo.IndexModel
}

// InitMongoDatabase initializes datas
func InitMongoDatabase() {
	mongodb := NewMongoClient()
	createCollections(mongodb, mongo_COLLECTIONS[:])
	addOccitanDialects(mongodb, data.OCCITAN)
	addSecretQuestions(mongodb, data.SECRET_QUESTIONS)
}

// createCollections creates all mongo collections
func createCollections(db *mongo.Database, collections []mongoCollection) {
	ctx := context.Background()

	for _, collection := range collections {
		err := db.CreateCollection(ctx, collection.Name)
		if err != nil {
			if strings.Contains(err.Error(), "(NamespaceExists)") {
				continue
			}
			log.Fatalln(err)
		}

		coll := db.Collection(collection.Name, nil)
		coll.Indexes().CreateMany(ctx, collection.Indexes)
	}
}

// addOccitanDialects adds all occitan dialects to a specified Collection
func addOccitanDialects(db *mongo.Database, occitan []data.Occitan) {
	occitanColl := db.Collection("Occitan")
	count, err := occitanColl.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return
	}

	var totalDialect int
	for _, dialect := range occitan {
		totalDialect += len(dialect.Subdialects)
	}

	log.Printf("Total Occitan Dialects: %d (JSON) | %d (Documents)\n", totalDialect, count)

	ctx := context.Background()
	for _, dialect := range occitan {
		for _, subdialect := range dialect.Subdialects {
			result := occitanColl.FindOne(ctx, bson.D{{Key: "dialectName", Value: dialect.Dialect}, {Key: "subdialectName", Value: subdialect}})
			if result.Err() != nil {
				if result.Err() == mongo.ErrNoDocuments {
					// insert occitan dialect if it's not in documents
					if _, err := occitanColl.InsertOne(ctx, bson.D{
						{Key: "dialectName", Value: dialect.Dialect}, {Key: "subdialectName", Value: subdialect},
					}); err != nil {
						log.Println(err)
					}
				}

				continue
			}
		}
	}
}

// addSecretQuestions adds all secret questions to a specified Collection
func addSecretQuestions(db *mongo.Database, secretQuestions []string) {
	secretQuestionsColl := db.Collection("SecretQuestions")
	count, err := secretQuestionsColl.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Total secret questions: %d (JSON) | %d (Documents)\n", len(secretQuestions), count)

	ctx := context.Background()
	for _, secretQuestion := range secretQuestions {
		result := secretQuestionsColl.FindOne(ctx, bson.D{{Key: "question", Value: secretQuestion}})
		if result.Err() != nil {
			if result.Err() == mongo.ErrNoDocuments {
				// insert secret question if it's not in documents
				if _, err := secretQuestionsColl.InsertOne(ctx, bson.D{{Key: "question", Value: secretQuestion}}); err != nil {
					log.Println(err)
				}
			}

			continue
		}
	}
}
