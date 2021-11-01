package mongodb

import "testing"

func TestNewMongoClient(t *testing.T) {
	db := NewMongoClient()
	t.Log(db.Name())
}
