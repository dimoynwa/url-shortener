package redis

import (
	"fmt"
	"strconv"

	"github.com/dimoynwa/url-shortener/shortener"
	"github.com/go-redis/redis"
	errs "github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisUrl string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, errs.Wrap(err, "redisRepository.NewClient")
	}

	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	return client, err
}

func NewRedisRepository(redisUrl string) (shortener.RedirectRepository, error) {
	client, err := newRedisClient(redisUrl)
	if err != nil {
		return nil, errs.Wrap(err, "redisRepository.New")
	}
	return &redisRepository{
		client: client,
	}, nil
}

func (repo *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%v", code)
}

func (repo *redisRepository) Find(code string) (*shortener.Redirect, error) {
	key := repo.generateKey(code)
	redirect := &shortener.Redirect{}

	data, err := repo.client.HGetAll(key).Result()
	if err != nil {
		return nil, errs.Wrap(err, "redisRepository.Find")
	}
	if len(data) == 0 {
		return nil, errs.Wrap(shortener.ErrRedirectNotFount, "redisRepository.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errs.Wrap(err, "redisRepository.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil
}

func (repo *redisRepository) Store(redirect *shortener.Redirect) error {
	key := repo.generateKey(redirect.Code)

	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := repo.client.HMSet(key, data).Result()
	if err != nil {
		return errs.Wrap(err, "redisRepository.Store")
	}
	return nil
}
