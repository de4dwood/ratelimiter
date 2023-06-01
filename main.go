package main

import (
	"fmt"
	"os"
	"os/signal"
	"rt/data"
	"rt/domain"
	"rt/presentation"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisAddr     = os.Getenv("REDIS_ADDR")
	redisUsername = os.Getenv("REDIS_USERNAME")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDB       = os.Getenv("REDIS_DB")
	listen        = os.Getenv("RL_ADDRESS")
	ttl           = os.Getenv("RL_TLL")
)

func main() {
	rdb, err := strconv.Atoi(redisDB)
	if err != nil {
		rdb = 0
	}
	ttL, err := time.ParseDuration(ttl)
	if err != nil {
		ttL = time.Duration(time.Hour * 24)
	}
	redisOpt := &redis.Options{
		Addr:     redisAddr,
		Username: redisUsername,
		Password: redisPassword,
		DB:       rdb,
	}

	store := data.NewDataRedis(redisOpt, ttL)
	service := domain.NewDomain(store)
	present := presentation.NewHttp(service)
	sigs := make(chan os.Signal, 1)
	errChan := make(chan int, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {

		if listen == "" {
			present.Start(":9000")
		} else {
			present.Start(listen)
		}

		errChan <- 0
	}()
	for {
		select {
		case <-sigs:
			fmt.Println("shouting down the application")
			os.Exit(0)

		case <-errChan:
			fmt.Println("there is an error in application")
			os.Exit(1)
		}

	}

}
