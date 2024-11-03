package repository

import (
	"context"

	"github.com/griggsjared/getsit/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// urlEntryDatabase is the name of the database that will repository the url entries
	urlEntryDatabase = "getsit"
	// urlEntryCollection is the name of the collection that will repository the url entries
	urlEntryCollection = "url_entries"
)

// MongoDBUrlEntryRepository is a mongodb repository that will repository the url entries
type MongoDBUrlEntryRepository struct {
	client *mongo.Client
}

// NewMongoDBUrlEntryRepository will create a new mongodb repository from the the mongodb client
func NewMongoDBUrlEntryRepository(c *mongo.Client) *MongoDBUrlEntryRepository {
	return &MongoDBUrlEntryRepository{
		client: c,
	}
}

// mongoDBUrlEntrySchema is the schema for the url entry in the mongodb repository
type mongoDBUrlEntrySchema struct {
	Token      entity.UrlToken `bson:"_id"`
	Url        entity.Url      `bson:"url"`
	VisitCount int             `bson:"visit_count"`
}

// Save will save the url entry to the repository
func (s *MongoDBUrlEntryRepository) SaveUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return nil, err
	}

	entry, err := s.GetFromUrl(ctx, url)
	if err == nil {
		return entry, nil
	}

	newEntry := mongoDBUrlEntrySchema{
		Token:      entity.NewUrlToken(),
		Url:        url,
		VisitCount: 0,
	}

	_, err = coll.InsertOne(ctx, newEntry)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Url:        newEntry.Url,
		Token:      newEntry.Token,
		VisitCount: newEntry.VisitCount,
	}, nil
}

// SaveVisit will increment the number of times the url has been visited
func (s *MongoDBUrlEntryRepository) SaveVisit(ctx context.Context, token entity.UrlToken) error {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return err
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": token.String()}, bson.M{"$inc": bson.M{"visit_count": 1}})
	if err != nil {
		return err
	}

	return nil
}

// GetFromToken will get the url entry from the token
func (s *MongoDBUrlEntryRepository) GetFromToken(ctx context.Context, token entity.UrlToken) (*entity.UrlEntry, error) {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return nil, err
	}

	var entry mongoDBUrlEntrySchema
	err = coll.FindOne(ctx, bson.M{"_id": token.String()}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Url:        entry.Url,
		Token:      entry.Token,
		VisitCount: entry.VisitCount,
	}, nil
}

// GetFromUrl will get the url entry from the url
func (s *MongoDBUrlEntryRepository) GetFromUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

	coll, err := s.newUrlEntryCollection()
	if err != nil {
		return nil, err
	}

	var schema mongoDBUrlEntrySchema
	err = coll.FindOne(ctx, bson.M{"url": url.String()}).Decode(&schema)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Url:        schema.Url,
		Token:      schema.Token,
		VisitCount: schema.VisitCount,
	}, nil
}

// newUrlEntryCollection will create a new collection for the url entries if it does not exist yet
// it will also create an index for the url field
func (s *MongoDBUrlEntryRepository) newUrlEntryCollection() (*mongo.Collection, error) {

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
