package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"os"
	"reflect"
	"runtime/debug"
	"time"
)

var logger zerolog.Logger
var ctx = context.Background()

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	w := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	logger = zerolog.New(w).With().Timestamp().Caller().Logger()
}

type dbRedisStore struct {
	dbRedis *redis.Client
}

func dbRedis() *dbRedisStore {
	dbr := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "",
		DB:       9,
	})
	return &dbRedisStore{
		dbRedis: dbr,
	}
}

func (d *dbRedisStore) rdGet(key string) (result string, err error) {
	result, err = d.dbRedis.Get(ctx, key).Result()
	return
}

func main() {
	logger.Info().Msg("Start...")
	go run()

	select {}
}

func run() {
	rdb := dbRedis()
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		test(rdb)
		logger.Info().Str("test", "5s").Msg("Next time after")
		time.Sleep(5 * time.Second)
	}
}

func test(rdb *dbRedisStore) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			logger.Err(err).Str("debug", string(debug.Stack())).Msg("Error")
		}
	}()
	a := divide(1, 0)
	fmt.Println(a)
	val, err := rdb.rdGet("pair_setting:coin06:usdt")

	if err != nil {
		logger.Err(err).Msg("error redis")
	}
	var out map[string]any
	err = json.Unmarshal([]byte(val), &out)
	if err != nil {
		return
	}
	fmt.Println(reflect.TypeOf(out))
	fmt.Println(out["max_amount"])
	logger.Info().Msg(val)
}

func divide(a, b int) int {
	return a / b
}
