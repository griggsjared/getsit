package repository

import (
	"context"

	"github.com/griggsjared/getsit/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// urlEntryDatabase is the name of the database that will store the url entries
	urlEntryDatabase = "getsit"
	// urlEntryCollection is the name of the collection that will store the url entries
	urlEntryCollection = "url_entries"
)

type MongoDBUrlEntryStore struct {
	client *mongo.Client
}

func NewMongoDBUrlEntryStore(c *mongo.Client) *MongoDBUrlEntryStore {
	return &MongoDBUrlEntryStore{
		client: c,
	}
}

type urlEntrySchema struct {
	Token      entity.UrlToken `bson:"_id"`
	Url        entity.Url      `bson:"url"`
	VisitCount int             `bson:"visit_count"`
}

func (s *MongoDBUrlEntryStore) Save(ctx context.Context, url entity.Url) (urlEntry *entity.UrlEntry, new bool, err error) {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return nil, false, err
	}

	entry, err := s.GetFromUrl(ctx, string(url))
	if err == nil {
		return entry, false, nil
	}

	newEntry := urlEntrySchema{
		Token:      entity.NewUrlToken(),
		Url:        url,
		VisitCount: 0,
	}

	_, err = coll.InsertOne(ctx, newEntry)
	if err != nil {
		return nil, false, err
	}

	return &entity.UrlEntry{
		Url:        newEntry.Url,
		Token:      newEntry.Token,
		VisitCount: newEntry.VisitCount,
	}, true, nil
}

func (s *MongoDBUrlEntryStore) SaveVisit(ctx context.Context, token string) error {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return err
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": token}, bson.M{"$inc": bson.M{"visit_count": 1}})
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoDBUrlEntryStore) GetFromToken(ctx context.Context, token string) (*entity.UrlEntry, error) {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return nil, err
	}

	var entry urlEntrySchema
	err = coll.FindOne(ctx, bson.M{"_id": token}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Url:        entry.Url,
		Token:      entry.Token,
		VisitCount: entry.VisitCount,
	}, nil
}

func (s *MongoDBUrlEntryStore) GetFromUrl(ctx context.Context, url string) (*entity.UrlEntry, error) {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return nil, err
	}

	var schema urlEntrySchema
	err = coll.FindOne(ctx, bson.M{"url": url}).Decode(&schema)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Url:        schema.Url,
		Token:      schema.Token,
		VisitCount: schema.VisitCount,
	}, nil
}

func (s *MongoDBUrlEntryStore) newUrlEntryCollection() (*mongo.Collection, error) {

	coll := s.client.Database(urlEntryDatabase).Collection(urlEntryCollection)

	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"url": 1,
		},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return nil, err
	}

	return coll, nil
}
