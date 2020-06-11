package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SemyonL95/filesaver/internals/cache"
	"github.com/SemyonL95/filesaver/internals/config"
	"github.com/SemyonL95/filesaver/internals/filestorage"
	"github.com/SemyonL95/filesaver/internals/handlers"
	"github.com/go-redis/redis/v8"
)

func main() {
	conf := config.NewConfig()

	log.Printf("app conf: %v", conf)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		Password: "",
		DB:       0,
	})

	cache := cache.NewRedisCache(redisClient)
	storage := filestorage.NewLocalStorage(conf.StoragePath)

	api := handlers.NewAPI(storage, cache)
	http.HandleFunc("/upload", api.Upload)
	http.HandleFunc("/files", api.GetFile)

	err := http.ListenAndServe(fmt.Sprintf(":%s", conf.AppPort), nil)
	if err != nil {
		log.Fatal("Cannot start http server")
	}
}
