package main

import (
	"context"
	"fmt"
	"io/ioutil"
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
var cacheKey = "weather_mashhad"

func found() (string, error) {
	checkData, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		fmt.Println("data found in chache")
		return checkData, nil
	}

	req, err := http.NewRequest("GET", "https://api.meteomatics.com/2025-03-09T00:00:00Z/t_2m:C/36.3,59.6/json", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("iran_salehi_mohsen", "FtB6c4eSE4")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	weatherData := string(body)

	err = rdb.Set(ctx, cacheKey, weatherData, 1*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return weatherData, nil
}

func Weather(w http.ResponseWriter, c *http.Request) {
	data, err := found()
	if err != nil {
		http.Error(w, "not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func main() {
	rdb = redisChack("localhost:6379")
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("اتصال به Redis موفقیت‌آمیز نبود: %v", err)
	} else {
		fmt.Println("اتصال به Redis برقرار شد!")
	}

	http.HandleFunc("/weather", Weather)
	err = http.ListenAndServe(":8082", nil)

}
