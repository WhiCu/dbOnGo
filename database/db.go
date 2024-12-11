package database

import (
	"context"
	"errors"
	"log"

	"github.com/WhiCu/mdb/config"
	"github.com/WhiCu/mdb/database/types"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrLoginTaken     = errors.New("database: a user with this login already exists")
	ErrTokenNotUnique = errors.New("database: the token must be unique")
	ErrTokenEmpty     = errors.New("database: the token cannot be empty")
)

type DB struct {
	Redis       *redis.Client
	DB          *mongo.Database
	Collections map[string]*mongo.Collection
}

func Client(uri string) *mongo.Client {

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		//TODO: add logging
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client

}

func NewDB(redis *redis.Client, client *mongo.Client, dbName string) *DB {
	return &DB{
		Redis:       redis,
		DB:          client.Database(dbName),
		Collections: make(map[string]*mongo.Collection),
	}
}

func (db *DB) Collection(name string) *mongo.Collection {
	//TODO: if value, ok := db.Collections[name]; !ok -?
	if db.Collections[name] == nil {
		//TODO: add logging
		log.Println("New collection")
		db.Collections[name] = db.DB.Collection(name)
	}
	return db.Collections[name]
}

// TODO: add context
func (db *DB) InsertOneUser(ctx context.Context, collectionName string, user *types.User) (string, error) {
	//TODO: add redis

	res, err := db.Collection(config.MustGet("MONGODB_USERS_COLLECTION")).InsertOne(ctx, *user)

	return res.InsertedID.(bson.ObjectID).Hex(), err
}

// TODO: add context
func (db *DB) FindOneUser(ctx context.Context, filter bson.D) (*types.User, bool) {

	log.Println("FindOneUser")

	user := new(types.User)

	err := db.Collection(config.MustGet("MONGODB_USERS_COLLECTION")).FindOne(ctx, filter).Decode(user)

	if err == mongo.ErrNoDocuments {
		return nil, false
	}

	if err != nil {
		panic("FindOne: " + err.Error())
	}

	log.Println("/FindOneUser")
	return user, true
}

func (db *DB) FindOneUserLogin(ctx context.Context, login string) (*types.User, bool) {

	log.Println("FindOneUserLogin")

	if u, _ := db.FindUserInRedis(context.TODO(), login); u != nil {
		return u, true
	}

	filter := bson.D{
		{Key: "login", Value: login},
	}

	user, have := db.FindOneUser(ctx, filter)

	if have {
		db.AddUserInRedis(ctx, user)
	}

	log.Println("/FindOneUserLogin")
	return user, have
}

func (db *DB) AddUser(ctx context.Context, user *types.User) (string, error) {

	log.Println("AddUser")

	if err := db.CorrectUser(user); err != nil {
		//TODO: add logging
		return "", err
	}

	//TODO: add redis

	_, err := db.InsertOneUser(ctx, config.MustGet("MONGODB_USERS_COLLECTION"), user)

	if err != nil {
		panic("AddUser: " + err.Error())
	}

	db.AddUserInRedis(ctx, user)

	return user.TemporaryToken, nil
}

func (db *DB) CorrectUser(user *types.User) error {
	log.Println("CorrectUser")

	if _, have := db.FindOneUserLogin(context.TODO(), user.Login); have {
		return ErrLoginTaken
	}

	filter := bson.D{
		{Key: "token", Value: user.Token},
	}

	_, have := db.FindOneUser(context.TODO(), filter)

	for have {
		user.GenerateTokens()

		filter = bson.D{
			{Key: "token", Value: user.Token},
		}
		_, have = db.FindOneUser(context.TODO(), filter)
	}

	return nil

}
