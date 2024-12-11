package main

import (
	"context"
	"log"
	"time"

	"github.com/WhiCu/mdb/config"
	"github.com/WhiCu/mdb/database"
	"github.com/WhiCu/mdb/database/types"
)

func main() {

	user := types.New(
		"email",
		"phone",
		"username",
		"NewLogin",
		"password",
	)

	uri := "mongodb://user:password@localhost:27017/fiberServer"

	client := database.Client(uri)

	db := database.NewDB(database.NewRedisLayer(), client, config.MustGet("MONGO_DB"))

	cont, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	token, err := db.AddUser(cont, user)
	cancel()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully added token: ", token)
}
