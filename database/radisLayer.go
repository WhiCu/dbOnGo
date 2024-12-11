package database

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/WhiCu/mdb/config"
	"github.com/WhiCu/mdb/database/types"
	"github.com/redis/go-redis/v9"
)

func NewRedisLayer() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.DefaultGet("REDIS_HOST", "localhost"), config.MustGet("REDIS_PORT")),
		Password: config.MustGet("REDIS_PASSWORD"),
		DB:       config.MustGetInt("REDIS_DB_ID"),
	})
}

func (db *DB) AddUserInRedis(ctx context.Context, user *types.User) error {
	log.Println("AddUserInRedis")

	err := db.Redis.Set(ctx, "user:"+user.Login, user.Json(), 1*time.Hour).Err()

	log.Println("/AddUserInRedis")
	return err
}

func (db *DB) FindUserInRedis(ctx context.Context, login string) (*types.User, error) {
	log.Println("FindUserInRedis")
	// Создаем новую переменную для пользователя
	// Пытаемся получить пользователя из Redis

	res, err := db.Redis.Get(ctx, "user:"+login).Result()

	if err != nil {
		// Если пользователь не найден, можно вернуть nil и обработать это в вызывающем коде
		if err == redis.Nil {
			log.Println("User not found in Redis")
			return nil, nil // Возвращаем nil, если пользователь не найден
		}
		return nil, err // Возвращаем ошибку в случае других проблем
	}

	fetchedUser := new(types.User)
	// Если пользователь был найден, декодируем его JSON в структуру
	if err := json.Unmarshal([]byte(res), fetchedUser); err != nil {
		return nil, err // Возвращаем ошибку в случае проблем с декодированием
	}

	log.Println("/FindUserInRedis")
	return fetchedUser, nil // Возвращаем найденного пользователя
}
