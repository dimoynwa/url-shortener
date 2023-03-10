package mongodb

import (
	"context"
	"time"

	errs "github.com/pkg/errors"

	"github.com/dimoynwa/url-shortener/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoUrl string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout))
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoRepository(mongoUrl, mongoDb string, mongoTimeout int) (shortener.RedirectRepository, error) {
	client, err := newMongoClient(mongoUrl, mongoTimeout)
	if err != nil {
		return nil, errs.Wrap(err, "mongo.Repository.New")
	}

	repo := &mongoRepository{
		client:   client,
		database: mongoDb,
		timeout:  time.Duration(mongoTimeout) * time.Second,
	}

	return repo, nil
}

func (repo *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	redirect := &shortener.Redirect{}
	collection := repo.client.Database(repo.database).Collection("redirects")
	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.Wrap(shortener.ErrRedirectNotFount, "mongoRepository.Find")
		}
		return nil, errs.Wrap(err, "mongoRepository.Find")
	}
	return redirect, nil
}

func (repo *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	collection := repo.client.Database(repo.database).Collection("redirects")
	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	})
	if err != nil {
		return errs.Wrap(err, "mongoRepository.Store")
	}
	return nil
}
