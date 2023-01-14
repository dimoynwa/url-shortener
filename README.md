Golang example project to shorten urls.
You can store url in database (MongoDb or Redis). It will generate unique code and store the URL.
Then you can call ```\{code}``` and it will redirect you to the stored URL.

Environment variables needed:

- PORT - integer, if not presented 8000 will be used as default
- URL_DB - (required) redis or mongodb
Redis:
    - REDIS_URL - url to redis database
Mongo DB:
    - MONGO_URL - url to mongo database
    - MONGO_DB - database name
    - MONGO_TIMEOUT - connection timeout in seconds

To start the project:
```go run .\main.go```