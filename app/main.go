package main

import (
	"log"
	"net/http"

	"github.com/SemyonL95/filesaver/internals/cache"
	"github.com/SemyonL95/filesaver/internals/filestorage"
	"github.com/SemyonL95/filesaver/internals/handlers"
	"github.com/go-redis/redis/v8"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	cache := cache.NewRedisCache(redisClient)
	storage := filestorage.NewLocalStorage("/storage")

	api := handlers.NewAPI(storage, cache)
	http.HandleFunc("/upload", api.Upload)
	http.HandleFunc("/files", api.GetFile)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Cannot start http server")
	}
}
