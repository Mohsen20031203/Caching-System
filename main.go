package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

func redisChack(addres string) *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr: addres,
	})
	return rds
}

var rdb *redis.Client
var ctx = context.Background()
var numberRequest int64 = 5

const windowSize = time.Minute

func apii(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")

	if !checkRateLimit(userID) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	fmt.Fprintf(w, "Request successful")

}

func checkRateLimit(userID string) bool {
	key := fmt.Sprintf("user:%s:requests", userID)

	pipe := rdb.Pipeline()

	pipe.Incr(ctx, key)

	pipe.Expire(ctx, key, windowSize)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error executing Redis pipeline: %v", err)
		return false
	}
	requestCount := cmds[0].(*redis.IntCmd).Val()
	if requestCount < numberRequest {
		return false
	}
	return true

}

func main() {
	rdb = redisChack("localhost:6379")
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("not connect %v", err)
	} else {
		fmt.Println("connect in server")
	}

	http.HandleFunc("/api/", apii)
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		panic(err)
	}

}
