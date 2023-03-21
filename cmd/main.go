package main

import (
	"context"
	"encoding/json"
	"github.com/erkkke/technodom_test/api"
	cache2 "github.com/erkkke/technodom_test/db/cache"
	db "github.com/erkkke/technodom_test/db/sqlc"
	"github.com/erkkke/technodom_test/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	redirectsList, err := parseRedirectsDataToStruct("./links.json")
	if err != nil {
		log.Fatalf("cannot parse redirects data: %s", err)
	}

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sqlx.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	database := db.New(conn)
	cache := cache2.NewInMemoryCache(1000)
	server, err := api.NewServer(database, cache)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	for _, redirect := range redirectsList {
		_, err = database.CreateRedirect(context.Background(), redirect)
		if err != nil {
			log.Fatal("cannot add to database:", err)
		}
	}

	server.Start(config.ServerAddress)
}

func parseRedirectsDataToStruct(filename string) ([]db.CreateRedirectParams, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var redirects []db.CreateRedirectParams
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&redirects)
	if err != nil {
		return nil, err
	}

	return redirects, nil
}
