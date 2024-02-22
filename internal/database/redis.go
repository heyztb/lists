package database

import "github.com/redis/go-redis/v9"

var Redis *redis.Client

var RedisSessionKeyPrefix string = "session_key:%d"
