package entities

type Dataset struct {
	ID       string `json:"id" bson:"_id"`
	Sentence string `json:"sentence" bson:"sentence"`
}
