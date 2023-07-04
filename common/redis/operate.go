package redis

import (
	"github.com/gomodule/redigo/redis"
)

// lua脚本保证锁原子性
const SetLua = `
if redis.call('set', KEYS[1]) == ARGV[1] then
  return redis.call('expire', KEYS[1],ARGV[2]) 				
 else
   return '0' 					
end`

const GetLua = `
if redis.call('get', KEYS[1]) == ARGV[1] then
  return redis.call('expire', KEYS[1],ARGV[2]) 				
 else
   return '0' 					
end`

// 设置一个有过期时间的键值对
func (m *redisObj) SetTTL(key string, value string, ttl int64) error {
	conn := m.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	if _, err := conn.Do("setex", key, ttl, value); err != nil {
		return err
	}
	return nil
}

// 设置一个键值对
func (m *redisObj) Set(key string, value string) error {
	conn := m.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	if _, err := conn.Do("set", key, value); err != nil {
		return err
	}
	return nil
}

// 获取一个键的值
func (m *redisObj) GetString(key string) (string, error) {
	conn := m.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return "", err
	}
	return redis.String(conn.Do("get", key))
}

// 删除一个键
func (m *redisObj) Remove(key string) error {
	conn := m.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	if _, err := conn.Do("del", key); err != nil {
		return err
	}
	return nil
}

// 判断一个键是否存在
func (m *redisObj) Exists(key string) (bool, error) {
	conn := m.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, err
	}
	return redis.Bool(conn.Do("exists", key))
}

// 获取匹配的所有键
func (m *redisObj) GetKeys(patter string) ([]string, error) {
	conn := m.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return redis.Strings(conn.Do("key", patter))
}

// 删除匹配的所有键
func (m *redisObj) DelKeys(patter string) error {
	conn := m.pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	keys, err := redis.Strings(conn.Do("key", patter))
	if err != nil {
		return err
	}
	for _, key := range keys {
		conn.Do("del", key)
	}
	return nil
}
