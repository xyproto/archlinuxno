package main

import (
	//"fmt"

	"github.com/garyburd/redigo/redis"
)

// Functions for dealing with a short list of string values in Redis

type RedisDatastructure struct {
	pool  *redis.Pool // connection pool
	id string
}

type RedisList RedisDatastructure
type RedisKeyValue RedisDatastructure
type RedisHashMap RedisDatastructure
type RedisSet RedisDatastructure

func NewRedisList(pool *redis.Pool, id string) *RedisList {
	return &RedisList{pool, id}
}

func NewRedisKeyValue(pool *redis.Pool, id string) *RedisKeyValue {
	return &RedisKeyValue{pool, id}
}

func NewRedisHashMap(pool *redis.Pool, id string) *RedisHashMap {
	return &RedisHashMap{pool, id}
}

func NewRedisSet(pool *redis.Pool, id string) *RedisSet {
	return &RedisSet{pool, id}
}

// Connect to the local instance of Redis at port 6379
func newRedisConnection() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func NewRedisConnectionPool() *redis.Pool {
	// The second argument is the maximum number of idle connections
	return redis.NewPool(newRedisConnection, 3)
}

func (rs *RedisSet) Add(value string) error {
	conn := rs.pool.Get()
	_, err := conn.Do("SADD", rs.id, value)
	return err
}

func (rs *RedisSet) Has(value string) (bool, error) {
	conn := rs.pool.Get()
	//fmt.Println("--- Has ---")
	//fmt.Println("command: SISMEMBER")
	//fmt.Println("fieldname:", rs.id)
	//fmt.Println("value:", value)
	retval, err := conn.Do("SISMEMBER", rs.id, value)
	//fmt.Println("retval:", retval)
	//fmt.Println("err:", err)
	if err != nil {
		//fmt.Println("noo")
		panic(err)
	}
	return redis.Bool(retval, err)
}

func (rs *RedisSet) GetAll() ([]string, error) {
	conn := rs.pool.Get()
	result, err := redis.Values(conn.Do("SMEMBERS", rs.id))
	strs := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		strs[i] = getString(result, i)
	}
	return strs, err

}

func (rs *RedisSet) Del(value string) error {
	conn := rs.pool.Get()
	_, err := conn.Do("SREM", rs.id, value)
	return err
}

func (rl *RedisList) Store(value string) error {
	conn := rl.pool.Get()
	_, err := conn.Do("RPUSH", rl.id, value)
	return err
}

func (rm *RedisKeyValue) Set(key, value string) error {
	conn := rm.pool.Get()
	_, err := conn.Do("SET", rm.id+":"+key, value)
	return err
}

func (rm *RedisKeyValue) Get(key string) (string, error) {
	conn := rm.pool.Get()
	result, err := redis.String(conn.Do("GET", rm.id+":"+key))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (rh *RedisHashMap) Set(hashkey, key, value string) error {
	conn := rh.pool.Get()
	_, err := conn.Do("HSET", rh.id+":"+hashkey, key, value)
	return err
}

func (rh *RedisHashMap) Get(hashkey, key string) (string, error) {
	conn := rh.pool.Get()
	result, err := redis.String(conn.Do("HGET", rh.id+":"+hashkey, key))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (rh *RedisHashMap) Del(hashkey, key string) error {
	conn := rh.pool.Get()
	_, err := conn.Do("HDEL", rh.id+":"+hashkey, key)
	return err
}

func bytes2string(b []uint8) string {
	return string(b)
	//return bytes.NewBuffer(b).String()
}

func getString(bi []interface{}, i int) string {
	return string(bi[i].([]uint8))
}

func (rl *RedisList) GetAll() ([]string, error) {
	conn := rl.pool.Get()
	result, err := redis.Values(conn.Do("LRANGE", rl.id, "0", "-1"))
	strs := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		strs[i] = getString(result, i)
	}
	return strs, err
}

func (rl *RedisList) GetLast() (string, error) {
	conn := rl.pool.Get()
	result, err := redis.Values(conn.Do("LRANGE", rl.id, "-1", "-1"))
	if len(result) == 1 {
		return getString(result, 0), err
	}
	return "", err
}

func (rl *RedisList) DelAll() error {
	conn := rl.pool.Get()
	_, err := conn.Do("DEL", rl.id)
	return err
}
